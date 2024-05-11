fn main() {
    // Linker arguments for building from Windows
    if cfg!(target_env = "msvc") {
        println!("cargo:rustc-link-arg=/ENTRY:WinMain");
        println!("cargo:rustc-link-arg=/SUBSYSTEM:windows");
    }
    // Linker arguments for GNU MinGW cross-compilation
    else {
        println!("cargo:rustc-link-arg=-e WinMain");
    }
}
