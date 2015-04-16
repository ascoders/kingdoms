package models

import (
	"labix.org/v2/mgo/bson"
	"strconv"
	"strings"
)

type Card struct {
	Id           uint16 `json:"_id" bson:"_id" form:"id" valid:"Required"`         //主键
	Name         string `json:"n" bson:"n" form:"name" valid:"Required"`           //姓名
	Ping         string `json:"p" bson:"p" form:"ping" valid:"Required;Alpha"`     //拼音
	Zi           string `json:"z" bson:"z" form:"zi"`                              //字
	Title        string `json:"ti" bson:"ti" form:"title"`                         //封号名
	Type         uint8  `json:"t" bson:"t" form:"type" valid:"Required"`           //类型 0~5 小兵 普通 虎将 神将 封号 主公
	Attribute    uint8  `json:"a" bson:"a" form:"attribute" valid:"Required"`      //属性
	Hp           uint8  `json:"h" bson:"h" form:"hp" valid:"Required"`             //体力
	Attack       uint8  `json:"at" bson:"at" form:"attack" valid:"Required"`       //武技
	AttackNumber uint8  `json:"an" bson:"an" form:"attacknumber" valid:"Required"` //攻击次数
	AttackSpeed  uint8  `json:"as" bson:"as" form:"attackspeed" valid:"Required"`  //行动速度
	Defense      uint8  `json:"d" bson:"d" form:"defense" valid:"Required"`        //防御力
	Lead         uint8  `json:"l" bson:"l" form:"lead" valid:"Required"`           //统帅力
	Cost         uint8  `json:"c" bson:"c" form:"cost" valid:"Required"`           //消耗珠宝数量
	DieSpeak     string `json:"ds" bson:"ds" form:"diespeak"`                      //亡语
	KillSpeak    string `json:"k" bson:"k" form:"killspeak"`                       //杀敌话语
	Skill        string `json:"s" bson:"s" form:"skill"`                           //技能话语描述
	SkillDetail  string `json:"sd" bson:"sd" form:"skilldetail"`                   //技能具体描述
}

var (
	base Base //基础数据库
)

func init() {
	//初始化数据库
	base.Init("card")
}

/* 是否含有某个属性 */
func (this *Card) HasAttribute(attribute uint8) bool {
	switch attribute {
	case 1:
		if this.Attribute&^31 == 32 {
			return true
		} else {
			return false
		}
	case 2:
		if this.Attribute&^47 == 16 {
			return true
		} else {
			return false
		}
	case 3:
		if this.Attribute&^55 == 8 {
			return true
		} else {
			return false
		}
	case 4:
		if this.Attribute&^59 == 4 {
			return true
		} else {
			return false
		}
	case 5:
		if this.Attribute&^61 == 2 {
			return true
		} else {
			return false
		}
	case 6:
		if this.Attribute&^62 == 1 {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

/* 查询所有卡牌信息 */
func (this *Card) FindAllInfo() []*Card {
	var result []*Card
	base.Conn.Find(bson.M{}).All(&result)
	return result
}

/* 查询全部卡牌信息 */
func (this *Card) FindAll() []Card {
	var result []Card
	err := base.Conn.Find(bson.M{}).Select(bson.M{"_id": 1, "n": 1, "p": 1}).All(&result)
	if err != nil {
		return nil
	}
	return result
}

/* 查询某个卡牌的某个技能 */
func (this *Card) FindSkill(id uint16, skillName string) string {
	result := Card{}
	err := base.Conn.Find(bson.M{"_id": id}).Select(bson.M{"sd": 1}).One(&result)
	if err != nil {
		return ""
	}
	skillArray := strings.Split(result.SkillDetail, ";")
	for k, _ := range skillArray {
		if strings.Contains(skillArray[k], skillName) {
			return skillArray[k]
		}
	}
	return ""
}

/* 插入新卡牌 */
func (this *Card) Insert() {
	base.Conn.Insert(this)
}

/* 下载数据库时提供的信息 */
func (this *Card) DownLoad() []*Card {
	var result []*Card
	err := base.Conn.Find(bson.M{}).Select(bson.M{"p": 0, "sd": 0}).All(&result)
	if err != nil {
		return nil
	}
	return result
}

/* 查询总数 */
func (this *Card) Count() int {
	number, err := base.Conn.Find(bson.M{}).Count()
	if err != nil {
		return 0
	}
	return number
}

/* 查询一定数目的卡牌
 * @params from 起始位置
 * @params number 查询数量
 */
func (this *Card) Find(from int, number int) []*Card {
	var result []*Card
	err := base.Conn.Find(bson.M{}).Sort("_id").Skip(from).Limit(number).All(&result)
	if err != nil {
		return nil
	}
	return result
}

/* 保存卡牌信息 */
func (this *Card) Update(id uint16) {
	base.Conn.Update(bson.M{"_id": id}, &this)
}

/* 查询某个id是否存在 */
func (this *Card) IdExist(id uint16) bool {
	var card *Card
	err := base.Conn.Find(bson.M{"_id": id}).One(&card)
	if err != nil {
		return false
	}
	return true
}

/* 删除卡牌 */
func (this *Card) Delete(id string) {
	cardId, _ := strconv.Atoi(id)
	base.Conn.Remove(bson.M{"_id": uint16(cardId)})
}

/*
属性对应表：
1: -> 32
2: -> 16
3: -> 8
4: -> 4
5: -> 2
6: -> 1
3,6-> 9
2,5->18
2,3->24

*/
