# ansible-playbook

This repository contains a collection of Ansible playbooks and roles for various infrastructure setups and service management tasks.

## Playbook 角色说明

1. playbooks/roles/docker：适用于简单的、单机环境的部署，主要使用 Docker 和 Docker Compose 进行容器化管理。
2. playbooks/roles/charts：面向大规模的 Kubernetes 集群，使用 Helm 和标准化 Chart 部署模式进行高可用和可扩展的管理。
3. playbooks/roles/vhosts：传统的非容器化部署方式，通常涉及手动配置服务器和虚拟主机，适用于不使用容器的应用场景。


## Role Summary

| Role Name               | Description                                           | Docker | Charts | VHosts | CICD    | Validate | Last Update  |
|-------------------------|-------------------------------------------------------|--------|--------|--------|---------|----------|--------------|
| `common`                | 通用角色，包含一些常用的功能，如日志记录、监控等。      |        |        |   ✔    |         |   yes    | 2025-02-14   |
| `keycloak`              | 用于管理身份认证和授权服务。                            |   ✔    |        |        | github  |   yes    | 2024-11-10   |
| `harbor`                | 容器镜像仓库角色，用于存储和管理容器镜像。              |   ✔    |        |        | github  |   yes    | 2024-11-14   |
| `app`                   | 参考模板。                                              |        |        |        |         |          |              |
| `nginx`                 | 用于设置 Nginx                                          |        |   ✔    |   ✔    |         |          |              |
| `grafana`               | 用于设置 Grafana                                        |        |   ✔    |   ✔    |         |          |              |
| `grafana-loki`          | 用于设置 Grafana-loki                                   |        |   ✔    |   ✔    |         |          |              |
| `Grafana-tempo`         | 用于设置 Grafana-tempo                                  |        |   ✔    |   ✔    |         |          |              |
| `prometheus`            | 用于设置 Prometheus                                     |        |   ✔    |   ✔    |         |          |              |
| `prometheus-transfer`   | 用于 Prometheus 数据传输设置。                          |        |        |   ✔    |         |          |              |
| `vector`                | 用于配置日志收集代理。                                  |        |        |   ✔    |         |          |              |
| `node-exporter`         | 用于导出系统和硬件的监控数据。                          |        |   ✔    |        |         |          |              |
| `observability-agent`   | 用于管理 Observability 代理。                           |        |   ✔    |   ✔    |         |          |              |
| `observability-server`  | 用于设置 Observability 服务端。                         |        |   ✔    |   ✔    |         |          |              |
| `wireguard-client`      | 用于设置 WireGuard 客户端。                             |        |        |   ✔    |         |          |              |
| `wireguard-gateway`     | 用于设置 WireGuard 网关。                               |        |        |   ✔    |         |          |              |
| `vault`                 | 用于管理敏感数据和密钥。                                |        |        |   ✔    |         |          |              |
| `postgresql`            | PostgreSQL 数据库角色，用于提供 PostgreSQL 数据库服务。 |        |   ✔    |        |         |          |              |
| `redis`                 | Redis 数据库角色，用于提供 Redis 数据库服务。           |        |   ✔    |        |         |          |              |
| `chartmuseum`           | 图表仓库角色，用于存储和管理 Kubernetes 图表。          |        |   ✔    |        |         |          |              |
| `gitlab`                | 代码仓库角色，用于存储和管理代码。                      |        |   ✔    |        |         |          |              |
| `mysql`                 | MySQL 数据库角色，用于提供 MySQL 数据库服务。           |        |   ✔    |        |         |          |              |
| `argo-server`           | 用于设置和管理 Argo Server。                            |        |   ✔    |        |         |          |              |
| `deepflow`              | 用于流量监控与网络性能分析的 DeepFlow 服务。            |        |   ✔    |        |         |          |              |
| `jenkins`               | Jenkins 自动化构建工具角色，用于 CI/CD 管道。           |        |   ✔    |        |         |          |              |
| `chaos-mesh`            | 用于 Chaos Engineering 测试的 Chaos Mesh 角色。         |        |   ✔    |        |         |          |              |
| `flagger-loadtester`    | 用于负载测试的 Flagger Loadtester 角色。                |        |   ✔    |        |         |          |              |
| `splunk-otel-collector` | 用于配置 Splunk OpenTelemetry Collector。               |        |   ✔    |        |         |          |              |
| `openldap`              | 用于设置和管理 OpenLDAP 身份认证服务。                  |        |   ✔    |        |         |          |              |
| `alerting`              | 用于设置和管理警报系统。                                |        |        |   ✔    |         |          |              |
| `k3s`                   | 用于创建 Kubernetes 集群。                              |        |        |   ✔    |         |          |              |
| `k3s-reset`             | 用于重置 Kubernetes 集群。                              |        |        |   ✔    |         |          |              |
| `k3s-addon`             | 用于安装 Kubernetes 集群插件。                          |        |        |   ✔    |         |          |              |
| `secret-manger`         | 密钥管理角色，用于管理密钥。                            |        |        |   ✔    |         |          |              |
| `cert-manager`          | 证书管理角色，用于管理证书。                            |        |        |   ✔    |         |          |              |

表格说明
- Docker：是否属于 Docker 角色。
- Charts：是否属于 Helm Chart 角色。
- VHosts：是否属于虚拟主机管理相关角色。
- CICD：是否启用 CICD 管道，标明是否集成了自动化流程。
- Validate：是否经过验证测试。
- Last Update：最后更新时间。

##  Usage Examples

- Linux OS Setup

ansible-playbook -i inventory/hosts/all playbooks/common -D -C
ansible-playbook -i inventory/hosts/all playbooks/common -D

- Gather Network Information

ansible-playbook -i inventory gather_network_info.yml -e target_group=master

- Display network information on all nodes

ansible -i inventory all -m script -a 'roles/network_info/tasks/files/display_network_info.sh'

- Deploy Keycloak Server

ansible-playbook -i inventory/hosts/core playbooks/keycloak_server -D

- Set up WireGuard Gateway

ansible-playbook -i inventory/hosts/vpn playbooks/wireguard_gateway.yaml -D

- Set up Grafana Alloy

ansible-playbook -i inventory/k3s-cluster playbooks/init_grafana_alloy -D -C -l cn-k3s-server.svc.plus -e @playbooks/roles/alloy/files/loki_journal_sources_k3s_server.yml -e "ansible_become_pass='xxxx'"
