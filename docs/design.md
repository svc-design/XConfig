# CraftWeave Design Overview

This document summarizes the current codebase layout and the main functions available in each module. The goal is to provide an easy reference when extending CraftWeave with new features (e.g. copy, command, package modules).

## Directory Structure

```
建议的目录结构
bash
复制
编辑
core/
├── parser/                 # 解析 playbook、task、role 结构
│   └── parser.go
├── executor/               # Playbook 执行逻辑（核心 orchestrator）
│   ├── playbook.go         # ExecutePlaybook 主入口
│   ├── condition.go        # when 条件处理
│   └── variable.go         # register、set_fact、变量作用域解析

internal/
├── modules/                # 内建模块实现，如 shell, copy, yum 等
│   ├── registry.go         # 模块注册机制（注册名 → HandlerFunc）
│   ├── shell.go
│   ├── command.go
│   ├── copy.go
│   ├── yum.go
│   ├── apt.go
│   ├── systemd.go
│   ├── service.go
│   ├── fail.go
│   ├── debug.go
│   └── types.go            # 模块接口定义（如 TaskHandler、ModuleResult）

├── ssh/                    # 底层 SSH 操作封装
│   ├── runner.go           # SSH 执行
│   ├── result.go           # 执行结果封装
│   ├── script.go           # 远程脚本
│   ├── template.go         # 模板渲染
│   ├── copy.go             # 文件传输
│   ├── facts.go            # setup 收集主机信息
│   └── pkg.go              # 包管理工具抽象层（apt/yum）

├── inventory/              # 解析 hosts.yaml/ini 文件
│   └── inventory.go

cmd/
├── ansible.go
├── playbook.go
├── plugin.go
└── root.go

example/
└── playbooks/


```

### Subcommands (`cmd/`)
- **root.go** – sets up the CLI and registers subcommands.
- **ansible.go** – ad-hoc task execution. Supports `shell` and `script` modules.
- **playbook.go** – runs a YAML playbook with `shell`, `script` and `template` tasks.
- **vault.go** – placeholder for future encryption features.
- **cmdb.go** – placeholder to export topology graphs.
- **plugin.go** – placeholder for loading plugins.

### Core Modules (`core/`)
- **executor/playbook.go** – iterates over plays and tasks and executes them using the SSH helpers. Supports roles by loading tasks from `roles/<role>/tasks/main.yaml`.
- **parser/parser.go** – defines YAML structures (`Play`, `Task`, `Template`) and parses playbooks.
- **cmdb/** – reserved for CMDB graph generation (not yet implemented).

### Internal Libraries (`internal/`)
- **inventory/inventory.go** – parses INI style inventory files and returns a list of `Host` structures.
- **ssh/runner.go** – runs remote shell commands via SSH with key or password authentication.
- **ssh/script.go** – uploads a local script using base64 and executes it remotely.
- **ssh/template.go** – renders a Go template and uploads the result to the remote host.
- **ssh/formatter.go** – utilities for aggregated output.
- **ssh/result.go** – defines the `CommandResult` struct used across the executor.

✅ 模块系统设计（internal/modules）
types.go
go
复制
编辑
type TaskContext struct {
  Host       inventory.Host
  Vars       map[string]string
  Task       parser.Task
  WorkingDir string
}

type ModuleResult struct {
  Host       string
  Output     string
  ReturnCode int
  ReturnMsg  string
}

type TaskHandler func(TaskContext) ModuleResult
registry.go
go
复制
编辑
var moduleRegistry = map[string]TaskHandler{}

func RegisterModule(name string, handler TaskHandler) {
  moduleRegistry[name] = handler
}

func GetHandler(name string) (TaskHandler, bool) {
  h, ok := moduleRegistry[name]
  return h, ok
}
每个模块如 copy.go
go
复制
编辑
func init() {
  RegisterModule("copy", CopyHandler)
}

func CopyHandler(ctx TaskContext) ModuleResult {
  // 使用 ssh.CopyFile，渲染变量，处理权限
  ...
}
✅ 执行器改动：executor/playbook.go
在每个任务类型分支前，尝试从 moduleRegistry 获取 Handler：

go
复制
编辑
if handler, ok := modules.GetHandler(task.Type()); ok {
  res = handler(TaskContext{
    Host: h,
    Vars: hv,
    Task: task,
    WorkingDir: baseDir,
  })
} else {
  res = ssh.CommandResult{
    Host:       h.Name,
    ReturnMsg:  "FAILED",
    ReturnCode: 1,
    Output:     fmt.Sprintf("Unsupported task type in '%s'", task.Name),
  }
}
新增 task.Type() 方法判断类型，如：

go
复制
编辑
func (t Task) Type() string {
  switch {
  case t.Shell != "":
    return "shell"
  case t.Command != "":
    return "command"
  case t.Script != "":
    return "script"
  case t.Copy != nil:
    return "copy"
  case t.Yum != nil:
    return "yum"
  case t.Apt != nil:
    return "apt"
  ...
  default:
    return "unknown"
  }
}
✅ 扩展功能设计建议
功能模块	建议位置	实现细节
set_fact	executor/variable.go	动态变量注入 map[string]string，生命周期限定
register	executor/playbook.go	将 res.Output 保存为变量
when	executor/condition.go	支持 ==, !=, bool 值判断
gather_facts	ssh/facts.go	支持 uname, lsb_release, hostname 等
become	ssh/runner.go	在 command/shell 执行前添加 sudo 前缀（可选实现）
delegate_to	暂不实现	复杂跨主机跳转，初期可忽略
