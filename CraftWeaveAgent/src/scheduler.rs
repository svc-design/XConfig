// File: src/scheduler.rs

use crate::{config, executor, result_store};
use crate::config::AgentConfig;
use crate::models::Play;
use tokio::time::{sleep, Duration};

pub async fn run_schedule(agent_config: &AgentConfig) -> anyhow::Result<()> {
    loop {
        let content = config::fetch_git_and_load_playbook(&agent_config.repo, "sync/playbook.yaml").await?;
        let parsed: Vec<Play> = serde_yaml::from_str(&content)?;
        let results = executor::run(parsed).await?;
        result_store::persist(results).await?;

        let interval = agent_config.interval.unwrap_or(60);
        sleep(Duration::from_secs(interval)).await;
    }
}
