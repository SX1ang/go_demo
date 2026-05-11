# go_demo
go web demo project

## wire依赖注入
```shell
go install github.com/google/wire/cmd/wire@latestc
```
进入cmd目录，执行：
```wire wire.go```

## GORM Gen
```bash
go get -u gorm.io/gen
```
提前在数据库中创建数据表。
```bash
cd cmd/gen
go run gen.go
```

## 生成Swagger文档
swagger注解：https://swaggo.github.io/swaggo.io/declarative_comments_format/
```
go install github.com/swaggo/swag/cmd/swag@latest
```
```bash
swag init -d cmd,internal/api --parseInternal --parseDependency
// --parseInternal 允许 swag 解析 internal/ 目录下的包（包括其中的 struct、注解、类型定义），可以在注解里引用 internal/api/dto.SignUpReq 等类型
// --parseDependency 如果你的注解里引用了其他包（依赖包）里的类型，swag 会继续去解析这些依赖包里的 struct 定义，把它们展开到 swagger 的 definitions/schemas 
```

## Run
```bash
go mod tidy
go run ./cmd  
```
or
```bash
go build -o bin/app ./cmd
```

# 配置
app dev模式下，无需认证即可访问接口，方便测试。