package logic

import (
	"container/list"
	"kingdoms/models"
	"strconv"
	"strings"
)

/* BUFF结构体 */
type BuffStruct struct {
	Effect map[uint8][2]uint8 //buff种类,数值,增益或者减益
	Last   int8               //buff剩余回合数 -1表示永久有效
}

/* 技能释放时机 */
const (
	Start               = "1"  //战斗开始前√
	ActionBefore        = "2"  //行动前√
	ActionAfter         = "3"  //行动后√
	AttackDelegate      = "4"  //攻击前√
	AttackHitDelegate   = "5"  //攻击命中√
	GetHitFrontDelegate = "6"  //受到攻击前√
	GetHitAfterDelegate = "7"  //受到攻击后√
	GetCritDelegate     = "8"  //受到暴击√
	MissDelegate        = "9"  //闪避后√
	DieFrontDelegate    = "10" //死亡前√
	DieAfterDelegate    = "11" //死亡后√
	AttackToDieDelegate = "12" //斩杀√
	AddHpFrontDelegate  = "13" //回血前√
	AddHpAfterDelegate  = "14" //回血后√
	ReLiveDelegate      = "15" //复活√
	RemoveTauntDelegate = "15" //移除嘲讽委托
	AfterRunSkill       = "16" //触发技能后√
	BeforeRunSkill      = "17" //触发技能前√
)

/* 技能结构体 */
type SkillStruct struct {
	Where  string  //技能释放时机 含有-表示释放技能 含有:表示初始属性 统一用;分隔
	Name   string  //技能名称
	Params []uint8 //参数
	player *Ai     //技能使用者
}

type Ai struct {
	Main          *Maingame              //游戏总控器地址
	AllSkill      *AllSkill              //技能释放器地址
	Base          models.Card            //卡牌信息
	Buff          map[string]*BuffStruct //BUFF信息
	Skill         list.List              //技能链表
	SkillTemp     map[string]int8        //技能状态临时变量
	SkillTempAi   map[string]*Ai         //技能状态临时ai指针
	PositionId    uint8                  //卡牌在队列中的位置
	LimitHP       uint8                  //血量上限
	MaxHP         uint8                  //当前最大血量
	AttackType    uint8                  //攻击类型
	Crit          uint8                  //暴击率
	CritDamage    uint8                  //暴击伤害
	HitRatio      uint8                  //命中率
	Miss          uint8                  //闪避率 1 = 1%
	MissType      uint8                  //闪避类型
	Dizzy         uint8                  //重击（晕眩）率 1 = 1%
	DizzyTime     float32                //重击持续时间
	Sputter       uint8                  //溅射率 1 = 1%
	SputterDamage uint8                  //溅射伤害 1 = 1%
	Volley        uint8                  //乱射率（无视嘲讽） 1 = 1%
	Reflect       uint8                  //反伤率 1 = 1%
	ReflectType   uint8                  //反射种类(0:百分比，1:固定)
	ReflectDamage uint8                  //反伤伤害 1 = 1% / 1
	IsMy          bool                   //是否为我方单位
	Taunt         bool                   //是否拥有嘲讽
	DizzyLeft     float32                //晕眩状态剩余时间（被重击）
	Hidden        bool                   //是否为埋伏类型
	Hide          bool                   //是否处于埋伏状态
	Chaos         uint8                  //混乱状态回合数 攻击己方单位
	CostStar      uint8                  //每种属性消耗个数
	ActionNumber  uint8                  //行动次数
	HuDun         uint8                  //护盾个数 免疫伤害次数
	SkillDamage   uint8                  //技能伤害
	ExtraAddHP    uint8                  //回复效果
	Focus         uint8                  //被集火剩余数
	Alive         bool                   //是否存活

	/* 技能信息 */

	/* 攻击前存储的信息 */
	Adamage     uint8   //攻击伤害
	Acritdamage uint8   //暴击比率
	Adizzytime  float32 //重击晕眩时间
	Aissputter  bool    //是否溅射
	Ahitratio   uint8   //命中修正率
	Atarget     *Ai     //此次攻击目标

	/* 受到攻击存储的信息 */
	Rdamage        uint8   //受到的伤害
	Rcritdamage    uint8   //暴击比率
	Rdizzytime     float32 //重击持续时间
	Rattacker      *Ai     //攻击者指针
	Rhitratio      uint8   //被命中修正率
	Rcanmiss       bool    //是否能被闪避
	Rcanreduce     bool    //是否能被防御抵消
	Rcanreflect    bool    //是否能反伤
	Rreflectdamage uint8   //当前反戈系数
	Rreflectbase   uint8   //当前反伤的基础伤害
	Rispoint       bool    //是否为指向性单体攻击
	Rtype          uint8   //攻击类型 0:普通攻击 1:智击 2:技能攻击

	/* 受到回复效果存储的信息 */
	Hadder *Ai   //为自己加血的目标
	Hvalue uint8 //加血量

	/* 总的统计信息 */
	AllAttack   uint16 //统计总伤害
	AllDamage   uint16 //统计总受到伤害
	AllAdd      uint16 //统计总治疗量
	AllRunTime  uint8  //统计总行动次数
	AllShowTime uint8  //统计总出牌次数
	AllKill     uint8  //统计总杀敌数

	/* 当前行动 */
	AlreadyRun float32 //已行动时间
}

/* 初始化创建单位 */
func (this *Ai) Create(info []int, allSkill *AllSkill) {
	//赋值技能释放器指针
	this.AllSkill = allSkill
	//设置额外属性
	this.Base.Hp += uint8(info[1])
	this.Base.AttackSpeed += uint8(info[2])
	this.Base.Attack += uint8(info[3])
	this.Base.Defense += uint8(info[4])
	//初始状态为存活
	this.Alive = true
	//初始化buff map
	this.Buff = make(map[string]*BuffStruct)
	//初始化skilltemp map
	this.SkillTemp = make(map[string]int8)
	//初始化skilltempai map
	this.SkillTempAi = make(map[string]*Ai)
	//初始化一些基础属性
	this.LimitHP = this.Base.Hp
	this.MaxHP = this.Base.Hp
	//记录出场次数
	this.AllShowTime = 1
	//解析技能
	skillArray := strings.Split(this.Base.SkillDetail, ";")
	for _, v := range skillArray {
		if strings.Index(v, ":") > 0 {
			//基础状态
			statusDetail := strings.Split(v, ":")
			switch statusDetail[0] {
			case "aty": //攻击类型
				ty, _ := strconv.Atoi(statusDetail[1])
				this.AttackType = uint8(ty)
			case "crit": //暴击
				t1, _ := strconv.Atoi(statusDetail[1])
				t2, _ := strconv.Atoi(statusDetail[2])
				this.Crit = uint8(t1)
				this.CritDamage = uint8(t2)
			case "dizzy": //重击
				t1, _ := strconv.Atoi(statusDetail[1])
				t2, _ := strconv.ParseFloat(statusDetail[2], 10)
				this.Dizzy = uint8(t1)
				this.DizzyTime = float32(t2)
			case "miss": //闪避
				m, _ := strconv.Atoi(statusDetail[1])
				this.Miss = uint8(m)
			case "volley": //精准（乱射，优先攻击非坚守目标）
				v, _ := strconv.Atoi(statusDetail[1])
				this.Volley = uint8(v)
			case "reflect": //反戈（反伤）
				r1, _ := strconv.Atoi(statusDetail[1])
				r2, _ := strconv.Atoi(statusDetail[2])
				r3, _ := strconv.Atoi(statusDetail[3])
				this.Reflect = uint8(r1)
				this.ReflectType = uint8(r2)
				this.ReflectDamage = uint8(r3)
			case "hidden": //埋伏
				this.Hidden = true
			case "hitratio": //命中率
				ratio, _ := strconv.Atoi(statusDetail[1])
				this.HitRatio = uint8(ratio)
			case "sputter": //溅射
				s1, _ := strconv.Atoi(statusDetail[1])
				s2, _ := strconv.Atoi(statusDetail[2])
				this.Sputter = uint8(s1)
				this.SputterDamage = uint8(s2)
			}
		} else if strings.Index(v, "-") > 0 {
			//添加技能
			this.AddSkill(v)
		}
	}
}

/* 释放关卡开始前技能 */
func (this *Ai) StartSkill() {
	//循环技能,只调用开场释放的技能
	for e := this.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == Start {
			this.AllSkill.RunSkill(e.Value.(*SkillStruct))
		}
	}
	extra := strings.Split(this.Base.SkillDetail, ";")
	for _, v := range extra {
		//承志
		if strings.Index(v, "+") > 0 {
			detail := strings.Split(v, "+")
			for _, v := range this.FindTeam(true) {
				if v.Base.Name == detail[0] {
					this.AddSkill(this.Base.FindSkill(v.Base.Id, detail[1]))
					this.Main.AllStep += "4:" + strconv.Itoa(int(this.PositionId)) + "," + strconv.Itoa(int(this.PositionId)) + "," + detail[1] + ",-1,0|"
				}
			}
		}
		if strings.Index(v, "*") > 0 {
			detail := strings.Split(v, "*")
			cArray := detail[1:]
			count := len(cArray) / 2
			for i := 0; i < count; i++ {
				this.Main.AllStep += "6:" + strconv.Itoa(int(this.PositionId)) + "," + cArray[i*2] + "," + cArray[i*2+1] + "|"
			}
		}
		//下回合减消耗
	}
	//有山属性则有坚守效果
	if this.Base.HasAttribute(4) {
		this.AddTaunt()
	}
}

/* 开始本单位的回合 */
func (this *Ai) StartAction() {
	//记录 统计总行动次数
	this.AllRunTime++
	//开始回合的时候，若处于埋伏状态，将会现形
	if this.Hide {
		this.RemoveHide()
	}
	//循环技能,只调用回合开始前释放的技能
	for e := this.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == ActionBefore {
			this.AllSkill.RunSkill(e.Value.(*SkillStruct))
		}
	}
	//开始攻击
	for i := uint8(0); i < this.Base.AttackNumber; i++ {
		var target []*Ai //设置一个空的当前目标
		//是否处于混乱状态
		targetMy := false
		if this.Chaos > 0 {
			targetMy = true
		}
		//是否触发乱射
		if this.Main.Rands.Int31n(100) < int32(this.Volley) {
			//触发乱射，寻找敌人非坚守目标
			target = this.FindTauntArray(targetMy, false, 1)
			if target == nil {
				//没找到，则随机选择目标
				target = this.FindRandArray(targetMy, 1)
			}
		} else {
			//没有乱射，则寻找敌人坚守目标
			target = this.FindTauntArray(targetMy, true, 1)
			if target == nil {
				//对方没有坚守单位，则随机寻找目标
				target = this.FindRandArray(targetMy, 1)
			}
		}
		if target != nil { //当前目标不为空，则进行攻击
			this.Attack(target[0])
		} else {
			//目标为空，跳过
			continue
		}
		//计算混乱时间是否过去
		if this.Chaos > 0 {
			targetMy = true
			this.Chaos--
			if this.Chaos == 0 {
				this.RemoveChaos()
			}
		}
	}
	//回合结束,考虑是否有buff要结束
	for k, _ := range this.Buff {
		if this.Buff[k].Last != -1 && this.Buff[k].Last == 1 {
			this.RemoveBuffByName(k)
		} else if this.Buff[k].Last != -1 {
			this.Buff[k].Last--
		}
	}
	//循环技能,只调用回合开始后释放的技能
	for e := this.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == ActionAfter {
			this.AllSkill.RunSkill(e.Value.(*SkillStruct))
		}
	}
}

/* 攻击某个目标 */
func (this *Ai) Attack(target *Ai) {
	this.Atarget = target
	this.Adamage = this.Base.Attack
	if this.Main.Rands.Int31n(100) < int32(this.Crit) { //是否暴击
		this.Acritdamage = this.CritDamage
	} else {
		this.Acritdamage = 0
	}
	if this.Main.Rands.Int31n(100) < int32(this.Dizzy) { //是否重击
		this.Adizzytime = this.DizzyTime
	} else {
		this.Adizzytime = 0
	}
	if this.Main.Rands.Int31n(100) < int32(this.Sputter) { //是否溅射
		this.Aissputter = true
	} else {
		this.Aissputter = false
	}
	this.Ahitratio = this.HitRatio
	//循环技能,只调用攻击后释放的技能
	for e := this.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == AttackDelegate {
			this.AllSkill.RunSkill(e.Value.(*SkillStruct))
		}
	}
	if this.Adamage == 0 {
		return
	}
	//造成伤害   攻击者 伤害 命中修正 暴击率 晕眩时间 是否能被闪避 是否能被防御抵消 是否能被反伤
	if this.AttackType == 2 {
		//智击，无视防御和闪避
		this.Atarget.ReduceHp(this, this.Adamage, this.Ahitratio, this.Acritdamage, this.Adizzytime, false, false, true, true, 1)
	} else {
		this.Atarget.ReduceHp(this, this.Adamage, this.Ahitratio, this.Acritdamage, this.Adizzytime, true, true, true, true, 0)
	}
	//如果溅射，则造成溅射伤害
	if this.Aissputter {
		team := this.FindTeam(false)
		for _, v := range team {
			if v != this.Atarget {
				damage := uint8(uint16(this.Adamage) * uint16(this.SputterDamage) / 100)
				if damage <= 0 {
					damage = 1
				}
				v.ReduceHp(this, damage, this.Ahitratio, this.Acritdamage, this.Adizzytime, false, false, true, false, 1)
			}
		}
	}
}

/* 受到伤害 */
func (this *Ai) ReduceHp(attacker *Ai, damage uint8, hitRatio uint8, critDamage uint8, dizzyTime float32, canMiss bool, canReduce bool, canReflect bool, isPoint bool, attackType uint8) {
	//存储受到攻击的数值
	this.Rattacker = attacker
	this.Rdamage = damage
	this.Rhitratio = hitRatio
	this.Rcritdamage = critDamage
	this.Rdizzytime = dizzyTime
	this.Rcanmiss = canMiss
	this.Rcanreduce = canReduce
	this.Rcanreflect = canReflect
	this.Rreflectdamage = this.ReflectDamage
	this.Rreflectdamage = this.ReflectDamage
	this.Rreflectbase = damage
	this.Rispoint = isPoint
	this.Rtype = attackType
	//属性相克 1 2 3 4 5 6
	for i := 1; i <= 6; i++ {
		if this.Base.HasAttribute(uint8(i)) {
			k := i + 1
			if k == 7 {
				k = 1
			}
			if this.Rattacker.Base.HasAttribute(uint8(k)) {
				this.Rdamage += 2
				break
			}
		}
	}
	//如果是技能攻击，增加技能伤害修正
	this.Rdamage += this.Rattacker.SkillDamage
	//循环技能,只调用受到攻击前释放的技能
	for e := this.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == GetHitFrontDelegate {
			this.AllSkill.RunSkill(e.Value.(*SkillStruct))
		}
	}
	//如果暴击，计算伤害
	if this.Rcritdamage > 0 {
		//循环技能,只调用受到暴击后释放的技能
		for e := this.Skill.Front(); e != nil; e = e.Next() {
			if e.Value.(*SkillStruct).Where == GetCritDelegate {
				this.AllSkill.RunSkill(e.Value.(*SkillStruct))
			}
		}
		this.Rdamage = uint8(float32(this.Rdamage) * (1 + (float32(this.Rcritdamage)+100)/100))
	}
	//如果可以闪避，并且闪避率大于命中修正，判断是否闪避
	if this.Rcanmiss && this.Miss > this.Rhitratio && this.Main.Rands.Int31n(100) < int32(this.Miss) { //闪避了
		//循环技能,只调用闪避后释放的技能
		for e := this.Skill.Front(); e != nil; e = e.Next() {
			if e.Value.(*SkillStruct).Where == MissDelegate {
				this.AllSkill.RunSkill(e.Value.(*SkillStruct))
			}
		}
		this.Main.AllStep += "1:" + strconv.Itoa(int(this.Rattacker.PositionId)) + "," + strconv.Itoa(int(this.PositionId)) + ",-1|"
	} else { //命中了
		//如果可以被防御抵消伤害，计算抵消防御后的伤害，防御力最多将伤害抵消到1
		if this.Rcanreduce {
			if this.Base.Defense >= this.Rdamage {
				if this.Rdamage > 0 {
					this.Rdamage = 1
				} else {
					this.Rdamage = 0
				}
			} else {
				this.Rdamage -= this.Base.Defense
			}
		}
		//如果有护盾，抵消一个护盾，并且不会被晕眩
		if this.HuDun > 0 {
			this.HuDun--
			this.Rdamage = 0
			this.Rdizzytime = 0
			this.Main.AllStep += "10:" + strconv.Itoa(int(this.PositionId)) + ",1|"
		}
		//暴击没有特殊效果(忽略) 重击则增加晕眩时间
		if this.Rdizzytime > 0 {
			this.DizzyLeft += this.Rdizzytime
		}
		this.Main.AllStep += "1:" + strconv.Itoa(int(this.Rattacker.PositionId)) + "," + strconv.Itoa(int(this.PositionId)) + "," + strconv.Itoa(int(this.Rdamage)) + "," + strconv.Itoa(int(this.Rcritdamage)) + "," + strconv.Itoa(int(this.Rdizzytime)) + "|"
		//如果可以反伤并且触发了反伤，对敌方造成伤害（反伤不会触发反伤），晕眩状态无法反伤
		if this.DizzyLeft == 0 && this.Rcanreflect && this.Main.Rands.Int31n(100) < int32(this.Reflect) {
			switch this.ReflectType {
			case 0: //百分比伤害
				this.Rattacker.ReduceHp(this, uint8(uint16(this.Rreflectbase)*uint16(this.Rreflectdamage)/100), 0, 0, 0, false, false, false, true, 0)
			case 1: //固定伤害
				this.Rattacker.ReduceHp(this, this.Rreflectdamage, 0, 0, 0, false, false, false, true, 0)
			}
		}
		//统计攻击者的造成伤害
		this.Rattacker.AllAttack += uint16(this.Rdamage)
		//统计被攻击者受到的伤害
		this.AllDamage += uint16(this.Rdamage)
		//扣血
		if this.Rdamage >= this.Base.Hp {
			this.Base.Hp = 0
			//循环技能,只调用攻击者的死亡前技能
			for e := this.Skill.Front(); e != nil; e = e.Next() {
				if e.Value.(*SkillStruct).Where == DieFrontDelegate {
					this.AllSkill.RunSkill(e.Value.(*SkillStruct))
				}
			}
		} else {
			this.Base.Hp -= this.Rdamage
			//受到攻击后,如果是埋伏单位,进入埋伏或者取消埋伏
			if this.Hidden {
				if this.Hide {
					this.RemoveHide()
				} else {
					this.AddHide()
				}
			}
		}
		//循环技能,只调用攻击者命中后释放的技能
		for e := this.Rattacker.Skill.Front(); e != nil; e = e.Next() {
			if e.Value.(*SkillStruct).Where == AttackHitDelegate {
				this.AllSkill.RunSkill(e.Value.(*SkillStruct))
			}
		}
	}
	//循环技能,只调用受到攻击后释放的技能
	for e := this.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == GetHitAfterDelegate {
			this.AllSkill.RunSkill(e.Value.(*SkillStruct))
		}
	}
}

/* 死亡 */
func (this *Ai) Die() {
	this.Alive = false
	//循环技能,只调用攻击者的死亡后技能
	for e := this.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == DieAfterDelegate {
			this.AllSkill.RunSkill(e.Value.(*SkillStruct))
		}
	}
	//循环技能,只调用攻击者的斩杀技能
	for e := this.Rattacker.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == AttackToDieDelegate {
			this.AllSkill.RunSkill(e.Value.(*SkillStruct))
		}
	}
	//记录攻击者杀敌数
	this.Rattacker.AllKill++
	this.Main.AllStep += "7:" + strconv.Itoa(int(this.PositionId)) + "|"
}

/* 复活 */
func (this *Ai) BackToLife() {
	//循环技能,只调用复活后释放的技能
	for e := this.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == ReLiveDelegate {
			this.AllSkill.RunSkill(e.Value.(*SkillStruct))
		}
	}
}

/* 加血 */
func (this *Ai) AddHp(adder *Ai, value uint8) {
	if value == 0 { //加血最低为1
		value = 1
	}
	this.Hadder = adder
	this.Hvalue = value
	this.Hvalue += this.ExtraAddHP
	//循环技能,只调用回血前释放的技能
	for e := this.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == AddHpFrontDelegate {
			this.AllSkill.RunSkill(e.Value.(*SkillStruct))
		}
	}
	if this.Base.Hp+value > this.MaxHP {
		this.Base.Hp = this.MaxHP
	} else {
		this.Base.Hp += value
	}
	this.Main.AllStep += "2:" + strconv.Itoa(int(this.Hadder.PositionId)) + "," + strconv.Itoa(int(this.PositionId)) + "," + strconv.Itoa(int(this.Hvalue)) + "|"
	//循环技能,只调用回血后释放的技能
	for e := this.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == AddHpAfterDelegate {
			this.AllSkill.RunSkill(e.Value.(*SkillStruct))
		}
	}
	//统计加血量
	this.Hadder.AllAdd += uint16(value)
}

/* 添加技能 */
func (this *Ai) AddSkill(skillString string) {
	if skillString == "" {
		return
	}
	skillDetail := strings.Split(skillString, "-")
	skill := &SkillStruct{}
	skill.player = this
	skill.Where = skillDetail[0]
	skill.Name = skillDetail[1]
	skillParam := skillDetail[2:]
	result := make([]uint8, len(skillParam))
	for k, _ := range skillParam {
		t, _ := strconv.Atoi(skillParam[k])
		result[k] = uint8(t)
	}
	skill.Params = result
	this.Skill.PushBack(skill)
}

/* ---BUFF方法--- */
/* 添加buff 释放者指针,buff名,效果种类,数值,增益/减益,持续回合,是否可叠加 */
func (this *Ai) AddBuff(adder *Ai, name string, last int8, add bool, params ...uint8) {
	//buff施加者和buff名称是否匹配检测这个buff是否存在
	key := strconv.Itoa(int(adder.PositionId)) + name
	_, ok := this.Buff[key]
	if ok { //已存在
		//如果可以叠加
		if add {
			//持续回合覆盖
			this.Buff[key].Last = last
			//效果叠加
			k := len(params) / 3
			for i := 0; i < k; i++ {
				//0-key 1-number 2-good
				if params[i*3+2] != this.Buff[key].Effect[params[i*3]][1] {
					//如果效果相反，数字减少
					if params[i*3+1] > this.Buff[key].Effect[params[i*3]][0] {
						//如果抵消的数字比原有的大，效果反向
						this.Buff[key].Effect[params[i*3]] = [2]uint8{params[i*3+1] - this.Buff[key].Effect[params[i*3]][0], this.Buff[key].Effect[params[i*3]][1]}
						if this.Buff[key].Effect[params[i*3]][1] == 0 {
							this.Buff[key].Effect[params[i*3]] = [2]uint8{this.Buff[key].Effect[params[i*3]][0], 1}
						} else {
							this.Buff[key].Effect[params[i*3]] = [2]uint8{this.Buff[key].Effect[params[i*3]][0], 0}
						}
					} else {
						this.Buff[key].Effect[params[i*3]] = [2]uint8{this.Buff[key].Effect[params[i*3]][0] - params[i*3+1], this.Buff[key].Effect[params[i*3]][1]}

					}
				} else { //效果相同，直接加
					this.Buff[key].Effect[params[i*3]] = [2]uint8{this.Buff[key].Effect[params[i*3]][0] + params[i*3+1], this.Buff[key].Effect[params[i*3]][1]}
				}
				//处理函数
				this.handlebuff(params[i*3], params[i*3+1], params[i*3+2])
			}
		} else {
			//不能叠加则退出
			return
		}
	} else { //不存在则赋值
		//false->0 true->1
		canAdd := "0"
		if add {
			canAdd = "1"
		}
		this.Main.AllStep += "4:" + strconv.Itoa(int(adder.PositionId)) + "," + strconv.Itoa(int(this.PositionId)) + "," + name + "," + strconv.Itoa(int(last)) + "," + canAdd + "|"
		//创建技能效果map
		effect := make(map[uint8][2]uint8)
		//赋值技能效果
		k := len(params) / 3
		for i := 0; i < k; i++ {
			//0-key 1-number 2-good
			effect[params[i*3]] = [2]uint8{params[i*3+1], params[i*3+2]}
			//处理函数
			this.handlebuff(params[i*3], params[i*3+1], params[i*3+2])
		}
		this.Buff[key] = &BuffStruct{effect, last}
	}
}

/* 移除buff */
func (this *Ai) RemoveBuff(buffer *Ai, key string) {
	key = strconv.Itoa(int(buffer.PositionId)) + key
	_, ok := this.Buff[key]
	if !ok {
		return
	}
	//调用处理函数
	this.Main.AllStep += "9:" + strconv.Itoa(int(this.PositionId)) + "," + key + "|"
	for k, v := range this.Buff[key].Effect {
		if v[1] == 0 {
			v[1] = 1
		} else {
			v[1] = 0
		}
		this.handlebuff(k, v[0], v[1])
	}
	delete(this.Buff, key)
}

/* 根据名称移除buff */
func (this *Ai) RemoveBuffByName(key string) {
	_, ok := this.Buff[key]
	if !ok {
		return
	}
	//调用处理函数
	this.Main.AllStep += "9:" + strconv.Itoa(int(this.PositionId)) + "," + key + "|"
	for k, v := range this.Buff[key].Effect {
		if v[1] == 0 {
			v[1] = 1
		} else {
			v[1] = 0
		}
		this.handlebuff(k, v[0], v[1])
	}
	delete(this.Buff, key)
}

/* 增益或者减益处理 */
func (this *Ai) handlebuff(effect uint8, number uint8, good uint8) {
	//定义一个统一修改的变量的指针
	var b [3]*uint8
	//根据效果做相应加成
	switch effect {
	case 1: //体力
		b[0] = &this.Base.Hp
		b[1] = &this.MaxHP
		b[2] = &this.LimitHP
	case 2: //武技
		b[0] = &this.Base.Attack
	case 3: //攻击次数
		b[0] = &this.Base.AttackNumber
	case 4: //防御力
		b[0] = &this.Base.Defense
	case 5: //行动力
		b[0] = &this.Base.AttackSpeed
	case 6: //暴击率
		b[0] = &this.Crit
	case 7: //暴击伤害
		b[0] = &this.CritDamage
	case 8: //闪避率
		b[0] = &this.Miss
	case 9: //技能效果
		b[0] = &this.SkillDamage
	case 10: //回复效果
		b[0] = &this.ExtraAddHP
	case 11: //护盾数量
		b[0] = &this.HuDun
	case 12: //再动次数
		b[0] = &this.ActionNumber
	case 13: //命中率
		b[0] = &this.HitRatio
	case 14: //毒火伤害
	case 15: //统帅力
		b[0] = &this.Base.Lead
	default:
	}
	//处理数据
	for _, v := range b {
		if v == nil {
			continue
		}
		if good == 1 {
			if 255-*v < number {
				*v = 255
			} else {
				*v += number
			}
		} else {
			if *v < number {
				*v = 0
			} else {
				*v -= number
			}
		}
	}
	//false->0 true->1
	isGood := "0"
	if good == 1 {
		isGood = "1"
	}
	this.Main.AllStep += "8:" + strconv.Itoa(int(this.PositionId)) + "," + strconv.Itoa(int(effect)) + "," + strconv.Itoa(int(number)) + "," + isGood + "|"
}

/* 获得坚守 */
func (this *Ai) AddTaunt() {
	this.Taunt = true
	this.Main.AllStep += "11:" + strconv.Itoa(int(this.PositionId)) + "|"
}

/* 移除坚守 */
func (this *Ai) RemoveTaunt() {
	this.Taunt = false
	this.Main.AllStep += "12:" + strconv.Itoa(int(this.PositionId)) + "|"
}

/* 进入埋伏 */
func (this *Ai) AddHide() {
	this.Hide = true
	this.Main.AllStep += "13:" + strconv.Itoa(int(this.PositionId)) + "|"
}

/* 移除埋伏 */
func (this *Ai) RemoveHide() {
	this.Hide = false
	this.Main.AllStep += "14:" + strconv.Itoa(int(this.PositionId)) + "|"
}

/* 行动力清零 */
func (this *Ai) ActionZero() {
	this.AlreadyRun = this.Main.Time
	this.Main.AllStep += "15:" + strconv.Itoa(int(this.PositionId)) + "|"
}

/* 进入混乱 */
func (this *Ai) AddChaos(number uint8) {
	this.Chaos += number
	this.Main.AllStep += "16:" + strconv.Itoa(int(this.PositionId)) + "," + strconv.Itoa(int(number)) + "|"
}

/* 移除混乱 */
func (this *Ai) RemoveChaos() {
	this.Chaos = 0
	this.Main.AllStep += "17:" + strconv.Itoa(int(this.PositionId)) + "|"
}

/* ---几种寻找目标方法--- */
/* 寻找具有/不具有嘲讽属性的目标集合 */
func (this *Ai) FindTauntArray(we bool, taunt bool, number uint8) []*Ai {
	team := this.FindTeam(we)
	if team == nil {
		return nil
	}
	//返回的数组
	var returnArray []*Ai
	//生成数组大小的随机数
	r := this.Main.Rands.Perm(len(team))
	//模拟随机访问数组
	for _, v := range r {
		//如果目标具有/不具有嘲讽属性
		if team[v].Taunt == taunt && number > 0 {
			number--
			returnArray = append(returnArray, team[v])
		}
	}
	return returnArray
}

/* 随机查找目标集合 */
func (this *Ai) FindRandArray(we bool, number int) []*Ai {
	team := this.FindTeam(we)
	if team == nil {
		return nil
	}
	//返回的数组
	var returnArray []*Ai
	//生成数组大小的随机数
	r := this.Main.Rands.Perm(len(team))
	//模拟随机访问数组
	for _, v := range r {
		if number > 0 {
			number--
			returnArray = append(returnArray, team[v])
		}
	}
	return returnArray
}

/* 查找拥有某个属性的集合 */
func (this *Ai) FindAttributeArray(we bool, attribute uint8, number int) []*Ai {
	team := this.FindTeam(we)
	if team == nil {
		return nil
	}
	//返回的数组
	var returnArray []*Ai
	//生成数组大小的随机数
	r := this.Main.Rands.Perm(len(team))
	//模拟随机访问数组
	for _, v := range r {
		if number > 0 && team[v].Base.HasAttribute(attribute) {
			number--
			returnArray = append(returnArray, team[v])
		}
	}
	return returnArray
}

/* 查找血量最少的目标 */
func (this *Ai) FindLessHp(we bool) *Ai {
	team := this.FindTeam(we)
	if team == nil {
		return nil
	}
	returnAi := team[0]
	for _, v := range team {
		if v.Base.Hp < returnAi.Base.Hp {
			returnAi = v
		}
	}
	return returnAi
}

/* 寻找有/没有某个buff的目标集合 */
func (this *Ai) FindBuffArray(we bool, adder *Ai, name string, has bool, number uint8) []*Ai {
	key := strconv.Itoa(int(adder.PositionId)) + name
	team := this.FindTeam(we)
	if team == nil {
		return nil
	}
	//返回的数组
	var returnArray []*Ai
	//生成数组大小的随机数
	r := this.Main.Rands.Perm(len(team))
	//模拟随机访问数组
	for _, v := range r {
		_, ok := team[v].Buff[key]
		if number > 0 && ok == has {
			number--
			returnArray = append(returnArray, team[v])
		}
	}
	return returnArray
}

/* 寻找所有队友(除了自己) */
func (this *Ai) FindWeExceptMe() []*Ai {
	team := this.FindTeam(true)
	//返回的数组
	var returnArray []*Ai
	//模拟随机访问数组
	for _, v := range team {
		if v != this {
			returnArray = append(returnArray, v)
		}
	}
	return returnArray
}

/* 返回埋伏状态数组 */
func (this *Ai) FindHide(we bool, number int) []*Ai {
	//返回的指针数组
	var returnArray []*Ai
	//选择的阵营
	var teamArray []*Ai
	if we {
		//是我方
		if this.IsMy {
			teamArray = this.Main.MyTeam
		} else {
			teamArray = this.Main.EnemieTeam
		}
	} else {
		//是敌方
		if this.IsMy {
			teamArray = this.Main.EnemieTeam
		} else {
			teamArray = this.Main.MyTeam
		}
	}
	for _, v := range teamArray {
		//如果存活，则添加到返回数组里
		if number > 0 && v.Alive && v.Hide == true {
			number--
			returnArray = append(returnArray, v)
		}
	}
	return returnArray
}

/* 返回我方或者敌方队伍（存活的） 如果是敌人，只能选择非埋伏的 */
func (this *Ai) FindTeam(we bool) []*Ai {
	//返回的指针数组
	var returnArray []*Ai
	//选择的阵营
	var teamArray []*Ai
	if we {
		//是我方
		if this.IsMy {
			teamArray = this.Main.MyTeam
		} else {
			teamArray = this.Main.EnemieTeam
		}
	} else {
		//是敌方
		if this.IsMy {
			teamArray = this.Main.EnemieTeam
		} else {
			teamArray = this.Main.MyTeam
		}
	}
	for _, v := range teamArray {
		//如果存活，则添加到返回数组里
		if v.Alive {
			//如果是敌人，不能选择敌方埋伏的目标
			if we == true || v.Hide == false {
				returnArray = append(returnArray, v)
			}
		}
	}
	return returnArray
}
