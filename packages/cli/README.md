# Cli

- [cobra](https://cobra.dev/)
- [viper](https://github.com/spf13/viper)

## 运行

```sh
go run ./main.go
go run ./main.go ui
go build -o yiu-ops
```

## vscode

```sh
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/fatih/gomodifytags@latest
go install github.com/cweill/gotests/gotests@latest
go install github.com/josharian/impl@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
```

## 发布

```sh
# 添加修改日志 只选择cli
pnpm run pkg:add
# 将有修改记录的包提升版本 只选择cli
pnpm run pkg:version
# push 修改记录
git add .
pnpm run vcs:commit
# yiu cli 添加修改日志
git push
# 生成对应的git tag
pnpm run pkg:tag
# 推送tag
git push origin --tags

# 下载二进制文件之后
chmod +x ./yiu-ops-darwin-arm64
# macOS Gatekeeper 临时本地运行可以去掉隔离属性
xattr -d com.apple.quarantine ./yiu-ops-darwin-arm64
./yiu-ops-darwin-arm64
```

> 删除tag的命令 `git push origin --delete cli@1.0.0-beta.1`
