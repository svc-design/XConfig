#!/bin/bash
set -e

# ============================================================
# 🧩 k3s-role-init.sh
# Version: v1.2.7
# Last Updated: 2025-03-13
#
# 🔄 Change Log:
# - v1.0.0: 初始版本
# - v1.1.0: 精简agent参数
# - v1.1.2: master允许调度pod，taint可选
# - v1.1.3: 修复Cilium Helm冲突
# - v1.1.4: 加入fixed参数清理旧环境
# - v1.1.5: 最小化Cilium部署配置
# ✅ v1.1.6: Cilium 调整为可选安装，通过 --with-cilium 开启
# ============================================================

ROLE=$1
INSTALL_CILIUM=false

# 解析额外参数
for arg in "$@"; do
  if [[ "$arg" == "--with-cilium" ]]; then
    INSTALL_CILIUM=true
  fi
done

print_usage() {
  echo "Usage:"
  echo "  $0 init"
  echo "  $0 fixed"
  echo "  $0 server <EGRESS_EXTERNAL_IP> [SERVER_NODE_IP] [FLANNEL_IFACE] [K3S_TOKEN] [CLUSTER_CIDR] [SERVICE_CIDR] [ADD_TAINT=true|false] [--with-cilium]"
  echo "  $0 agent <SERVER_NODE_IP> <K3S_TOKEN>"
  exit 1
}

if [[ "$ROLE" != "init" && "$ROLE" != "server" && "$ROLE" != "agent" && "$ROLE" != "fixed" ]]; then
  print_usage
fi

### FIXED 模式 ###
if [[ "$ROLE" == "fixed" ]]; then
  echo "🛠️  正在清理旧环境（k3s和Cilium）..."

  /usr/local/bin/k3s-uninstall.sh || true
  /usr/local/bin/k3s-agent-uninstall.sh || true
  rm -rf /etc/rancher /opt/rancher ~/.kube || true

  helm uninstall cilium -n kube-system || true
  helm uninstall cilium-crds -n kube-system || true
  kubectl delete namespace cilium-secrets --ignore-not-found
  kubectl delete crd $(kubectl get crd | grep cilium | awk '{print $1}') --ignore-not-found || true
  kubectl taint nodes $(kubectl get nodes -o name) node.cilium.io/agent-not-ready:NoSchedule- || true

  for iface in $(ip link | awk '/flannel|cilium/ {print $2}' | sed 's/@.*//'); do
    echo "🔥 清理网卡 $iface"
    ip link set $iface down || true
    ip link delete $iface || true
  done

  echo "[3] 清理 flannel 和 cilium 网络接口"
  for iface in $(ip -o link show | awk -F': ' '{print $2}' | grep -E '^(flannel|cilium)'); do
    echo "🔥 删除网卡 $iface"
    ip link set $iface down || true
    ip link delete $iface || true
  done

  echo "[4] 清理可能残留的 cilium 相关接口"
  ip link | grep -E 'cilium_|cilium@|cilium_vxlan' | awk -F': ' '{print $2}' | sed 's/@.*//' | while read -r iface; do
    echo "🔥 删除额外的 cilium 网卡：$iface"
    ip link set "$iface" down || true
    ip link delete "$iface" || true
  done

  echo "✅ 环境清理完成"
  exit 0
fi

### INIT ###
if [[ "$ROLE" == "init" ]]; then
  echo "⚙️ 系统优化启动"
  fallocate -l 1G /swapfile || dd if=/dev/zero of=/swapfile bs=1M count=1024
  chmod 600 /swapfile && mkswap /swapfile && swapon /swapfile
  grep -q swapfile /etc/fstab || echo '/swapfile none swap sw 0 0' >> /etc/fstab

  cat <<EOF >/etc/sysctl.d/k3s.conf
vm.swappiness=10
vm.vfs_cache_pressure=50
net.ipv4.ip_forward=1
EOF
  sysctl --system

  systemctl disable --now snapd motd-news.service rsyslog apport ufw || true
  apt purge -y cloud-init lxd lxc unattended-upgrades || yum remove -y cloud-init || true

  echo "✅ 系统优化完成"
  exit 0
fi

### 下载 K3s 安装器 ###
curl -sfL https://get.k3s.io >install_k3s.sh && chmod +x install_k3s.sh

### SERVER ###
if [ "$ROLE" == "server" ]; then
  EGRESS_EXTERNAL_IP=$2
  SERVER_NODE_IP=${3:-$(hostname -I | awk '{print $1}')}
  FLANNEL_IFACE=${4:-""}
  K3S_TOKEN=${5}
  CLUSTER_CIDR=${6}
  SERVICE_CIDR=${7}
  ADD_TAINT=${8:-false}

  [[ -z "$EGRESS_EXTERNAL_IP" ]] && { echo "❌ 缺少EGRESS_EXTERNAL_IP"; print_usage; }

  echo "🔧 部署参数："
  echo "  SERVER_NODE_IP=${SERVER_NODE_IP}"
  echo "  EGRESS_EXTERNAL_IP=${EGRESS_EXTERNAL_IP}"
  echo "  K3S_TOKEN=${K3S_TOKEN:-自动生成}"
  echo "  ADD_TAINT=${ADD_TAINT}"
  echo "  INSTALL_CILIUM=${INSTALL_CILIUM}"

  INSTALL_K3S_EXEC="server --disable=traefik,servicelb,local-storage \
    --data-dir=/opt/rancher/k3s \
    --node-ip=${SERVER_NODE_IP} \
    --node-external-ip=${EGRESS_EXTERNAL_IP} \
    --advertise-address=${SERVER_NODE_IP}"

  [[ -n "$FLANNEL_IFACE" ]] && INSTALL_K3S_EXEC+=" --flannel-iface=${FLANNEL_IFACE}"
  [[ -n "$K3S_TOKEN" ]] && INSTALL_K3S_EXEC+=" --token=${K3S_TOKEN}"
  [[ -n "$CLUSTER_CIDR" ]] && INSTALL_K3S_EXEC+=" --cluster-cidr=${CLUSTER_CIDR}"
  [[ -n "$SERVICE_CIDR" ]] && INSTALL_K3S_EXEC+=" --service-cidr=${SERVICE_CIDR}"

  INSTALL_K3S_EXEC="$INSTALL_K3S_EXEC" ./install_k3s.sh

  until kubectl get pods -A | grep -q "coredns.*Running"; do sleep 3; done
  mkdir -p ~/.kube && cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
  export KUBECONFIG=~/.kube/config

  if [[ "$INSTALL_CILIUM" == "true" ]]; then
    echo "🚀 开始安装 Cilium..."
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
    helm repo add cilium https://helm.cilium.io && helm repo update

    cat <<EOF >cilium-egress-values.yaml
routingMode: native
ipv4NativeRoutingCIDR: "10.42.0.0/16"
kubeProxyReplacement: false
enableIPv4Masquerade: true
nodePort:
  enabled: true
bpf:
  masquerade: true
ipam:
  mode: kubernetes
egressGateway:
  enabled: true
  installRoutes: true
endpointRoutes:
  enabled: true
cni:
  exclusive: false
envoy:
  enabled: false
proxy:
  enabled: false
l7Proxy: false
hubble:
  enabled: false
operator:
  enabled: true
  replicas: 1
  resources:
    requests:
      cpu: 20m
      memory: 30Mi
    limits:
      cpu: 100m
      memory: 128Mi
resources:
  requests:
    cpu: 20m
    memory: 50Mi
  limits:
    cpu: 100m
    memory: 128Mi
EOF

    helm upgrade --install cilium cilium/cilium \
      -n kube-system \
      --set installCRDs=true \
      -f cilium-egress-values.yaml \
      --wait

    kubectl label node $(hostname) egress-gateway=true --overwrite
    echo "✅ Cilium 安装完成"
  else
    echo "🚫 未启用 Cilium 安装，跳过..."
  fi

  [[ "$ADD_TAINT" == "true" ]] && kubectl taint node $(hostname) node-role.kubernetes.io/master=:NoSchedule --overwrite

  echo "✅ Server 安装完成"
  echo "================================================="
  echo "🌟 当前 Kubernetes & Cilium 环境信息 🌟"
  echo "================================================="

  echo -e "\n📌 Kubernetes 版本："
  kubectl version

  echo -e "\n📌 集群 CIDR (Pod subnet)："
  kubectl cluster-info dump | grep -m 1 cluster-cidr || echo "未显式设置，默认: 10.42.0.0/16"

  echo -e "\n📌 服务 CIDR (Service subnet)："
  kubectl cluster-info dump | grep -m 1 service-cluster-ip-range || echo "未显式设置，默认: 10.43.0.0/16"

  echo -e "\n📌 Helm 版本："
  helm version --short || echo "未安装"

  echo -e "\n📌 Cilium Helm Chart 版本："
  helm list -n kube-system | grep cilium | awk '{print $9}' || echo "未安装"

  echo -e "\n📌 Cilium Pod 版本："
  kubectl -n kube-system get pods -l k8s-app=cilium -o jsonpath='{.items[0].spec.containers[0].image}' 2>/dev/null | cut -d':' -f2 || echo "未安装"

  echo "================================================="
  echo "🎯 环境信息显示完毕!"
  echo "================================================="
fi

### AGENT ###
if [ "$ROLE" == "agent" ]; then
  SERVER_NODE_IP=$2
  K3S_TOKEN=$3
  [[ -z "$SERVER_NODE_IP" || -z "$K3S_TOKEN" ]] && print_usage

  NODE_IP=$(hostname -I | awk '{print $1}')
  INSTALL_K3S_EXEC="agent --server=https://${SERVER_NODE_IP}:6443 --node-ip=${NODE_IP} --token=${K3S_TOKEN}"
  INSTALL_K3S_EXEC="$INSTALL_K3S_EXEC" ./install_k3s.sh

  echo "✅ Agent 安装完成"
fi
