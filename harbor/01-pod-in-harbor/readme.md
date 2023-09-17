修改containerd中的配置文件：
```toml
    [plugins."io.containerd.grpc.v1.cri".registry]
      config_path = ""

      [plugins."io.containerd.grpc.v1.cri".registry.auths]

    # 新增
      [plugins."io.containerd.grpc.v1.cri".registry.configs]
        [plugins."io.containerd.grpc.v1.cri".registry.configs."harbor.go189.cn".tls]
                insecure_skip_verify = true # 跳过https验证
        [plugins."io.containerd.grpc.v1.cri".registry.configs."harbor.go189.cn".auth]
                username = "admin"
                password = "Harbor12345"
    # 以上为新增
      [plugins."io.containerd.grpc.v1.cri".registry.headers]

      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
          endpoint = ["https://p3ubfc8j.mirror.aliyuncs.com"]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."k8s.gcr.io"]
          endpoint = ["https://registry.cn-hangzhou.aliyuncs.com/google_containers"]

```

重启containerd

验证
```bash
# 登录 由于是自签名证书，所以需要加上 --insecure-registr参数 
nerdctl login -u admin harbor.go189.cn  --insecure-registr

#尝试上传镜像
nerdctl tag nginx:alpine harbor.go189.cn/library/nginx:alpine  
nerdctl push --insecure-registry harbor.go189.cn/library/nginx:alpine

# 拉取镜像
nerdctl push --insecure-registry harbor.go189.cn/library/nginx:alpine

```

在K8S集群中使用
需要将 Harbor 的认证信息以 Secret 的形式添加到集群中去：[参考连接](https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/pull-image-private-registry/#inspecting-the-secret-regcred)
```bash
kubectl create secret docker-registry regcred \
  --docker-server=<你的镜像仓库服务器> \
  --docker-username=<你的用户名> \
  --docker-password=<你的密码> \
  --docker-email=<你的邮箱地址>
  -n namespace
```
编写yaml文件时，需要将secret写在pod.spec.imagePullSecrets
```yaml
    spec:
      containers:
      - name: nginx
        image: harbor.go189.cn/library/nginx:latest
      imagePullSecrets:
      - name: harbor-auth
```
