// File: src/config.rs
use serde::Deserialize;
use std::process::Command;
use std::fs;
use std::path::Path;
use tokio::fs as tokio_fs;

#[derive(Debug, Deserialize)]
pub struct TaskConfig {
    pub name: String,
    pub command: String,
}

pub async fn fetch() -> anyhow::Result<Vec<TaskConfig>> {
    fetch_git_and_load_yaml(
        "https://github.com/svc-design/gitops.git",
        "sync/config.yaml",
    ).await
}

/// 克隆仓库并解析指定路径的 YAML 配置
pub async fn fetch_git_and_load_yaml(repo: &str, subpath: &str) -> anyhow::Result<Vec<TaskConfig>> {
    let tmp_dir = "/tmp/cw-agent-sync";
    let full_path = format!("{}/{}", tmp_dir, subpath);

    // 清理旧目录（或优化为 pull）
    let _ = fs::remove_dir_all(tmp_dir);
    fs::create_dir_all(tmp_dir)?;

    let status = Command::new("git")
        .args(["clone", "--depth", "1", repo, tmp_dir])
        .status()?;

    if !status.success() {
        anyhow::bail!("git clone failed with code: {:?}", status.code());
    }

    if !Path::new(&full_path).exists() {
        anyhow::bail!("config file not found: {}", full_path);
    }

    let content = tokio_fs::read_to_string(&full_path).await?;
    let tasks: Vec<TaskConfig> = serde_yaml::from_str(&content)?;
    Ok(tasks)
}

/// 本地模式，从文件加载 YAML 任务配置
pub async fn load_from_file<P: AsRef<Path>>(path: P) -> anyhow::Result<Vec<TaskConfig>> {
    let content = tokio_fs::read_to_string(path).await?;
    let tasks: Vec<TaskConfig> = serde_yaml::from_str(&content)?;
    Ok(tasks)
}
