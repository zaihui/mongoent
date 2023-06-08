package go_mongo

import "go.mongodb.org/mongo-driver/bson"

type UserInfoQuery struct {
	config
	Conditions bson.M
}
func (uq *UserInfoQuery) Where(ps ...bson.M)*UserInfoQuery{
	for _, p := range ps {
		for s, v := range p {
			uq.Conditions[s] = v
		}
	}
	return uq
}
