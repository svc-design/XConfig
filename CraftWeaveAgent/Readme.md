```
cw-agent/
├── Cargo.toml
├── README.md
├── cw-agent.service              # 可选：systemd 单元文件
└── src/
    ├── main.rs                  # 入口，CLI 参数解析 + 调度器入口
    ├── scheduler.rs             # 定时/触发式拉取、执行、保存
    ├── config.rs                # 拉取并解析远程配置（Git/HTTP 等）
    ├── executor.rs              # 执行任务（command/copy/service 等）
    ├── result_store.rs          # 存储执行结果（JSON / local DB）
    └── models.rs                # 配置结构体 / 执行结果结构体
```

# 扩展

cw-agent run --oneshot             # 单次运行
cw-agent daemon                    # 持续运行
cw-agent status                    # 输出本地最新执行结果
cw-agent apply --file config.json  # 本地执行一次（不拉取）
