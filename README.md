# my-go-zero-looklook

本项目是主要参考了 [go-zero-looklook](https://github.com/Mikaelemmmm/go-zero-looklook)  的学习与实践项目。

## 项目简介

本项目采用微服务架构，基于 `go-zero` 框架开发。主要包含以下功能模块：

*   **API 网关**: 使用 Nginx 作为外部网关。
*   **业务服务**:
    *   用户服务 (Usercenter)
    *   民宿服务 (Travel)
    *   订单服务 (Order)
    *   支付服务 (Payment)
    *   定时任务服务 (Mqueue)
*   **可观测性**:
    *   日志: Filebeat + Kafka + Go-stash + Elasticsearch + Kibana
    *   监控: Prometheus + Grafana
    *   链路追踪: Jaeger

## 目录结构

*   `app`: 业务代码 (API, RPC, MQ)
*   `common` / `pkg`: 通用组件
*   `deploy`: 部署相关配置 (Filebeat, Nginx, Prometheus 等)

## 本地部署

1. docker安装软件依赖, docker-compose up -d
2. 启动微服务 make run SERVICE="xxx"

## Kubernetes 部署

本项目的 Kubernetes 部署配置文件托管在独立的仓库中，配合 Jenkins 实现自动化 CI/CD。

**K8s 配置仓库**: [https://github.com/wujunhui99/looklook-pro-conf](https://github.com/wujunhui99/looklook-pro-conf) 仓库也包含Jenkinsfile

该仓库包含了各微服务的 K8s 配置 YAML 文件：
*   `usercenter`
*   `travel`
*   `order`
*   `payment`
*   `mqueue`

### 部署流程简述

1.  **代码提交**: 提交代码到 Git 仓库。
2.  **Jenkins 构建**: Jenkins 拉取代码，构建 Docker 镜像并推送至 Harbor 仓库。
3.  **配置同步**: Jenkins 拉取 `looklook-pro-conf` 中的 K8s 配置。
4.  **自动部署**: 使用 `kubectl` 将服务发布到 Kubernetes 集群。

## 致谢

*   [go-zero](https://github.com/zeromicro/go-zero)
*   [go-zero-looklook](https://github.com/Mikaelemmmm/go-zero-looklook)
