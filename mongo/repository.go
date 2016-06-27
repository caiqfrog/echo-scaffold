package mongo

import (
	"reflect"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type mongoRepository struct {
	seed MongoModal
	seedType reflect.Type
}

// create
func NewRepository(seed MongoModal) *mongoRepository {
	return &mongoRepository{
		seed: seed,
		seedType: reflect.TypeOf(seed),
	}
}

func (m *mongoRepository) Get() MongoModal {
	return m.seed
}

func (m *mongoRepository) Reset(seed MongoModal) {
	m.seed = seed
	m.seedType = reflect.TypeOf(seed)
}

func (m *mongoRepository) Insert() error {
	return Do(m.seed.MongoCollection(), func(c *mgo.Collection) error {
		id := m.seed.MongoID()
		if _, ok := id.(bson.ObjectId); ok && bson.ObjectId("") == id {
			bid := bson.NewObjectId()
			m.seed.SetObjectId(bid)
		}
		return c.Insert(m.seed)
	})
}

func (m *mongoRepository) FindId(id ...interface{}) error {
	var _id interface{}
	if len(id) == 0 {
		_id = m.seed.MongoID()
	} else {
		_id = id[0]
	}
	return Do(m.seed.MongoCollection(), func(c *mgo.Collection) error {
		query := c.FindId(_id)
		return query.One(m.seed)
	})
}

func (m *mongoRepository) Find(query bson.M) (data interface{}, err error) {
	err = Do(m.seed.MongoCollection(), func(c *mgo.Collection) error {
		query := c.Find(query)
		if _, e := query.Count(); nil != e {
			return e
		} else {
			t := reflect.SliceOf(m.seedType)

			v := reflect.New(t)

			if e := query.All(v.Interface()); nil != e {
				return e
			} else {
				data = v.Elem().Interface()
				return nil
			}
		}
	})
	return
}

func (m *mongoRepository) UpdateId(set bson.M) error {
	return Do(m.seed.MongoCollection(), func(c *mgo.Collection) error {
		return c.UpdateId(m.seed.MongoID(), set)
	})
}

func (m *mongoRepository) Update(query, set bson.M) error {
	return Do(m.seed.MongoCollection(), func(c *mgo.Collection) error {
		return c.Update(query, set)
	})
}

func (m *mongoRepository) RemoveId(id ...interface{}) error {
	var _id interface{}
	if len(id) == 0 {
		_id = m.seed.MongoID()
	} else {
		_id = id[0]
	}
	return Do(m.seed.MongoCollection(), func(c *mgo.Collection) error {
		return c.RemoveId(_id)
	})
}

func (m *mongoRepository) Remove(query bson.M) error {
	return Do(m.seed.MongoCollection(), func(c *mgo.Collection) error {
		return c.Remove(query)
	})
}
