.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@set GOOS=linux
	@set GOARCH=arm
	@go build -o webook .
	@docker rmi -f yangyiliy/webook:v0.0.1
	@docker build -t yangyiliy/webook:v0.0.1 .


