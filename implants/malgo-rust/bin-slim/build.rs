fn main() {
    // Linker arguments for building from Windows
    if cfg!(target_os = "windows") {
        println!("cargo:rustc-link-arg=/ENTRY:_start");
        println!("cargo:rustc-link-arg=/SUBSYSTEM:windows");
    }
    // Linker arguments for GNU MinGW cross-compilation
    else {
        println!("cargo:rustc-link-arg=-e _start");
        println!("cargo:rustc-link-arg=-Tlink.x");
    }
}
