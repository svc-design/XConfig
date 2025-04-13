#!/bin/bash
set -e

export INSTALL_K3S_EXEC="server --disable=traefik,servicelb,local-storage --data-dir=/opt/rancher/k3s --kube-apiserver-arg=service-node-port-range=0-50000 --flannel-iface=br0"
curl -sfL https://get.k3s.io | sh -

export INSTALL_K3S_EXEC="server --data-dir=/mnt/opt/rancher/k3s --disable=traefik,servicelb,local-storage --kube-apiserver-arg=service-node-port-range=0-50000 --system-default-registry=registry.cn-hangzhou.aliyuncs.com --flannel-iface=br0"
curl -sfL https://rancher-mirror.rancher.cn/k3s/k3s-install.sh | sh -

# === 设置本地 kubeconfig ===
mkdir -p ~/.kube
cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
chmod 600 ~/.kube/config
export KUBECONFIG=~/.kube/config

curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# === 等待 CoreDNS 启动 ===
echo "⏳ 等待 CoreDNS 启动..."
until kubectl get pods -A 2>/dev/null | grep -q "coredns.*Running"; do
  sleep 3
done
echo "✅ K3s 安装完成，kubectl/helm 已就绪"
