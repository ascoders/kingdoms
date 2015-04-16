package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ascoders/upyun"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"kingdoms/logic"
	"kingdoms/models"
	"labix.org/v2/mgo/bson"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type IndexController struct {
	beego.Controller
	member models.Member
	isOk   bool //用户是否存在且合法
}

var (
	allskill *logic.AllSkill        //技能
	cards    map[uint16]models.Card //卡牌
	images   map[string]string      //所有卡牌图片详细信息
	log      *logs.BeeLogger        //打印日志
	ku       *upyun.UpYun           //又拍云
	newRand  *rand.Rand             //随机

	cardRanks string //武将进阶所需经验
	chuanqi   []int  //传奇神将
	shenjiang []int  //神将
	hujiang   []int  //虎将
	putong    []int  //普通武将
	xiaobing1 []int  //初阶小兵
	xiaobing2 []int  //进阶小兵
	xiaobing3 []int  //高阶小兵
	RliveCost int    //复活并且秒杀敌军消耗金币数量
)

func init() {
	//根据随机时间种子实例化随机数
	newRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	//如果没有日志目录则创建日志目录
	_, err := os.Open("log")
	if err != nil && os.IsNotExist(err) {
		os.Mkdir("log", 0777)
	}
	//初始化日志
	log = logs.NewLogger(10000)
	log.SetLogger("file", `{"filename":"log/kingdoms.log"}`)
	//初始化又拍云
	/// 初始化空间
	ku = upyun.NewUpYun("kingdoms", beego.AppConfig.String("UpyunAccount"), beego.AppConfig.String("UpyunPassword"))
	ku.SetApiDomain("v0.api.upyun.com")
	//初始化技能
	allskill = &logic.AllSkill{}
	//初始化所有图片json值
	updateImage()
	//初始化游戏数据
	cardRanks = "100,500,2000,10000"
	chuanqi = []int{16, 132, 181, 230, 278, 308, 501, 539, 566}
	shenjiang = []int{17, 18, 24, 26, 57, 84, 91, 98, 180, 190, 192, 196, 233, 294, 395, 575}
	hujiang = []int{14, 20, 25, 27, 29, 31, 32, 33, 129, 130, 131, 183, 191, 215, 220, 225, 229, 237, 283, 284, 285, 288, 293, 296, 297, 352, 383, 450, 451, 479, 495, 530, 546, 574, 576, 577}
	putong = []int{20}
	xiaobing1 = []int{1001, 1002, 1003, 1004, 1005, 1006}
	xiaobing2 = []int{1007, 1008, 1009, 1010, 1011, 1012}
	xiaobing3 = []int{1013, 1014, 1015, 1016, 1017, 1018}
	RliveCost = 10
	//赋值所有卡牌
	SetAllCards()
}

func (this *IndexController) Prepare() {
	//查找用户数据
	this.isOk = this.member.CheckOne(this.GetString("id"), this.GetString("pass"))
	if this.isOk { //执行统一逻辑
		//如果是第二天第一次登陆
		if this.member.LastOperateTime.Year() < bson.Now().Year() || this.member.LastOperateTime.YearDay() < bson.Now().YearDay() { //产生时间（天）跨度
			//疲劳恢复满
			this.member.NowStamina = this.member.AllStamina
			//刷新购买商店
			this.RefreshShop()
		}
		//每20秒恢复一点能量，计算经过的秒数
		second := math.Floor(time.Now().Sub(this.member.LastOperateTime).Seconds())
		if this.member.NowEnergy+second/20 > this.member.AllEnergy {
			this.member.NowEnergy = this.member.AllEnergy
		} else {
			this.member.NowEnergy += second / 20
		}
		//更新最后操作时间
		this.member.LastOperateTime = bson.Now()
	}
}

/* -------- 游戏服务 -------- */
/* 用户注册 */
func (this *IndexController) Register() {
	this.member.Name = this.GetString("name")
	//获得当前时间的md5作为原密码
	t := md5.New()
	t.Write([]byte(time.Now().String()))
	pass := hex.EncodeToString(t.Sum(nil))
	this.member.Password = pass
	//用户初始信息
	this.member.Team = "默认卡包-0,1,3"
	this.member.AddCard(395, 539, 288, 183)
	this.RefreshShop()
	//插入用户
	id := this.member.Insert()
	//处理返回数据
	type result struct {
		Id   string `json:"id"`   //用户ID
		Pass string `json:"pass"` //密码
	}
	r := result{id, pass}
	this.Data["json"] = r
	this.ServeJson()
}

/* 用户普通登陆游戏 */
func (this *IndexController) Login() {
	if !this.isOk {
		this.StopRun()
	}
	//登陆返回信息结构体
	type loginInfo struct {
		Name          string  `json:"n"`   //账号名称
		Team          string  `json:"t"`   //出场队伍
		SwitchTeam    int     `json:"s"`   //选择哪一只队伍
		BackPack      string  `json:"b"`   //背包卡牌
		Gold          int     `json:"g"`   //金币数
		MaxLevel      uint16  `json:"m"`   //当前达到的最大关卡数
		Friend        string  `json:"f"`   //好友列表
		FriendRequest string  `json:"fr"`  //好友申请列表
		CardVersion   string  `json:"cv"`  //卡牌数据库版本
		LevelVersion  string  `json:"lv"`  //关卡数据库版本
		ImageVersion  string  `json:"i"`   //图片信息版本
		AllStamina    uint16  `json:"asa"` //最大疲劳
		NowStamina    uint16  `json:"ns"`  //当前疲劳
		AllEnergy     float64 `json:"ae"`  //最大能量（pk消耗）
		NowEnergy     float64 `json:"ne"`  //当前能量（pk消耗）
		Power         uint32  `json:"po"`  //最大战斗力
		AttackNumber  uint16  `json:"an"`  //总pk次数
		DefenseNumber uint16  `json:"dn"`  //总被他人pk次数
		Reward        int     `json:"r"`   //奖励
		CardRanks     string  `json:"cr"`  //武将进阶所需经验
		Shop          string  `json:"sp"`  //每日商店
		RliveCost     string  `json:"rc"`  //复活并秒杀敌军消耗金币数
	}
	//获得版本信息
	version := &models.Version{}
	version.GetVersion()
	//计算是否获得奖励
	//初始化返回信息
	r := loginInfo{this.member.Name, this.member.Team, this.member.SwitchTeam, this.member.BackPack, this.member.Gold, this.member.MaxLevel, this.member.Friend, this.member.FriendRequest, version.CardVersion.Hex(), version.LevelVersion.Hex(), version.ImageVersion.Hex(), this.member.AllStamina, this.member.NowStamina, this.member.AllEnergy, this.member.NowEnergy, this.member.Power, this.member.AttackNumber, this.member.DefenseNumber, -1, cardRanks, this.member.Shop, strconv.Itoa(RliveCost)}
	this.Data["json"] = r
	this.ServeJson()
}

/* 处理用户特殊操作的枚举方法 */
func (this *IndexController) Special() {
	if !this.isOk {
		this.StopRun()
	}
	//获取操作类型
	switch this.GetString("type") {
	case "updateTeam": //更新卡组队伍配置
		this.member.Team = this.GetString("team")
		this.Data["json"] = true
	case "randCard": //招贤（随机获取一张卡牌）
		if this.member.Gold < 100 { //金币不足退出
			this.StopRun()
		} else {
			this.member.Gold -= 100 //扣除金币
			//30概率小兵 30概率普通 30概率虎将 7概率神将 3概率传奇
			randNumber := newRand.Int31n(1000)
			//定义抽到卡牌的id
			var selectId int
			//判断落在哪个区间
			if randNumber < 300 { //小兵
				randNumber = newRand.Int31n(1000)
				if randNumber < 700 { //初阶
					selectId = xiaobing1[newRand.Intn(len(xiaobing1))]
				} else if randNumber < 900 { //进阶
					selectId = xiaobing2[newRand.Intn(len(xiaobing2))]
				} else { //高阶
					selectId = xiaobing3[newRand.Intn(len(xiaobing3))]
				}
			} else if randNumber < 600 { //普通
				selectId = putong[newRand.Intn(len(putong))]
			} else if randNumber < 900 { //虎将
				selectId = hujiang[newRand.Intn(len(hujiang))]
			} else if randNumber < 970 { //神将
				selectId = shenjiang[newRand.Intn(len(shenjiang))]
			} else { //传奇
				selectId = chuanqi[newRand.Intn(len(chuanqi))]
			}
			this.member.AddCard(selectId)
			this.Data["json"] = []string{strconv.Itoa(selectId), this.member.BackPack}
		}
	case "findOpponent": //匹配相近战斗力的玩家
		//计算当前pk消耗的活力值
		cost := this.member.AttackNumber/100*5 + 3
		if this.member.NowEnergy < float64(cost) { //如果活力值不够立即停止执行
			this.StopRun()
		} else {
			this.member.NowEnergy -= float64(cost)
		}
		//寻找对手
		opponent := this.member.FindOpponent()
		if opponent == nil { //如果没有可以作战的目标
			this.Data["json"] = []string{"三国军争", "57-0-0-0-0&190-0-0-0-0"}
		} else {
			//记录当前对手id
			this.member.PkOpponentId = opponent.Id.Hex()
			this.Data["json"] = []string{opponent.Name, opponent.PkTeam}
		}
	case "startPk": //开始和某位玩家pk
		//验证
		valid := validation.Validation{}
		valid.Required(this.GetString("teamId"), "teamId")
		if valid.HasErrors() { //没有通过验证则退出
			return
		}
		//获得当前出场卡牌位置数组
		teamIdArray := strings.Split(this.GetString("teamId"), ",")
		//获得当前选择卡组的卡牌数组
		teamArray := strings.Split(strings.Split(strings.Split(this.member.Team, "&")[this.member.SwitchTeam], "-")[1], ",")
		//解析背包卡牌为数组
		backPackArray := strings.Split(this.member.BackPack, ";")
		//我方出场卡牌ID
		myId := make([]string, len(teamIdArray))
		//我方即将保存的卡牌ID列表
		var myIdString string
		for k, _ := range teamIdArray {
			//获取出战队伍位置
			teamposition, _ := strconv.Atoi(teamIdArray[k])
			//获取出场卡牌在背包中位置
			position, _ := strconv.Atoi(teamArray[teamposition])
			//获取卡牌真实id
			cardInfo := strings.Split(backPackArray[position], ":")[1]
			if myIdString == "" {
				myIdString = cardInfo
			} else {
				myIdString += "&" + cardInfo
			}
			myId[k] = cardInfo
		}
		//查询对手信息
		oe := this.member.FindAccount(this.member.PkOpponentId)
		//获取对方出场id
		enemieId := strings.Split(oe.PkTeam, "&")
		//模拟战斗
		result := this.SimulateFight(myId, "0;0", enemieId, "0,0")
		if strings.Split(result, "*")[3] == "1" { //在pk场中获胜
			//战斗胜场增加
			this.member.AttackNumber++
			//更新获胜队伍列表
			this.member.PkTeam = myIdString
			//计算随机数
			randPower := newRand.Int31n(5)
			randGold := newRand.Int31n(5)
			//增加总战斗力
			this.member.Power += uint32(randPower + 20)
			//随机获得金币
			this.member.Gold += int(randGold + 20)
			//如果战斗胜场是10的倍数，增加3点活力值
			if this.member.AttackNumber%10 == 0 {
				this.member.AllEnergy += 3
			}
			this.Data["json"] = []string{result, strconv.Itoa(int(randPower + 20)), strconv.Itoa(int(randGold + 20))}
		} else {
			//计算随机数
			randNumber := newRand.Int31n(3)
			//随机获得金币
			this.member.Gold += int(randNumber + 1)
			this.Data["json"] = []string{result, strconv.Itoa(int(randNumber + 1)), "0"}
		}
	case "getGold": //获取当前金币数
		this.Data["json"] = []string{strconv.Itoa(this.member.Gold)}
	case "updateCard": //卡牌进阶
		uniqueId := this.GetString("uniqueId")
		//拆分背包
		teamArray := strings.Split(this.member.BackPack, ";")
		for k, _ := range teamArray {
			detail := strings.Split(teamArray[k], ":")
			if detail[0] == uniqueId { //操作此卡牌
				info := strings.Split(detail[1], ",")
				nowRank, _ := strconv.Atoi(info[6])       //当前等级
				nowExperience, _ := strconv.Atoi(info[5]) //当前经验值
				//拆分出经验数组
				cardRanksArray := StringToIntArray(strings.Split(cardRanks, ","))
				if nowRank+1 > len(cardRanksArray) { //超过最大等级
					this.StopRun()
				}
				//fmt.Println(nowExperience, cardRanksArray[nowRank])
				if nowExperience < cardRanksArray[nowRank] { //经验值不足
					this.StopRun()
				}
				//卡牌进阶
				nowRank++
				nowExperience = 0
				info[6] = strconv.Itoa(nowRank)
				info[5] = strconv.Itoa(nowExperience)
				//增加属性
				hp, _ := strconv.Atoi(info[1])
				speed, _ := strconv.Atoi(info[2])
				attack, _ := strconv.Atoi(info[3])
				defense, _ := strconv.Atoi(info[4])
				addhp, addspeed, addattack, adddefense := this.RandBase(nowRank)
				hp += addhp
				speed += addspeed
				attack += addattack
				defense += adddefense
				info[1] = strconv.Itoa(hp)
				info[2] = strconv.Itoa(speed)
				info[3] = strconv.Itoa(attack)
				info[4] = strconv.Itoa(defense)
				//重组detail[1]
				detail[1] = strings.Join(info, ",")
				//重组teamArray
				teamArray[k] = detail[0] + ":" + detail[1]
				this.Data["json"] = []string{info[1], info[2], info[3], info[4]}
				break
			}
		}
		//重组背包
		this.member.BackPack = strings.Join(teamArray, ";")
	case "buyShopCard": //购买每日商店的卡牌
		position, _ := strconv.Atoi(this.GetString("position"))
		shopArray := strings.Split(this.member.Shop, "&")
		detail := strings.Split(shopArray[position], ",")
		id, _ := strconv.Atoi(detail[0])   //卡牌id
		cost, _ := strconv.Atoi(detail[1]) //价格
		//购买此卡牌
		if this.member.Gold >= cost { //金币足够
			this.member.Gold -= cost //扣除购买金币
			this.member.AddCard(id)  //新增卡牌
			//剔除商店中的此位置
			shopArrayBefore := shopArray[:position]
			shopArrayAfter := shopArray[position+1:]
			shopArray = append(shopArrayBefore, shopArrayAfter...)
			this.member.Shop = strings.Join(shopArray, "&")
		}
		this.Data["json"] = []string{this.member.BackPack}
	case "refreshShop": //刷新每日商店
		if this.member.Gold > 30 {
			this.RefreshShop()
			//扣除金币
			this.member.Gold -= 30
			this.Data["json"] = []string{this.member.Shop}
		}
	case "backToLife": //战斗复活并秒杀敌人
		if this.member.Gold >= RliveCost {
			//消耗金币
			this.member.Gold -= RliveCost
			//战斗过程状态从失败变成成功
			info := strings.Split(this.member.LastGet, "&")
			info[6] = "1"
			//战斗步骤自增1
			turn, _ := strconv.Atoi(info[1])
			turn++
			//必须不是最后一关
			if turn >= len(strings.Split(info[0], "|")) {
				this.StopRun()
			}
			info[1] = strconv.Itoa(turn)
			//重组字符串
			this.member.LastGet = strings.Join(info, "&")
		}
	case "getCoins": //购买金币

	}
	this.ServeJson()
}

/* 用户开始游戏 */
func (this *IndexController) StartGame() {
	if !this.isOk {
		this.StopRun()
	}
	//验证
	postLevel, _ := this.GetInt("level")
	postSwitchTeam, _ := this.GetInt("switchTeam")
	valid := validation.Validation{}
	valid.Required(postLevel, "level")
	valid.Required(postSwitchTeam, "switchteam")
	if valid.HasErrors() { //没有通过验证则退出
		return
	}
	this.member.SwitchTeam = int(postSwitchTeam)
	//获取当前关卡信息
	level := models.Level{}
	ok := level.FindOne(int(postLevel))
	if ok == false { //关卡不存在，则退出
		return
	}
	//获得每一轮敌人信息
	eachLevel := strings.Split(level.Information, "|")
	//判断疲劳是否足够
	if this.member.NowStamina < uint16(len(eachLevel)) { //疲劳不够，停止运行
		this.StopRun()
	} else {
		this.member.NowStamina -= uint16(len(eachLevel))
	}
	//获得经验
	var experience uint16
	//敌人信息
	enemieInfo := make([]string, len(eachLevel))
	for k, v := range eachLevel {
		if v == "" {
			continue
		}
		detail := strings.Split(v, ":")
		//确定敌人数量
		var count int
		if strings.Contains(detail[0], "~") { //敌人数是个范围
			countArray := strings.Split(detail[0], "~")
			s, _ := strconv.Atoi(countArray[0])
			e, _ := strconv.Atoi(countArray[1])
			count = newRand.Intn(e-s+1) + s
		} else {
			c, _ := strconv.Atoi(detail[0])
			count = c
		}
		//获取敌人备选队列，eg []array{'181-0','416-0'}
		descrip := strings.Split(detail[1], ",")
		//从敌人备选队列中选择选定的数量
		sel := len(descrip) - count
		if sel == 0 {
			sel = 1
		}
		randStart := newRand.Intn(sel) //获取随机开始位置
		enemieArray := descrip[randStart : randStart+count]
		//判断敌人给的经验
		for _, enemie := range enemieArray {
			if enemie == "" {
				continue
			}
			//这张敌人卡牌的具体信息
			enemieDetail := strings.Split(enemie, "-")
			id, _ := strconv.Atoi(enemieDetail[0])
			//判断给的经验
			quality, _ := strconv.Atoi(enemieDetail[1])
			if id <= 1000 { //武将
				switch quality {
				case 0:
					experience += 10
				case 1:
					experience += 20
				case 2:
					experience += 25
				case 3:
					experience += 45
				case 4:
					experience += 70
				}
			} else { //小兵
				switch quality {
				case 0:
					experience += 5
				case 1:
					experience += 10
				case 2:
					experience += 15
				case 3:
					experience += 25
				case 4:
					experience += 45
				}
			}
			var hp, speed, attack, defense int
			//根据敌人阶级给相应增益属性
			switch quality {
			case 1:
				hp, speed, attack, defense = this.RandBase(1)
			case 2:
				hp, speed, attack, defense = this.RandBase(1)
				hp1, speed1, attack1, defense1 := this.RandBase(2)
				hp += hp1
				speed += speed1
				attack += attack1
				defense += defense1
			case 3:
				hp, speed, attack, defense = this.RandBase(1)
				hp1, speed1, attack1, defense1 := this.RandBase(2)
				hp2, speed2, attack2, defense2 := this.RandBase(3)
				hp += (hp1 + hp2)
				speed += (speed1 + speed2)
				attack += (attack1 + attack2)
				defense += (defense1 + defense2)
			case 4:
				hp, speed, attack, defense = this.RandBase(1)
				hp1, speed1, attack1, defense1 := this.RandBase(2)
				hp2, speed2, attack2, defense2 := this.RandBase(3)
				hp3, speed3, attack3, defense3 := this.RandBase(4)
				hp += (hp1 + hp2 + hp3)
				speed += (speed1 + speed2 + speed3)
				attack += (attack1 + attack2 + attack3)
				defense += (defense1 + defense2 + defense3)
			}
			//将敌人加入敌人队列中
			if enemieInfo[k] == "" {
				enemieInfo[k] = strconv.Itoa(id) + "," + strconv.Itoa(hp) + "," + strconv.Itoa(speed) + "," + strconv.Itoa(attack) + "," + strconv.Itoa(defense)
			} else {
				enemieInfo[k] += ";" + strconv.Itoa(id) + "," + strconv.Itoa(hp) + "," + strconv.Itoa(speed) + "," + strconv.Itoa(attack) + "," + strconv.Itoa(defense)
			}
		}
	}
	//判断是可以否会开启下一关
	levelCount := uint16(level.GetCount())
	//总闯关关卡小于总关卡数且当前关卡等于总关卡数（最新关）
	canOpenNext := "0"
	if this.member.MaxLevel < levelCount && uint16(postLevel) == this.member.MaxLevel {
		canOpenNext = "1"
	}
	//随机奖励的金币
	randGold := newRand.Intn(4) + 20
	//初始化我方卡牌出场次数
	actionTime := make([]string, len(strings.Split(strings.Split(strings.Split(this.member.Team, "&")[this.member.SwitchTeam], "-")[1], ",")))
	for k, _ := range actionTime {
		actionTime[k] = "0"
	}
	//保存到lastget 敌人信息 & 从第一关(0)开始打 & 是否开启下关 & 奖励金币 & 奖励经验 & 各卡牌出场次数 & 战斗是否结束了(0:结束)
	this.member.LastGet = strings.Join(enemieInfo, "|") + "&0" + "&" + canOpenNext + "&" + strconv.Itoa(randGold) + "&" + strconv.Itoa(int(experience)) + "&" + strings.Join(actionTime, ",") + "&1"
	//查询该关卡的父级关卡珠子类型
	starArray := level.FindParentStarArray()
	//返回数据
	this.Data["json"] = []string{this.member.LastGet, starArray}
	this.ServeJson()
}

/* 根据敌人阶级给相应增益属性 返回值依次是 体力 行动力 攻击力 防御力 */
func (this *IndexController) RandBase(rand int) (int, int, int, int) {
	var hp, speed, attack, defense int
	switch rand {
	case 1:
		hp = newRand.Intn(6) + 1
		speed = newRand.Intn(4) + 1
		attack = newRand.Intn(2)
		defense = newRand.Intn(2)
	case 2:
		hp = newRand.Intn(10) + 5
		speed = newRand.Intn(5) + 5
		attack = newRand.Intn(2) + 1
		defense = newRand.Intn(2) + 1
	case 3:
		hp = newRand.Intn(15) + 10
		speed = newRand.Intn(10) + 5
		attack = newRand.Intn(4) + 1
		defense = newRand.Intn(3) + 1
	case 4:
		hp = newRand.Intn(15) + 20
		speed = newRand.Intn(10) + 10
		attack = newRand.Intn(8) + 2
		defense = newRand.Intn(6) + 2
	}
	//血量最高增量：82
	return hp, speed, attack, defense
}

/* 刷新购买商店列表 */
func (this *IndexController) RefreshShop() {
	//传奇：1 神将：1 虎将：5 普通：1(10) 小兵：2
	//初始化商店数组
	shopArray := make([]string, 10)
	shopArray[0] = strconv.Itoa(chuanqi[newRand.Intn(len(chuanqi))]) + "," + strconv.Itoa(1900+newRand.Intn(200))
	shopArray[1] = strconv.Itoa(shenjiang[newRand.Intn(len(shenjiang))]) + "," + strconv.Itoa(950+newRand.Intn(100))
	for i := 2; i <= 6; i++ {
		shopArray[i] = strconv.Itoa(hujiang[newRand.Intn(len(hujiang))]) + "," + strconv.Itoa(200+newRand.Intn(50))
	}
	for i := 7; i <= 7; i++ {
		shopArray[i] = strconv.Itoa(putong[newRand.Intn(len(putong))]) + "," + strconv.Itoa(120+newRand.Intn(30))
	}
	for i := 8; i <= 9; i++ {
		randNumber := newRand.Int31n(1000)
		if randNumber < 700 { //初阶
			shopArray[i] = strconv.Itoa(xiaobing1[newRand.Intn(len(xiaobing1))]) + "," + strconv.Itoa(100+newRand.Intn(5))
		} else if randNumber < 900 { //进阶
			shopArray[i] = strconv.Itoa(xiaobing2[newRand.Intn(len(xiaobing2))]) + "," + strconv.Itoa(200+newRand.Intn(10))
		} else { //高阶
			shopArray[i] = strconv.Itoa(xiaobing3[newRand.Intn(len(xiaobing3))]) + "," + strconv.Itoa(500+newRand.Intn(100))
		}
	}
	this.member.Shop = strings.Join(shopArray, "&")
}

/* 开始剧情战斗 */
func (this *IndexController) StartFight() {
	if !this.isOk {
		this.StopRun()
	}
	//验证
	valid := validation.Validation{}
	valid.Required(this.GetString("teamId"), "teamId")
	if valid.HasErrors() { //没有通过验证则退出
		return
	}
	//获得当前出场卡牌位置数组
	teamIdArray := strings.Split(this.GetString("teamId"), ",")
	//获得当前选择卡组的卡牌数组
	teamArray := strings.Split(strings.Split(strings.Split(this.member.Team, "&")[this.member.SwitchTeam], "-")[1], ",")
	//解析背包卡牌为数组
	backPackArray := strings.Split(this.member.BackPack, ";")
	//我方出场卡牌ID
	myId := make([]string, len(teamIdArray))
	for k, _ := range teamIdArray {
		//获取出战队伍位置
		teamposition, _ := strconv.Atoi(teamIdArray[k])
		//获取出场卡牌在背包中位置
		position, _ := strconv.Atoi(teamArray[teamposition])
		//获取卡牌真实id
		myId[k] = strings.Split(backPackArray[position], ":")[1]
	}
	info := strings.Split(this.member.LastGet, "&")
	//如果有结束标识，立刻退出
	if info[6] == "0" {
		this.StopRun()
	}
	//根据lastget获取战斗回合数
	turn, _ := strconv.Atoi(info[1])
	//根据lastget获取此轮敌方卡牌数组
	enemieId := strings.Split(strings.Split(strings.Split(this.member.LastGet, "&")[0], "|")[turn], ";")
	//模拟战斗
	result := this.SimulateFight(myId, "0;0", enemieId, "0;0")
	//如果赢了，保存我方战斗步骤
	if strings.Split(result, "*")[3] == "1" {
		//步骤+1
		turn++
		info[1] = strconv.Itoa(turn)
	} else { //输了则设置结束标识
		info[6] = "0"
		//奖励0金币和10%经验
	}
	//增加位置卡牌出场次数
	eachActionTime := strings.Split(info[5], ",")
	for _, v := range teamIdArray {
		position, _ := strconv.Atoi(v)
		intTime, _ := strconv.Atoi(eachActionTime[position])
		intTime++
		eachActionTime[position] = strconv.Itoa(intTime)
	}
	//重组出场次数
	info[5] = strings.Join(eachActionTime, ",")
	//如果战斗全部结束并且胜利，累积奖励
	if turn == len(strings.Split(strings.Split(this.member.LastGet, "&")[0], "|")) {
		//如果可以开启下一关，说明是第一次通关
		if info[2] == "1" {
			this.member.MaxLevel++
			//增加1点疲劳值
			this.member.AllStamina += 1
			this.member.NowStamina += 1
		}
		//获得金币
		gold, _ := strconv.Atoi(info[3])
		this.member.Gold += gold
		//武将根据出场次数分配经验
		experience, _ := strconv.Atoi(info[4])
		actionArray := strings.Split(info[5], ",") //获取各个位置武将出场次数
		//计算总出场次数
		allAction := 0
		for _, v := range actionArray {
			number, _ := strconv.Atoi(v)
			allAction += number
		}
		//拆分出经验数组
		cardRanksArray := StringToIntArray(strings.Split(cardRanks, ","))
		//分别保存各个卡牌的经验
		for k, v := range actionArray {
			if v == "0" {
				continue //出场次数为0则跳过
			}
			actionTime, _ := strconv.Atoi(v) //出场次数
			//获取当前背包中卡片信息 4经验 5等级
			position, _ := strconv.Atoi(teamArray[k])
			detail := strings.Split(backPackArray[position], ":")
			cardInfo := strings.Split(detail[1], ",")
			//如果卡牌是满级，不累计经验
			rank, _ := strconv.Atoi(cardInfo[6]) //当前等级
			exp, _ := strconv.Atoi(cardInfo[5])  //当前经验
			addExp := int(math.Ceil(float64(actionTime) / float64(allAction) * float64(experience)))
			if rank < len(cardRanksArray) { //还没有满级
				if exp+addExp > int(cardRanksArray[rank]) { //超过最大经验
					exp = int(cardRanksArray[rank])
				} else {
					exp += addExp
				}
			}
			cardInfo[5] = strconv.Itoa(exp)
			//保存到此卡牌信息
			backPackArray[position] = detail[0] + ":" + strings.Join(cardInfo, ",")
		}
		//保存背包卡牌信息
		this.member.BackPack = strings.Join(backPackArray, ";")
		//设置结束标识
		info[6] = "0"
	}
	//重组lastget并保存
	this.member.LastGet = strings.Join(info, "&")
	//输出结果
	this.Ctx.WriteString(result)
}

/* 模拟一场战斗 武将：id,体力,行动力,武技,防御力 军队：id,数量;id,数量 */
func (this *IndexController) SimulateFight(my []string, myArmy string, enemie []string, enemieArmy string) string {
	//初始化游戏总控器
	main := &logic.Maingame{}
	main.MyTeam = make([]*logic.Ai, len(my))
	main.EnemieTeam = make([]*logic.Ai, len(enemie))
	main.Time = 80
	main.Rands = newRand
	//武将不能超过5个
	if len(my) > 5 {
		this.StopRun()
	}
	if len(enemie) > 5 {
		this.StopRun()
	}
	//一维数组转化为二维信息数组
	myArray := make([][]int, len(my))
	for k, _ := range myArray {
		myArray[k] = StringToIntArray(strings.Split(my[k], ","))
	}
	enemieArray := make([][]int, len(enemie))
	for k, _ := range enemieArray {
		enemieArray[k] = StringToIntArray(strings.Split(enemie[k], ","))
	}
	//判断是否有重复
	for k, v := range myArray {
		if v[0] > 1000 { //小兵可以重复
			continue
		}
		for key, val := range myArray {
			if key != k && val[0] == v[0] {
				this.StopRun() //重复
			}
		}
	}
	for k, v := range enemieArray {
		if v[0] > 1000 { //小兵可以重复
			continue
		}
		for key, val := range enemieArray {
			if key != k && val[0] == v[0] {
				this.StopRun() //重复
			}
		}
	}
	//初始化我方卡牌信息
	for k, v := range myArray {
		//实例化卡牌，让指针指向它
		main.MyTeam[k] = &logic.Ai{}
		main.MyTeam[k].Base = cards[uint16(v[0])]
		main.MyTeam[k].Create(v, allskill)
		//赋值全局控制指针
		main.MyTeam[k].Main = main
		//行动时间初始化
		main.MyTeam[k].AlreadyRun = main.Time
		main.MyTeam[k].PositionId = uint8(k)
		main.MyTeam[k].IsMy = true
	}
	//初始化敌方卡牌信息
	for k, v := range enemieArray {
		//实例化卡牌，让指针指向它`
		main.EnemieTeam[k] = &logic.Ai{}
		main.EnemieTeam[k].Base = cards[uint16(v[0])]
		main.EnemieTeam[k].Create(v, allskill)
		//赋值全局控制指针
		main.EnemieTeam[k].Main = main
		//行动时间初始化
		main.EnemieTeam[k].AlreadyRun = main.Time
		main.EnemieTeam[k].PositionId = uint8(k + len((my)))
		main.EnemieTeam[k].IsMy = false
	}
	//战斗状态 -1失败 0表示没结束 1胜利
	status := int8(0)
	//组合我方和地方单位指针为一个大的slice
	all := append(main.MyTeam, main.EnemieTeam...)
	//每个卡牌释放关卡开始技能
	for _, v := range all {
		main.AllStep += strconv.Itoa(int(v.PositionId)) + "&"
		v.StartSkill()
		main.AllStep += "#"
	}
	main.AllStep += "*"
	//开始循环战斗，最多循环100次
	for i := 0; i < 100; i++ {
		//定义最短时间(100秒),当前行动单位
		shortime := float32(100)
		var runAi *logic.Ai
		//寻找用时最短的
		for _, v := range all {
			//若已死亡则不在考虑范围
			if v.Alive == false {
				continue
			}
			//行动用时等于 已经运动数量/行动速度+晕眩时间
			usetime := v.AlreadyRun/float32(v.Base.AttackSpeed) + v.DizzyLeft
			if usetime < shortime {
				shortime = usetime
				runAi = v
			}
		}
		// 开始行动
		for _, v := range all {
			//若已死亡则不在考虑范围
			if v.Alive == false {
				continue
			}
			if v == runAi { //单位开始行动
				main.AllStep += strconv.Itoa(int(v.PositionId)) + "&" //开始
				v.StartAction()                                       //行动
				//让所有HP为0的活着的单位死亡
				for _, a := range all {
					if a.Alive == true && a.Base.Hp == 0 {
						a.Die()
					}
				}
				main.AllStep += "#" //结束
				//行动结束后，重置行动时间
				v.AlreadyRun = main.Time
				//晕眩时间重置为0
				v.DizzyLeft = 0
			} else { //单位行动剩余时间减少
				if v.DizzyLeft > 0 {
					//有晕眩时间优先抵消晕眩时间
					if v.DizzyLeft > shortime {
						//晕眩时间还未结束
						v.DizzyLeft -= shortime
					} else {
						//保存抵消过晕眩时间后的剩余时间
						backShortTime := shortime - v.DizzyLeft
						//行动剩余时间减少
						v.AlreadyRun -= backShortTime * float32(v.Base.AttackSpeed)
						//晕眩时间为0
						v.DizzyLeft = 0
					}
				} else {
					//没有晕眩时间，直接减少剩余行动时间(一定会有剩余或者等于0，浮点计算会有小误差)
					v.AlreadyRun -= shortime * float32(v.Base.AttackSpeed)
				}
			}
		}
		//判断是否败走
		myAlive := 0
		for _, v := range main.MyTeam {
			if v.Alive {
				myAlive++
			}
		}
		if myAlive == 0 {
			//我方全军覆没
			status = -1
			goto Next
		}
		//判断是否胜利
		enemieAlive := 0
		for _, v := range main.EnemieTeam {
			if v.Alive {
				enemieAlive++
			}
		}
		if enemieAlive == 0 {
			//敌方全军覆没
			status = 1
			goto Next
		}
	}
Next:
	//计算总评价
	main.AllStep += "*"
	//记录所有评价
	for _, v := range main.MyTeam {
		main.AllStep += strconv.Itoa(int(v.AllAttack)) + "," + strconv.Itoa(int(v.AllDamage)) + "," + strconv.Itoa(int(v.AllAdd)) + "," + strconv.Itoa(int(v.AllRunTime)) + "," + strconv.Itoa(int(v.AllShowTime)) + "," + strconv.Itoa(int(v.AllKill)) + ";"
	}
	//记录结果状态
	main.AllStep += "*" + strconv.Itoa(int(status))
	//返回过程
	return main.AllStep
}

//下载资料
func (this *IndexController) Download() {
	if !this.isOk {
		this.StopRun()
	}
	switch this.GetString("path") {
	case "card": //下载卡牌信息
		card := &models.Card{}
		allCard := card.DownLoad()
		this.Data["json"] = allCard
	case "level": //下载关卡信息
		level := &models.Level{}
		allLevel := level.DownLoad()
		this.Data["json"] = allLevel
	case "image": //下载图片信息
		this.Data["json"] = images
	default:
	}
	this.ServeJson()
}

//返回可以下载的token
func (this *IndexController) DownloadImage() {
	if !this.isOk {
		this.StopRun()
	}
	cards := strings.Split(this.GetString("cards"), ",")
	urls := make([]string, len(cards))
	for k, _ := range cards {
		//生成token
		path := "/download/" + cards[k]
		etime := strconv.Itoa(int(time.Now().Unix()) + 300*(k+1))
		t := md5.New()
		t.Write([]byte(beego.AppConfig.String("uKingdomsToken") + "&" + etime + "&" + path))
		sign := hex.EncodeToString(t.Sum(nil))
		urls[k] = "http://kingdoms.wokugame.com/download/" + cards[k] + "?_upt=" + sign[12:20] + etime
	}
	this.Ctx.WriteString(strings.Join(urls, ","))
}

//更新图片信息
func updateImage() {
	/// 读取目录
	dirs, err := ku.ReadDir("/download")
	if err != nil {
		log.Error("更新图片信息出错：", err)
	}
	//重定义信息json大小
	images = make(map[string]string, len(dirs))
	for _, d := range dirs {
		//设置images
		images[d.Name] = strconv.Itoa(int(d.Time)) + "," + strconv.Itoa(int(d.Size))
	}
}

/* 付款成功接收post */
func (this *IndexController) PaySuccess() {
	log.Info("异步接收", this.Input().Encode())
}

func (this *IndexController) Finish() {
	//保存更新后数据
	this.member.Update()
}
