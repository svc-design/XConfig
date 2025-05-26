// File: src/config.rs

use serde::Deserialize;
use std::fs;
use std::path::Path;
use std::process::Command;
use tokio::fs as tokio_fs;

#[derive(Debug, Deserialize)]
pub struct AgentConfig {
    /// Git 仓库地址
    pub repo: String,

    /// 轮询间隔（仅在 daemon 模式生效）
    pub interval: Option<u64>,

    /// 要执行的 playbook 路径列表（相对于仓库根目录）
    pub playbook: Vec<String>,
}

/// 加载本地 agent 配置文件（例如 /etc/cw-agent.conf）
pub async fn load_agent_config(path: &str) -> anyhow::Result<AgentConfig> {
    let content = tokio_fs::read_to_string(path).await?;
    let config: AgentConfig = serde_yaml::from_str(&content)?;
    Ok(config)
}

/// 克隆 Git 仓库并读取指定路径的 playbook.yaml 内容（作为字符串返回）
pub async fn fetch_git_and_load_playbook(repo: &str, subpath: &str) -> anyhow::Result<String> {
    let tmp_dir = "/tmp/cw-agent-sync";
    let full_path = format!("{}/{}", tmp_dir, subpath);

    // 每次拉取前清理临时目录
    let _ = fs::remove_dir_all(tmp_dir);
    fs::create_dir_all(tmp_dir)?;

    // 克隆 repo（浅拷贝）
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

