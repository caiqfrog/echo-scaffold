package mongo

import (
	"fmt"
	"log"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TestType struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Name  string        `bson:"name"`
	Value string        `bson:"value"`
}

func (self *TestType) MongoID() interface{} {
	return self.Id
}

func (self *TestType) MongoCollection() string {
	return "test"
}

func (self *TestType) SetObjectId(bid bson.ObjectId) {
	self.Id = bid
}

func Test_Mgo(t *testing.T) {
	if e := Dial("mongodb://localhost:27017", "colledge", 10); nil != e {
		log.Fatal("dial: ", e)
		return
	}
	defer Close()
	data := &TestType{
		Name:  "hello",
		Value: "value",
	}
	//
	var err error
	r := NewRepository(data)
	// insert
	err = r.Insert()
	// ci, e := Upsert(thesis.MgoDatabase, _t, _t)
	fmt.Println("insert: ", err, data, data.Id)
	// find
	_data, err := r.Find(bson.M{"name": "hello"})
	fmt.Println("find: ", err, _data)
	fmt.Println("find array: ", _data.([]*TestType)[0])
	// find by id
	data3 := new(TestType)
	data3.Id = bson.ObjectIdHex("574e5451046b15f76ba90c23")
	r.Reset(data3)
	err = r.FindId()
	fmt.Println("findById: ", err, data3)
	// update
	err = r.Update(
		bson.M{
			"_id": bson.ObjectIdHex("574e5451046b15f76ba90c23"),
		},
		bson.M{
		"value": time.Now().String(),
	})
}
