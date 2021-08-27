package init_sch

import (
	"Sch_engine/rely"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoInit() *mongo.Collection  {
	clientOptions := options.Client().ApplyURI(rely.MongoAddr)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		Error.Println(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		Error.Println(err)
	}
	Info.Println("Connected to MongoDB!")

	postCol := client.Database(rely.DataBase).Collection(rely.Collection)
	return postCol
}

func FindOne(collection *mongo.Collection, id string) rely.PidInfo {
	filter := bson.D{
		{"_id", id},
	}

	original := make(map[string]interface{})
	var postData rely.PidInfo
	err := collection.FindOne(context.TODO(), filter).Decode(&original)
	if err != nil {
		Warning.Printf("%s:no pid !!!", id)
	} else {
		pid, _ := original["_id"].(string)
		pType, _ := original["ptype"].(string)
		tid, _ := original["tid"].(string)
		tname, _ := original["tname"].(string)
		content, _ := original["content"].(string)
		comments, _ := original["top100_reviews"].(string)

		postData.Pid = pid
		postData.Ptype = pType
		postData.Tid = tid
		postData.Tname = tname
		postData.Content = content
		postData.Comment = comments
	}
	return postData
}

