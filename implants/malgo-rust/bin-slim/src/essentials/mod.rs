#[panic_handler]
pub fn panic(_info: &core::panic::PanicInfo) -> ! {
    unsafe {
        windows_sys::Win32::System::Threading::ExitProcess(2137);
    }
}

#[no_mangle]
pub extern "C" fn memcpy(dest: *mut u8, src: *mut u8, n: usize) {
    unsafe {
        for i in 0..n {
            *dest.add(i) = *src.add(i);
        }
    }
}

#[no_mangle]
pub extern "C" fn memset(s: *mut u8, c: u8, n: usize) {
    unsafe {
        for i in 0..n {
            *s.add(i) = c;
        }
    }
}

#[no_mangle]
pub extern "C" fn memcmp(s1: *const u8, s2: *const u8, n: usize) -> i32 {
    unsafe {
        for i in 0..n {
            if *s1.add(i) != *s2.add(i) {
                return *s1.add(i) as i32 - *s2.add(i) as i32;
            }
        }
    }
    return 0;
}

#[no_mangle]
#[allow(non_upper_case_globals)]
pub static _fltused: u32 = 0;
