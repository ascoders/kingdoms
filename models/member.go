package models

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strconv"
	"strings"
	"time"
)

//用户
type Member struct {
	Id              bson.ObjectId `bson:"_id"` //主键
	WokuId          string        `bson:"a"`   //官网关联id string类型的ObjectId
	Name            string        `bson:"n"`   //昵称
	Password        string        `bson:"p"`   //密码
	Uniqid          int           `bson:"u"`   //不重复的卡牌id
	Team            string        `bson:"t"`   //出场队伍
	SwitchTeam      int           `bson:"s"`   //选择哪一只队伍
	BackPack        string        `bson:"b"`   //背包卡牌
	Gold            int           `bson:"g"`   //金币数
	MaxLevel        uint16        `bson:"m"`   //当前达到的最大关卡数
	Friend          string        `bson:"f"`   //好友列表
	FriendRequest   string        `bson:"fr"`  //好友申请列表
	LastGet         string        `bson:"l"`   //上一句战斗情况
	LastOperateTime time.Time     `bson:"lo"`  //最后操作时间
	AllStamina      uint16        `bson:"as"`  //最大疲劳
	NowStamina      uint16        `bson:"ns"`  //当前疲劳
	AllEnergy       float64       `bson:"ae"`  //最大能量（pk消耗）
	NowEnergy       float64       `bson:"ne"`  //当前能量（pk消耗）
	Power           uint32        `bson:"po"`  //最大战斗力
	AttackNumber    uint16        `bson:"an"`  //总pk次数
	DefenseNumber   uint16        `bson:"dn"`  //总被他人pk次数
	PkOpponentId    string        `bson:"poi"` //当前pk对手的id
	PkTeam          string        `bson:"pt"`  //PK赛中最后一次获胜的队伍id列表
	Shop            string        `bson:"sp"`  //卡牌商店，列举了当天可购买的卡牌
}

var (
	memberC *mgo.Collection //数据库连接
)

func init() {
	//获取数据库连接
	session, err := mgo.Dial(beego.AppConfig.String(beego.RunMode + "::MongoDb"))
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	memberC = session.DB("kingdoms").C("member")
}

/* 插入用户 */
func (this *Member) Insert() string {
	this.Id = bson.NewObjectId()
	this.LastOperateTime = bson.Now()
	this.MaxLevel = 1     //最大关卡数初始为1
	this.AllStamina = 100 //最大疲劳值初始值为100
	this.NowStamina = 100 //当前疲劳值初始值为100
	this.AllEnergy = 20   //最大能量初始值为20
	this.NowEnergy = 20   //当前能量初始值为20
	this.Gold = 100       //金币初始值为100

	err := memberC.Insert(this)
	if err != nil {
		return ""
	}

	return bson.ObjectId.Hex(this.Id)
}

/* 根据ObjectId和pass查询某个用户信息 */
func (this *Member) CheckOne(id string, pass string) bool {
	if !bson.IsObjectIdHex(id) {
		return false
	}
	//查询用户
	err := memberC.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&this)
	if err != nil {
		return false //账号不存在
	}
	//对比
	if pass != this.Password { //验证出错
		return false
	} else { //通过验证
		return true
	}
}

/* 根据ObjectId查询用户信息 */
func (this *Member) FindAccount(id string) *Member {
	//查询用户
	var member *Member
	err := memberC.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&member)
	if err != nil {
		return nil //账号不存在
	} else {
		return member
	}
}

/* 更新用户信息 */
func (this *Member) Update() {
	err := memberC.UpdateId(this.Id, this)
	if err != nil {
		//处理错误
	}
}

/* 插入一张卡牌 */
func (this *Member) AddCard(ids ...int) {
	//背包拆分为数组
	backPackArray := strings.Split(this.BackPack, ";")
	//遍历背包数组
	for _, id := range ids {
		exist := false
		if len(backPackArray) > 0 {
			for k, v := range backPackArray {
				if v == "" {
					continue
				}
				info := strings.Split(v, ":")
				detail := strings.Split(info[1], ",")
				//判断这张卡牌是否已存在
				if detail[0] == strconv.Itoa(id) {
					exist = true
					//存在数量+1
					this.StringFloatMath(&detail[1], 1)
					info[1] = strings.Join(detail, ",")
					backPackArray[k] = strings.Join(info, ":")
					break
				}
			}
		}
		if !exist {
			//卡牌不重复id : 卡牌id +体力 +行动力 +攻击力 +防御力 当前经验 卡牌等级 杀敌数 伤害 加血
			cardString := strconv.Itoa(this.Uniqid) + ":" + strconv.Itoa(id) + ",0,0,0,0,0,0,0,0,0"
			this.Uniqid++
			//追加这张卡牌
			if backPackArray[0] == "" {
				backPackArray[0] = cardString
			} else {
				backPackArray = append(backPackArray, cardString)
			}
		}
	}
	//保存修改后的背包
	this.BackPack = strings.Join(backPackArray, ";")
}

/* 查询战斗力比自己高的对手（如果自己是最高，则查找第二高的对手） */
func (this *Member) FindOpponent() *Member {
	//查询用户
	var member *Member
	//找战斗力比他大的最小的一个
	err := memberC.Find(bson.M{"po": bson.M{"$gt": this.Power}}).Sort("po").One(&member)
	if err != nil { //没有战斗力比他高的对手
		//找战斗力比他小的最大的一个
		err = memberC.Find(bson.M{"po": bson.M{"$lt": this.Power}}).Sort("-po").One(&member)
		if err != nil {
			return nil
		} else {
			return member
		}
	} else {
		return member
	}
}

/* 让字符串转为数组进行运算，再转换为字符串 */
func (this *Member) StringFloatMath(value *string, number float64) {
	f, err := strconv.ParseFloat(*value, 10)
	if err == nil {
		*value = strconv.FormatFloat(f+number, 'f', 0, 64)
	}
}

/* 生成签名 md5(id&expire&pay&token) */
func (this *Member) CreateSign(expire int64, pay int64) string {
	token := this.Id.Hex() + "&" + strconv.Itoa(int(expire)) + "&" + strconv.Itoa(int(pay)) + "&" + this.Password
	m := md5.New()
	m.Write([]byte(token))
	token = hex.EncodeToString(m.Sum(nil))
	return token
}
