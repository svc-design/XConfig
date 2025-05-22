// File: src/result_store.rs
// -------------------------
use chrono::Utc;
use serde::Serialize;
use std::fs;
use std::path::Path;

#[derive(Debug, Serialize)]
pub struct CommandResult {
    pub task: String,
    pub stdout: String,
    pub stderr: String,
    pub success: bool,
    pub return_code: i32,
}

pub async fn persist(results: Vec<CommandResult>) -> anyhow::Result<()> {
    let json = serde_json::to_string_pretty(&results)?;
    let ts = Utc::now().format("%Y%m%d%H%M%S");
    let path = format!("/var/lib/cw-agent/status-{}.json", ts);
    fs::create_dir_all(Path::new("/var/lib/cw-agent"))?;
    fs::write(path, json)?;
    Ok(())
}
