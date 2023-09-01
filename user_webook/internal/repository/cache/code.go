package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("验证码错误")
	ErrCodeTimeOut            = errors.New("验证码超时")
)

type CodeCacheInter interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type localCodeCache struct {
	lclcache sync.Map
}

type CodeCache struct {
	phone string
	code  string
	cnt   int
	ttl   int64
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewLocalCodeCache() CodeCacheInter {
	var m sync.Map
	return &localCodeCache{
		lclcache: m,
	}
}

func (lc *localCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	key := lc.key(biz, phone)

	value, ok := lc.lclcache.Load(key)

	if ok {
		val, _ := value.(CodeCache)

		fmt.Println("ok", val.phone, val.code, val.ttl)
		if time.Now().Unix()-val.ttl < 10 || val.cnt == 3 {
			return ErrCodeSendTooMany
		}

	} else {
		data := CodeCache{
			phone: phone,
			code:  code,
			cnt:   3,
			ttl:   time.Now().Unix(),
		}

		fmt.Println("data", data.phone, data.code, data.ttl, key)
		lc.lclcache.Store(key, data)
	}

	return nil
}

func (lc *localCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	fmt.Printf("params: %s,%s,%s", biz, phone, inputCode)
	key := lc.key(biz, phone)
	value, ok := lc.lclcache.Load(key)

	if !ok {
		return false, ErrUnknownForCode
	}

	val, ok := value.(CodeCache)

	if !ok {
		return false, errors.New("系统错误2")
	}
	cnt := val.cnt - 1

	fmt.Println(val)

	if time.Now().Unix()-val.ttl > 60 {
		return false, ErrCodeTimeOut
	}

	if cnt < 0 {
		lc.lclcache.Delete(key)
		return false, ErrCodeVerifyTooManyTimes
	}

	if val.code != inputCode {
		data := CodeCache{
			phone: phone,
			code:  inputCode,
			cnt:   cnt,
			ttl:   val.ttl,
		}

		lc.lclcache.Store(key, data)
		return false, nil
	}

	return true, nil
}

func NewCodeCache(client redis.Cmdable) CodeCacheInter {
	return &RedisCodeCache{
		client: client,
	}
}

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		// 毫无问题
		return nil
	case -1:
		// 发送太频繁
		return ErrCodeSendTooMany
	//case -2:
	//	return
	default:
		// 系统错误
		return errors.New("系统错误")
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	fmt.Println(res, err)
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		// 正常来说，如果频繁出现这个错误，你就要告警，因为有人搞你
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
		//default:
		//	return false, ErrUnknownForCode
	}
	return false, ErrUnknownForCode
}

func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (lc *localCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
