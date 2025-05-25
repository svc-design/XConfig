// File: src/config.rs

use serde::Deserialize;
use std::fs;
use std::path::Path;
use std::process::Command;
use tokio::fs as tokio_fs;

#[derive(Debug, Deserialize)]
pub struct AgentConfig {
    pub repo: String,
    pub interval: Option<u64>, // in seconds
}

pub async fn load_agent_config(path: &str) -> anyhow::Result<AgentConfig> {
    let content = tokio_fs::read_to_string(path).await?;
    let config: AgentConfig = serde_yaml::from_str(&content)?;
    Ok(config)
}

/// 克隆 Git 仓库并读取 playbook.yaml 内容（作为字符串返回）
pub async fn fetch_git_and_load_playbook(repo: &str, subpath: &str) -> anyhow::Result<String> {
    let tmp_dir = "/tmp/cw-agent-sync";
    let full_path = format!("{}/{}", tmp_dir, subpath);

    // 清理旧目录
    let _ = fs::remove_dir_all(tmp_dir);
    fs::create_dir_all(tmp_dir)?;

    let status = Command::new("git")
        .args(["clone", "--depth", "1", repo, tmp_dir])
        .status()?;

    if !status.success() {
        anyhow::bail!("git clone failed with code: {:?}", status.code());
    }

    if !Path::new(&full_path).exists() {
        anyhow::bail!("playbook not found: {}", full_path);
    }

    let content = tokio_fs::read_to_string(&full_path).await?;
    Ok(content)
}

