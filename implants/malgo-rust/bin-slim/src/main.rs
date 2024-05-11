#![no_std]
#![no_main]

#[cfg(not(windows))]
compile_error!("This program is only designed to run on Windows");

#[allow(unused)]
use bin_slim::essentials::*;

#[no_mangle]
pub extern "C" fn WinMain(
    _h_instance: u32,
    _h_prev_instance: u32,
    _lp_cmd_line: *const u8,
    _n_cmd_show: i32,
) -> u32 {
    panic!()
}
