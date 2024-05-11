fn main() {
    if !cfg!(target_os = "windows") {
        panic!("This build script is intended to run only on Windows.");
    }
}
