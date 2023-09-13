# 13_pv&pvc
## PV和PVC

真正的存储是PV（或者说PV关联的资源，如hostpath，nfs等），pvc只是一个对象
###  hostpath
#### 创建PV
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-hostpath  #全局的pv，没有namespace限制
  labels:
    type: local
spec:
  hostPath:
    path: /root/k8s/11pv_pvc/01hostpath/html
  capacity: 
    storage: 10Gi
  accessModes:
  - ReadWriteOnce
  storageClassName: manual
```
#### 对应的PVC
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-hostpath
spec:
  resources:
    requests:
      storage: 3Gi
  accessModes:
    - ReadWriteOnce
  storageClassName: manual
  selector:
    matchLabels:
      type: local
```
PV和PVC对应，需要storageClassName和accessModes一致
如果有多个pv满足要求，这个时候随机绑定一个满足条件的PV。如果需要指定绑定PV,如上使用selector中的matchLabes
#### 在pod中使用pvc
注意，pod被调度的节点需要有pv中的路径，示例中pv的路径在master中存在，所以pod如下：
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pv-pod
  labels:
    name: pv-pod
spec:
  containers:
  - name: pv-pod
    image: nginx:latest
    resources:
      limits:
        memory: "128Mi"
        cpu: "500m"
    ports:
      - containerPort: 80
    volumeMounts:
      - name:  pvc-hostpath-storage
        mountPath:  /usr/share/nginx/html
  tolerations:
    - key: node-role.kubernetes.io/control-plane
      operator: Exists
      effect: NoSchedule
    - key: node-role.kubernetes.io/master
      operator: Exists
      effect: NoSchedule
  nodeSelector:
    kubernetes.io/hostname: master

  volumes:
    - name:  pvc-hostpath-storage
      persistentVolumeClaim:
        claimName: pvc-hostpath
        readOnly: false
```
以上，hostpath具有局限性，pod不能飘逸，必须固定到一个节点上，一旦飘逸到其他宿主机上就没有了对应的数据，所以在使用hostpath 的pv时都会搭配nodeSeletor指定宿主机。优点是hostpath的pv直接使用的本地磁盘，读写性能相比于大多数远程存储来说要好很多。不过对应正常PV来说，使用了hostpath的节点一旦宕机，数据有可能会丢失，所以要求使用hostpath的应用必须具备数据备份和恢复的能力。
### Local PV
在hostpath的基础上，K8S依靠pv和pvc实现了一个新的特性，叫Local PV[#local](https://kubernetes.io/zh-cn/docs/concepts/storage/volumes/#local)
LocalPV实现的功能类似hostPath加上nodeAffinity。
在实际使用过程中，不能将一个宿主机上的目录当作PV来使用，因为本地目录的存储行为不可控，所以在的磁盘可能会被应用写满，甚至造成宿主机宕机。所以一般LocalPV对应的存储介质时一块额外挂载在宿主上的磁盘或者块设备，认为“一个pv一个盘”

local 卷所代表的是某个被挂载的本地存储设备，例如磁盘、分区或者目录。

local 卷只能用作静态创建的持久卷。不支持动态配置。

与 hostPath 卷相比，local 卷能够以持久和可移植的方式使用，而无需手动将 Pod 调度到节点。系统通过查看 PersistentVolume 的节点亲和性配置，就能了解卷的节点约束。
对于普通PV来说，K8S都是先调度Pod到某个节点上，再持久化节点上的volume目录，完成Volume目录与容器的绑定挂载，但是对于localPV来说，节点上可供使用的磁盘必须先准备好，因为不同节点上的挂载情况可能不能，甚至有些节点没有这个磁盘，所以调度器必须知道所有节点和localPv对于磁盘的关联关系，然后根据这个信息来调度Pod。

然而，local 卷仍然取决于底层节点的可用性，并不适合所有应用程序。 如果节点变得不健康，那么 local 卷也将变得不可被 Pod 访问。使用它的 Pod 将不能运行。 使用 local 卷的应用程序必须能够容忍这种可用性的降低，以及因底层磁盘的耐用性特征而带来的潜在的数据丢失风险。

#### 创建PV
下面测试Local PV，暂时将node1上的/root/k8s/localpv/看作独立的磁盘，下面来声明对应的local PV
```yaml
#pv-local.yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: local-storage
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-local
  labels:
    type: local
spec:
  capacity:
    storage: 10Gi
  volumeMode: Filesystem
  accessModes:
    - ReadWriteOnce
  storageClassName: local-storage
  persistentVolumeReclaimPolicy: Delete
  local:
    path: /root/k8s/localpv
  nodeAffinity:
    required:
      nodeSelectorTerms:
       - matchExpressions:
           - key: kubernetes.io/hostname
             operator: In
             values: 
              - node1
```
因为pvc一旦创建，就会和对应的pv绑定，所以使用 StorageClass，将volumeBindingMode设置成WaitForFirstConsumer，当消费pvc的时候才会绑定
#### 创建pvc
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-local
spec:
  resources:
    requests:
      storage: 3Gi
  volumeMode: Filesystem
  accessModes:
    - ReadWriteOnce
  storageClassName: local-storage
  selector:
    matchLabels:
      type: local
```
#### 创建pod
```
apiVersion: apps/v1
kind: Deployment
metadata: 
  name: pv-local-dep
  labels:
    app: pv-local
spec:
  replicas: 2
  selector:
    matchLabels:
      app: pv-local
  template:
    metadata:
      name: pv-local-pod
      labels:
        app: pv-local
    spec:
      containers:
      - name: pv-local-container
        image: nginx
        ports:
          - containerPort: 80
        volumeMounts:
          - mountPath: /usr/share/nginx/html
            name: pv-local-storage
      volumes:
        - name: pv-local-storage
          persistentVolumeClaim:
            claimName: pvc-local
```

使用localPV的时候，默认就会调度到pv指定的节点上

注意：删除PV时的流程：
- 删除使用这个pvc的pod
- 从宿主机上移除本地磁盘
- 删除pvc
- 删除pv


###常用localPV的使用场景
kafka、es、etcd