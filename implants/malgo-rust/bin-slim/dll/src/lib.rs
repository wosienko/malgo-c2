#![no_std]
#![no_main]

#[allow(unused)]
use bin_slim::essentials::*;

use windows_sys::Win32::UI::WindowsAndMessaging::MessageBoxA;
use windows_sys::Win32::Foundation::HINSTANCE;
use windows_sys::Win32::System::SystemServices::*;

fn attach() {
    unsafe {
        MessageBoxA(0, "Hello, Slim!\0".as_ptr(), "Hello\0".as_ptr(), 0);
    }
}

#[no_mangle]
#[allow(non_snake_case)]
pub extern "C" fn _DllMainCRTStartup(
    _hinstDLL: HINSTANCE,
    fdwReason: u32,
    _: *mut core::ffi::c_void,
) -> bool {
    match fdwReason {
        DLL_PROCESS_ATTACH => attach(),
        _ => (),
    }
    true
}