use std::io::{self, BufRead, BufReader};
use sha1::{Sha1, Digest};
use anyhow::{Context, Result};
use std::fs::File;

pub fn extract_paths(joined_paths: &String) -> Vec<&str> {
    let extracted: Vec<&str> = joined_paths.split(";").collect();
    extracted
}

fn get_first_line(path: &str) -> Result<String> {
    // try opening claim file
    let file: File = File::open(path).with_context(|| format!("Failed to open file {}", path))?;
    let mut reader: BufReader<File> = BufReader::new(file);
    let mut first_line = String::new();

    // return first line of claim file if it is not empty
    if reader.read_line(&mut first_line)? > 0 {
        Ok(first_line.trim_end().to_string())
    } else {
        Ok("".to_string())
    }
}

pub fn generate_challenge_hash(path: &str) -> Result<String> {
    // generate the challenge plaintext:
    let claim_content: String = get_first_line(path)?;
    if claim_content == "" {
        return Ok("".to_string());
    }
    let challenge_plaintext = claim_content + ":" + path;
    println!("{}", challenge_plaintext);

    // hash challenge with sha1
    let mut hasher = Sha1::new();
    hasher.update(challenge_plaintext.as_bytes());
    let challenge_hash = hasher.finalize();

    Ok(format!("{:x}", (&challenge_hash)))
}
