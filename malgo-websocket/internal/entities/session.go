package entities

type SessionKeyValue struct {
	SessionId string `json:"session_id" db:"session_id"`
	Key       string `json:"key" db:"key"`
	Value     string `json:"value" db:"value"`
}

type SessionKeyValueSentToOperator struct {
	MessageType string `json:"message_type"`
	SessionId   string `json:"session_id" db:"session_id"`
	Key         string `json:"key" db:"key"`
	Value       string `json:"value" db:"value"`
}

type SessionName struct {
	SessionId string `json:"session_id" db:"session_id"`
	Name      string `json:"name" db:"name"`
}

type SessionNameSentToOperator struct {
	MessageType string `json:"message_type"`
	SessionId   string `json:"session_id" db:"session_id"`
	Name        string `json:"name" db:"name"`
}
