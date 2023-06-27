package go_mongo

import (
	"context"
	"cc/go-mongo/userinfo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserInfoQuery struct {
	config
	Predicates []userinfo.UserInfoPredicate
	limit  *int64
	offset *int64
	dbName string
	options bson.D

}

func (uq *UserInfoQuery) Limit(limit int64) *UserInfoQuery{
	uq.limit = &limit
	return uq
}

func (uq *UserInfoQuery) Offset(offset int64) *UserInfoQuery{
	uq.offset = &offset
	return uq
}

func (uq *UserInfoQuery) Order(o ...OrderFunc) *UserInfoQuery {
	for _, fn := range o {
		fn(&uq.options)
	}
	return uq
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

	o := options.Find()
	if uq.limit != nil && *uq.limit != 0 {
		o = o.SetLimit(*uq.limit)
	}
	if uq.offset != nil && *uq.offset != 0 {
		o = o.SetSkip(*uq.offset)
	}
	o.SetSort(uq.options)
	cur, err := uq.Database(uq.dbName).Collection(userinfo.UserInfoMongo).Find(ctx, filter,o)
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
