use core::net::{SocketAddr};

use alloc::format;

use hickory_resolver::{Resolver};
use hickory_resolver::config::{NameServerConfig, ResolverConfig, ResolverOpts, Protocol};
use rand::Rng;

use crate::transport::Transport;
use crate::session::Session;

pub struct DNS<'a> {
    pub server: &'a str,
    pub domain: &'a str,

    resolver: Resolver,
}

// Generate a random blob of data of the same length as the target
fn generate_blob(mut target: alloc::vec::Vec<u8>) -> alloc::vec::Vec<u8> {
    let mut rng = rand::thread_rng();
    for i in 0..target.len() {
        match rng.gen_range(0..3) {
            0 => target[i] = rng.gen_range(48..58), // fill with numbers
            1 => target[i] = rng.gen_range(65..91), // fill with uppercase letters
            _ => target[i] = rng.gen_range(97..123), // fill with lowercase letters
        }
    }
    target
}

impl Transport for DNS<'_> {
    fn register_session(&self, session: Session) -> Result<(), &str> {
        // Generate random value so that the DNS query is not cached
        let mut blob = alloc::vec![0; 4];
        blob = generate_blob(blob);
        let blob = core::str::from_utf8(&blob).unwrap();

        let query = format!(
            "{}.{}.{}.{}",
            session.project_id,
            session.id,
            blob,
            self.domain
        );

        return match self.resolver.ipv6_lookup(query) {
            Ok(_) => {
                Ok(())
            }
            Err(_) => {
                Err("Could not register session")
            }
        }
    }
}

impl DNS<'_> {
    pub fn new<'a>(server: &'a str, domain: &'a str) -> DNS<'a> {
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