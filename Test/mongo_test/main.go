package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestStruct struct {
	Name string
	ID int32
}

func main() {
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar?readPreference=primary&ssl=false"))
	if err != nil {
		panic(err)
	}
	defer mc.Disconnect(c)
	col := mc.Database("coolcar").Collection("account")

	insertRows(c, col)
	change()
}
func change()  {
	data, err := bson.Marshal(&TestStruct{Name: "Bob"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Marshal:%q", data)
	value := TestStruct{}
	err2 := bson.Unmarshal(data, &value)
	if err2 != nil {
		panic(err)
	}
	fmt.Println("value:", value)

	mmap := bson.M{}
	err3 := bson.Unmarshal(data, mmap)
	if err3!=nil{
		panic(err3)
	}
	fmt.Println("mmap:", mmap)
}
func findRows(c context.Context, col *mongo.Collection) {
	cur, err := col.Find(c, bson.M{})
	if err != nil {
		panic(err)
	}
	for cur.Next(c) {
		var row struct {
			ID     primitive.ObjectID `bson:"_id"`
			OpenID string             `bson:"open_id"`
		}
		err = cur.Decode(&row)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", row)
	}
}
func createIndex(c context.Context, col *mongo.Collection) {

}
func insertRows(c context.Context, col *mongo.Collection) {
	res, err := col.InsertMany(c, []interface{}{
		bson.M{
			"open_id": "123",
			"name":    "123",
		},
		bson.M{
			"open_id": "456",
			"name":    "456",
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", res)
}
