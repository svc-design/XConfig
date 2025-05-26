use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Play {
    pub name: String,
    pub tasks: Vec<Task>,
}

#[derive(Debug, Deserialize)]
pub struct Task {
    pub name: String,
    pub shell: Option<String>,
    pub script: Option<String>,
}

