fn main() {
    // Panic if not building for Windows
    if !cfg!(target_os = "windows") {
        panic!("This build script is intended to run only on Windows.");
    }

    // Linker arguments for MSVC
    if cfg!(target_os = "windows") && cfg!(target_env = "msvc") {
        println!("cargo:rustc-link-arg=/ENTRY:_start");
        println!("cargo:rustc-link-arg=/SUBSYSTEM:windows");
    }
    // Linker arguments for GNU MinGW
    else if cfg!(target_os = "windows") && cfg!(target_env = "msvc") {
        println!("cargo:rustc-link-arg=-e _start");
        println!("cargo:rustc-link-arg=-Tlink.x");
    } else {
        panic!("Unsupported target environment.");
    }
}
