#!/bin/bash

#https://github.com/containerd/nerdctl/releases/download/v2.0.2/nerdctl-2.0.2-linux-amd64.tar.gz
#https://github.com/containerd/nerdctl/releases/download/v2.0.2/nerdctl-full-2.0.2-linux-amd64.tar.gz
#wget https://github.com/containernetworking/plugins/releases/download/v1.6.2/cni-plugins-linux-amd64-v1.6.2.tgz

#!/bin/bash
set -e

# =============================================
# ✅ 环境变量检查（可配置）
# =============================================
: "${REGISTRY_DOMAIN:=kube.registry.local}"
: "${REGISTRY_PORT:=5000}"
: "${NERDCTL_VERSION:=v2.0.2}"
: "${CNI_VERSION:=v1.6.2}"
: "${CNI_DIR:=/opt/cni/bin}"
: "${CERT_DIR:=/opt/registry/certs}"
: "${CONFIG_DIR:=/opt/registry/config}"
: "${REGISTRY_DATA:=/var/lib/registry}"
: "${REGISTRY_YAML:=registry-config.yaml}"
: "${COMPOSE_YAML:=compose.yaml}"
: "${TAR_FILE:=registry.tar}"

# =============================================
# ✅ 自动检测 containerd.sock
# =============================================
if [[ -S "/run/k3s/containerd/containerd.sock" ]]; then
  export CONTAINERD_ADDRESS="/run/k3s/containerd/containerd.sock"
elif [[ -S "/run/containerd/containerd.sock" ]]; then
  export CONTAINERD_ADDRESS="/run/containerd/containerd.sock"
elif [[ -S "/var/run/containerd/containerd.sock" ]]; then
  export CONTAINERD_ADDRESS="/var/run/containerd/containerd.sock"
else
  echo "❌ 未检测到有效的 containerd.sock，请确认 containerd 是否正常运行。"
  exit 1
fi

export NERDCTL_NAMESPACE="k8s.io"

# =============================================
echo "📦 准备 nerdctl 全功能版..."
if ! command -v nerdctl &>/dev/null; then
  if [ ! -f /tmp/nerdctl-full.tgz ]; then
    echo "⬇️ 下载 nerdctl..."
    wget -O /tmp/nerdctl-full.tgz \
      "https://github.com/containerd/nerdctl/releases/download/${NERDCTL_VERSION}/nerdctl-full-${NERDCTL_VERSION#v}-linux-amd64.tar.gz"
  else
    echo "📦 已存在 nerdctl-full.tgz，跳过下载"
  fi

  echo "📦 解压 nerdctl 到 /usr/local..."
  sudo tar -C /usr/local -xzf /tmp/nerdctl-full.tgz
  echo "✅ nerdctl 安装完成: $(nerdctl --version)"
else
  echo "✅ nerdctl 已存在: $(nerdctl --version)"
fi

# =============================================
echo "📦 安装 CNI 插件..."
if [ ! -f "${CNI_DIR}/bridge" ]; then
  if [ ! -f /tmp/cni.tgz ]; then
    echo "⬇️ 下载 CNI 插件..."
    wget -O /tmp/cni.tgz \
      "https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-linux-amd64-${CNI_VERSION}.tgz"
  else
    echo "📦 已存在 cni.tgz，跳过下载"
  fi

  sudo mkdir -p "${CNI_DIR}"
  sudo tar -C "${CNI_DIR}" -xzf /tmp/cni.tgz
  echo "✅ CNI 插件已安装到: ${CNI_DIR}"
else
  echo "✅ CNI 插件已存在: ${CNI_DIR}/bridge"
fi

# =============================================
echo "📦 解压 SSL 证书..."
if [ -d "$CERT_DIR" ] && [ -f "${CERT_DIR}/kube.registry.local.cert" ]; then
  echo "✅ 证书目录已存在: $CERT_DIR"
else
  if [ -f "ssl_certificates.tar.gz" ]; then
    mkdir -p "$CERT_DIR"
    tar -xvpf ssl_certificates.tar.gz -C "$CERT_DIR"
    echo "✅ 证书已解压至: $CERT_DIR"
  else
    echo "⚠️ 未找到 ssl_certificates.tar.gz，跳过证书解压"
  fi
fi

# =============================================

# ============ 生成 registry-config ============
echo "⚙️ 准备 registry 配置..."
sudo mkdir -p "$COMPOSE_DIR"
echo "📝 写入 registry-config.yaml..."
sudo tee "$REGISTRY_CONFIG" > /dev/null <<EOF
version: 0.1
log:
  fields:
    service: registry
storage:
  cache:
    blobdescriptor: inmemory
  filesystem:
    rootdirectory: /var/lib/registry
  delete:
    enabled: true
http:
  addr: :5000
  headers:
    X-Content-Type-Options: [nosniff]
  tls:
    certificate: ${CERT_DIR}/kube.registry.local.cert
    key: ${CERT_DIR}/kube.registry.local.key
health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3
EOF

sudo cp "$COMPOSE_YAML" "$CONFIG_DIR/compose.yaml"
sudo mkdir -p "$REGISTRY_DATA"
echo "✅ 写入完成: $REGISTRY_CONFIG"

# =============================================
echo "📦 导入本地 registry 镜像..."
if [ -f "/usr/local/deepflow/$TAR_FILE" ]; then
  sudo nerdctl --namespace $NERDCTL_NAMESPACE load -i "/usr/local/deepflow/$TAR_FILE"
else
  echo "⚠️ 本地镜像文件不存在：/usr/local/deepflow/$TAR_FILE"
fi

# =============================================
echo "🔁 重启 registry 服务..."
sudo nerdctl --namespace $NERDCTL_NAMESPACE compose -f "$CONFIG_DIR/compose.yaml" down || true
sudo nerdctl --namespace $NERDCTL_NAMESPACE compose -f "$CONFIG_DIR/compose.yaml" up -d

# =============================================
echo "🔗 添加 hosts 映射..."
if ! grep -q "$REGISTRY_DOMAIN" /etc/hosts; then
  echo "127.0.0.1 $REGISTRY_DOMAIN" | sudo tee -a /etc/hosts
  echo "✅ /etc/hosts 已添加 $REGISTRY_DOMAIN"
else
  echo "✅ hosts 中已存在 $REGISTRY_DOMAIN"
fi

echo "✅ Registry 启动成功: https://$REGISTRY_DOMAIN:$REGISTRY_PORT"

