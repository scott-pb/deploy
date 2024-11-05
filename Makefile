# 管理后台构建工具

# 自动生成接口
.PHONY: init-api
init-api:
	go run main.go --scene init-api --ignore "ignore"

# 生成接口文档（自动引入依赖项）
.PHONY: swag
swag:
	swag init --parseDependency --parseInternal


# 生成接口&接口文档
.PHONY: doc
doc:
	make swag
	make init-api