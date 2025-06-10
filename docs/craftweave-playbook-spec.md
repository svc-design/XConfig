# CraftWeave Playbook YAML 语法规范（v0.1）

# ✅ 顶层为一个或多个 Play（支持并发）
# 每个 Play 对应一个 hosts 主机组 + 一组 tasks

- name: Run system checks         # ✅ 可选：描述 Play
  hosts: all                     # ✅ 必需：支持 inventory 中定义的组名或 all
  gather_facts: false            # ✅ 可选：默认 true，设置为 false 可跳过自动收集系统信息
  vars:                          # ✅ 可选：play 级变量，可在 shell/template 中渲染
    message: "hello world"

  tasks:
    - name: Show hostname        # ✅ 可选：描述任务
      shell: hostname            # ✅ shell 模块，执行远程 shell 命令

    - name: Run CPU count script
      script: ./example/nproc.sh # ✅ script 模块：上传本地脚本并远程执行

    - name: Render config file
      template:
        src: ./templates/nginx.conf.j2
        dest: /tmp/nginx.conf

    - name: Show welcome message
      shell: echo "{{ message }}"

    - name: Render remote MOTD
      template:
        src: ./templates/motd.tmpl
        dest: /tmp/motd.txt

  roles:
    - role: common

---

# 🚀 TODO 支持（版本 roadmap）
# - copy: src= dest= mode=
# - when / tags / loop 等语法糖

---

# 🧪 示例测试 Playbook 1（echo + script）

- name: Simple echo
  hosts: all
  tasks:
    - name: Echo message
      shell: echo Hello from CraftWeave

    - name: Show OS
      script: ./example/uname.sh

---

# ❌ 错误示例（用于 parser 校验测试）

# - name: Missing hosts
#   tasks:
#     - shell: echo "invalid"

# - name: Invalid task type
#   hosts: all
#   tasks:
#     - foo: bar

---

# 🔍 Parser 校验建议：
# 1. 检查 play 是否为 list（YAML 顶层）
# 2. 检查每个 play 是否包含 hosts + tasks 字段
# 3. tasks 中每项是否为 map，包含合法的模块字段（shell/script）
# 4. 如果指定了模块字段，值必须是字符串
# 5. 错误信息应带行号与 task 名称提示

# ✅ 合法模块 key（暂支持）: shell, script, template
# 🚫 不合法的 key：除上述外都报错（为后续模块保留）

# CraftWeave Playbook 元素定义表格

| 元素名   | 类型   | 是否必要 | 示例说明                           |
|----------|--------|----------|------------------------------------|
| `name`   | string | ✅ 是     | Play 或 task 的描述                 |
| `hosts`  | string | ✅ 是     | 当前 play 作用的 inventory 主机组  |
| `gather_facts` | bool | 可选      | 是否自动收集系统信息，默认 true |
| `tasks`  | list   | ✅ 是     | 每条任务可以是 shell、script、template 等 |
| `shell`  | string | 可选      | 执行单条远程命令                   |
| `script` | string | 可选      | 执行本地脚本并上传远程运行         |
| `template` | map  | 可选      | 渲染本地模板并上传至远程           |
| `vars`   | map    | 可选（V1）| 支持在 shell 和 template 中引用     |
| `roles`  | list  | 可选      | 引用 role 名称，自动加载其 tasks    |
