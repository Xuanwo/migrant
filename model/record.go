package model

// Record is the record for a migration.
type Record struct {
	ID        string `db:"id"`
	Type      string `db:"type"`
	AppliedAt int64  `db:"applied_at"`
}

// ByID sorts migrations by ID.
type ByID []*Record

func (b ByID) Len() int           { return len(b) }
func (b ByID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByID) Less(i, j int) bool { return b[i].ID < b[j].ID }
