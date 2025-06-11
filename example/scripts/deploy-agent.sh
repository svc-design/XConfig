#!/bin/bash

set -e

####################################
# 🌐 配置区
####################################

IP_LIST="./ip.list"
SERVICE_NAME="deepflow-agent"
PKG_DIR="deepflow-agent-for-linux"
MAX_PARALLEL=5  # 可调：最大并发数

# === 默认值，可通过参数覆盖 ===
CONTROLLER_IP=""
VTAP_GROUP_ID=""

# === SSH 通用选项（含超时 15 秒）
SSH_OPTS="-o StrictHostKeyChecking=no -o ConnectTimeout=15"

####################################
# 参数解析
####################################

if [[ $# -eq 0 ]]; then
  echo "用法: $0 {deploy|upgrade|verify} --controller <ip> --group <id>"
  exit 1
fi

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


if [[ "$ACTION" != "deploy" && "$ACTION" != "upgrade" && "$ACTION" != "verify" ]]; then
  echo "用法: $0 {deploy|upgrade|verify} --controller <ip> --group <id>"
  exit 1
fi

if [[ "$ACTION" != "verify" && ( -z "$CONTROLLER_IP" || -z "$VTAP_GROUP_ID" ) ]]; then
  echo "❗ deploy/upgrade 必须传入 --controller 和 --group 参数"
  exit 1
fi

####################################
# 核心函数
####################################

worker() {
  local ip="$1"
  local user="$2"
  local pass="$3"

  echo "🔧 [$ACTION] 处理主机 $ip ($user)"

  if [[ "$ACTION" == "verify" ]]; then
    verify_agent "$ip" "$user" "$pass"
    return
  fi

  local remote_info arch init pkg_type

  remote_info=$(fetch_remote_info "$ip" "$user" "$pass") || {
    echo "❌ $ip 获取远程信息失败"
    return
  }

  arch=$(echo "$remote_info" | cut -d'|' -f1)
  init=$(echo "$remote_info" | cut -d'|' -f2)
  pkg_type=$(echo "$remote_info" | cut -d'|' -f3)

  if [[ "$init" == "unknown" || "$pkg_type" == "unknown" ]]; then
    echo "❌ $ip 不支持的初始化或包管理器: $init/$pkg_type"
    return
  fi

  pkg_path=$(choose_agent_package "$arch" "$init" "$pkg_type")

  if [[ "$pkg_path" == "UNSUPPORTED" ]]; then
    echo "❌ $ip 无匹配安装包: $arch/$init/$pkg_type"
    return
  fi

  install_agent "$ip" "$user" "$pass" "$pkg_path" && update_config "$ip" "$user" "$pass"
  echo "✅ $ip $ACTION 完成"
  echo "-------------------------------------------"
}

fetch_remote_info() {
  local ip="$1"
  local user="$2"
  local pass="$3"

  sshpass -p "$pass" ssh $SSH_OPTS "$user@$ip" bash <<'EOF'
arch=$(uname -m)
if command -v systemctl >/dev/null; then init=systemd;
elif command -v initctl >/dev/null; then init=upstart;
else init=unknown; fi
if command -v rpm >/dev/null; then pkg=rpm;
elif command -v dpkg >/dev/null; then pkg=deb;
else pkg=unknown; fi
echo "${arch}|${init}|${pkg}"
EOF
}

choose_agent_package() {
  local arch="$1" init="$2" pkg_type="$3"

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
  local ip="$1" user="$2" pass="$3" pkg_path="$4"
  local remote_pkg="/tmp/agent.${pkg_path##*.}"

  sshpass -p "$pass" scp $SSH_OPTS "$pkg_path" "$user@$ip:$remote_pkg"

  sshpass -p "$pass" ssh $SSH_OPTS "$user@$ip" bash <<EOF
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
  local ip="$1" user="$2" pass="$3"
  sshpass -p "$pass" ssh $SSH_OPTS "$user@$ip" bash <<EOF
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
  local ip="$1" user="$2" pass="$3"
  echo "🔍 $ip 状态检查："
  sshpass -p "$pass" ssh $SSH_OPTS "$user@$ip" "
    systemctl is-active $SERVICE_NAME 2>/dev/null || \
    service $SERVICE_NAME status || \
    initctl status $SERVICE_NAME
  "
}

####################################
# 控制并发执行主逻辑
####################################

# 简单并发控制函数 (纯 Bash 无需 parallel)
sem(){
  while [[ $(jobs -r | wc -l) -ge $MAX_PARALLEL ]]; do
    sleep 0.5
  done
}

while read -r ip user pass; do
  sem
  worker "$ip" "$user" "$pass" &
done < "$IP_LIST"

wait
echo "🎯 全部任务执行完成"
