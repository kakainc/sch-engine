module Sch_engine

go 1.16

require (
	github.com/adamzy/cedar-go v0.0.0-20170805034717-80a9c64b256d // indirect
	github.com/go-nlp/bm25 v1.0.0
	github.com/go-nlp/tfidf v1.1.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/huichen/sego v0.0.0-20180617034105-3f3c8a8cfacc
	github.com/issue9/assert v1.4.1 // indirect
	github.com/onsi/gomega v1.15.0 // indirect
	go.mongodb.org/mongo-driver v1.7.1
)

replace github.com/huichen/sego v0.0.0-20180617034105-3f3c8a8cfacc => github.com/kakainc/sego v1.0.0
