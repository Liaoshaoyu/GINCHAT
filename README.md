修改为国内可用的代理地址
在命令提示符输入：go env -w GOPROXY=https://goproxy.cn

### 笔记
> 1. 修改为国内可用的代理地址，在命令提示符输入：go env -w GOPROXY=https://goproxy.cn
> 2. 安装包：go get <pkg name>; 删除没有用到的包：go mod tidy
> 