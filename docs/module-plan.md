# Module Classification Plan

该文档梳理 CraftWeave 现有目录，规划哪些组件作为核心模块，哪些属于内置模块，以及后续可扩展的模块类型，便于 ansible 与 playbook 两个子命令共同复用。

## 核心模块 (core/)

这些代码对 CLI 功能至关重要，所有任务执行都会依赖：

- `parser/`：解析 YAML playbook、task 与 role 结构。
- `executor/`：执行 playbook 的主要逻辑，包含 `Playbook` 入口、`when` 条件判断及变量作用域处理等。
- `inventory/`：解析 INI 风格的 inventory 文件，提供主机信息。
- `ssh/`：封装底层 SSH 连接、脚本上传与模板渲染。

## 内置模块 (internal/modules)

开箱即用的任务模块，由 `registry.go` 注册后可直接在 playbook/ansible 中调用：

| 模块名   | 说明                     |
|---------|------------------------|
| `shell`    | 在目标主机执行 shell 命令，支持模板渲染 |
| `script`   | 上传并执行本地脚本             |
| `template` | 渲染 Go 模板并上传到远端          |

未来计划加入：`copy`、`command`、`apt`、`yum`、`systemd`、`service` 等。

### DeepFlow Agent 角色验证反馈

在运行 `example/deploy_deepflow_agent` 作为基准角色时，CLI 输出大量 `Unsupported task type` 错误，说明以下模块尚未实现：

- `setup`/`gather_facts` – 自动收集远端主机信息。
- `set_fact`/`register` – 在任务之间传递变量。
- `fail`、`debug` – 显式失败或调试输出。
- `copy`、`stat` – 文件复制与检查。
- `command` – 直接执行命令（无 shell 解析）。
- `apt`、`yum` – 包管理安装。
- `systemd`、`service` – 管理系统服务。

这些模块将根据优先级逐步纳入内置实现，以保证 Ansible Playbook 与 CraftWeave Playbook 的兼容性。

## 扩展模块

通过插件或自定义实现的模块，利用注册机制插入到执行流程，可覆盖以下场景：

- 额外的文件/包管理功能，如 `copy`、`apt` 等（若未内置）。
- 主机信息收集 `gather_facts`、云平台 API 等。
- 基于 `plugin` 子命令的第三方扩展或 WASM 模块。

核心与内置模块都由 ansible 与 playbook 子命令共用，扩展模块则按需加载，保持架构灵活。
