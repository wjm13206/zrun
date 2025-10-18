# zrun - 跨平台脚本语言

zrun 是一个简单的跨平台脚本语言。

[![go](https://img.shields.io/badge/Go-1.24.5+-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)


## 特性

- 自动检测操作系统类型 (Windows/Linux/macOS)
- 根据操作系统执行对应的命令块
- 支持默认命令块 (`@default`)
- 支持Unix通用平台 (`@unix` 适用于Linux和macOS)
- 支持命令回显控制 (`@echo on`/`@echo off`)

## 为什么会制作项目？

1.闲

2.制作懒人包

## 安装

确保你已经安装了 Go 环境。

```bash
# 克隆仓库
git clone https://github.com/wjm13206/zrun.git
cd zrun
# 构建
go build -o zrun main.go
```

或

```bash
go build -o zrun .
```

## 使用方法

创建一个 `.zr` 扩展名的脚本文件：

```zr
@echo off

@windows {
    echo "Hello Windows!"
    dir
}

@linux {
    echo "Hello Linux!"
    ls -la
}

@macos {
    echo "Hello macOS!"
    ls -lG
}

@unix {
    echo "Hello Unix"
    uname -a
}

@default {
    echo "Hello!"
}
```

执行脚本：

```bash
./zrun script.zr
```

## 语法规则

1. 使用 `@平台 { }` 来定义平台特定的命令块
2. 支持的平台标识符：
   - `@windows` - Windows系统
   - `@linux` - Linux系统
   - `@macos` - macOS系统
   - `@unix` - Unix类系统 (包括Linux和macOS)
   - `@default` - 默认块
3. 在大括号内编写需要执行的系统命令，每行一个命令
4. 使用 `@echo off` 禁止回显（默认是on，即显示命令）
5. 使用 `@echo on` 重新开启命令回显

## 项目结构

```
zrun/
├── main.go                 # 程序入口
├── go.mod                  # Go模块文件
└── src/                    # 内部
    ├── types/              # 类型定义
    ├── parser/             # 脚本解析器
    ├── executor/           # 命令执行器
    └── utils/              # 工具函数
```


## 工作原理

zrun 会：
1. 解析 `.zr` 脚本文件
2. 检测当前运行的操作系统
3. 按顺序执行所有匹配当前操作系统的命令块，和`@default` 块

