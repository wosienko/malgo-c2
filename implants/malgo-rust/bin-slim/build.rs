fn main() {
    if cfg!(target_os = "windows") {
        println!("cargo:rustc-link-arg=/ENTRY:_start");
        println!("cargo:rustc-link-arg=/SUBSYSTEM:windows");
    }
}