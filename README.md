# ansible-playbook

# Playook 角色说明

- Common 部分
  - common	        通用角色，包含一些常用的功能，如日志记录、监控等。
  - secret-manger	密钥管理角色，用于管理密钥。
  - cert-manager	证书管理角色，用于管理证书。
  - k3s	                用于创建 Kubernetes 集群。
  - k3s-reset	        用于重置 Kubernetes 集群。
  - k3s-addon	        用于安装 Kubernetes 集群插件。
  - app	                应用程序服务角色，提供应用程序运行所需的服务，如 Nginx、App_backend 等。
  - redis	        Redis 数据库角色，用于提供 Redis 数据库服务。
  - postgresql	        PostgreSQL 数据库角色，用于提供 PostgreSQL 数据库服务。

- DevOPSPlatform 使用的角色列表：
  - chartmuseum	        图表仓库角色，用于存储和管理 Kubernetes 图表。
  - gitlab	        代码仓库角色，用于存储和管理代码。
  - harbor	        容器镜像仓库角色，用于存储和管理容器镜像。
  - mysql	        MySQL 数据库角色，用于提供 MySQL 数据库服务。
- Federated-IdentityProvider 使用的角色列表：
  - openldap
  - keycloak
- ObservabilityPlatform 使用的角色列表：
  - alerting
  - clickhouse
  - node-exporter
  - observability-agent
  - observability-server
  - prometheus-transfer
  - promtail-agent

ansible-playbook  -i inventory gather_network_info.yml -e target_group=master
ansible -i inventory all -m script -a 'roles/network_info/tasks/files/display_network_info.sh'


ansible-playbook -i inventory/hosts/core playbooks/keycloak_server -D


- Setup Wireguard Gateway: 

``` ansible-playbook -i inventory/hosts/vpn playbooks/wireguard_gateway.yaml -D ```


# all_in_one scripts

- curl -sfL https://mirrors.onwalk.net/public/k3s_setup.sh | bash -

https://www.henryxieblogs.com/2021/12/how-to-expose-kube-api-server-via-nginx.html

