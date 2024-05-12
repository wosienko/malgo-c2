use crate::session::Session;
pub trait Transport {
    fn register_session(&self, session: Session) -> Result<(), String>;
}

pub mod dns;