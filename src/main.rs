mod network;
mod constant;
mod respond;
mod util;

fn main() {
    println!("tally started, listening on {}", format!("{}:{}", constant::ADDR, constant::PORT));
    let _ = network::listen();
}