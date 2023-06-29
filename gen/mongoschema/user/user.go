package user

import "go.mongodb.org/mongo-driver/bson"

const (
	UserMongo     = "user"
	UserNameField = "user_name"
	AgeField      = "age"
)

func UserName(v string) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   UserNameField,
			Value: bson.M{"$eq": v},
		})
	}
}

func UserNameEQ(v string) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   UserNameField,
			Value: bson.M{"$eq": v},
		})
	}
}

func UserNameNE(v string) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   UserNameField,
			Value: bson.M{"$ne": v},
		})
	}
}

func UserNameRegex(v string) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   UserNameField,
			Value: bson.M{"$regex": v},
		})
	}
}

func UserNameIN(v ...string) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   UserNameField,
			Value: bson.M{"$in": v},
		})
	}
}

func UserNameNIN(v ...string) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   UserNameField,
			Value: bson.M{"$nin": v},
		})
	}
}

func Age(v int) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$eq": v},
		})
	}
}

func AgeEQ(v int) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$eq": v},
		})
	}
}

func AgeNE(v int) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$ne": v},
		})
	}
}

func AgeGT(v int) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$gt": v},
		})
	}
}

func AgeLT(v int) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$lt": v},
		})
	}
}

func AgeGTE(v int) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$gte": v},
		})
	}
}

func AgeLTE(v int) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$lte": v},
		})
	}
}

func AgeIN(v ...int) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$in": v},
		})
	}
}

func AgeNIN(v ...int) UserPredicate {
	return func(d *bson.D) {
		*d = append(*d, bson.E{
			Key:   AgeField,
			Value: bson.M{"$nin": v},
		})
	}
}

type UserPredicate func(*bson.D)
