package rely

// LogFile to write
const LogFile = "./rely/log.txt"

// local_redis
const (
IndexRedisAddr = "127.0.0.1:6379"
IndexRedisAuth = ""
)

// local_mongo
const (
	MongoAddr = "mongodb://127.0.0.1:27017"
	DataBase = "newpost"
	Collection = "postinfo"
)


// 结构体

type PidScore struct {
	Pid   int64
	Score float64
}

type PidInfo struct {
	Pid      string
	Ptype    string
	Tid      string
	Tname    string
	Content  string
	Comment  string
}

type RetPids struct {
	Pids []int64
}


// 参数
const(
	CONTENTWEIGHT = 1.0
	COMMENTWEIGHT = 1.0
	TNAMEWEIGHT = 1.0
	REPEATWEIGHT = 1.0
)

const (
	BMWEIGHT = 1.0
	HOTWEIGHT = 1.0
	CTRWEIGHT = 1.0
	MODELWEIGHT = 1.0
)


