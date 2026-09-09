# zrun - 跨平台任务脚本


[![go](https://img.shields.io/badge/Go-1.21+-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

## 特性

- 以 `task` 为执行单元，默认顺序执行、遇错即停
- `--on` 按平台过滤（`windows/linux/macos/unix` + `amd64/386/arm/arm64`，`64` 是 `amd64` 别名）
- `--deps` 声明依赖，自动拓扑排序并检出成环
- `var/env` + `$VAR/${VAR}` 插值，内置 `$OS/$ARCH/$ARGS`
- 回显按任务作用域控制（`--echo on/off`，默认继承全局 `@echo`）
- `env/workdir/shell` 按任务配置，`--dry-run/--list` 便于预览

## 安装

要求 Go 1.21+。

```bash
git clone https://github.com/wjm13206/zrun.git
cd zrun
go build -o zrun .
```

## 使用方法

```bash
zrun script.zr                # 按文件顺序运行全部可执行任务
zrun script.zr build          # 只运行 build 及其依赖
zrun --list script.zr         # 列出任务
zrun --dry-run script.zr build
zrun script.zr run -- --port 8080   # -- 后面的内容进 $ARGS
```

## 语法规则（v2，不兼容 v1）

```zr
@syntax 2.0
@echo on

var APP = "demo"
env GREETING = "hello"

task hello --desc "打招呼" --on [windows, linux, macos] {
  echo $GREETING from $OS/$ARCH app=$APP args=$ARGS
}

task build --desc "构建" --deps [hello] --on [linux, windows/amd64] --workdir ./ --echo on {
  go build -o dist/$APP .
}
```

1. 文件必须以 `@syntax 2.x` 开头；`v1` 的 `@windows { }` 等写法已废弃，会直接报错。
2. 顶层只允许 `@syntax/@echo/var/env/task`；任务外出现命令直接报错并带行号。
3. `task 名称 --选项 { ... }`，选项仅允许 `--desc/--on/--deps/--shell/--workdir/--echo`。
4. `--on` 为空表示全平台；`unix` 指 `linux/macos`；`default/any/*` 恒真；大小写不敏感。
5. 注释只认整行 `#` 或 `//`，避免误伤命令里的 `#`。
6. 变量用 `$NAME/${NAME}`，`$$` 转义为 `$`；查找顺序 `var > env > $OS/$ARCH/$ARGS > 系统环境变量`。

## 项目结构

```
zrun/
├── main.go                 # 程序入口（CLI）
├── go.mod                  # Go模块文件（go 1.21+）
└── src/
    ├── types/              # v2 类型定义
    ├── parser/             # v2 解析器
    ├── executor/           # 顺序执行器
    └── utils/              # 平台归一化/插值/更新检查
```

## 工作原理

1. 解析 `.zr` 脚本（带行号报错）
2. 按请求任务展开依赖并拓扑排序，检出成环
3. 跳过 `--on` 不匹配的任务
4. 顺序执行，命令先插值再交给 `cmd /C` 或 `sh -c`
