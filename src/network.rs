use std::{io::{Read, Write}, net::{TcpListener, TcpStream}};
use crate::{constant, respond::generate_response};
use eyre::Result;

pub fn listen() -> Result<()> {
    let listener: TcpListener = TcpListener::bind(format!("{}:{}", constant::ADDR, constant::PORT))?;
    for stream in listener.incoming() {
        handle_connection(&mut stream?);
    };
    Ok(())
}

pub fn handle_connection(stream: &mut TcpStream) {
    let mut buffer = [0; 1024];
    while let Ok(n) = stream.read(&mut buffer) {
        if n == 0 { break; }
        let received_string = String::from_utf8_lossy(&buffer[..n]);
        // println!("Received: {}", received_string);

        let response = generate_response(received_string.trim().to_string());
        // println!("{}",response);
        let _ = stream.write_all(response.as_bytes());
    }
}