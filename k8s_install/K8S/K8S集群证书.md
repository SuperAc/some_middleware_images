# K8S集群证书

## 重新生成证书：
- kubeadmin

### 检查证书时间
```shell
新版本(1.15+)：kubeadm certs check-expiration
或
openssl x509 -in /etc/kubernetes/pki/apiserver.crt -noout -text |grep ' Not '
其他同理
```

### 证书备份
```shell
cp -rp /etc/kubernetes /etc/kubernetes.bak
```
### 移除过期证书
```shell
rm -f /etc/kubernetes/pki/apiserver*
rm -f /etc/kubernetes/pki/front-proxy-client.*
rm -rf /etc/kubernetes/pki/etcd/healthcheck-client.*
rm -rf /etc/kubernetes/pki/etcd/server.*
rm -rf /etc/kubernetes/pki/etcd/peer.*
备注：可以使用命令openssl x509 -in [证书全路径] -noout -text查看证书详情。
```
### 重新生成证书
```shell
老版本：kubeadm alpha certs renew all
或
新版本(1.15+)：kubeadm certs renew all 使用该命令不用提前删除过期证书
重新生成配置文件
```
### 重新生成配置
```shell
mv /etc/kubernetes/*.conf /tmp/
老版本：kubeadm alpha phase kubeconfig all
或
新版本(1.15+)：kubeadm init phase kubeconfig all
```
### 更新kubectl配置
```shell
cp /etc/kubernetes/admin.conf ~/.kube/config
重启kubelet
systemctl restart kubelet
证书过期时间确认
openssl x509 -in /etc/kubernetes/pki/apiserver.crt -noout -text |grep ' Not '
其他同理
集群确认
kubectl get no
```
# 如果发现集群，能读不能写，则请重启一下组件：
```shell
docker ps | grep apiserver 
docker ps  | grep scheduler
docker ps  | grep controller-manager

docker restart 容器标识
```


