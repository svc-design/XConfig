#!/bin/bash

set -e

####################################
# 🌐 配置区：需根据实际环境修改
####################################

IP_LIST="./ip.list"  # 定义主机清单文件路径，每行格式为：IP USER PASSWORD
SERVICE_NAME="deepflow-agent"  # 定义要操作的服务名称（deepflow-agent）
PKG_DIR="deepflow-agent-for-linux"  # 存放各平台 RPM 包的目录

# === 默认值，可通过参数覆盖 ===
CONTROLLER_IP=""
VTAP_GROUP_ID=""

# === 参数解析 ===
ACTION="$1"
shift

while [[ $# -gt 0 ]]; do
  case "$1" in
    --controller)
      CONTROLLER_IP="$2"
      shift 2
      ;;
    --group)
      VTAP_GROUP_ID="$2"
      shift 2
      ;;
    *)
      echo "未知参数: $1"
      exit 1
      ;;
  esac
done

# 检查参数完整性
if [[ "$ACTION" != "deploy" && "$ACTION" != "upgrade" && "$ACTION" != "verify" ]]; then
  echo "用法: $0 {deploy|upgrade|verify} --controller <ip> --group <id>"
  exit 1
fi

if [[ "$ACTION" != "verify" && ( -z "$CONTROLLER_IP" || -z "$VTAP_GROUP_ID" ) ]]; then
  echo "❗ deploy/upgrade 必须传入 --controller 和 --group 参数"
  exit 1
fi

choose_agent_package() {
  local arch="$1"
  local init=""
  local pkg=""

  init=$(sshpass -p "$pass" ssh -o StrictHostKeyChecking=no "$user@$ip" '
    if command -v systemctl >/dev/null; then echo systemd;
    elif command -v initctl >/dev/null; then echo upstart;
    else echo unknown; fi')

  if [[ "$init" == "unknown" ]]; then
    echo "UNSUPPORTED"
    return
  fi

  pkg_type=$(sshpass -p "$pass" ssh -o StrictHostKeyChecking=no "$user@$ip" '
    if command -v rpm >/dev/null; then echo rpm;
    elif command -v dpkg >/dev/null; then echo deb;
    else echo unknown; fi')

  if [[ "$pkg_type" == "unknown" ]]; then
    echo "UNSUPPORTED"
    return
  fi

  # 查找匹配初始化系统和包格式的文件，优先考虑带架构字段的，降级用通用版
  pkg=$(find "$PKG_DIR" -type f \( \
    -name "deepflow-agent-*.$init-*.$pkg_type" -o \
    -name "deepflow-agent-*.$init.$pkg_type" \) | sort -V | tail -1)

  if [[ -n "$pkg" ]]; then
    echo "$pkg"
  else
    echo "UNSUPPORTED"
  fi
}

install_agent() {
  local ip="$1"
  local user="$2"
  local pass="$3"
  local pkg_path="$4"

  local remote_pkg="/tmp/agent.${pkg_path##*.}"

  sshpass -p "$pass" scp -o StrictHostKeyChecking=no "$pkg_path" "$user@$ip:$remote_pkg"

  sshpass -p "$pass" ssh -o StrictHostKeyChecking=no "$user@$ip" bash <<EOF
set -e

if [[ "$remote_pkg" == *.rpm ]]; then
  rpm -Uvh --replacepkgs "$remote_pkg"
elif [[ "$remote_pkg" == *.deb ]]; then
  dpkg -i "$remote_pkg" || apt-get install -f -y
else
  echo "❌ 不支持的安装包格式"
  exit 1
fi

if command -v systemctl &>/dev/null; then
  systemctl enable $SERVICE_NAME
  systemctl restart $SERVICE_NAME
elif command -v service &>/dev/null; then
  service $SERVICE_NAME restart
  chkconfig $SERVICE_NAME on
elif command -v initctl &>/dev/null; then
  initctl restart $SERVICE_NAME || initctl start $SERVICE_NAME
else
  echo "❌ 无法识别服务管理方式"
fi
EOF
}

update_config() {
  local ip="$1"
  local user="$2"
  local pass="$3"

  sshpass -p "$pass" ssh -o StrictHostKeyChecking=no "$user@$ip" bash <<EOF
set -e
CONFIG_FILE="/etc/deepflow-agent.yaml"
mkdir -p \$(dirname \$CONFIG_FILE)
cat > "\$CONFIG_FILE" <<CFG
controller-ips:
  - $CONTROLLER_IP
vtap-group-id: "$VTAP_GROUP_ID"
CFG
chmod 644 "\$CONFIG_FILE"
chown root:root "\$CONFIG_FILE"
EOF
}

verify_agent() {
  local ip="$1"
  local user="$2"
  local pass="$3"
  echo "🔍 $ip 状态检查："
  sshpass -p "$pass" ssh -o StrictHostKeyChecking=no "$user@$ip" "
    systemctl is-active $SERVICE_NAME 2>/dev/null || \
    service $SERVICE_NAME status || \
    initctl status $SERVICE_NAME
  "
}

while read -r ip user pass; do
  echo "🔧 [$ACTION] 处理主机 $ip ($user)"

  if [[ "$ACTION" == "verify" ]]; then
    verify_agent "$ip" "$user" "$pass"
    continue
  fi

  arch=$(sshpass -p "$pass" ssh -o StrictHostKeyChecking=no "$user@$ip" "uname -m")
  pkg_path=$(choose_agent_package "$arch")

  if [[ "$pkg_path" == "UNSUPPORTED" ]]; then
    echo "❌ 不支持的系统架构或未找到匹配包: $arch"
    continue
  fi

  install_agent "$ip" "$user" "$pass" "$pkg_path"
  update_config "$ip" "$user" "$pass"
  echo "✅ $ip $ACTION 完成"
  echo "-------------------------------------------"
done < "$IP_LIST"
