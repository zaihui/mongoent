package mongoschema

type UserInfo struct {
	UserName string `bson:"user_name"`
	Age      int    `bson:"age"`
}
