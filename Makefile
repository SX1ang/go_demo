.PHONY: all build run tidy wire gen swagger clean help

# 默认目标
all: tidy wire gen swagger build

# ================================
# 依赖管理
# ================================

## 整理依赖
tidy:
	go mod tidy

## 安装所有工具依赖
install-tools:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go get -u gorm.io/gen

# ================================
# 代码生成
# ================================

## Wire 依赖注入代码生成
wire:
	cd cmd && wire wire.go

## GORM Gen 生成数据库模型
gen:
	cd cmd/gen && go run gen.go

## 生成 Swagger 文档
swagger:
	swag init -d cmd,internal/api --parseInternal --parseDependency

# ================================
# 构建 & 运行
# ================================

## 编译项目
build:
	go build -o bin/app ./cmd

## 直接运行（dev 模式）
run:
	go run ./cmd

## dev 模式运行（无需认证）
run-dev:
	APP_ENV=dev go run ./cmd

# ================================
# 清理
# ================================

## 清理编译产物
clean:
	rm -rf bin/

# ================================
# 帮助
# ================================
help:
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^##' Makefile | sed 's/## /  /'
	@echo ""