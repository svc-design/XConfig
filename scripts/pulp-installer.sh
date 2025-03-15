#!/bin/bash
set -e

echo "🚀 开始离线安装 Pulp Operator..."

# 安装 nerdctl（如存在）
if [ -f nerdctl.tar.gz ]; then
  echo "📦 解压 nerdctl..."
  tar xzvf nerdctl.tar.gz -C /usr/local/bin/
fi

# 导入镜像
echo "🚀 导入 pulp-operator 镜像..."
if command -v docker &>/dev/null && docker info &>/dev/null; then
  docker load -i images/pulp-operator.tar
elif [ -S /run/k3s/containerd/containerd.sock ]; then
  export CONTAINERD_ADDRESS=/run/k3s/containerd/containerd.sock
  nerdctl --namespace k8s.io load -i images/pulp-operator.tar
elif [ -S /run/containerd/containerd.sock ]; then
  export CONTAINERD_ADDRESS=/run/containerd/containerd.sock
  nerdctl --namespace k8s.io load -i images/pulp-operator.tar
else
  echo "❌ 没有可用的容器运行时"
  exit 1
fi

# 创建命名空间
kubectl create namespace pulp || true

# 安装 chart
echo "📦 安装本地 Helm Chart..."
helm upgrade --install pulp-operator ./charts/pulp-operator -n pulp --create-namespace

# 等待 CRD 注册
sleep 10

# 生成默认 CR yaml（可改为 values 覆盖渲染）
echo "📝 生成 CR manifests/pulp-cr.yaml..."
mkdir -p manifests
cat > manifests/pulp-cr.yaml <<EOF
apiVersion: repo-manager.pulpproject.org/v1beta2
kind: Pulp
metadata:
  name: pulp
  namespace: pulp
spec:
  api:
    replicas: 1
    ingress:
      enabled: true
      tls:
        enabled: true
        secretName: pulp-tls-secret
  content:
    replicas: 1
  worker:
    replicas: 2
  plugins:
    - pulp-container
    - pulp-rpm
    - pulp-deb
    - pulp-helm
    - pulp-file
    - pulp-nuget
  storage:
    type: s3
    s3:
      bucket: pulp-repo-bucket
      accessKeyID: <your-access-key>
      secretAccessKey: <your-secret-key>
      endpointURL: https://oss-cn-beijing.aliyuncs.com
      region: cn-beijing
      tls:
        insecure: false
EOF

# 应用 CR
echo "✅ 应用 Pulp CR"
kubectl apply -f manifests/pulp-cr.yaml

echo "🎉 Pulp 安装完成，查看状态：kubectl -n pulp get pods"

