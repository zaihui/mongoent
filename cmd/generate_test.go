package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/zaihui/mongoent/gen/mongoschema"
	"github.com/zaihui/mongoent/gen/mongoschema/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGenerate(t *testing.T) {
	// GetStructsFromGoFile("/Users/zh/sdk/go1.16/demo/go-mongo/spec/model.go")

	cc := []string{"a", "b", "c"}
	for s := range cc {
		fmt.Println(s)
	}
}

func TestQueryUserInfo(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100000*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	newClient := mongoschema.NewClient(mongoschema.Driver(*client))

	all, err := newClient.User.SetDBName("my_mongo").Query().
		Where(user.UserNameRegex("c*")).
		All(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(all[0], all[1])
}
