# CraftWeave

**CraftWeave** 是一个使用 Go 编写的 Ansible-like 工具，支持任务执行、架构编织、图模型导出与插件扩展。

## 🧩 特性

- 🛠️ `craftweave ansible`：执行单条远程命令
- 📜 `craftweave playbook`：运行 YAML Playbook
- 🔐 `craftweave vault`：加解密配置
- 🧠 `craftweave cmdb`：输出图数据库模型
- 🧩 `craftweave plugin`：加载并运行插件（支持 WASM）

## 🚀 快速开始

```bash
go build -o craftweave
./craftweave
./craftweave playbook deploy.yaml

项目结构

CraftWeave/
├── cmd/                  # Cobra 命令定义
│   ├── root.go           # 根命令
│   ├── ansible.go        # 类 ansible 子命令
│   ├── playbook.go       # 执行 playbook
│   ├── vault.go          # 加解密相关
│   ├── cmdb.go           # 输出图模型
│   └── plugin.go         # 插件运行
├── core/                 # 核心逻辑模块
│   ├── executor/         # 执行器引擎
│   ├── parser/           # playbook/拓扑解析
│   ├── cmdb/             # 图模型构建与导出
│   └── plugin/           # 插件接口定义与加载
├── plugins/              # 插件目录（WASM/Go 可选）
├── internal/             # 内部工具库
├── go.mod
├── main.go
├── README.md
└── banner.txt            # CLI 启动 ASCII 图标

🔮 愿景
CraftWeave 旨在成为下一代 DevOps 工具平台 —— 融合任务调度、架构可视化与智能插件。
