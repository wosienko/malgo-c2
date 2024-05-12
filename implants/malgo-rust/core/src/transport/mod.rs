use crate::command::CommandInfo;
pub type LastOne = bool;
pub trait Transport {
    fn register_session(&self) -> Result<(), String>;
    fn command_info(&self) -> Result<Option<CommandInfo>, String>;
    fn command_details(&self, command: &mut CommandInfo, offset: &mut u64) -> Result<LastOne, String>;
}

pub mod dns;