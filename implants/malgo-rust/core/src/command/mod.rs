use serde::Deserialize;

#[derive(Deserialize, Debug)]
pub struct CommandInfo {
    command_id: String,
    #[serde(rename(deserialize = "type"))]
    command_type: String, // TODO: switch to Enum
    command_length: u64
}