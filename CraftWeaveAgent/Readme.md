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

目标是构建一个完全独立、本地执行的 CraftWeaveAgent，支持从 Git 同步 playbook 并执行 shell/script 任务，不依赖 controller，也不进行远程 SSH。

✅ 功能目标（本地 Playbook 执行器）
✅ 支持从 Git 仓库拉取 playbook.yaml
✅ 支持任务类型：shell, script（本地执行）
✅ 支持 --oneshot 或 daemon 模式定期同步执行
✅ 所有任务限定运行在本机（无 SSH）


# 🧩 支持命令说明

命令格式	功能说明
cw-agent --mode oneshot	一次性从 /etc/cw-agent.conf 拉取 Git 仓库并执行配置
cw-agent --mode daemon	持续运行，从 Git 周期性拉取并执行（定期刷新）
cw-agent playbook <path>	直接执行本地指定 Playbook 文件（只作用于本机，不依赖 Git）
cw-agent status	输出最近一次任务的 JSON 执行结果
cw-agent version	显示版本号
