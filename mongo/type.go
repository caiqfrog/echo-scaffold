package mongo

import "gopkg.in/mgo.v2/bson"

type Data interface {
	ID() bson.ObjectId
}

type MongoModal interface {
	MongoID() interface{}
	SetObjectId(id bson.ObjectId)
	MongoCollection() string
}
