package userinfo

import "go.mongodb.org/mongo-driver/bson"

const (
	UserInfoMongo = "user_info"
	UserNameField = "user_name"
	AgeField      = "age"
)

func UserName(v string) UserInfoPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   UserNameField,
			Value: bson.M{"$eq": v},
		})
	}
}

func UserNameEQ(v string) UserInfoPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   UserNameField,
			Value: bson.M{"$eq": v},
		})
	}
}

func UserNameNE(v string) UserInfoPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   UserNameField,
			Value: bson.M{"$ne": v},
		})
	}
}

func UserNameRegex(v string) UserInfoPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   UserNameField,
			Value: bson.M{"$regex": v},
		})
	}
}

func Age(v int) UserInfoPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$eq": v},
		})
	}
}

func AgeEQ(v int) UserInfoPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$eq": v},
		})
	}
}

func AgeNE(v int) UserInfoPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$ne": v},
		})
	}
}

func AgeGT(v int) UserInfoPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$gt": v},
		})
	}
}

func AgeLT(v int) UserInfoPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$lt": v},
		})
	}
}

func AgeGTE(v int) UserInfoPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$gte": v},
		})
	}
}

func AgeLTE(v int) UserInfoPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$lte": v},
		})
	}
}

type UserInfoPredicate func(*bson.D)
