.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@set GOOS=linux
	@set GOARCH=arm
	@go build -o webook .
	@docker rmi -f yangyiliy/webook:v0.0.1
	@docker build -t yangyiliy/webook:v0.0.1 .
.PHONY: mock
mock:
# windows中-有特殊意义，需要加上特殊转义
	@mockgen `-source=./internal/service/code.go `-package=svcmocks `-destination=./internal/service/mocks/code.mock.go
	@mockgen `-source=./internal/service/user.go `-package=svcmocks `-destination=./internal/service/mocks/user.mock.go

	@mockgen `-source=./internal/repository/code.go `-package=repomocks `-destination=./internal/repository/mocks/code.mock.go
    @mockgen `-source=./internal/repository/user.go `-package=repomocks `-destination=./internal/repository/mocks/user.mock.go

	@mockgen `-source=./internal/repository/dao/user.go `-package=daomocks `-destination=./internal/repository/dao/mocks/user.mock.go

    @mockgen `-source=./internal/repository/cache/user.go `-package=cachemocks `-destination=./internal/repository/cache/mocks/user.mock.go
	@mockgen `-source=./internal/repository/cache/code.go `-package=cachemocks `-destination=./internal/repository/cache/mocks/code.mock.go

	@go mod tidy
