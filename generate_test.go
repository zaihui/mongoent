package go_mongo

import (
	user1 "cc/go-mongo/user"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	GetStructsFromGoFile("/Users/zh/sdk/go1.16/demo/go-mongo/model.go")
}
func TestQueryUserInfo(t *testing.T) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	newClient := NewClient(Driver(*client))

	all, err := newClient.User.SetDBName("my_mongo").Query().Where(
		user1.FindUserByUserName("ccc"),
		user1.FindUserByAge(2),
	).All(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(all[0])

}
