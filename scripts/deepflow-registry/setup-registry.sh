#!/bin/bash

#https://github.com/containerd/nerdctl/releases/download/v2.0.2/nerdctl-2.0.2-linux-amd64.tar.gz
#https://github.com/containerd/nerdctl/releases/download/v2.0.2/nerdctl-full-2.0.2-linux-amd64.tar.gz
#https://github.com/containernetworking/plugins/releases/download/v1.6.2/cni-plugins-linux-amd64-v1.6.2.tgz
#wget https://github.com/containernetworking/plugins/releases/download/v1.6.2/cni-plugins-linux-amd64-v1.6.2.tgz
sudo mkdir -pv /opt/deepflow-registry/config/
sudo cp compose.yaml /opt/deepflow-registry/config/compose.yaml
sudo cp deepflow-registry.yaml /opt/deepflow-registry/config/deepflow-registry.yaml
#运行时为Containerd
sudo CONTAINERD_ADDRESS=/run/k3s/containerd/containerd.sock nerdctl --namespace k8s.io load -i /usr/local/deepflow/registry.tar
sudo CONTAINERD_ADDRESS=/run/k3s/containerd/containerd.sock nerdctl --namespace k8s.io compose -f /opt/deepflow-registry/config/compose.yaml down
sudo CONTAINERD_ADDRESS=/run/k3s/containerd/containerd.sock nerdctl --namespace k8s.io compose -f /opt/deepflow-registry/config/compose.yaml up -d
sudo echo "127.0.0.1 repo.onwalk.net" >> /etc/hosts
