package models

import (
	"github.com/astaxie/beego"
	"labix.org/v2/mgo"
)

var (
	Db *mgo.Database //数据库
)

type Base struct {
	Conn *mgo.Collection //数据库连接
}

func init() {
	//获取数据库连接
	session, err := mgo.Dial(beego.AppConfig.String(beego.RunMode + "::MongoDb"))
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	Db = session.DB("kingdoms")
}

/* 初始化 */
func (this *Base) Init(dbName string) {
	this.Conn = Db.C(dbName)
}

/* 增 */
func (this *Base) Insert(docs interface{}) error {
	return this.Conn.Insert()
}

/* 删 */
func (this *Base) Delete(id interface{}) error {
	return this.Conn.RemoveId(id)
}

/* 改 */
func (this *Base) Update(id interface{}, update interface{}) error {
	return this.Conn.UpdateId(id, update)
}

/* 查 */
func (this *Base) Find(id interface{}, result interface{}) error {
	return this.Conn.FindId(id).One(&result)
}
