package entities

type ResultChunk struct {
	CommandId string `db:"command_id"`
	Chunk     []byte `db:"chunk"`
	Offset    int    `db:"offset"`
}
