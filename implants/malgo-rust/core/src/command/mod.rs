use serde::Deserialize;

#[derive(Deserialize, Debug)]
pub struct CommandInfo {
    pub command_id: String,
    #[serde(rename(deserialize = "type"))]
    pub command_type: String, // TODO: switch to Enum
    pub command_length: u64,
    #[serde(skip_deserializing)]
    pub command: String
}

#[derive(Deserialize)]
pub struct CommandDetails {
    pub data: String,
    pub is_last_chunk: bool
}