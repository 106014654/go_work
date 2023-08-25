.PHONY:docker
docker:
	@#rm user_webook || true
	@GOOS=linux GOARCH=arm go build  -o user_webook .
	@docker rmi -f chrisch/user_webook:v0.0.1
	@docker build -t chrisch/user_webook:v0.0.1 .