// File: src/config.rs
// -------------------
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct TaskConfig {
    pub name: String,
    pub command: String,
}

pub async fn fetch() -> anyhow::Result<Vec<TaskConfig>> {
    // Placeholder: fetch from Git or HTTP
    Ok(vec![
        TaskConfig {
            name: "echo hello".into(),
            command: "echo hello world".into(),
        },
    ])
}
