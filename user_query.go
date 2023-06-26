package go_mongo

import (
	"cc/go-mongo/user"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type UserQuery struct {
	config
	Predicates []user.UserPredicate
	dbName     string
}

func (uq *UserQuery) Where(ps ...user.UserPredicate) *UserQuery {
	for _, p := range ps {
		uq.Predicates = append(uq.Predicates, p)
	}
	return uq
}
func (uq *UserQuery) All(ctx context.Context) ([]*User, error) {
	filter := bson.D{}
	for _, p := range uq.Predicates {
		p(&filter)
	}
	cur, err := uq.Database(uq.dbName).Collection(user.UserMongo).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	temp := make([]*User, 0, 0)
	for cur.Next(ctx) {
		var u User
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
