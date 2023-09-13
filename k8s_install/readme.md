# K8S some note
# 使用HA-Proxy和KeepAlived安装高可用K8S 注意事项

1. 使用`kubeadm config print init-defaults --component-configs KubeletConfiguration > kubeadm.yaml`生成的kubeadm.yaml中，如下配置需要修改
```yaml
localAPIEndpoint:
  advertiseAddress: 1.2.3.4  //KeepAlived的VIP
  bindPort: 6443  //HA-Proxy的端口
```
其他内容可以参考文档