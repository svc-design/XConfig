// File: src/scheduler.rs
// ----------------------
use crate::{config, executor, result_store};
use tokio::time::{sleep, Duration};

pub async fn run_schedule(once: bool) -> anyhow::Result<()> {
    loop {
        let cfg = config::fetch().await?;
        let results = executor::apply(cfg).await?;
        result_store::persist(results).await?;

        if once {
            break;
        }
        sleep(Duration::from_secs(60)).await;
    }
    Ok(())
}
