#!/bin/bash
set -e

INGRESS_IP="${1:-$(hostname -I | awk '{print $1}')}"
NODE_LABEL="$2"

echo "🚀 Ingress离线部署开始，IP: ${INGRESS_IP}"

# 解压 nerdctl 并安装
echo "📦 安装nerdctl..."
tar xzvf nerdctl.tar.gz -C /usr/local/bin/

# 导入镜像
echo "🚀 导入镜像到本地containerd..."
nerdctl load -i images/nginx-ingress.tar
nerdctl load -i images/kube-webhook-certgen.tar

# 创建命名空间
kubectl create namespace ingress || true

# 生成 Helm values.yaml
cat > values.yaml <<EOF
controller:
  ingressClass: nginx
  ingressClassResource:
    enabled: true
  replicaCount: 2
  image:
    registry: docker.io
    image: nginx/nginx-ingress
    tag: "2.4.0"
  service:
    enabled: true
    type: NodePort
    externalIPs:
      - $INGRESS_IP
    nodePorts:
      http: 80
      https: 443
EOF

# 节点标签
if [[ -n "$2" ]]; then
cat >> values.yaml <<EOF
  nodeSelector:
    ${NODE_LABEL%%=*}: "${NODE_LABEL#*=}"
EOF
fi

# 安装 Helm Chart（使用本地chart）
helm upgrade --install nginx ./charts/nginx-ingress \
  --namespace ingress -f values.yaml

# 配置 ConfigMap 优化参数
kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-nginx-ingress
  namespace: ingress
data:
  proxy-connect-timeout: "10"
  proxy-read-timeout: "10"
  client-header-buffer-size: 64k
  client-body-buffer-size: 64k
  client-max-body-size: 1000m
  proxy-buffers: "8 32k"
  proxy-buffer-size: 32k
EOF

echo "✅ 离线安装完成，Ingress IP: $INGRESS_IP"
