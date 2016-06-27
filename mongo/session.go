package mongo

import (
	"github.com/name5566/leaf/db/mongodb"
	"fmt"
	"gopkg.in/mgo.v2"
)

/**
 * 连接会话
 */
var (
	gContext *mongodb.DialContext
	gDatabase string
)

/**
 * 开启会话
 * uri: [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
 */
func Dial(url, database string, sessionNum int) (err error) {
	gDatabase = database
	gContext, err = mongodb.Dial(url, sessionNum)
	return
}

/**
 * 关闭会话
 */
func Close() error {
	if nil == gContext {
		return fmt.Errorf("关闭未连接会话")
	}
	gContext.Close()
	gContext = nil
	gDatabase = ""
	return nil
}

/***
 * 执行操作
 */
func Do(collection string, do func(c *mgo.Collection) error) error {
	s := gContext.Ref()
	defer gContext.UnRef(s)

	c := s.DB(gDatabase).C(collection)

	if e := do(c); nil != e {
		return e
	} else {
		return nil
	}
}