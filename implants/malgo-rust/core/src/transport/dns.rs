use core::net::{SocketAddr};

use hickory_resolver::{Resolver};
use hickory_resolver::config::{NameServerConfig, ResolverConfig, ResolverOpts, Protocol};
use rand::Rng;

use crate::transport::Transport;
use crate::session::Session;

pub struct DNS {
    pub server: String,
    pub domain: String,

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
    fn register_session(&self, session: Session) -> Result<(), String> {
        let query = format!(
            "{}.{}.{}.{}",
            session.project_id,
            session.id,
            generate_blob(4),
            self.domain
        );

        return match self.resolver.ipv6_lookup(query) {
            Ok(_) => {
                Ok(())
            }
            Err(_) => {
                Err("Could not register session".to_string())
            }
        }
    }
}

impl DNS {
    pub fn new<'a>(server: String, domain: String) -> DNS {
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
            resolver,
        }
    }
}