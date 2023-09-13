# containerdK8s
```
30  vi /etc/hosts
   31  systemctl stop firewalld
   32  systemctl disable firewalld
   33  setenforce 0
   34  vi /etc/selinux/
   35  vi /etc/selinux/config
    SELINUX=disabled
   36  modprobe br_netfilter
   37  vi /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
# 下面的内核参数可以解决ipvs模式下长连接空闲超时的问题
net.ipv4.tcp_keepalive_intvl = 30
net.ipv4.tcp_keepalive_probes = 10
net.ipv4.tcp_keepalive_time = 600
vm.swappiness=0

   

   38  sysctl -p /etc/sysctl.d/k8s.conf
   39  cat > /etc/sysconfig/modules/ipvs.modules <<EOFmodprobe -- ip_vs
modprobe -- ip_vs_rr
modprobe -- ip_vs_wrr
modprobe -- ip_vs_sh
modprobe -- nf_conntrack_ipv4
EOF

   40  chmod 755 /etc/sysconfig/modules/ipvs.modules && bash /etc/sysconfig/modules/ipvs.modules && lsmod | grep -e ip_vs -e nf_conntrack_ipv4chmod 755 /etc/sysconfig/modules/ipvs.modules && bash /etc/sysconfig/modules/ipvs.modules && lsmod | grep -e ip_vs -e nf_conntrack_ipv4
   41  chmod 755 /etc/sysconfig/modules/ipvs.modules && bash /etc/sysconfig/modules/ipvs.modules && lsmod | grep -e ip_vs -e nf_conntrack_ipv4

   43  yum install ipset ipvsadm chrony  -y
   44  systemctl enable chronydsystemctl enable chronyd
   45  systemctl enable chronyd
   46  systemctl start chronyd
   47  chronyc sources
   48  data
   49  date
   50  swapoff -a
   51  vi /etc/fstab
   注释掉swap分区
   52  vi /etc/sysctl.d/k8s.conf
   53  sysctl -p /etc/sysctl.d/k8s.conf

   55  mkdir -p /etc/containerd
   56  containerd config default > /etc/containerd/config.toml
   57  vi /etc/containerd/config.toml
   加镜像
   58  systemctl daemon-reload
   59  systemctl enable containerd --now
   60  ctr version
   61  cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF

   62  setenforce 0
   63  yum install -y kubelet kubeadm kubectl
   64  systemctl enable kubelet && systemctl start kubelet
```
```
由于开启内核 ipv4 转发需要加载 br_netfilter 模块，所以加载下该模块：

➜  ~ modprobe br_netfilter

最好将上面的命令设置成开机启动，因为重启后模块失效，下面是开机自动加载模块的方式。首先新建 /etc/rc.sysinit 文件，内容如下所示：

#!/bin/bash
for file in /etc/sysconfig/modules/*.modules ; do
[ -x $file ] && $file
done

然后在 /etc/sysconfig/modules/ 目录下新建如下文件：

➜  ~ cat /etc/sysconfig/modules/br_netfilter.modules
modprobe br_netfilter

增加权限：

➜  ~ chmod 755 br_netfilter.modules
```
k8s安装
- 初始化集群
`kubeadm config print init-defaults --component-configs KubeletConfiguration > kubeadm.yaml`
如果自定义的` KubeletConfiguration API` 对象使用像 `kubeadm ... --config some-config-file.yaml` 这样的配置文件进行传递，则可以配置 kubeadm 启动的 kubelet。
通过调用` kubeadm config print init-defaults --component-configs KubeletConfiguration`， 你可以看到此结构中的所有默认值。
也可以在基础 KubeletConfiguration 上应用实例特定的补丁

主机修改 主机节点，主机名称、加上Pod网段等
如：
```
apiVersion: kubeadm.k8s.io/v1beta3
bootstrapTokens:
- groups:
  - system:bootstrappers:kubeadm:default-node-token
  token: abcdef.0123456789abcdef
  ttl: 24h0m0s
  usages:
  - signing
  - authentication
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: 192.168.127.100 #master主机ip
  bindPort: 6443
nodeRegistration:
  criSocket: unix:///var/run/containerd/containerd.sock
  imagePullPolicy: IfNotPresent
  name: master #主机hostname
  taints: null
---
#因为用了ipvs，所以加了一个ipvs的配置
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: ipvs  # kube-proxy 模式

---
apiServer:
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta3
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
controllerManager: {}
dns: {}
etcd:
  local:
    dataDir: /var/lib/etcd
imageRepository: registry.cn-hangzhou.aliyuncs.com/google_containers #加速镜像地址
kind: ClusterConfiguration
kubernetesVersion: 1.27.0
networking:
  dnsDomain: cluster.local
  serviceSubnet: 10.96.0.0/12
  podSubnet: 10.244.0.0/16 #pod网段
scheduler: {}
---
apiVersion: kubelet.config.k8s.io/v1beta1
authentication:
  anonymous:
    enabled: false
  webhook:
    cacheTTL: 0s
    enabled: true
  x509:
    clientCAFile: /etc/kubernetes/pki/ca.crt
authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: 0s
    cacheUnauthorizedTTL: 0s
cgroupDriver: systemd  #确认cgroupDriver是ETCD
clusterDNS:
- 10.96.0.10
clusterDomain: cluster.local
containerRuntimeEndpoint: ""
cpuManagerReconcilePeriod: 0s
evictionPressureTransitionPeriod: 0s
fileCheckFrequency: 0s
healthzBindAddress: 127.0.0.1
healthzPort: 10248
httpCheckFrequency: 0s
imageMinimumGCAge: 0s
kind: KubeletConfiguration
logging:
  flushFrequency: 0
  options:
    json:
      infoBufferSize: "0"
  verbosity: 0
memorySwap: {}
nodeStatusReportFrequency: 0s
nodeStatusUpdateFrequency: 0s
rotateCertificates: true
runtimeRequestTimeout: 0s
shutdownGracePeriod: 0s
shutdownGracePeriodCriticalPods: 0s
staticPodPath: /etc/kubernetes/manifests
streamingConnectionIdleTimeout: 0s
syncFrequency: 0s
volumeStatsAggPeriod: 0s

```
开始安装
```
kubeadmin init --config kubeadm.yaml

期间安装失败，查看日志发现需要registry.k8s.io/pause:3.8这个镜像，使用ctr下载，需要指定一下namespace

[root@master k8s]# ctr -n k8s.io i pull  registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.8
registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.8:                    resolved       |++++++++++++++++++++++++++++++++++++++|
index-sha256:9001185023633d17a2f98ff69b6ff2615b8ea02a825adffa40422f51dfdcde9d:    done           |++++++++++++++++++++++++++++++++++++++|
manifest-sha256:f5944f2d1daf66463768a1503d0c8c5e8dde7c1674d3f85abc70cef9c7e32e95: done           |++++++++++++++++++++++++++++++++++++++|
layer-sha256:9457426d68990df190301d2e20b8450c4f67d7559bdb7ded6c40d41ced6731f7:    done           |++++++++++++++++++++++++++++++++++++++|
config-sha256:4873874c08efc72e9729683a83ffbb7502ee729e9a5ac097723806ea7fa13517:   done           |++++++++++++++++++++++++++++++++++++++|
elapsed: 0.7 s                                                                    total:  2.7 Ki (3.8 KiB/s)
unpacking linux/amd64 sha256:9001185023633d17a2f98ff69b6ff2615b8ea02a825adffa40422f51dfdcde9d...
done: 50.643599ms
[root@master k8s]# ctr -n k8s.io  image tag registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.8 registry.k8s.io/pause:3.8
registry.k8s.io/pause:3.8

```


常见问题
node节点无法使用kubectl
把master节点下`/etc/kubernetes/admin.conf` 复制到node节点相同的目录下
```
mkdir -p $HOME/.kube
#保存 admin.conf 到 $HOME/.kube/config这个目录下
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
#授予权限
sudo chown $(id -u):$(id -g) $HOME/.kube/config

```

k8s拉不下来
使用ctr命令将镜像拉下来


安装nerdctl工具
```
mkdir -p /usr/local/containerd/bin/ && tar -zxvf nerdctl-1.4.0-linux-amd64.tar.gz nerdctl && mv nerdctl /usr/local/containerd/bin/
``