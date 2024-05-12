#![no_std]
#![no_main]

#[cfg(not(windows))]
compile_error!("This program is only designed to run on Windows");

#[allow(unused)]
use bin_slim::essentials::*;

use malgo_rust::session::Session;
use malgo_rust::transport::dns::DNS;
use malgo_rust::transport::Transport;

const SESSION: Session = Session {
    id: "336df889-c0e9-453a-8675-8bff4176e1b0",
    project_id: "38b6d3a1-4373-4202-b9d7-b11d399d0ebf"
};

#[no_mangle]
pub extern "C" fn WinMain(
    _h_instance: u32,
    _h_prev_instance: u32,
    _lp_cmd_line: *const u8,
    _n_cmd_show: i32,
) -> u32 {
    let dns = DNS::new("127.0.0.1:53", "a.example.com");

    match dns.register_session(SESSION) {
        Ok(_) => 0,
        Err(_) => panic!(),
    }
    return 0;
}
