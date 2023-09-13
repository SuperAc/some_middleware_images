## 创建分支
`git checkout -b newBranch`
……
修改
……
## 修改完成后
```bash
# 查看状态
git status

# 增加修改内容
git add ***

# 提交
# -s Signed-off-by  
git commit git -s -m "修改内容"

# push
git push --set-upstream orgin newBranch

# 进入github 确认，没问题就可以提交 并申请合并
# 合并后同步
git checkout main

# 加一次就行
git remote add upstream https://github.com/****（仓库地址） 
git remote set-url --push upstream no-pushing


# 拉取主仓库的main分支
git pull upstream main


# 本地和远端的同步
git push

# 如果分支合并了，后续需要修改，切换分支需要rebase, 不然两次修改会重复
git checkout newBranch

git rebase main
```