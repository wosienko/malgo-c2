use malgo_rust::session::Session;
use malgo_rust::transport::dns::*;
use malgo_rust::transport::Transport;

const SESSION: Session = Session {
    id: "336df889-c0e9-453a-8675-8bff4176e1b0",
    project_id: "38b6d3a1-4373-4202-b9d7-b11d399d0ebf"
};

fn main() {
    let dns = DNS::new("127.0.0.1:53".to_string(), "a.example.com".to_string());

    match dns.register_session(SESSION) {
        Ok(_) => println!("Session registered successfully"),
        Err(e) => println!("Error registering session: {}", e),
    }
}
