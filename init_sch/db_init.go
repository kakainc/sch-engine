package init_sch

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

func ServiceDbInit() (*redis.Client, *mongo.Collection){
	rd := RedisInit()
	mg := MongoInit()
	return rd, mg
}