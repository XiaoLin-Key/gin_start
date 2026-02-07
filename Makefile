.PHONY: all build run clean test help

# 变量定义
BINARY_NAME=bluebell
# 你的配置文件路径在 ./conf/config.yaml
CONF_PATH=./conf/config.yaml

all: build

build:
	@echo "正在编译 Go 二进制文件..."
	@go build -o $(BINARY_NAME) main.go

run:
	@echo "正在启动项目..."
	@go build -o $(BINARY_NAME) main.go
	@./$(BINARY_NAME) $(CONF_PATH)

clean:
	@echo "正在清理二进制文件..."
	@if [ -f $(BINARY_NAME) ] ; then rm $(BINARY_NAME) ; fi

test:
	@echo "正在执行单元测试..."
	@go test -v ./...

# 专门为 air 准备的指令
watch:
	@air -c .air.toml

help:
	@echo "使用方法:"
	@echo "  make build  - 编译项目生成二进制文件"
	@echo "  make run    - 编译并直接运行项目"
	@echo "  make watch  - 使用 Air 进行热重载开发"
	@echo "  make clean  - 移除编译生成的二进制文件"
	@echo "  make help   - 查看帮助信息"

.DEFAULT_GOAL := help