package go_mongo

type UserInfo struct {
	UserName string `bson:"user_name"`
	Age      int    `bson:"age"`
}