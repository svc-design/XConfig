// File: src/main.rs

mod config;
mod executor;
mod result_store;
mod scheduler;
mod models;

use clap::{Parser, Subcommand};
use std::path::PathBuf;
use crate::executor::run as run_playbook;
use crate::config::load_agent_config;
use tokio::fs;

#[derive(Parser, Debug)]
#[command(name = "cw-agent", version)]
#[command(about = "CraftWeave Agent - lightweight local playbook runner")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand, Debug)]
enum Commands {
    /// Run once using playbook(s) from Git repo
    Oneshot,
    /// Run as daemon with interval from config file
    Daemon,
    /// Run full playbook from file (array of plays)
    Playbook {
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
            for path in &agent_config.playbook {
                let content = config::fetch_git_and_load_playbook(&agent_config.repo, path).await?;
                let parsed: Vec<models::Play> = serde_yaml::from_str(&content)?;
                let results = run_playbook(parsed).await?;
                result_store::persist(results).await?;
            }
        }
        Commands::Daemon => {
            scheduler::run_schedule(&agent_config).await?;
        }
        Commands::Playbook { file } => {
            let content = fs::read_to_string(file).await?;
            let parsed: Vec<models::Play> = serde_yaml::from_str(&content)?;
            let results = run_playbook(parsed).await?;
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
