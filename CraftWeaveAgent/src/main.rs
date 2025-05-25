// File: src/main.rs
mod config;
mod executor;
mod result_store;
mod scheduler;

use clap::{Parser, Subcommand};
use std::path::PathBuf;
use crate::scheduler::run_schedule;
use crate::config::{load_agent_config, load_from_file};

#[derive(Parser, Debug)]
#[command(name = "cw-agent", version)]
#[command(about = "CraftWeave Agent - lightweight remote playbook executor")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand, Debug)]
enum Commands {
    /// Run once using config from Git repo
    Oneshot,
    /// Run as daemon with interval from config file
    Daemon,
    /// Apply local config file
    Apply {
        #[arg(short, long)]
        file: PathBuf,
    },
    /// Print latest execution result from local store
    Status,
    /// Show version info
    Version,
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let args = Cli::parse();

    let agent_config = load_agent_config("/etc/cw-agent.conf").await.unwrap_or_else(|e| {
        eprintln!("⚠️ Failed to load config: {}", e);
        std::process::exit(1);
    });

    match args.command {
        Commands::Oneshot => {
            run_schedule(true, &agent_config).await?;
        }
        Commands::Daemon => {
            run_schedule(false, &agent_config).await?;
        }
        Commands::Apply { file } => {
            let cfg = load_from_file(file).await?;
            let results = executor::apply(cfg).await?;
            result_store::persist(results).await?;
        }
        Commands::Status => {
            result_store::print_latest().await?;
        }
        Commands::Version => {
            println!("cw-agent version {}", env!("CARGO_PKG_VERSION"));
        }
    }

    Ok(())
}
