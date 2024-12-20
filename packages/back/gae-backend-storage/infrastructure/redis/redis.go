package redis

import (
	"context"
	"encoding/json"
	"errors"
	"gae-backend-storage/constant/common"
	"github.com/redis/go-redis/v9"
	"time"
)

type Client interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, k string, v any) error
	SetExpire(ctx context.Context, k string, v any, ddl time.Duration) error
	Get(ctx context.Context, k string) (string, error)

	LRange(ctx context.Context, k string, start int, end int) ([]string, error)
	LRangeAll(ctx context.Context, k string) ([]string, error)
	LRem(ctx context.Context, k string, count int, v any) (int64, error)
	LPush(ctx context.Context, k string, v any) error
	RPop(ctx context.Context, k string) (string, error)

	ZRem(ctx context.Context, k string, vs ...any) (int64, error)
	// ZScore 获取指定元素的分数 Zset
	ZScore(ctx context.Context, k string, member string) (isExist bool, score int)
	ZRangeAll(ctx context.Context, k string) ([]string, error)
	ZRange(ctx context.Context, k string, start int, end int) ([]string, error)

	SetStruct(ctx context.Context, k string, vStruct any) error
	SetStructExpire(ctx context.Context, k string, vStruct any, ddl time.Duration) error

	//TODO 这个方法有问题！！！！
	GetStruct(ctx context.Context, k string, targetStruct any) error

	SAdd(ctx context.Context, k string, v ...string) error
	SAddExpire(ctx context.Context, ddl time.Duration, k string, v string) error
	SISMember(ctx context.Context, k string, v string) bool
	SCard(ctx context.Context, k string) int
	SMembers(ctx context.Context, k string) ([]string, error)

	HSet(ctx context.Context, k string, v ...any) error
	HGet(ctx context.Context, k string, field string) (any, error)

	Del(ctx context.Context, k string) error

	ExecuteLuaScript(ctx context.Context, luaScript string, k string) (any, error)
	ExecuteArgsLuaScript(ctx context.Context, luaScript string, keys []string, args ...interface{}) (err error, retValue []any)

	IsEmpty(err error) bool

	// Lru 此调用方式不太合理 哪有从redis里面调用lru的 有也是从lru里面调用redis实现
	//Lru(ctx context.Context, maxCapacity int, dataType int) lru.Lru
}

type redisClient struct {
	rcl *redis.Client
}

func NewRedisClient(r *InitRedisApplication) (Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     r.UserAddr,
		Password: r.Password,
	})

	return &redisClient{rcl: client}, nil
}

func (r *redisClient) Ping(ctx context.Context) error {
	_, err := r.rcl.Ping(ctx).Result()
	return err
}

func (r *redisClient) Set(ctx context.Context, k string, v any) error {
	return r.rcl.Set(ctx, k, v, 0).Err()
}

func (r *redisClient) SetExpire(ctx context.Context, k string, v any, ddl time.Duration) error {
	return r.rcl.Set(ctx, k, v, ddl).Err()
}

func (r *redisClient) Get(ctx context.Context, k string) (string, error) {
	return r.rcl.Get(ctx, k).Result()
}

func (r *redisClient) LRange(ctx context.Context, k string, start int, end int) ([]string, error) {
	return r.rcl.LRange(ctx, k, int64(start), int64(end)).Result()
}

func (r *redisClient) LRangeAll(ctx context.Context, k string) ([]string, error) {
	return r.LRange(ctx, k, 0, -1)
}

func (r *redisClient) LPush(ctx context.Context, k string, v any) error {
	return r.rcl.LPush(ctx, k, v).Err()
}

func (r *redisClient) RPop(ctx context.Context, k string) (string, error) {
	return r.rcl.RPop(ctx, k).Result()
}

func (r *redisClient) SetStruct(ctx context.Context, k string, vStruct any) error {
	vJsonData, _ := json.Marshal(vStruct)
	return r.Set(ctx, k, vJsonData)
}

func (r *redisClient) SetStructExpire(ctx context.Context, k string, vStruct any, ddl time.Duration) error {
	vJsonData, _ := json.Marshal(vStruct)
	return r.SetExpire(ctx, k, vJsonData, ddl)
}

// GetStruct 获取自定义结构体
func (r *redisClient) GetStruct(ctx context.Context, k string, targetStruct any) error {
	vJsonData, err := r.rcl.Get(ctx, k).Result() // 获取存储的 JSON 字符串
	if err != nil {
		return err
	}
	// 将 JSON 字符串反序列化为结构体
	err = json.Unmarshal([]byte(vJsonData), targetStruct)
	if err != nil {
		return err
	}
	return nil
}

func (r *redisClient) ZRem(ctx context.Context, k string, vs ...any) (int64, error) {
	return r.rcl.ZRem(ctx, k, vs).Result()
}

// ZScore 获取指定元素的分数 Zset
func (r *redisClient) ZScore(ctx context.Context, k string, member string) (isExist bool, score int) {
	result, err := r.rcl.ZScore(ctx, k, member).Result()
	if r.IsEmpty(err) {
		return false, common.FalseInt
	}
	return true, int(result)
}

func (r *redisClient) ZRangeAll(ctx context.Context, k string) ([]string, error) {
	return r.rcl.ZRange(ctx, k, 0, -1).Result()
}

func (r *redisClient) ZRange(ctx context.Context, k string, start int, end int) ([]string, error) {
	return r.rcl.ZRange(ctx, k, int64(start), int64(end)).Result()
}

func (r *redisClient) LRem(ctx context.Context, k string, count int, v any) (int64, error) {
	return r.rcl.LRem(ctx, k, int64(count), v).Result()
}

func (r *redisClient) SAddExpire(ctx context.Context, ddl time.Duration, k string, v string) error {
	err := r.SAdd(ctx, k, v)
	if err != nil {
		return err
	}
	return r.rcl.Expire(ctx, k, ddl).Err()
}

func (r *redisClient) SAdd(ctx context.Context, k string, v ...string) error {
	return r.rcl.SAdd(ctx, k, v).Err()
}

func (r *redisClient) SISMember(ctx context.Context, k string, v string) bool {
	result, _ := r.rcl.SIsMember(ctx, k, v).Result()
	return result
}

func (r *redisClient) SCard(ctx context.Context, k string) int {
	result, _ := r.rcl.SCard(ctx, k).Result()
	return int(result)
}

func (r *redisClient) SMembers(ctx context.Context, k string) ([]string, error) {
	return r.rcl.SMembers(ctx, k).Result()
}

// HSet 支持批量添加 但是kv必须成对出现
func (r *redisClient) HSet(ctx context.Context, k string, v ...any) error {
	return r.rcl.HSet(ctx, k, v).Err()
}

func (r *redisClient) HGet(ctx context.Context, k string, field string) (any, error) {
	return r.rcl.HGet(ctx, k, field).Result()
}

func (r *redisClient) Del(ctx context.Context, k string) error {
	return r.rcl.Del(ctx, k).Err()
}

// ExecuteLuaScript 执行lua脚本 保证操作原子性
func (r *redisClient) ExecuteLuaScript(ctx context.Context, luaScript string, k string) (any, error) {
	result, err := r.rcl.Eval(ctx, luaScript, []string{k}).Result()
	if err != nil {
		return nil, err
	}
	return result, err
}

func (r *redisClient) ExecuteArgsLuaScript(ctx context.Context, luaScript string, keys []string, args ...interface{}) (err error, retValue []any) {
	res, err := r.rcl.Eval(ctx, luaScript, keys, args).Result()
	if err != nil {
		return err, nil
	}

	// 如果 result 是一个切片，直接返回
	if res, ok := res.([]interface{}); ok {
		return err, res
	}

	// 否则，将 result 包装为切片并返回
	return err, []interface{}{res}
}

func (r *redisClient) IsEmpty(err error) bool {
	return errors.Is(err, redis.Nil)
}

type InitRedisApplication struct {
	UserAddr string
	Password string
}

func NewRedisApplication(addr string, password string) *InitRedisApplication {
	return &InitRedisApplication{
		UserAddr: addr,
		Password: password,
	}
}
