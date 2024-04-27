# 说明

本模块内容与`proxy` 模块类似，也是把集群内部的服务，通过网关暴露出去。

与 `proxy` 模块的区别是，本模块使用的集群，是局域网内部集群，而非供应商的集群。

所以 LoadBalancer 类型 service，没有 external ip，故先要利用 `metallb` 提供 external ip，详见 [metallb 入门说明](./gateway/metallb入门说明.md)

之后就与`proxy` 模块一样了。