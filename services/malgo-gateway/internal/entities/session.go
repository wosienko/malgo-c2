package entities

type RegisterNewSession struct {
	SessionId string `db:"session_id"`
	ProjectId string `db:"project_id"`
}
