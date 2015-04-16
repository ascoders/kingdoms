package models

import (
	"github.com/astaxie/beego"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Version struct {
	CardVersion  bson.ObjectId `bson:"c"` //卡牌数据库版本
	LevelVersion bson.ObjectId `bson:"l"` //关卡数据库版本
	ImageVersion bson.ObjectId `bson:"i"` //图片数据库版本
}

var (
	versionC *mgo.Collection //数据库连接
)

func init() {
	//获取数据库连接
	session, err := mgo.Dial(beego.AppConfig.String(beego.RunMode + "::MongoDb"))
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	versionC = session.DB("kingdoms").C("version")
}

/* 更新版本 */
func (this *Version) UpdateVersion(name string) {
	colQuerier := bson.M{"_id": 1}
	change := bson.M{"$set": bson.M{name: bson.NewObjectId()}}

	err := versionC.Update(colQuerier, change)
	if err != nil {
		//处理错误
	}
}

/* 获取版本信息 */
func (this *Version) GetVersion() {
	//查询用户
	err := versionC.Find(bson.M{"_id": 1}).One(&this)
	if err != nil {
		return
	}
}
