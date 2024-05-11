use std::panic::catch_unwind;
use malgo_rust::mathy::add::*;

fn main() {
    println!("EXE Fat: {}", add(2, 2));

    // panic and catch it, logging it
    catch_unwind(|| {
        panic!("Fatality!");
    }).ok();
}
