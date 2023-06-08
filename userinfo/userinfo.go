package userinfo

import "go.mongodb.org/mongo-driver/bson"

const (
	UserInfoMongo = "user_info"
	UserNameField = "user_name"
	AgeField      = "age"
)

func FindUserInfoByUserName(userName string) bson.M {
	return bson.M{UserNameField: userName}
}

func FindUserInfoByAge(age int) bson.M {
	return bson.M{AgeField: age}
}
