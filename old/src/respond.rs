use crate::util::{extract_paths, generate_challenge_hash};
use rand::{distributions::Alphanumeric, Rng};

pub fn generate_response(received_response: String) -> String {
    let paths = extract_paths(&received_response);
    let hashes: Result<Vec<String>, anyhow::Error> = paths.iter().map(|path| generate_challenge_hash(path) ).collect::<Result<Vec<_>, _>>();
    let hashes = hashes.unwrap_or_default();

    // If hashes is empty or any hash is empty, return a random string
    if hashes.is_empty() || hashes.iter().any(|h| h.is_empty()) {
        let rand_str: String = rand::thread_rng()
            .sample_iter(&Alphanumeric)
            .take(16)
            .map(char::from)
            .collect();
        rand_str
    } else {
        hashes.join(";")
    }
}
