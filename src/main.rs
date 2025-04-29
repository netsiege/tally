use std::io::{self, BufRead, BufReader};
use sha1::{Sha1, Digest};
use anyhow::Result;
use std::fs::File;

fn get_first_line(path: &str) -> Result<String> {
    // try opening control file
    let file: File = File::open(path)?;
    let mut reader: BufReader<File> = BufReader::new(file);
    let mut first_line = String::new();

    // return first line of control file if it is not empty
    if reader.read_line(&mut first_line)? > 0 {
        Ok(first_line.trim_end().to_string())
    } else {
        Err(io::Error::new(io::ErrorKind::Other, "File is empty").into())
    }
}

fn generate_challenge_hash(path: &str) -> Result<[u8; 20]> {
    // generate the challenge plaintext:
    let control_content: String = get_first_line(path)?;
    let challenge_plaintext = control_content + ":" + path;
    println!("{}", challenge_plaintext);

    // hash challenge with sha1
    let mut hasher = Sha1::new();
    hasher.update(challenge_plaintext.as_bytes());
    let challenge_hash = hasher.finalize();

    Ok(challenge_hash.into())
}

fn main() {
    println!("{:?}", get_first_line("control.txt"));
    println!("{:?}", generate_challenge_hash("control.txt"));
}
