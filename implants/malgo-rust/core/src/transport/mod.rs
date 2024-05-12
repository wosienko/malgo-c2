use crate::command::CommandInfo;
pub trait Transport {
    fn register_session(&self) -> Result<(), String>;
    fn command_info(&self) -> Result<Option<CommandInfo>, String>;
}

pub mod dns;