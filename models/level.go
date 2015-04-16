package models

import (
	"labix.org/v2/mgo/bson"
	"strconv"
)

type Level struct {
	Id          uint16 `json:"_id" bson:"_id" form:"id"`        //主键
	Name        string `json:"n" bson:"n" form:"name"`          //关卡名称
	Description string `json:"dn" bson:"dn" form:"description"` //关卡描述
	ParentName  string `json:"p" bson:"p" form:"parentname"`    //父级关卡名称
	RoomParent  uint8  `json:"r" bson:"r" form:"roomparent"`    //父级关卡id
	Information string `json:"i" bson:"i" form:"information"`   //敌人信息
	Dialogue    string `json:"d" bson:"d" form:"dialogue"`      //开场白
	Reward      string `json:"re" bson:"re" form:"reward"`      //奖励
	Special     string `json:"s" bson:"s" form:"special"`       //附加特殊字段
	StarArray   string `json:"sa" bson:"sa" form:"stararray"`   //出场珠子种类
}

var (
	basec Base //基础数据库
)

func init() {
	//初始化数据库
	basec.Init("level")
}

//根据id查询关卡
func (this *Level) FindOne(id int) bool {
	//查询关卡
	err := basec.Conn.Find(bson.M{"_id": id}).One(&this)
	if err != nil {
		return false
	}
	return true
}

//获取关卡总行数
func (this *Level) GetCount() int {
	//查询关卡
	n, err := basec.Conn.Count()
	if err != nil {
		return 0
	}
	return n
}

/* 下载数据库时提供的信息 */
func (this *Level) DownLoad() []*Level {
	var result []*Level
	err := basec.Conn.Find(bson.M{}).Sort("_id").Select(bson.M{"re": 0, "s": 0}).All(&result)
	if err != nil {
		return nil
	}
	return result
}

/* 获取父级关卡的出场 珠宝 字符串 */
func (this *Level) FindParentStarArray() string {
	var result *Level
	err := basec.Conn.Find(bson.M{"r": this.RoomParent}).Sort("_id").Select(bson.M{"sa": 1}).One(&result)
	if err != nil {
		return ""
	}
	return result.StarArray
}

/* 插入新关卡 */
func (this *Level) Insert() {
	basec.Conn.Insert(this)
}

/* 查询总数 */
func (this *Level) Count() int {
	number, err := basec.Conn.Find(bson.M{}).Count()
	if err != nil {
		return 0
	}
	return number
}

/* 查询一定数目的关卡
 * @params from 起始位置
 * @params number 查询数量
 */
func (this *Level) Find(from int, number int) []*Level {
	var result []*Level
	err := basec.Conn.Find(bson.M{}).Sort("_id").Skip(from).Limit(number).All(&result)
	if err != nil {
		return nil
	}
	return result
}

/* 保存关卡信息 */
func (this *Level) Update(id uint16) {
	basec.Conn.Update(bson.M{"_id": id}, &this)
}

/* 查询某个id是否存在 */
func (this *Level) IdExist(id uint16) bool {
	var card *Level
	err := basec.Conn.Find(bson.M{"_id": id}).One(&card)
	if err != nil {
		return false
	}
	return true
}

/* 删除关卡 */
func (this *Level) Delete(id string) {
	levelId, _ := strconv.Atoi(id)
	basec.Conn.Remove(bson.M{"_id": uint16(levelId)})
}
