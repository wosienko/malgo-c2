package entities

type CommandInfo struct {
	Type   string `db:"type"`
	ID     string `db:"id"`
	Length int64  `db:"length"`
}
