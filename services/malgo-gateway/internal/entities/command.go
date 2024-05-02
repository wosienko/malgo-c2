package entities

type CommandInfo struct {
	Type   string `db:"type"`
	ID     string `db:"id"`
	Length int64  `db:"length"`
}

type CommandChunkQuery struct {
	CommandID string `db:"command_id"`
	Offset    int    `db:"offset"`
	Length    int    `db:"length"`
}

type CommandChunk struct {
	CommandID string `db:"id"`
	SessionID string `db:"session_id"`
	Offset    int    `db:"offset"`
	Data      string `db:"data"`
	Length    int    `db:"length"`
	IsLast    bool   `db:"is_last"`
}
