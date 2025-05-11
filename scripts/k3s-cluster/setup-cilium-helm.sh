API_SERVER_IP=10.253.253.1
# Kubeadm default is 6443
API_SERVER_PORT=6443
helm install cilium cilium/cilium --version 1.17.3 \
    --namespace kube-system \
    --set kubeProxyReplacement=true \
    --set k8sServiceHost=${API_SERVER_IP} \
    --set k8sServicePort=${API_SERVER_PORT}
