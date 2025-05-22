// File: src/main.rs
// ------------------
mod config;
mod executor;
mod result_store;
mod scheduler;

use crate::scheduler::run_schedule;
use clap::Parser;

#[derive(Parser, Debug)]
#[command(name = "cw-agent", version)]
struct Cli {
    /// Run once and exit (for testing)
    #[arg(short, long)]
    oneshot: bool,
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let args = Cli::parse();
    run_schedule(args.oneshot).await
}
