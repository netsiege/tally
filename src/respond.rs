use crate::util::{extract_paths, generate_challenge_hash};

pub fn generate_response(received_response: String) -> String {
    let paths = extract_paths(&received_response);
    let hashes = paths.iter().map(|path| generate_challenge_hash(path) ).collect::<Result<Vec<_>, _>>();
    hashes.unwrap().join(";")
}
