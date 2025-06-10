# CraftWeave 设计概览

本文件结合当前代码库，对 CraftWeave 的架构模式、已实现模块以及后续可扩展点进行简要说明，便于后续开发。

## 1. 现在的架构模式

CraftWeave 控制端采用 Go 实现，整体架构遵循 "CLI -> 执行器 -> 模块 -> SSH" 的流水线：

1. **CLI 层 (`cmd/`)** 使用 [Cobra](https://github.com/spf13/cobra) 提供 `ansible` 与 `playbook` 等子命令，解析用户输入。
2. **解析与执行 (`core/`)**
   - `parser`：解析 YAML Playbook，转化为 `Play`、`Task` 等结构体。
   - `executor`：遍历任务，调用模块或内置逻辑执行，并支持并发、`when` 条件与日志收集。
3. **内部库 (`internal/`)**
   - `inventory`：解析 INI 格式 inventory，提供主机信息。
   - `ssh`：封装远程命令执行、脚本上传和模板渲染。
   - `modules`：注册并实现可扩展的任务模块。
4. **Agent (`CraftWeaveAgent/`)** 采用 Rust 编写，可定时拉取并在本地执行 Playbook，实现轻量化的边缘执行能力。

整体模式保持松耦合，模块通过注册表动态查找，方便后续按需扩展。

## 2. 已经支持的模块

当前 `internal/modules` 目录下提供以下内建模块：

| 模块名   | 功能简介                         |
|---------|---------------------------------|
| `shell`    | 在目标主机执行 shell 命令，支持变量渲染 |
| `script`   | 上传并运行本地脚本                   |
| `template` | 渲染 Go 模板并上传到远端               |

Playbook 任务若未找到注册模块，则回退到同名内置处理，因此上述三个模块也是 CLI 默认支持的功能。

## 3. 需要开发扩展的

为了接近完整的 Ansible 能力，仍有若干模块和功能待实现，示例包括：

- **文件与包管理**：`copy`、`command`、`apt`、`yum` 等模块。
- **变量相关**：`set_fact`、`register` 用于动态变量赋值及结果保存。
- **条件与权限**：完善 `when` 表达式，支持 `become` 以 sudo 身份执行。
- **信息收集**：`gather_facts` 模块，用于获取主机系统信息。
- **插件与拓扑**：预留的 `plugin`、`cmdb` 子命令，将来可提供 WASM 扩展和架构导出能力。

社区或开发者可根据模块注册机制，自行扩展以上能力或添加新的功能模块。
