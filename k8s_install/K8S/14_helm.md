# 14_helm
[安装](https://github.com/helm/helm/releases)

Helm 客户端准备成功后，我们就可以添加一个 chart 仓库，当然最常用的就是官方的 Helm stable charts 仓库，但是由于官方的 charts 仓库地址需要科学上网，我们可以使用微软的 charts 仓库代替：

```bash
helm repo add stable http://mirror.azure.cn/kubernetes/charts/
helm repo list
NAME            URL
stable          http://mirror.azure.cn/kubernetes/charts/
```
安装完成后可以用 search 命令来搜索可以安装的 chart 包：
```bash
helm search repo stable
NAME                                    CHART VERSION   APP VERSION                     DESCRIPTION
stable/acs-engine-autoscaler            2.2.2           2.1.1                           DEPRECATED Scales worker nodes within agent pools
stable/aerospike                        0.3.1           v4.5.0.5                        A Helm chart for Aerospike in Kubernetes
stable/airflow                          5.2.1           1.10.4                          Airflow is a platform to programmatically autho...
stable/ambassador                       5.1.0           0.85.0                          A Helm chart for Datawire Ambassador
stable/anchore-engine                   1.3.7           0.5.2                           Anchore container analysis and policy evaluatio...
stable/apm-server                       2.1.5           7.0.0                           The server receives data from the Elastic APM a...
......
```
## ex
首先从仓库中将可用的 charts 信息同步到本地，可以确保我们获取到最新的 charts 列表：
```bash
 helm repo update
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "stable" chart repository
Update Complete. ⎈ Happy Helming!⎈
```

我们现在安装一个 mysql 应用：
```bash
helm install stable/mysql --generate-name
WARNING: This chart is deprecated
NAME: mysql-1693711024
LAST DEPLOYED: Sun Sep  3 11:17:08 2023
NAMESPACE: default
STATUS: deployed
REVISION: 1
NOTES:
MySQL can be accessed via port 3306 on the following DNS name from within your cluster:
mysql-1693711024.default.svc.cluster.local
……
```
可以看到 stable/mysql 这个 chart 已经安装成功了，我们将安装成功的这个应用叫做一个 release，由于我们在安装的时候指定了--generate-name 参数，所以生成的 release 名称是随机生成的，名为mysql-1693711024。我们可以用下面的命令来查看 release 安装以后对应的 Kubernetes 资源的状态：
```bash
 kubectl get all -l release=mysql-1693711024
NAME                                   READY   STATUS    RESTARTS   AGE
pod/mysql-1693711024-5d8ddcfb5-fl9fn   0/1     Pending   0          102s

NAME                       TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/mysql-1693711024   ClusterIP   10.111.90.151   <none>        3306/TCP   102s

NAME                               READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/mysql-1693711024   0/1     1            0           102s

NAME                                         DESIRED   CURRENT   READY   AGE
replicaset.apps/mysql-1693711024-5d8ddcfb5   1         1         0       102s

```
可以通过 helm show chart 命令来了解 MySQL 这个 chart 包的一些特性：
helm show chart stable/mysql

如果想要了解更多信息，可以用 helm show all 命令：
 helm show all stable/mysql

需要注意的是无论什么时候安装 chart，都会创建一个新的 release，所以一个 chart 包是可以多次安装到同一个集群中的，每个都可以独立管理和升级。

同样我们也可以用 Helm 很容易查看到已经安装的 release：
```bash
helm ls
NAME                    NAMESPACE       REVISION        UPDATED                                 STATUS          CHART           APP VERSION
mysql-1693711024        default         1               2023-09-03 11:17:08.201915846 +0800 CST deployed        mysql-1.6.9     5.7.30
testsql                 default         2               2023-09-03 11:11:06.354158218 +0800 CST deployed        mysql-1.6.9     5.7.30

```

如果需要删除这个 release，也很简单，只需要使用 helm uninstall 命令即可

uninstall 命令会从 Kubernetes 中删除 release，也会删除与 release 相关的所有 Kubernetes 资源以及 release 历史记录。也可以在删除的时候使用 --keep-history 参数，则会保留 release 的历史记录，可以获取该 release 的状态就是 UNINSTALLED，而不是找不到 release了

```bash
[root@master 01mysql]# helm uninstall mysql-1693711320 --keep-history
release "mysql-1693711320" uninstalled
[root@master 01mysql]# helm ls
NAME    NAMESPACE       REVISION        UPDATED                                 STATUS          CHART           APP VERSION
testsql default         2               2023-09-03 11:11:06.354158218 +0800 CST deployed        mysql-1.6.9     5.7.30
[root@master 01mysql]# helm ls -a
NAME                    NAMESPACE       REVISION        UPDATED                                 STATUS          CHART           APP VERSION
mysql-1693711320        default         1               2023-09-03 11:22:02.896169112 +0800 CST uninstalled     mysql-1.6.9     5.7.30
testsql                 default         2               2023-09-03 11:11:06.354158218 +0800 CST deployed        mysql-1.6.9     5.7.30

```
## 定制
上面都是直接使用的 helm install 命令安装的 chart 包，这种情况下只会使用 chart 的默认配置选项，但是更多的时候，是各种各样的需求，索引我们希望根据自己的需求来定制 chart 包的配置参数。

我们可以使用 helm show values 命令来查看一个 chart 包的所有可配置的参数选项：
```bash
[root@master 01mysql]# helm show values stable/mysql
## mysql image version
## ref: https://hub.docker.com/r/library/mysql/tags/
##
image: "mysql"
imageTag: "5.7.30"

strategy:
  type: Recreate

busybox:
  image: "busybox"
  tag: "1.32"

testFramework:
  enabled: true
  image: "bats/bats"
  tag: "1.2.1"
  imagePullPolicy: IfNotPresent
  securityContext: {}

## Specify password for root user
##
## Default: random 10 character string
# mysqlRootPassword: testing

## Create a database user
##
# mysqlUser:
## Default: random 10 character string
# mysqlPassword:

## Allow unauthenticated access, uncomment to enable
##
# mysqlAllowEmptyPassword: true

## Create a database
##
# mysqlDatabase:

## Specify an imagePullPolicy (Required)
## It's recommended to change this to 'Always' if the image tag is 'latest'
## ref: http://kubernetes.io/docs/user-guide/images/#updating-images
##
imagePullPolicy: IfNotPresent

……
```


所有参数都是可以用自己的数据来覆盖的，可以在安装的时候通过 YAML 格式的文件来传递这些参数：
```bash

[root@master 01mysql]# cat MySQL.yaml
mysqlRootPassword: root
mysqlDatabase: helmtest
persistence:
  enabled: false
  
[root@master 01mysql]# helm install -f MySQL.yaml testsql stable/mysql
WARNING: This chart is deprecated
NAME: testsql
LAST DEPLOYED: Sun Sep  3 11:33:30 2023
NAMESPACE: default
STATUS: deployed
REVISION: 1
NOTES:
MySQL can be accessed via port 3306 on the following DNS name from within your cluster:
testsql-mysql.default.svc.cluster.local

```
release 安装成功后，可以查看对应的 Pod 信息：
```bash
[root@master 01mysql]# kubectl get pod -l release=testsql
NAME                            READY   STATUS    RESTARTS   AGE
testsql-mysql-c498cf466-pzrgd   1/1     Running   0          111s
[root@master 01mysql]# kubectl describe pod testsql-mysql-c498cf466-pzrgd  
……
    Environment:
      MYSQL_ROOT_PASSWORD:  <set to the key 'mysql-root-password' in secret 'testsql-mysql'>  Optional: false
      MYSQL_PASSWORD:       <set to the key 'mysql-password' in secret 'testsql-mysql'>       Optional: true
      MYSQL_USER:
      MYSQL_DATABASE:       helmtest
    Mounts:
……
```
可以看到环境变量 MYSQL_DATABASE=helmtest的值和我们上面配置的值是一致的。在安装过程中，有两种方法可以传递配置数据：
- `--values（或者 -f）`：指定一个 YAML 文件来覆盖 values 值，可以指定多个值，最后边的文件优先
- `--set`：在命令行上指定覆盖的配置
如果同时使用两种方式，则 --set 中的值会被合并到 --values 中，但是 --set 中的值优先级更高。使用 --set 指定的值将持久化在 ConfigMap 中
如：
```bash
helm upgrade  testsql -f MySQL.yaml --set persistence.enabled=true stable/mysql
[root@master 01mysql]# kubectl get pod -l release=testsql
NAME                             READY   STATUS    RESTARTS   AGE
testsql-mysql-645797fc4d-h4z4v   0/1     Pending   0          1s
```
因为`--set`优先级更高，所以此时由于没有持久卷，pod状态变为pending

对于给定的 release，可以使用 helm get values <release-name> 来查看已经设置的值，已设置的值也通过允许 helm upgrade 并指定 --reset 值来清除。
```bash
[root@master 01mysql]# helm get values testsql
USER-SUPPLIED VALUES:
mysqlDatabase: helmtest
mysqlRootPassword: root
persistence:
  enabled: true
```
## 更多安装方式
- chart 仓库（类似于上面我们提到的）
- 本地 chart 压缩包（helm install foo-0.1.1.tgz）
通过`helm fetch stable/mysql`可以将stable中的某个chart压缩
- 本地解压缩的 chart 目录（helm install foo path/to/foo）
- 在线的 URL（helm install fool https://example.com/charts/foo-1.2.3.tgz）

## 升级和回滚
当新版本的 chart 包发布的时候，或者要更改 release 的配置的时候，可以使用 helm upgrade 命令来操作
```bash
[root@master 01mysql]# helm upgrade  testsql -f MySQL.yaml --set persistence.enabled=true stable/mysql
WARNING: This chart is deprecated
^[[A^[[ARelease "testsql" has been upgraded. Happy Helming!
NAME: testsql
LAST DEPLOYED: Sun Sep  3 11:39:17 2023
NAMESPACE: default
STATUS: deployed
REVISION: 2
……
```
这里就是通过upgrade更新了persistence.enabled的参数
```bash
[root@master 01mysql]# helm list
NAME    NAMESPACE       REVISION        UPDATED                                 STATUS          CHART           APP VERSION
testsql default         2               2023-09-03 11:39:17.362568675 +0800 CST deployed        mysql-1.6.9     5.7.30
[root@master 01mysql]# helm history testsql
REVISION        UPDATED                         STATUS          CHART           APP VERSION     DESCRIPTION
1               Sun Sep  3 11:33:30 2023        superseded      mysql-1.6.9     5.7.30          Install complete
2               Sun Sep  3 11:39:17 2023        deployed        mysql-1.6.9     5.7.30          Upgrade complete
[root@master 01mysql]# helm rollback testsql 1
Rollback was a success! Happy Helming!
[root@master 01mysql]# kubectl get pod -l release=testsql
NAME                            READY   STATUS    RESTARTS   AGE
testsql-mysql-c498cf466-6mjpb   1/1     Running   0          14s
[root@master 01mysql]# helm list
NAME    NAMESPACE       REVISION        UPDATED                                 STATUS          CHART           APP VERSION
testsql default         3               2023-09-03 12:11:04.470435533 +0800 CST deployed        mysql-1.6.9     5.7.3
[root@master 01mysql]# helm history  testsql
REVISION        UPDATED                         STATUS          CHART           APP VERSION     DESCRIPTION
1               Sun Sep  3 11:33:30 2023        superseded      mysql-1.6.9     5.7.30          Install complete
2               Sun Sep  3 11:39:17 2023        superseded      mysql-1.6.9     5.7.30          Upgrade complete
3               Sun Sep  3 12:11:04 2023        deployed        mysql-1.6.9     5.7.30          Rollback to 1
   
```

可以看到 values 配置已经回滚到之前的版本，pod也运行了。上面的命令回滚到了 release 的第一个版本，每次进行安装、升级或回滚时，修订号都会加 1，第一个修订号始终为1，我们可以使用 helm history [RELEASE] 来查看某个版本的修订号


除此之外我们还可以指定一些有用的选项来定制 install/upgrade/rollback 的一些行为，要查看完整的参数标志，我们可以运行 helm <command> --help 来查看，这里我们介绍几个有用的参数：

- --timeout: 等待 Kubernetes 命令完成的时间，默认是 300（5分钟）
- --wait: 等待直到所有 Pods 都处于就绪状态、PVCs 已经绑定、Deployments 具有处于就绪状态的最小 Pods 数量（期望值减去 maxUnavailable）以及 Service 有一个 IP 地址，然后才标记 release 为成功状态。它将等待与 --timeout 值一样长的时间，如果达到超时，则 release 将标记为失败。注意：在 Deployment 将副本设置为 1 并且作为滚动更新策略的一部分，maxUnavailable 未设置为0的情况下，--wait 将返回就绪状态，因为它已满足就绪状态下的最小 Pod 数量
- --no-hooks: 将会跳过命令的运行 hooks
- --recreate-pods: 仅适用于 upgrade 和 rollback，这个标志将导致重新创建所有的 Pods。（Helm3 中启用了）

## charts模板
文件结构
```base
wordpress/
  Chart.yaml          # 包含当前 chart 信息的 YAML 文件
  LICENSE             # 可选：包含 chart 的 license 的文本文件
  README.md           # 可选：一个可读性高的 README 文件
  values.yaml         # 当前 chart 的默认配置 values
  values.schema.json  # 可选: 一个作用在 values.yaml 文件上的 JSON 模式
  charts/             # 包含该 chart 依赖的所有 chart 的目录
  crds/               # Custom Resource Definitions
  templates/          # 模板目录，与 values 结合使用时，将渲染生成 Kubernetes 资源清单文件
  templates/NOTES.txt # 可选: 包含简短使用使用的文本文件
```

### Chart.yaml
对于一个 chart 包来说 Chart.yaml 文件是必须的，它包含下面的这些字段：
```bash
apiVersion: chart API 版本 (必须)
name: chart 名 (必须)
version: SemVer 2版本 (必须)
kubeVersion: 兼容的 Kubernetes 版本 (可选)
description: 一句话描述 (可选)
type: chart 类型 (可选)
keywords:
  - 当前项目关键字集合 (可选)
home: 当前项目的 URL (可选)
sources:
  - 当前项目源码 URL (可选)
dependencies: # chart 依赖列表 (可选)
  - name: chart 名称 (nginx)
    version: chart 版本 ("1.2.3")
    repository: 仓库地址 ("https://example.com/charts")
maintainers: # (可选)
  - name: 维护者名字 (对每个 maintainer 是必须的)
    email: 维护者的 email (可选)
    url: 维护者 URL (可选)
icon: chart 的 SVG 或者 PNG 图标 URL (可选).
appVersion: 包含的应用程序版本 (可选). 不需要 SemVer 版本
deprecated: chart 是否已被弃用 (可选, boolean)
```
其他字段默认会被忽略。

### 使用 Helm 管理 Charts
#### 创建chart包
```bash
helm create mychart
```
#### 将chart打包到一个独立文件中
```bash
helm package mychart
```
#### 使用 helm 帮助查找 chart 包的格式要求方面或其他问题：
```bash
helm lint mychart
```
### Chart 仓库
chart 仓库实际上就是一个 HTTP 服务器，其中包含一个或多个打包的 chart 包，虽然可以使用 helm 来管理本地 chart 目录，但是在共享 charts 的时候，最好的还是使用 chart 仓库。

可以提供 YAML 文件和 tar文件并可以相应 GET 请求的任何 HTTP 服务器都可以作为 chart 仓库服务器。仓库的主要特征是存在一个名为 index.yaml 的特殊文件，该文件具有仓库中提供的所有软件包的列表以及允许检索和验证这些软件包的元数据。

在客户端，可以使用 helm repo 命令来管理仓库，但是 Helm 不提供用于将 chart 上传到远程 chart 仓库的工具。

### template
可以通过`helm template`获取渲染后的temlpate
```
helm template traefik ./traefik > traefik.yaml
```

