# zrun - 跨平台脚本语言

zrun 是一个简单的跨平台脚本语言，可以根据不同的操作系统执行相应的命令。

## 特性

- 自动检测操作系统类型 (Windows/Linux/macOS)
- 根据操作系统执行对应的命令块
- 支持默认命令块 (`@default`)
- 支持Unix通用平台 (`@unix` 适用于Linux和macOS)

## 安装

确保你已经安装了 Go 语言环境。

```bash
git clone <repository-url>
cd zrun
go build -o zrun main.go
```

## 使用方法

创建一个 `.zr` 扩展名的脚本文件：

```zr
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

## 工作原理

zrun 解释器会：
1. 解析 `.zr` 脚本文件
2. 检测当前运行的操作系统
3. 执行第一个匹配的操作系统命令块
4. 如果没有匹配的块，则尝试执行 `@default` 块

## 示例

查看 `example.zr` 和 `advanced_example.zr` 文件获取更多使用示例。