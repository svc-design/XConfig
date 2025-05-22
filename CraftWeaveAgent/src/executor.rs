// File: src/executor.rs
// ---------------------
use crate::config::TaskConfig;
use crate::result_store::CommandResult;
use tokio::process::Command;

pub async fn apply(tasks: Vec<TaskConfig>) -> anyhow::Result<Vec<CommandResult>> {
    let mut results = vec![];
    for task in tasks {
        let output = Command::new("sh")
            .arg("-c")
            .arg(&task.command)
            .output()
            .await?;

        results.push(CommandResult {
            task: task.name,
            stdout: String::from_utf8_lossy(&output.stdout).into(),
            stderr: String::from_utf8_lossy(&output.stderr).into(),
            success: output.status.success(),
            return_code: output.status.code().unwrap_or(-1),
        });
    }
    Ok(results)
}
