// File: src/scheduler.rs

use crate::{config, executor, result_store};
use crate::config::AgentConfig;
use tokio::time::{sleep, Duration};

pub async fn run_schedule(once: bool, agent_config: &AgentConfig) -> anyhow::Result<()> {
    if once {
        let cfg = config::fetch_git_and_load_yaml(&agent_config.repo, "sync/config.yaml").await?;
        let results = executor::apply(cfg).await?;
        result_store::persist(results).await?;
        return Ok(());
    }

    loop {
        let cfg = config::fetch_git_and_load_yaml(&agent_config.repo, "sync/config.yaml").await?;
        let results = executor::apply(cfg).await?;
        result_store::persist(results).await?;
        let interval = agent_config.interval.unwrap_or(60);
        sleep(Duration::from_secs(interval)).await;
    }
}
