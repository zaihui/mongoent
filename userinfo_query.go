package go_mongo

import (
	"context"
	"cc/go-mongo/userinfo"
	"go.mongodb.org/mongo-driver/bson"
)
type UserInfoQuery struct {
	config
	Predicates []userinfo.UserInfoPredicate
	dbName string

}
func (uq *UserInfoQuery) Where(ps ...userinfo.UserInfoPredicate)*UserInfoQuery{
	for _, p := range ps {
		uq.Predicates = append(uq.Predicates, p)
	}
	return uq
}
func (uq *UserInfoQuery) All(ctx context.Context)([]*UserInfo,error) {
	filter := bson.D{}
	for _, p := range uq.Predicates {
		p(&filter)
	}
	cur, err := uq.Database(uq.dbName).Collection(userinfo.UserInfoMongo).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	temp := make([]*UserInfo, 0, 0)
	for cur.Next(ctx) {
		var u UserInfo
		err = cur.Decode(&u)
		if err != nil {
			return nil, err
		}
		temp = append(temp, &u)
	}
	if err = cur.Err(); err != nil {
		return nil, err
	}
	return temp, nil
}
