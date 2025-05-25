// File: src/scheduler.rs
// ----------------------
use crate::{config, executor, result_store};
use tokio::time::{sleep, Duration};

pub async fn run_schedule(once: bool) -> anyhow::Result<()> {
    if once {
        let cfg = config::fetch().await?;
        let results = executor::apply(cfg).await?;
        result_store::persist(results).await?;
    } else {
        loop {
            let cfg = config::fetch().await?;
            let results = executor::apply(cfg).await?;
            result_store::persist(results).await?;
            sleep(Duration::from_secs(60)).await;
        }
    }
    Ok(())
}
