使用helm安装harbor,具体文件参考harbor-ingree.yaml
常用命令
```bash
helm upgrade --install name dir -f filename -n namespace


helm list -n namespace


helm history name -n namespace
helm rollback name -n namespace version
```