package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pengsrc/go-shared/check"
	"github.com/pengsrc/go-shared/convert"
	"gopkg.in/urfave/cli.v1"

	"github.com/Xuanwo/migrant/config"
	"github.com/Xuanwo/migrant/constants"
	"github.com/Xuanwo/migrant/migration"
)

func up(ctx *cli.Context) error {
	c, _ := config.New()
	err := c.LoadFromFilePath(ctx.GlobalString("config"))
	if err != nil {
		return err
	}

	err = migration.Setup(c)
	if err != nil {
		return err
	}

	id, err := migration.Up()
	if err != nil {
		return err
	}
	if id != "" {
		fmt.Printf("Run migration complete, %s.\n", id)
	} else {
		fmt.Println("No migration executed.")
	}

	return nil
}

func down(ctx *cli.Context) error {
	c, _ := config.New()
	err := c.LoadFromFilePath(ctx.GlobalString("config"))
	if err != nil {
		return err
	}

	err = migration.Setup(c)
	if err != nil {
		return err
	}

	id, err := migration.Down()
	if err != nil {
		return err
	}
	if id != "" {
		fmt.Printf("Revert migration complete, %s.\n", id)
	} else {
		fmt.Println("No migration executed.")
	}

	return nil
}

func status(ctx *cli.Context) error {
	c, _ := config.New()
	err := c.LoadFromFilePath(ctx.GlobalString("config"))
	if err != nil {
		return err
	}

	err = migration.Setup(c)
	if err != nil {
		return err
	}

	ids, err := migration.Status()
	if err != nil {
		return err
	}

	fmt.Println("Applied At               Migration")
	fmt.Println(strings.Repeat("=", 80))

	for _, v := range ids {
		if v.AppliedAt != 0 {
			appliedAt := convert.TimestampToString(v.AppliedAt, convert.ISO8601)
			fmt.Printf("%s  -  %s\n", appliedAt, v.ID)
			continue
		}
		fmt.Printf("Pending               -  %s\n", v.ID)
	}

	return nil
}

func sync(ctx *cli.Context) error {
	c, _ := config.New()
	err := c.LoadFromFilePath(ctx.GlobalString("config"))
	if err != nil {
		return err
	}

	err = migration.Setup(c)
	if err != nil {
		return err
	}

	ids, err := migration.Sync()
	if err != nil {
		return err
	}

	s := make([]string, len(ids))
	for k, v := range ids {
		s[k] = fmt.Sprintf("- %s\n", v.ID)
	}
	if len(ids) > 0 {
		fmt.Printf("Migrated to the latest database schema.\n%s", strings.Join(s, ""))
	} else {
		fmt.Println("Already has the latest database schema.")
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = constants.Name
	app.Usage = constants.Usage
	app.Version = constants.Version
	app.Commands = []cli.Command{
		{
			Name:   "up",
			Usage:  "Run one migration",
			Action: up,
		},
		{
			Name:   "down",
			Usage:  "Revert one migration",
			Action: down,
		},
		{
			Name:   "sync",
			Usage:  "Sync database schema",
			Action: sync,
		},
		{
			Name:   "status",
			Usage:  "Show database schema status",
			Action: status,
		},
	}

	// Setup flags.
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))

	check.ErrorForExit(constants.Name, app.Run(os.Args))
}
