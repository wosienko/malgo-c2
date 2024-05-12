use crate::session::Session;
pub trait Transport {
    fn register_session(&self, session: Session) -> Result<(), &str>;
}

pub mod dns;