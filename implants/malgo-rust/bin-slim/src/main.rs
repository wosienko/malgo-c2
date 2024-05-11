#![no_std]
#![no_main]

#[cfg(not(windows))]
compile_error!("This program is only designed to run on Windows");

#[allow(unused)]
use bin_slim::essentials::*;

#[no_mangle]
pub extern "C" fn _start() -> u32 {
    panic!()
}
