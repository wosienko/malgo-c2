#[cfg(not(windows))]
compile_error!("This program is only designed to run on Windows");

use malgo_rust::mathy::add::*;
use std::panic::catch_unwind;

fn main() {
    println!("EXE Fat: {}", add(2, 2));

    // panic and catch it, logging it
    catch_unwind(|| {
        panic!("Fatality!");
    })
    .ok();
}
