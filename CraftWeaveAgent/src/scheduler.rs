// File: src/scheduler.rs

use crate::{config, executor, result_store};
use crate::config::AgentConfig;
use crate::models::Play;
use tokio::time::{sleep, Duration};

pub async fn run_schedule(agent_config: &AgentConfig) -> anyhow::Result<()> {
    loop {
        let mut all_results = vec![];

        for path in &agent_config.playbook {
            println!("📦 Fetching and executing: {}", path);
            match config::fetch_git_and_load_playbook(&agent_config.repo, path).await {
                Ok(content) => {
                    match serde_yaml::from_str::<Vec<Play>>(&content) {
                        Ok(parsed) => {
                            match executor::run(parsed).await {
                                Ok(results) => all_results.extend(results),
                                Err(e) => eprintln!("❌ Executor error [{}]: {}", path, e),
                            }
                        }
                        Err(e) => eprintln!("❌ YAML parse error [{}]: {}", path, e),
                    }
                }
                Err(e) => eprintln!("❌ Failed to fetch [{}]: {}", path, e),
            }
        }

        result_store::persist(all_results).await?;

        let interval = agent_config.interval.unwrap_or(60);
        println!("🕒 Sleeping {}s before next run...\n", interval);
        sleep(Duration::from_secs(interval)).await;
    }
}
