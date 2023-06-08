package go_mongo

import "go.mongodb.org/mongo-driver/bson"

type UserQuery struct {
	config
	Conditions bson.M
}
func (uq *UserQuery) Where(ps ...bson.M)*UserQuery{
	for _, p := range ps {
		for s, v := range p {
			uq.Conditions[s] = v
		}
	}
	return uq
}
