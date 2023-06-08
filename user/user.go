package user

import "go.mongodb.org/mongo-driver/bson"

const (
	UserMongo     = "user"
	UserNameField = "user_name"
	AgeField      = "age"
)

func FindUserByUserName(userName string) bson.M {
	return bson.M{UserNameField: userName}
}

func FindUserByAge(age int) bson.M {
	return bson.M{AgeField: age}
}
