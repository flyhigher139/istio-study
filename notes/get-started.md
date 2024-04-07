# Istio Get Started

## 1. 如何在部署应用的时候，注入Envoy

最简单的方法是，给命名空间添加标签，指示 Istio 在部署应用的时候，自动注入 Envoy 边车代理：

```shell
kubectl label namespace default istio-injection=enabled

```