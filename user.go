package go_mongo

type User struct {
	UserName string `bson:"user_name"`
	Age      int    `bson:"age"`
}