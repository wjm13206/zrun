# zrun - 跨平台脚本语言

zrun 是一个简单的跨平台脚本语言，可以根据不同的操作系统执行相应的命令。

## 特性

- 自动检测操作系统类型 (Windows/Linux/macOS)
- 根据操作系统执行对应的命令块
- 支持默认命令块 (`@default`)
- 支持Unix通用平台 (`@unix` 适用于Linux和macOS)
- 支持命令回显控制 (`@echo on`/`@echo off`)

## 安装

确保你已经安装了 Go 语言环境。

```bash
git clone <repository-url>
cd zrun
go build -o zrun main.go
```

或者

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
    echo "Hello Unix-like systems!"
    uname -a
}

@default {
    echo "Hello Unknown OS!"
}
```

执行脚本：

```bash
./zrun script.zr
```

## 语法规则

1. 使用 `@platform { }` 来定义平台特定的命令块
2. 支持的平台标识符：
   - `@windows` - Windows系统
   - `@linux` - Linux系统
   - `@macos` - macOS系统
   - `@unix` - Unix类系统 (包括Linux和macOS)
   - `@default` - 默认块，当没有其他块匹配时执行
3. 在大括号内编写需要执行的系统命令，每行一个命令
4. 使用 `@echo off` 禁止显示将要执行的命令（默认是on，即显示命令）
5. 使用 `@echo on` 可以重新开启命令回显
6. `@echo` 指令在整个脚本中全局生效，会影响其后所有命令的回显状态

## 项目结构

```
zrun/
├── main.go                 # 程序入口
├── go.mod                  # Go模块文件
├── *.zr                    # 脚本示例文件
└── internal/               # 内部模块
    ├── types/              # 类型定义
    ├── parser/             # 脚本解析器
    ├── executor/           # 命令执行器
    └── utils/              # 工具函数
```

## 代码质量改进

根据代码分析报告，我们对项目进行了以下改进：

1. **增强代码注释** - 为所有模块添加了详细的中文注释
2. **简化复杂函数** - 重构了parser模块，将复杂函数拆分为更小的函数
3. **改善代码结构** - 保持了清晰的模块化架构

## 工作原理

zrun 解释器会：
1. 解析 `.zr` 脚本文件
2. 检测当前运行的操作系统
3. 按顺序执行所有匹配当前操作系统的命令块
4. 如果没有匹配的块，则尝试执行 `@default` 块

## 示例

查看 `example.zr`、`advanced_example.zr`、`echo_example.zr` 和 `full_demo.zr` 文件获取更多使用示例。