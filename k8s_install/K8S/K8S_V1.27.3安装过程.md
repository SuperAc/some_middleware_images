# K8S_V1.27.3
# 安装前准备 
准备三到五台服务器
参考文档：`https://kubernetes.io/zh-cn/docs/setup/production-environment/tools/kubeadm/install-kubeadm/`
说明：此次部署由DOCKER作为容器运行时

|  节点   | 资源需求 |           IP           |
| ------ | -------- | ---------------------- |
| Master | 2C2G80G  | 192.168.127.10、20、30 |
| Node   | 2C2G60G  | 192.168.127.40、50     |

# 安装前配置
## 相关环境修改配置
所有机器都需要修改
- vi /etc/hostname
```
master1 （对应自己的hostname）
```
- vi /etc/hosts
```
192.168.127.10 master1
192.168.127.20 master2
192.168.127.30 master3
192.168.127.40 node1
192.168.127.50 node2
```
- shutdonw -r
- 修改iptables
```
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
br_netfilter
EOF

cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sudo sysctl --system
```
- 关闭防火墙
```
systemctl stop firewalld
systemctl disable firewalld
systemctl status firewalld
```
- 关闭SELinux
```
vi /etc/selinux/config

# This file controls the state of SELinux on the system.
# SELINUX= can take one of these three values:
#     enforcing - SELinux security policy is enforced.
#     permissive - SELinux prints warnings instead of enforcing.
#     disabled - No SELinux policy is loaded.
#SELINUX=enforcing
# SELINUXTYPE= can take one of three values:
#     targeted - Targeted processes are protected,
#     minimum - Modification of targeted policy. Only selected processes are protected.
#     mls - Multi Level Security protection.
#SELINUXTYPE=targeted 
SELINUX=disabled 
```
- 关闭Swap分区
```
// temp turn off
 swapoff -a
// Permanently turn off
vi /etc/fstab

#
# /etc/fstab
# Created by anaconda on Sun Aug 22 00:58:48 2021
#
# Accessible filesystems, by reference, are maintained under '/dev/disk'
# See man pages fstab(5), findfs(8), mount(8) and/or blkid(8) for more info
#
/dev/mapper/centos-root /                       xfs     defaults        0 0
UUID=be6fcb6b-d425-4cc6-9e97-0bba2b1c7236 /boot                   xfs     defaults        0 0
/dev/mapper/centos-home /home                   xfs     defaults        0 0
#/dev/mapper/centos-swap swap                    swap    defaults        0 0  // Comment out the current line

shutdown -r
```
- 同步时区
```
yum install ntp
systemctl enable ntpd
systemctl start ntpd
timedatectl set-timezone Asia/Shanghai
timedatectl set-ntp yes
ntpq -p
```
## 安装docker
可以参考[官方教程](https://docs.docker.com/engine/install/centos/)
安装docker完成后，可以配置一下docker的镜像加速源
```
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": ["加速的镜像地址"]
}
EOF
sudo systemctl daemon-reload
sudo systemctl restart docker
```
## 安装cri-dockerd
由于高版本K8S不支持docker作为容器运行时，所以需要安装cri-dockerd
- 安装golang
可以参考[官方教程](https://go.dev/doc/install)
```
cd /usr/local
# 下载 go1.20.5.linux-amd64.tar.gz ,也可以本地下载完后，上传到服务器/usr/local目录下
wget https://go.dev/dl/go1.20.5.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin  # 可以把他放到类似/etc/profile文件中
go version
```
接着安装cri-dockerd
项目地址：https://github.com/Mirantis/cri-dockerd
```
git clone https://github.com/Mirantis/cri-dockerd.git
cd cri-dockerd
mkdir bin
go build -o bin/cri-dockerd   # 这一步比较耗时，耐心等待
mkdir -p /usr/local/bin
install -o root -g root -m 0755 bin/cri-dockerd /usr/local/bin/cri-dockerd
cp -a packaging/systemd/* /etc/systemd/system
sed -i -e 's,/usr/bin/cri-dockerd,/usr/local/bin/cri-dockerd,' /etc/systemd/system/cri-docker.service
systemctl daemon-reload
systemctl enable cri-docker.service
systemctl enable --now cri-docker.socket
```
## 安装kubelet、kubeadm、kubectl#
https://developer.aliyun.com/mirror/kubernetes/
```
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
setenforce 0
yum install -y kubelet kubeadm kubectl
systemctl enable kubelet && systemctl start kubelet
```
## 主节点部署
由于国内访问不了registry.k8s.io，所以手动拉一下前面所需要的镜像
```kubeadm config images pull --image-repository registry.cn-hangzhou.aliyuncs.com/google_containers --cri-socket=unix:///var/run/cri-dockerd.sock```
如果你是用containerd，最后的cri-socket可以不用指定，但是后续registry.k8s.io/pause:3.6这个镜像我没搞清楚怎么本地下载然后由containerd管理
下载完成后，建议把registry.k8s.io/pause:3.6也下载一下
```docker pull docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.6```
 k8s安装的时候会获取registry.k8s.io前缀的镜像，所以修改一下之前安装的镜像tag
 ```
 docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/kube-apiserver:v1.27.3 \
registry.k8s.io/kube-apiserver:v1.27.3

docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/kube-scheduler:v1.27.3 \
registry.k8s.io/kube-scheduler:v1.27.3

docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/kube-controller-manager:v1.27.3 \
registry.k8s.io/kube-controller-manager:v1.27.3

docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/kube-proxy:v1.27.3 \
registry.k8s.io/kube-proxy:v1.27.3

docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/coredns:v1.10.1 \
registry.k8s.io/coredns:v1.10.1

docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/etcd:3.5.7-0 \
registry.k8s.io/etcd:3.5.7-0 

docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.6 \
registry.k8s.io/pause:3.6

-------------------------- 这个pause3.6和3.9建议都下载，我部署的时候总出问题
docker tag registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.9 \
registry.k8s.io/pause:3.9
 ```
 使用kubeadm进行初始化
 ```
 kubeadm init \
--apiserver-advertise-address=192.168.127.10 \
--control-plane-endpoint=master1 \
--image-repository registry.cn-hangzhou.aliyuncs.com/google_containers \
--cri-socket=unix:///var/run/cri-dockerd.sock \
--pod-network-cidr=172.16.0.0/16 \
--v=5
```
`piserver-advertise-address`写master节点的IP
`control-plane-endpoint`写master节点的hostname
`pod-network-cidr`是pod节点的网段
安装完成后，有一个kubeadm join ***,这是其他节点加入k8s集群的命令
生产kubeconfig文件
```
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```
## 安装网络插件Calico（也可以安装flannel）
参考[官方教程](https://docs.tigera.io/calico/latest/getting-started/kubernetes/quickstart)
你可以把`tigera-operator.yaml`和`custom-resources.yaml`下载到本地上传到服务器上
注意记得修改`custom-resources.yaml`中的`cidr`，为你自己的POD网段


## 加入node节点
由于使用的是cri-dockerd，所以需要加入`--cri-socket=unix:///var/run/cri-dockerd.sock`
kubeadm join ***  --cri-socket=unix:///var/run/cri-dockerd.sock


## 加入master节点
master节点上输入
```
kubeadm init phase upload-certs --upload-certs
#这里会生成一个certificate key

kubeadm token create --print-join-command

```
在master节点上，最后有一条 kubeadm join***，之后master节点加入的时候，在这条命令后面加上
```--control-plane --certificate-key 之前生产的certificate key```

登录其他master节点
```
kubeadm join ***  --control-plane --certificate-key 之前生产的certificate key
```

Kubectl命令补全
```
yum install -y bash-completion
source /usr/share/bash-completion/bash_completion
echo 'source <(kubectl completion bash)' >> /etc/profile
source /etc/profile
```