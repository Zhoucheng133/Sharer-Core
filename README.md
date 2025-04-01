# Sharer (Core)

## 简介

![License](https://img.shields.io/badge/License-MIT-dark_green)

这是一个用于在PC/Mac上的文件分享工具，可以将PC/Mac作为文件共享服务器


> [!IMPORTANT]
> 这是一个二进制文件的仓库，如果你要查找有GUI界面的App，请前往[Sharer-App](https://github.com/Zhoucheng133/Sharer-App)，如果你要查找网页的源代码请前往[Sharer-Web](https://github.com/Zhoucheng133/Sharer-Web)

> [!NOTE]
> 你可以单独[使用](#使用)它，但是更建议使用[Sharer-App](https://github.com/Zhoucheng133/Sharer-App)

✅ 打包目录下载  
✅ 多文件下载  
✅ 上传本地目录  
✅ 添加访问用户名和密码  
✅ 在线预览各种常见格式  

## 目录
- [简介](#简介)
- [目录](#目录)
- [截图](#截图)
- [使用](#使用)
- [构建](#构建)

## 截图

![截图0](demo/demo0.png)

![截图1](demo/demo1.png)

## 使用

1. 你需要前往Release页找到合适自己设备的可执行程序
2. 如果你使用的是macOS系统，在下载到本地之后需要使用这个命令允许访问这个可执行文件:
   ```bash
   chmod 777 <可执行文件的位置>
   ```
3. 之后需要通过这样的命令来运行这个程序
   ```bash
   <可执行文件位置>
   -port <端口号>
   -d <需要分享的目录位置>
   -u <用户名>
   -p <密码>
   ```
4. 执行上述命令之后会有这样的提示
   ```
   服务运行在:
   ➜ http://192.168.124.22:8081
   ➜ http://127.0.0.1:808
   ```

### 需要用户名和密码登录

例如这样:
```bash
/Users/zhoucheng/Downloads/macOS -port 8081 -d /Users/zhoucheng/Downloads -u admin -p 123456
```
端口号为`8081`  
分享位置为`/Users/zhoucheng/Downloads`  
用户名为`admin`  
密码为`123456`

### 不需要用户名和密码登录

忽略用户名和密码，这种情况下所有局域网的用户都可以直接访问而不需要用户名和密码登录，例如这样：
```bash
/Users/zhoucheng/Downloads/macOS -port 8081 -d /Users/zhoucheng/Downloads
```

端口号为`8081`  
分享位置为`/Users/zhoucheng/Downloads`  

注意只输入用户名或者只输入密码不会运行服务

## 构建

### 准备

你需要在你的设备上安装配置好这些东西：
- bun
- go

### 生成二进制文件

1. 你需要先克隆或者下载本仓库
2. 在仓库中执行此命令下载子模块:
   ```bash
   git submodule update --init --recursive
   ```
3. 生成Web页面
   ```bash
   cd Sharer-Web
   bun run build
   ```
4. 生成二进制文件
   ```bash
   cd .. #回到仓库根目录
   go run . #运行程序
   go build #打包
   ```

### 生成动态库供[Sharer-App](https://github.com/Zhoucheng133/Sharer-App)使用或二次开发

1. 你需要先克隆或者下载本仓库
2. 在仓库中执行此命令下载子模块:
   ```bash
   git submodule update --init --recursive
   ```
3. 生成Web页面
   ```bash
   cd Sharer-Web
   bun run build
   ```
4. 生成动态库
   ```bash
   cd .. #回到仓库根目录
   go build -buildmode=c-shared -o libserver.dll .    # 生成macOS动态库
   go build -buildmode=c-shared -o libserver.dylib .  # 生成macOS动态库
   ```
