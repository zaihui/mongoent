package spec

type User struct {
	UserName string `bson:"user_name"`
	Age      int    `bson:"age"`
}

type UserInfo struct {
	UserName string `bson:"user_name"`
	Age      int    `bson:"age"`
}
