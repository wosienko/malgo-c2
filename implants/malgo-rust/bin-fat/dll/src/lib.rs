use windows_sys::Win32::UI::WindowsAndMessaging::MessageBoxA;
use windows_sys::Win32::Foundation::HINSTANCE;
use windows_sys::Win32::System::SystemServices::*;

fn attach() {
    unsafe {
        MessageBoxA(0, "Hello, Fat!\0".as_ptr(), "Hello\0".as_ptr(), 0);
    }
}

#[no_mangle]
#[allow(non_snake_case)]
fn DllMain(
    _hinstDLL: HINSTANCE,
    fdwReason: u32,
    _: *mut std::ffi::c_void,
) -> bool {
    match fdwReason {
        DLL_PROCESS_ATTACH => attach(),
        _ => (),
    }
    true
}