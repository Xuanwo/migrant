package sqlparse

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	sqlCmdPrefix = "-- +migrate "
)

// ParsedMigration is the migration that parsed.
type ParsedMigration struct {
	UpStatements   []string
	DownStatements []string
}

var (
	// LineSeparator can be used to split migrations by an exact line match. This line
	// will be removed from the output. If left blank, it is not considered. It is defaulted
	// to blank so you will have to set it manually.
	// Use case: in MSSQL, it is convenient to separate commands by GO statements like in
	// SQL Query Analyzer.
	LineSeparator = ""
)

func errNoTerminator() error {
	if len(LineSeparator) == 0 {
		return errors.New(`last statement must be ended by a semicolon`)
	}

	return fmt.Errorf(`last statement must be ended by a semicolon or a line whose contents are %q`, LineSeparator)
}

// Checks the line to see if the line has a statement-ending semicolon
// or if the line contains a double-dash comment.
func endsWithSemicolon(line string) bool {

	prev := ""
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		if strings.HasPrefix(word, "--") {
			break
		}
		prev = word
	}

	return strings.HasSuffix(prev, ";")
}

type migrationDirection int

const (
	directionNone migrationDirection = iota
	directionUp
	directionDown
)

type migrateCommand struct {
	Command string
	Options []string
}

func (c *migrateCommand) HasOption(opt string) bool {
	for _, specifiedOption := range c.Options {
		if specifiedOption == opt {
			return true
		}
	}

	return false
}

func parseCommand(line string) (*migrateCommand, error) {
	cmd := &migrateCommand{}

	if !strings.HasPrefix(line, sqlCmdPrefix) {
		return nil, errors.New("not a sql-migrate command")
	}

	fields := strings.Fields(line[len(sqlCmdPrefix):])
	if len(fields) == 0 {
		return nil, errors.New(`incomplete migration command`)
	}

	cmd.Command = fields[0]

	cmd.Options = fields[1:]

	return cmd, nil
}

// ParseMigration will split the given sql script into individual statements.
//
// The base case is to simply split on semicolons, as these
// naturally terminate a statement.
//
// However, more complex cases like pl/pgsql can have semicolons
// within a statement. For these cases, we provide the explicit annotations
// 'StatementBegin' and 'StatementEnd' to allow the script to
// tell us to ignore semicolons.
func ParseMigration(r io.ReadSeeker) (*ParsedMigration, error) {
	p := &ParsedMigration{}

	_, err := r.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	scanner := bufio.NewScanner(r)

	statementEnded := false
	currentDirection := directionNone

	for scanner.Scan() {
		line := scanner.Text()

		// handle any migrate-specific commands
		if strings.HasPrefix(line, sqlCmdPrefix) {
			cmd, err := parseCommand(line)
			if err != nil {
				return nil, err
			}

			switch cmd.Command {
			case "Up":
				if len(strings.TrimSpace(buf.String())) > 0 {
					return nil, errNoTerminator()
				}
				currentDirection = directionUp
				break

			case "Down":
				if len(strings.TrimSpace(buf.String())) > 0 {
					return nil, errNoTerminator()
				}
				currentDirection = directionDown
				break
			}
		}

		if currentDirection == directionNone {
			continue
		}

		isLineSeparator := len(LineSeparator) > 0 && line == LineSeparator

		if !isLineSeparator {
			if _, err := buf.WriteString(line + "\n"); err != nil {
				return nil, err
			}
		}

		// Wrap up the two supported cases: 1) basic with semicolon; 2) psql statement
		// Lines that end with semicolon that are in a statement block
		// do not conclude statement.
		if (endsWithSemicolon(line) || isLineSeparator) || statementEnded {
			statementEnded = false
			switch currentDirection {
			case directionUp:
				p.UpStatements = append(p.UpStatements, buf.String())

			case directionDown:
				p.DownStatements = append(p.DownStatements, buf.String())

			default:
				panic("impossible state")
			}

			buf.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if currentDirection == directionNone {
		return nil, errors.New(`no Up/Down annotations found, no statements were executed`)
	}

	if len(strings.TrimSpace(buf.String())) > 0 {
		return nil, errNoTerminator()
	}

	return p, nil
}
