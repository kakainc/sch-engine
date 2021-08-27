package init_sch

import (
	"Sch_engine/rely"
	"github.com/go-redis/redis"
	"strconv"
)

func RedisInit() *redis.Client {
	redisdb := redis.NewClient(&redis.Options{
		Addr:     rely.IndexRedisAddr,
		Password: rely.IndexRedisAuth,
		DB:       0,
	})
	_, err := redisdb.Ping().Result()
	if err != nil {
		Error.Printf("Can get redisdb! error: %s", err)
	}
	Info.Println("Connected to Redis!")
	return redisdb
}

// ZRange 有序集合ZSET
func ZRange(rdb *redis.Client, zsetKey string, nRecall int64) []rely.PidScore {
	var result []rely.PidScore
	ret, err := rdb.ZRevRangeWithScores(zsetKey, 0, nRecall).Result()
	if err != nil {
		panic(err)
	}

	for _, z := range ret {
		pid, _ := strconv.Atoi(z.Member.(string))
		result = append(result, rely.PidScore{
			Pid: int64(pid),
			Score: z.Score,
		})
	}
	return result
}

// ZGet 有序集合ZSET
func ZGet(rdb *redis.Client, zsetKey string, pid string) float64 {
	ret, err := rdb.ZScore(zsetKey, pid).Result()
	if err != nil {
		//Info.Printf("%s dont have comment!!!", pid)
		return 0
	}
	return ret
}

// HGet 哈希Hash
func HGet(rdb *redis.Client, hKey string, pid string) float64 {
	ret, err := rdb.HGet(hKey, pid).Result()
	if err != nil {
		return 0
	}
	ctr, _ := strconv.ParseFloat(ret, 64)
	return ctr
}
