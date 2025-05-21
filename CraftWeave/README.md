# CraftWeave

**CraftWeave** 是一个使用 Go 编写的 Ansible-like 工具，支持任务执行、架构编织、图模型导出与插件扩展。

---

## 🧩 特性

- 🛠️ `craftweave ansible`：执行单条远程命令（支持 shell 模块）
- 📜 `craftweave playbook`：运行 YAML Playbook
- 🔐 `craftweave vault`：加解密配置(Todo)
- 🧠 `craftweave cmdb`：输出图数据库模型(Todo)
- 🧩 `craftweave plugin`：加载并运行插件（Todo 支持 WASM）

---

## 🚀 快速开始

1. 编译项目 make
2. 执行远程 shell 命令（类似 ansible）
使用 INI 格式的 inventory 文件：

ini
[all]
deepflow-demo  ansible_host=192.168.124.77     ansible_ssh_user=shenlan
cn-hub         ansible_host=1.15.155.245       ansible_ssh_user=ubuntu
...

[all:vars]
ansible_port=22
ansible_ssh_private_key_file=~/.ssh/id_rsa
执行命令： ./craftweave ansible all -i example/inventory -m shell -a 'id'
输出示例：

🧶 欢迎使用：CraftWeave - 任务与架构编织工具
deepflow-demo | CHANGED | rc=0 >>
uid=1000(shenlan) gid=1000(shenlan) groups=1000(shenlan),10(wheel)

cn-hub | CHANGED | rc=0 >>
uid=1000(ubuntu) gid=1001(ubuntu) groups=1001(ubuntu),27(sudo),...

...
支持 dry-run 模式：

bash
./craftweave ansible all -i example/inventory -m shell -a 'id' -C

# 📁 项目结构

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
├── internal/             # 内部工具库（如 ssh 执行器、inventory 解析器）
│   ├── ssh/
│   └── inventory/
├── plugins/              # 插件目录（WASM/Go 可选）
├── example/              # 示例配置（inventory 等）
│   └── inventory
├── banner.txt            # CLI 启动 ASCII 图标
├── go.mod
├── go.sum
├── main.go
└── README.md

# 🔮 愿景

CraftWeave 旨在成为下一代 DevOps 工具 —— 融合任务调度、架构可视化与智能插件能力，支持轻量化、模块化和智能化的运维体验。
