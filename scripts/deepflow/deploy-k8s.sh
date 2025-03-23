
# 新建部署目录，并解压安装包到该目录

mkdir /opt/k8s-deploy && tar -xvpf sealos-amd64-k8s-1.25.16.tar.gz -C /opt/k8s-deploy
cd /opt/k8s-deploy/                        && \
cp sealos helm calicoctl nerdctl /usr/bin/ && \
chmod +x /usr/bin/sealos /usr/bin/helm /usr/bin/calicoctl /usr/bin/nerdctl

# 导入离线镜像
sealos load -i sealos-calico.tar
sealos load -i sealos-helm.tar
sealos load -i sealos-k8s-1.25.16.tar

# 单机部署(单机部署无需ssh密码，root用户本机直接执行即可)
sealos run \
    registry.cn-shanghai.aliyuncs.com/labring/kubernetes:v1.25.16  \
    registry.cn-shanghai.aliyuncs.com/labring/helm:v3.9.4          \
    registry.cn-shanghai.aliyuncs.com/labring/calico:v3.24.1       \
    --single
