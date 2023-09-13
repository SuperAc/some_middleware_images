# 08_RBAC

## RBAC和用户绑定
```
apiVersion: v1
kind: ServiceAccount
metadata:
  name: xaby-test
  namespace: xaby-test
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: xaby-test
  namespace: xaby-test
  labels:
    app: xaby-test
    env: xaby-test
rules:
  - apiGroups:
     - ""
    resources:
     - "*"
    verbs:
    - "*"
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  namespace: xaby-test
  name: xaby-test
  labels:
    app: xaby-test
    env: xaby-test
roleRef:
  name: xaby-test
  apiGroup: rbac.authorization.k8s.io
  kind: Role
subjects:
  - kind: User
    name: xaby
    apiGroup: ''
  - kind: ServiceAccount
    name: xaby-test
    apiGroup: ''
---
apiVersion: v1
kind: Secret
metadata:
  name: xaby
  namespace: xaby-test
  annotations:
    kubernetes.io/service-account.name: "xaby-test"
type: kubernetes.io/service-account-token
```
思路：
创建role和rolebinding（只有namespace角色的权限，如果是集群权限的话，就用clusterRole）
绑定之后，生成token，创建一个config文件
```
# config
apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://192.168.127.100:6443 #集群的serverIP地址
    certificate-authority-data: #这里可以复制已有的 LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvakNDQWVhZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJek1Ea3dOVEExTkRJd05Gb1hEVE16TURrd01qQTFOREl3TkZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTXF5Cm9sVE1DamFhMnE4eUdFRUpTdHE1K2dnakJScm9DOWNpM3k5KzB5azlJVkFQL1dsdWNwVE5uemdyaVBmY1hWcWwKUVYyM2dHTlNIcG9DYVhQaENhVEgrazlCSzhJTW82RjhYdkcvdDloL0dTU0R1TGVCQmlqNU9hRnMvU3VKeDU3Zgo4bEtaM2c0b0c4eW5SdzAraEFpcitBTkdsU01xZFdXekNIejllb1diSGRsRVcvYlR4MmJKSUMxalk4bStnOENpCm1zOEtRVlNnMWVQdzZYS2VrZ0xkcGdibXJ4VndPTEJVbldoYnUxVTdxaGpWQVR1bEhTMVBxRE5HZm15WUJTclcKRUR0TldJVytpV0piZ280T0RLbmVWeWVlK3Z0c2Q3eUJQbVErTjVFQ3lZeXdSRzNuaXBwSWgwSTJ5SitvSUhSUwptcXgyUHlJeGU3NmM4ZVZFNm1FQ0F3RUFBYU5aTUZjd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZCOEhaRGFVend0aXVIQThLeGZjQ2pvbytiNUFNQlVHQTFVZEVRUU8KTUF5Q0NtdDFZbVZ5Ym1WMFpYTXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBRklOTTVnU2lkeWFWR3pEMDRqVQpNRjdKM1pEQ0NwM2dreStxU1R6NVBlM0pRWXNQT25rNlBOSnJFc1FlMFMzZVV4MmVvU0NPWU9jTXkyOU96c1NvCkNBVlVIZW9PS3g2Vk9VYmtYS1V0b2lJd1BXV0dMaWRzTFRBRHltYnV4Yko2aXg5WlVNaXZsNGFXaktJN1p0OTAKQzZHNmJaT004NTdpV2IwVDc1RlRRVXRnRnNBeVJ6eGphVitQOWwwcjQ2Vi8vWVcyYWFlUC9sY09wb0ZzSk50cApHOUt4b053elk3MG9WRHA4d1d5amNnMi81YjBqdDJnRHZSeUM0SGVDUmE1MmRON091c2pGelJBUERYQVdPc3psCmtJUXdKbnMrUUVyenRmK0xxYTkxcVpmeEpQRVF1OFlJWXpxdStmL2d4WlcvUldIZytvZmo5WUdwRktwS0hJc3AKMkpzPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
  name: xaby-test
users:
- name: "xaby"
  user:
    token:  这个就是刚才的token "eyJhbGciOiJSUzI1NiIsImtpZCI6InpLQThoNHJCOGVEMUtzejlmMnVKaEdZX2NCbXZEMDJsZDZqc3lZcWh3Ym8ifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNjk0NDg4MTAxLCJpYXQiOjE2OTQ0ODQ1MDEsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJ4YWJ5LXRlc3QiLCJzZXJ2aWNlYWNjb3VudCI6eyJuYW1lIjoieGFieS10ZXN0IiwidWlkIjoiNDA2OGZiM2YtOWU1My00ZGU1LTlkMGYtZGM3NTM5MTVkMjJlIn19LCJuYmYiOjE2OTQ0ODQ1MDEsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDp4YWJ5LXRlc3Q6eGFieS10ZXN0In0.jbRLZ0Dc0recikGLgAYfaY90QIuvplB6ZXvu_D6iwSUkUuWfjPNLOLgROf2nlQHgToQ2GVZ4RZumIliES53ErS8QIQdmXrlfFBiMMhEp8Mwb4KAahDbLABYqVUGwZZpzTWj-378YLY-JAuQMiiWqhaBRmG_p6BL02YxR0xlSx8-qd2mdq6XN4Y0PioTSjHGjqH_IjIAxTDtgNZg_JoeWkmto1ZYhxCRRf9aqbHlkPV0Z_G2EP3RL5H57zDjgejdU0b7VcRbcDuS7mmj7qICzunFxxN36lIwNeLj_6i0by42xK7Pe_PaRzS5FLN1qHixl8HbwX53IAaGxOh4CUJox6g"
contexts:
- context:
    cluster: xaby-test
    user: "xaby"
  name: xaby
preferences: {}
current-context: xaby

```

其实最后一步的secret可以不用，直接用命令生成的token也行
```
kubectl -n xaby-test create token xaby
```
这样子创建的config是跟着sa走的，以上的yaml只是对某一个namespace，如果想看其他的，就需要clusterrole，不给高权限
```

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: xaby-test
  labels:
    app: xaby-test
    env: xaby-test
rules:
  - apiGroups:
     - ""
    resources: [pods,deployment,secrets,configmap,service]
    verbs: [get,list,watch,describe]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: xaby-test
  labels:
    app: xaby-test
    env: xaby-test
roleRef:
  name: xaby-test
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
subjects:
  - kind: ServiceAccount
    namespace: xaby-test
    name: xaby-test
    apiGroup: ''
---

```
在测试中需要权限了再更新就可以
