# 07_cm_secret
对于应用的配置管理、敏感信息的存储和使用，容器运行资源的配置、安全管理、身份认证等应用的可变配置，K8S通过configmap资源对象来实现的。日常使用中，应用经常会从配置文件，命令行参数或者环境变量中读取一些配置信息。这些配置信息一般不会写死在程序中，configmap提供了向容器注意配置信息的能力，可以保存单个属性，以及整个配置文件。
## ConfigMap
`ConfigMap`使用`key-value`的形式来配置数据
```yaml
apiVersion:v1
kind: ConfigMap
metadata:
  name: cm-demo
data:
  data.1: hello
  data.2: world
  config: |
    property.1=value-1
    property.2=value-2
    property.3=value-3
```
配置数据在data下进行配置，前两个用来保存单个属性，之后的config用来保存配置文件。config后面的竖线`|`，在yaml中表示保留换行，每行的缩进和行尾空白都会被去掉，而额外的缩进会被保留。
```yaml
lines: |
  我是第一行
  我是第二行
    我是吴彦祖
      我是第四行
  我是第五行

# JSON
{"lines": "我是第一行\n我是第二行\n  我是吴彦祖\n     我是第四行\n我是第五行"}
```
除了竖线`|`之外，还有`>`右尖括号，宝石折叠换行，只有空白行才会被识别为换行，原来的换行符都会被转换成空格。
```yaml
lines: >
  我是第一行
  我也是第一行
  我仍是第一行
  我依旧是第一行

  我是第二行
  这么巧我也是第二行

# JSON
{"lines": "我是第一行 我也是第一行 我仍是第一行 我依旧是第一行\n我是第二行 这么巧我也是第二行"}
```
还可以使用竖线和加号或者减号进行配合使用，+ 表示保留文字块末尾的换行，- 表示删除字符串末尾的换行。
```yaml
value: |
  hello

# {"value": "hello\n"}

value: |-
  hello

# {"value": "hello"}

value: |+
  hello

# {"value": "hello\n\n"} (有多少个回车就有多少个\n)
```
除了配置文件，还可以使用`kubectl create cm(configmap的简写)`来创建configmap。通过使用`kubectl create cm --help`查看相关的帮助信息
```sh
Examples:
  # Create a new config map named my-config based on folder bar 根据文件创建一个名为my-config的configmap
  kubectl create configmap my-config --from-file=path/to/bar

  # Create a new config map named my-config with specified keys instead of file basenames on disk 创建时使用磁盘上文件代替key的value
  kubectl create configmap my-config --from-file=key1=/path/to/bar/file1.txt --from-file=key2=/path/to/bar/file2.txt

  # Create a new config map named my-config with key1=config1 and key2=config2 直接指定key-value
  kubectl create configmap my-config --from-literal=key1=config1 --from-literal=key2=config2

  # Create a new config map named my-config from the key=value pairs in the file。根据文件中的key=value对创建名为my-config的新配置映射。
  kubectl create configmap my-config --from-file=path/to/bar

  # Create a new config map named my-config from an env file(会解析env文件中的key-value)
  kubectl create configmap my-config --from-env-file=path/to/foo.env --from-env-file=path/to/bar.env
```
比如下面有个`testcm`目录，里面有两个配置文件
```sh
ls testcm
redis.conf
mysql.conf

cat testcm/redis.conf
host=127.0.0.1
port=6379

cat testcm/mysql.conf
host=127.0.0.1
port=3306
```
使用`from-file`关键字创建该目录的configmap
```sh
~ kubectl create cm cm-demo-fromdir --from-file ./testcm
configmap/cm-demo-fromdir created

~ kubectl describe cm cm-demo-fromdir
Name:         cm-demo-fromdir
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
mysql.conf:
----
my.host=127.0.0.1
my.port=3306

redis.conf:
----
red.host=127.0.0.0
red.port=6379


BinaryData
====

Events:  <none>


~ kubectl get cm cm-demo-fromdir -o yaml
apiVersion: v1
data:
  mysql.conf: |
    my.host=127.0.0.1
    my.port=3306
  redis.conf: |
    red.host=127.0.0.0
    red.port=6379
kind: ConfigMap
metadata:
  creationTimestamp: "2023-07-18T05:27:00Z"
  name: cm-demo-fromdir
  namespace: default
  resourceVersion: "455871"
  uid: cfd96392-cfbd-4ce0-851f-4ecfd2ad7e4c
```
通过文件方式创建configmap时，还可以指定key值
```sh
[root@master ~]# kubectl create cm cm-demo-fromfilekey --from-file key1=./testcm/redis.conf --from-file redis=./testcm/mysql.conf
configmap/cm-demo-fromfilekey created
[root@master ~]# kubectl describe cm cm-demo-fromfilekey
Name:         cm-demo-fromfilekey
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
key1:
----
red.host=127.0.0.0
red.port=6379

redis:
----
my.host=127.0.0.1
my.port=3306


BinaryData
====

Events:  <none>
[root@master ~]# kubectl get cm cm-demo-fromfilekey -o yaml
apiVersion: v1
data:
  key1: |
    red.host=127.0.0.0
    red.port=6379
  redis: |
    my.host=127.0.0.1
    my.port=3306
kind: ConfigMap
metadata:
  creationTimestamp: "2023-07-18T05:39:42Z"
  name: cm-demo-fromfilekey
  namespace: default
  resourceVersion: "457014"
  uid: ee5d0436-1ae3-44eb-89c9-6c1eeece067c
```
可以看到key从默认的文件名，变成了指定的key值

还可以通过直接使用字符串的方式创建，通过`--from-literal`传递配置信息
```sh
[root@master ~]# kubectl create cm cm-demo-fromstring --from-literal key1=value1 --from-literal key2=number2
configmap/cm-demo-fromstring created
[root@master ~]# kubectl describe cm cm-demo-fromstring
Name:         cm-demo-fromstring
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
key1:
----
value1
key2:
----
number2

BinaryData
====

Events:  <none>
[root@master ~]# kubectl get cm !$ -o yaml
kubectl get cm cm-demo-fromstring -o yaml
apiVersion: v1
data:
  key1: value1
  key2: number2
kind: ConfigMap
metadata:
  creationTimestamp: "2023-07-18T05:45:17Z"
  name: cm-demo-fromstring
  namespace: default
  resourceVersion: "457518"
  uid: 35de58fd-334c-444d-9927-fb3b4150fabb
```
`ConfigMap`创建成功了，但是如何在Pod中使用?`ConfigMap`配置的数据，有很多种方式在pod中使用， 主要为：
- 设置环境变量的值
- 在容器中设置命令行参数
- 在数据卷挂载配置文件

1. 环境变量
使用configMap作为环境变量

```
apiVersion: v1
kind: Pod
metadata:
  name: env-cm-demo
  labels:
    name: env-cm-demo
spec:
  containers:
  - name: env-cm-demo
    image: busybox
    command: [ "/bin/sh", "-c", "env" ]
    #command: [ "/bin/sh", "-c", "echo $(DB_HOST) $(DB_PORT)" ] 这样可以输出cm注入的环境变量
    env:
      - name: DB_HOST
        valueFrom: 
          configMapKeyRef:
            key: db.host
            name: cm-demo2
      - name: redis_config
        valueFrom:
          configMapKeyRef:
            key: redis.conf
            name: cm-demo1
    envFrom:
      - configMapRef:
          name: cm-env-demo
      - configMapRef:
          name: cm-demo
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-env-demo
  labels:
    app: myapplication
data:
  my-key: my-value
```
运行POD后，查看一下日志信息以及对应cm
```sh
[root@master configmap]# kubectl logs -f pods/env-cm-demo
data.2=world

redis_config=red.host=127.0.0.0
red.port=6379

my-key=my-value

config=property.1=value-1
property.2=value-2
property.3=value-3


DB_HOST=127.0.0.1
data.1=hello

Name:         cm-demo
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
config:
----
property.1=value-1
property.2=value-2
property.3=value-3

data.1:
----
hello
data.2:
----
world

BinaryData
====

Events:  <none>


Name:         cm-demo2
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
db.host:
----
127.0.0.1
db.port:
----
2379

BinaryData
====

Events:  <none>


Name:         cm-env-demo
Namespace:    default
Labels:       app=myapplication
Annotations:  <none>

Data
====
my-key:
----
my-value

BinaryData
====

Events:  <none>
```
2. 通过数据卷的方式，即将文件填入数据卷，键就是文件名，value就是文件内容

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: vol-cm-demo
  labels:
    name: vol-cm-demo
spec:
  volumes:
    - name:  config-cm-redis-mysql
      configMap:
        name: cm-demo1
    - name: cm-config
      configMap:
        name: cm-demo2
  containers:
  - name: vol-cm-demo
    image: busybox
    command:
     - /bin/sh
     - "-c"
     - "ls /etc/config/  /etc/data/; sleep 3600"
    volumeMounts:
      - name:  config-cm-redis-mysql
        mountPath: /etc/config
      - name : cm-config
        mountPath: /etc/data
```

```sh
[root@master configmap]# kubectl logs -f pods/vol-cm-demo
/etc/config/:
mysql.conf
redis.conf

/etc/data/:
db.host
db.port

[root@master configmap]# kubectl exec -it pods/vol-cm-demo -- bin/sh
/ #
/ # cat /etc/config/mysql.conf
my.host=127.0.0.1
/ # cat /etc/config/redis.conf
red.host=127.0.0.0
red.port=6379/ #
/ # cat /etc/config/mysql.conf
my.host=127.0.0.1
my.port=3306/ #
/ # cat /etc/config/redis.conf
red.host=127.0.0.0
red.port=6379/ #
/ # cat /etc/data/db.port
2379/ #
/ # cat /etc/data/db.host
127.0.0.1/ #
/ #
```

> 只有通过 Kubernetes API 创建的 Pod 才能使用 ConfigMap，其他方式创建的（比如静态 Pod）不能使用；ConfigMap 文件大小限制为 1MB（ETCD 的要求）。
## Secret
Secret和ConfigMap类似，但是一般Secret用来保存敏感信息，例如密码、OAuth 令牌和 ssh key 等等，将这些信息放在 Secret 中比放在 Pod 的定义中或者 Docker 镜像中要更加安全和灵活。
Secret主要使用[类型](https://kubernetes.io/zh-cn/docs/concepts/configuration/secret/#secret-types)：
- Opaque: base64编码格式的Secret，用来存储密码，密钥等，但是数据可以通过`base64 -decode`解码得到原始数据，加密性很弱
- `kubernetes.io/dockercfg`：`~/.dockercfg`文件的序列化形式
- `kubernetes.io/dockerconfigjson`：用来存储私有`docker registry`的认证信息，`~/.docker/config.json` 文件的序列化形式
- `kubernetes.io/service-account-token`：用于 S`erviceAccount`,` ServiceAccount` 创建时 Kubernetes 会默认创建一个对应的 Secret 对象，Pod 如果使用了 ServiceAccount，对应的 Secret 会自动挂载到 **Pod** 目录 `/run/secrets/kubernetes.io/serviceaccount` 中
- `kubernetes.io/ssh-auth`：用于 SSH 身份认证的凭据
- `kubernetes.io/basic-auth`：用于基本身份认证的凭据
- `bootstrap.kubernetes.io/token`：用于节点接入集群的校验的 Secret
>上面是 Secret 对象内置支持的几种类型，通过为 Secret 对象的 type 字段设置一个非空的字符串值，也可以定义并使用自己 Secret 类型。如果 type 值为空字符串，则被视为 Opaque 类型。Kubernetes 并不对类型的名称作任何限制，不过，如果要使用内置类型之一， 则你必须满足为该类型所定义的所有要求。

### Opaque Secret
Secret资源包含两个键值对：`data`和`stringData`，`data`字段用来存储 base64 编码的任意数据，提供 `stringData` 字段是为了方便，它允许 Secret 使用未编码的字符串。 `data` 和` stringData` 的键必须由字母、数字、`-`，`_`或`.` 组成。

创建Secret对象时，首先要将保存到信息做base64编码
```sh
echo -n "admin" | base64
YWRtaW4=
echo -n "admin321" | base64
YWRtaW4zMjE=
```
通过编码后的数据，就可以呃编写yaml文件了
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: YWRtaW4zMjE=

---
kubectl get secret
NAME                  TYPE                                  DATA      AGE
mysecret              Opaque                                2         40s

kubectl describe secret mysecret
Name:         mysecret
Namespace:    default
Labels:       <none>
Annotations:  <none>

Type:  Opaque

Data
====
password:  8 bytes
username:  5 bytes

---
kubectl get secret mysecret -o yaml
apiVersion: v1
data:
  password: YWRtaW4zMjE=
  username: YWRtaW4=
kind: Secret
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","data":{"password":"YWRtaW4zMjE=","username":"YWRtaW4="},"kind":"Secret","metadata":{"annotations":{},"name":"mysecret","namespace":"default"},"type":"Opaque"}
  creationTimestamp: "2023-07-18T07:27:16Z"
  name: mysecret
  namespace: default
  resourceVersion: "466809"
  uid: 7ca18e59-bab2-4dee-b626-349cd061bea8
type: Opaque
```
某些场景下，可能希望使用`stringData`字段，这个字段可以将一个非base64编码的字符串直接放入secret中，当创建或更新该Secret时，此字段将被编码
比如应用程序需要使用一下配置
```yaml
apiUrl: "https://my.api.com/api/v1"
username: "<user>"
password: "<password>"
```
对应的secret yaml 如下
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: string-data
type: Opaque
stringData: 
  config.yaml: |
    apiUrl: "https://my.api.com/api/v1"
    username: "<user>"
    password: "<password>"
```
```sh
[root@master secret]# kubectl apply -f stringData.yaml
secret/string-data created
[root@master secret]# kubectl describe secrets string-data
Name:         string-data
Namespace:    default
Labels:       <none>
Annotations:  <none>

Type:  Opaque

Data
====
config.yaml:  78 bytes
[root@master secret]# kubectl get secrets string-data -o yaml
apiVersion: v1
data:
  config.yaml: YXBpVXJsOiAiaHR0cHM6Ly9teS5hcGkuY29tL2FwaS92MSIKdXNlcm5hbWU6ICI8dXNlcj4iCnBhc3N3b3JkOiAiPHBhc3N3b3JkPiIK
kind: Secret
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Secret","metadata":{"annotations":{},"name":"string-data","namespace":"default"},"stringData":{"config.yaml":"apiUrl: \"https://my.api.com/api/v1\"\nusername: \"\u003cuser\u003e\"\npassword: \"\u003cpassword\u003e\"\n"},"type":"Opaque"}
  creationTimestamp: "2023-07-18T07:39:01Z"
  name: string-data
  namespace: default
  resourceVersion: "467882"
  uid: 61333d4d-0b6c-449a-96b9-30c395bb5714
type: Opaque
```

同样的，secret和configmap类似，也有两种方式来使用
- 环境变量
- volume挂载
### 下面以环境变量和挂载的形式同时演示
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: secret-demo
  labels:
    name: secret-demo
spec:
  containers:
  - name: secret-demo
    image: busybox
    command: ["/bin/sh", "-c", "env"]
    #command: ["/bin/sh", "-c", "sleep 3600"]
    env:
      - name: mysecrty-pass
        valueFrom:
          secretKeyRef:
            name: mysecrty
            key: password
      - name: mysecrty-user
        valueFrom:
          secretKeyRef:
            name: mysecrty
            key: username
    envFrom:
      - secretRef:
          name: secretenvfrom
    volumeMounts:
      - name:  volsecret
        mountPath:  /etc/secrets
  volumes:
    - name: volsecret
      secret:
        secretName: secret-vol
----
apiVersion: v1
kind: Secret
metadata:
  name:  secrtyen
data:
   username: YWRtaW4=
   password: YWRtaW4zMjE=
type: Opaque
---
apiVersion: v1
kind: Secret
metadata:
  name:  secretenvfrom
data:
   message: dGhpcyBpcyBzZWNyZXRlbnZmcm9t
type: Opaque

---
apiVersion: v1
kind: Secret
metadata:
  name:  secret-vol
data:
   message:  dGhpcyBpcyBzZWNyZXQgdm9sIG1vdW50
type: Opaque

---

[root@master secret]# kubectl logs -f pods/secret-demo
……
mysecrty-pass=admin321

mysecrty-user=admin

message=this is secretenvfrom
……
----
[root@master secret]# kubectl exec -it pods/secret-demo -- /bin/sh

~ # cat /etc/secrets/message
this is secret vol mount~ #
~ #

```
### kubernetes.io/dockerconfigjson
`dockerconfigjson`来创建用户`docker registry`认证的`secret`，可以用`kubectl create secret -h`来查看相关信息
```sh
[root@master secret]# kubectl create secret --help
Create a secret using specified subcommand.

Available Commands:
  docker-registry   创建一个给 Docker registry 使用的 Secret
  generic           Create a secret from a local file, directory, or literal value
  tls               创建一个 TLS secret

Usage:
  kubectl create secret [flags] [options]

Use "kubectl <command> --help" for more information about a given command.
Use "kubectl options" for a list of global command-line options (applies to all commands).
```
可以使用
`kubectl create secret docker-registry NAME --docker-username=user --docker-password=password --docker-email=email [--docker-server=string] [--from-file=[key=]source] [--dry-run=server|client|none] [options]`
也可以使用
`kubectl create secret generic myregistry --from-file=.dockerconfigjson=/root/.docker/config.json --type=kubernetes.io/dockerconfigjson`
如果我们需要拉取私有仓库中的 Docker 镜像的话就需要使用到上面的 myregistry 这个 Secret：
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: foo
spec:
  containers:
  - name: foo
    image: 192.168.1.100:5000/test:v1
  imagePullSecrets:
  - name: myregistry
```
>ImagePullSecrets 与 Secrets 不同，因为 Secrets 可以挂载到 Pod 中，但是 ImagePullSecrets 只能由 Kubelet 访问。
我们需要拉取私有仓库镜像 192.168.1.100:5000/test:v1，我们就需要针对该私有仓库来创建一个如上的 Secret，然后在 Pod 中指定 imagePullSecrets。

除了设置 Pod.spec.imagePullSecrets 这种方式来获取私有镜像之外，我们还可以通过在 ServiceAccount 中设置 imagePullSecrets，然后就会自动为使用该 SA 的 Pod 注入 imagePullSecrets 信息：
```
apiVersion: v1
kind: ServiceAccount
metadata:
  creationTimestamp: "2019-11-08T12:00:04Z"
  name: default
  namespace: default
  resourceVersion: "332"
  selfLink: /api/v1/namespaces/default/serviceaccounts/default
  uid: cc37a719-c4fe-4ebf-92da-e92c3e24d5d0
secrets:
- name: default-token-5tsh4
imagePullSecrets:
- name: myregistry
```
### ServiceAccount


## Secret和ConfigMap的区别
### 相同点
- key/value的形式
- 属于某个特定的命名空间
- 可以导出到环境变量
- 可以通过目录/文件形式挂载
- 通过 volume 挂载的配置信息均可热更新
### 不同点
- Secret 可以被 `ServerAccount` 关联
- Secret 可以存储 `docker register` 的鉴权信息，用在 `- ImagePullSecret` 参数中，用于拉取私有仓库的镜像
- `Secret `支持 Base64 加密
- `Secret` 分为 `kubernetes.io/service-account-token`、- `kubernetes.io/dockerconfigjson`、`Opaque `等多种[类型](https://kubernetes.io/zh-cn/docs/concepts/configuration/secret/#secret-types)，而 `Configmap` 不区分类型