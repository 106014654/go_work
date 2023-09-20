.PHONY: mock
mock:
	@mockgen -source=user_webook/internal/service/user.go -package=svcmocks -destination=user_webook/internal/service/mocks/user.mock.go
	@mockgen -source=user_webook/internal/service/code.go -package=svcmocks -destination=user_webook/internal/service/mocks/code.mock.go
	@mockgen -source=user_webook/internal/repository/code.go -package=repomocks -destination=user_webook/internal/repository/mocks/code.mock.go
	@mockgen -source=user_webook/internal/repository/user.go -package=repomocks -destination=user_webook/internal/repository/mocks/user.mock.go
	@mockgen -source=user_webook/internal/repository/dao/user.go -package=daomocks -destination=user_webook/internal/repository/dao/mocks/user.mock.go
	@mockgen -source=user_webook/internal/repository/cache/user.go -package=cachemocks -destination=user_webook/internal/repository/cache/mocks/user.mock.go
	@mockgen -source=user_webook/internal/service/sms/type.go -package=smsmocks -destination=user_webook/internal/service/sms/mocks/sms.mock.go
	#@mockgen -source=user_webook/pkg/ratelimit/types.go -package=limitmocks -destination=user_webook/pkg/ratelimit/mocks/ratelimit.mock.go
	@mockgen -package=redismocks -destination=user_webook/internal/repository/cache/redismocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy