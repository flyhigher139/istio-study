# 背景
在供应商提供的云主机上安装 Kubernetes 集群时，如果 Service 的类型是 LoadBalancer，会自动配置 external ip，但在自己的私有局域网中，external ip 始终会显示 pending 状态。
External ip 用于从集群外访问集群服务，如果没有它，只能把 Service 改为 NodePort 类型，直接在集群上把端口暴露出去。
很多场景下，比如配合 istio 的 gateway 时，还是 external ip 更方便，如何让私有局域网中的 k8s 集群具备 external ip，现在有一种生产环境可用的解决方案是，使用 [Metallb](https://metallb.universe.tf/)

# Metallb 使用
## 安装
Metallb 提供了多种[安装方式](https://metallb.universe.tf/installation/)，我们采用最简单的 Installation By Manifest 方式：
```bash
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.14.5/config/manifests/metallb-native.yaml
```

## 使用
要先配置，再使用

### 配置

核心是配置 ip 池，可以配公网 ip，也可以配局域网 ip
我们使用最简单的Layer 2 Configuration即可，配置完成后，创建 Service 时，会自动分配 external ip
本配置只有创建 2 个对象：`IPAddressPool`，`L2Advertisement`

```yaml
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: first-pool
  namespace: metallb-system
spec:
  addresses:
  - 192.168.1.240-192.168.1.250

```
```yaml
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: example
  namespace: metallb-system
spec:
  ipAddressPools:
  - first-pool

```

更多配置，参考[官方文档](https://metallb.universe.tf/configuration/)

### 使用
Metallb 配置完成后，创建 LoadBalancer 类型的 Service 时，会自动分配 external ip，如果想指定 external ip 的值，可以加一个 annotation：`metallb.universe.tf/loadBalancerIPs: 192.168.1.100`，如：
```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx
  annotations:
    metallb.universe.tf/loadBalancerIPs: 192.168.1.100
spec:
  ports:
  - port: 80
    targetPort: 80
  selector:
    app: nginx
  type: LoadBalancer

```
### Samples
[MetalLB, bare metal load-balancer for Kubernetes](https://metallb.universe.tf/usage/example/)

我个人在 istio gateway 中的使用案例为：
配置：
```yaml
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: first-pool
  namespace: metallb-system
spec:
  addresses:
#    - 198.19.249.34/24
#    - 198.19.249.34/24
    - 198.19.249.35-198.19.249.40
#    - 198.19.249.1-198.19.249.2

---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: first-advertisement
  namespace: metallb-system
spec:
  ipAddressPools:
    - first-pool
```
应用：
```yaml
apiVersion: v1
kind: Service
metadata:
  name: middleware-gateway-mysql-svc
  namespace: istio-system
  annotations:
    metallb.universe.tf/loadBalancerIPs: 198.19.249.35
spec:
  type: LoadBalancer
  selector:
    app: istio-ingressgateway
  ports:
    - port: 3310
      targetPort: 3310
      name: mysql-middleware
```
结果：
```bash
➜  gateway git:(main) rkectl -n istio-system get svc
NAME                           TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                                      AGE
istio-ingressgateway           LoadBalancer   10.43.185.90    198.19.249.36   15021:30570/TCP,80:31676/TCP,443:30602/TCP   5h14m
istiod                         ClusterIP      10.43.136.114   <none>          15010/TCP,15012/TCP,443/TCP,15014/TCP        5h15m
middleware-gateway-mysql-svc   LoadBalancer   10.43.34.15     198.19.249.35   3310:30590/TCP                               4h9m

```
