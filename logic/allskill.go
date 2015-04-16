package logic

import (
	"strconv"
	"strings"
)

type AllSkill struct {
}

/* 技能列表 */
const (
	JiJiang            = "激将"
	JiJiangEnd         = "激将结束"
	BaGua              = "八卦"
	WenZhan            = "温斩"
	YiJue              = "义绝"
	YiJueStart         = "义绝开始"
	YiJueEnd           = "义绝结束"
	DuanHou            = "断吼"
	LongDan            = "龙胆"
	LongTi             = "龙体"
	LianHuan           = "连环"
	NiePan             = "涅槃"
	TieQi              = "铁骑"
	TuXi               = "突袭"
	MiJian             = "弥坚"
	MiJianEnd          = "弥坚结束"
	XiangLe            = "享乐"
	KuangGu            = "狂骨"
	AoGu               = "傲骨"
	EnYuan             = "恩怨"
	EnYuanRun          = "恩怨执行"
	DuLiang            = "度量"
	ShengXi            = "生息"
	KuangJian          = "匡谏"
	ZhongZheng         = "忠正"
	LianNu             = "连弩"
	NanMan             = "南蛮"
	TengJia            = "藤甲"
	TengJiaRun         = "藤甲发动"
	TiaoXin            = "挑衅"
	ShiCheng           = "师承"
	FuZheng            = "辅政"
	WuDi               = "武帝"
	GuiXin             = "归心"
	GuiXinRun          = "归心发动"
	GangLie            = "刚烈"
	QiangXi            = "强袭"
	WenDi              = "文帝"
	WenDiEnd           = "文帝结束"
	ChengXiang         = "称象"
	ChengXiangRun      = "称象发动"
	HuBao              = "虎豹"
	HuBaoEnd           = "虎豹结束"
	MingDi             = "明帝"
	MingDiRun          = "明帝发动"
	CaoZu              = "曹祖"
	CaoZuEnd           = "曹祖结束"
	LuoYi              = "裸衣"
	RenJie             = "忍戒"
	RenJieRun          = "忍戒发动"
	QuHu               = "驱虎"
	MouZhu             = "谋主"
	QiCe               = "奇策"
	LuanWu             = "乱舞"
	WeiMu              = "帷幕"
	WeiMu1             = "帷幕1"
	WeiMu2             = "帷幕2"
	FuJi               = "伏计"
	GangLi             = "刚戾"
	HuiXue             = "回血"
	Buff               = "增益"
	LiaoHua            = "廖化"
	AddAttribute       = "加力"
	NotTauntDamage     = "非坚破"
	ZhuRong            = "祝融"
	XianXing           = "现形"
	TongGong           = "同功"
	TongGongWrite      = "同功记录"
	TongGongRun        = "同功发动"
	ZhuGeZhan          = "诸葛瞻"
	YiJi               = "伊籍"
	LiuBa              = "刘巴"
	MaZhong            = "马忠"
	WangPing           = "王平"
	LiHui              = "李恢"
	CaoXiu             = "曹休"
	CaoXiuStart        = "曹休开始"
	CaoXiuRun          = "曹休技能开始"
	CaoAng             = "曹昂"
	CaoAngEnd          = "曹昂结束"
	CaoZhi             = "曹植"
	CaoZhiRun          = "曹植发动"
	CaoSong            = "曹嵩"
	CaoShuang          = "曹爽"
	CaoShuangRun       = "曹爽发动"
	CaoDe              = "曹德"
	QingTengJia        = "轻藤甲"
	QingTengJiaRun     = "轻藤甲回合"
	QingTengJiaDefense = "轻藤甲受到攻击"
	ShiDun             = "石盾"
	PoJia              = "破甲"
)

/* 执行技能 */
func (this *AllSkill) RunSkill(skill *SkillStruct) {
	//循环技能,只调用技能触发前释放的技能
	for e := skill.player.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == "17" {
			//指针thisSkill 代替 指针e.Value.(*SkillStruct)
			thisSkill := e.Value.(*SkillStruct)
			switch thisSkill.Name {
			case WeiMu2: //遗技失效
				if skill.Where == "11" && thisSkill.player.SkillTempAi[WeiMu].Alive { //如果技能是遗技，帷幕技能释放者没有死亡则退出
					this.WriteStepSkill(thisSkill.player.SkillTempAi[WeiMu], WeiMu)
					return
				}
			}
		}
	}
	switch skill.Name {
	case JiJiang: //激将 本方所有武将加武技增加
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindTeam(true) {
			v.AddBuff(skill.player, "激将", -1, false, 2, skill.Params[0], 1)
		}
		//死亡时候调用结束技能
		newskillend := &SkillStruct{"11", JiJiangEnd, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskillend)
	case JiJiangEnd: //激将结束
		for _, v := range skill.player.FindTeam(true) {
			v.RemoveBuff(skill.player, "激将")
		}
	case BaGua: //八卦 每回合触发下面一个效果：【生】本方1武将+体力 【休】本方所有武将+体力【开】本方1武将永久+武技
		//【杜】本方1武将获得护盾【景】敌方1武将永久-武技【惊】敌方1武将眩晕秒【伤】敌方所有武将点伤害【死】敌方1武将点伤害
		randNumber := skill.player.Main.Rands.Int31n(8)
		switch randNumber {
		case 0:
			this.WriteStepSkill(skill.player, skill.Name+"生")
			skill.player.FindLessHp(true).AddHp(skill.player, skill.Params[0])
		case 1:
			this.WriteStepSkill(skill.player, skill.Name+"休")
			for _, v := range skill.player.FindTeam(true) {
				v.AddHp(skill.player, skill.Params[1])
			}
		case 2:
			target := skill.player.FindRandArray(true, 1)
			if target != nil {
				this.WriteStepSkill(skill.player, skill.Name+"开")
				target[0].AddBuff(skill.player, "开", -1, true, 2, skill.Params[2], 1)
			}
		case 3:
			target := skill.player.FindRandArray(true, 1)
			if target != nil {
				this.WriteStepSkill(skill.player, skill.Name+"杜")
				target[0].AddBuff(skill.player, "杜", -1, true, 11, skill.Params[3], 1)
			}
		case 4:
			target := skill.player.FindRandArray(false, 1)
			if target != nil {
				this.WriteStepSkill(skill.player, skill.Name+"景")
				target[0].AddBuff(skill.player, "景", -1, true, 2, skill.Params[4], 0)
			}
		case 5:
			target := skill.player.FindRandArray(false, 1)
			if target != nil {
				this.WriteStepSkill(skill.player, skill.Name+"惊")
				target[0].ReduceHp(skill.player, 0, 0, 0, float32(skill.Params[5]), false, false, false, true, 2)
			}
		case 6:
			this.WriteStepSkill(skill.player, skill.Name+"伤")
			for _, v := range skill.player.FindTeam(false) {
				v.ReduceHp(skill.player, skill.Params[6], 0, 0, 0, false, false, false, false, 2)
			}
		case 7:
			target := skill.player.FindRandArray(false, 1)
			if target != nil {
				this.WriteStepSkill(skill.player, skill.Name+"死")
				target[0].ReduceHp(skill.player, skill.Params[7], 0, 0, 0, false, false, false, true, 2)
			}
		}
	case WenZhan: //温斩 斩杀后永久+武技
		this.WriteStepSkill(skill.player, skill.Name)
		skill.player.AddBuff(skill.player, "温斩", -1, true, 2, skill.Params[0], 1)
	case YiJue: //义绝 每有一个友方武将，+暴击率 为每个武将绑定死亡后触发事件
		this.WriteStepSkill(skill.player, skill.Name)
		skill.player.AddBuff(skill.player, "义绝", -1, true, 6, skill.Params[0]*uint8(len(skill.player.FindTeam(true))-1), 1)
		for _, v := range skill.player.FindTeam(true) {
			newskillstart := &SkillStruct{"15", YiJueStart, skill.Params, skill.player}
			v.Skill.PushBack(newskillstart)
			newskillend := &SkillStruct{"11", YiJueEnd, skill.Params, skill.player}
			v.Skill.PushBack(newskillend)
		}
	case YiJueStart: //义绝开始 +暴击率
		if skill.player.Alive {
			skill.player.AddBuff(skill.player, "义绝", -1, true, 6, skill.Params[0], 1)
		}
	case YiJueEnd: //义绝结束 -暴击率
		if skill.player.Alive {
			skill.player.AddBuff(skill.player, "义绝", -1, true, 6, skill.Params[0], 0)
		}
	case DuanHou: //断吼 如果我方没有坚守武将，敌方晕眩秒，防御减少点，只能触发一次
		if skill.player.SkillTemp[DuanHou] == 0 && skill.player.FindTauntArray(true, true, 1) == nil { //我方没有坚守武将
			this.WriteStepSkill(skill.player, skill.Name)
			skill.player.SkillTemp[DuanHou] = 1
			for _, v := range skill.player.FindTeam(false) {
				v.ReduceHp(skill.player, 0, 0, 0, float32(skill.Params[0]), false, false, false, false, 2)
				v.AddBuff(skill.player, "断吼", -1, false, 4, skill.Params[1], 0)
			}
		}
	case LongDan: //龙胆 闪现：攻击次数+
		skill.player.AddBuff(skill.player, "龙胆", 1, false, 3, skill.Params[0], 1)
	case LongTi: //龙体 闪避率+
		this.WriteStepSkill(skill.player, skill.Name)
		skill.player.AddBuff(skill.player, "龙体", -1, true, 8, skill.Params[0], 1)
	case LianHuan: //连环 对地方全体造成伤害
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindTeam(false) {
			v.ReduceHp(skill.player, skill.Params[0], 0, 0, 0, false, false, false, false, 2)
		}
	case NiePan: //涅槃 复活一次
		if skill.player.SkillTemp[NiePan] == 0 {
			this.WriteStepSkill(skill.player, skill.Name)
			skill.player.SkillTemp[NiePan] = 1
			skill.player.AddHp(skill.player, skill.player.MaxHP-skill.player.Base.Hp)
		}
	case TieQi: //铁骑 行动力增加
		this.WriteStepSkill(skill.player, skill.Name)
		skill.player.AddBuff(skill.player, "铁骑", -1, true, 5, skill.Params[0], 1)
	case TuXi: //突袭 本次攻击锁定对方体力最低的单位，并对其造成1秒眩晕（仅1次）
		if skill.player.SkillTemp[NiePan] == 0 {
			skill.player.SkillTemp[NiePan] = 1
			this.WriteStepSkill(skill.player, skill.Name)
			skill.player.Atarget = skill.player.FindLessHp(false)
			skill.player.Adizzytime = float32(skill.Params[0])
		}
	case MiJian: //弥坚 火属性武将武技+
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindAttributeArray(true, 3, 255) {
			v.AddBuff(skill.player, "弥坚", -1, false, 2, skill.Params[0], 1)
		}
		//死亡时候调用结束技能
		newskillend := &SkillStruct{"11", MiJianEnd, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskillend)
	case MiJianEnd: //弥坚结束
		for _, v := range skill.player.FindAttributeArray(true, 3, 255) {
			v.RemoveBuff(skill.player, "坚")
		}
	case XiangLe: //享乐 所有武将体力-
		this.WriteStepSkill(skill.player, skill.Name)
		all := skill.player.FindTeam(true)
		all = append(all, skill.player.FindTeam(false)...)
		for _, v := range all {
			v.AddBuff(skill.player, "享乐", -1, false, 1, v.Base.Hp/2, 0)
		}
	case KuangGu: //狂骨 每次攻击回复%血
		this.WriteStepSkill(skill.player, skill.Name)
		skill.player.AddHp(skill.player.Atarget, skill.player.Atarget.Rdamage/skill.Params[0])
	case AoGu: //傲骨 每回合增加个护盾
		this.WriteStepSkill(skill.player, skill.Name)
		skill.player.AddBuff(skill.player, "傲骨", 1, false, 11, 1, 1)
	case EnYuan: //恩怨 本方所有回复效果变为对对方的伤害效果
		//为所有我方单位绑定回血时技能
		for _, v := range skill.player.FindTeam(true) {
			newskillrun := &SkillStruct{"13", EnYuanRun, skill.Params, v}
			v.SkillTempAi[EnYuan] = skill.player
			v.Skill.PushBack(newskillrun)
		}
	case EnYuanRun: //恩怨执行
		if skill.player.SkillTempAi[EnYuan].Alive {
			//如果恩怨武将还活着，回复变伤害
			target := skill.player.FindRandArray(false, 1)
			if target != nil {
				this.WriteStepSkill(skill.player.SkillTempAi[EnYuan], EnYuan)
				target[0].ReduceHp(skill.player.SkillTempAi[EnYuan], skill.player.Hvalue, 0, 0, 0, false, false, false, true, 2)
				skill.player.Hvalue = 0
			}
		}
	case DuLiang: //度量 我方体力没有上线（255）
		for _, v := range skill.player.FindTeam(true) {
			v.LimitHP = 255
		}
	case ShengXi: //生息 回复友方1名武将当前体力的%
		this.WriteStepSkill(skill.player, skill.Name)
		target := skill.player.FindLessHp(true)
		target.AddHp(skill.player, target.Base.Hp*skill.Params[0]/100)
	case KuangJian: //匡谏 本方技能伤害+
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindTeam(true) {
			v.AddBuff(skill.player, "匡谏", -1, false, 9, skill.Params[0], 1)
		}
	case ZhongZheng: //忠正 本方回复效果+
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindTeam(true) {
			v.AddBuff(skill.player, "忠正", -1, false, 10, skill.Params[0], 1)
		}
	case LianNu: //连弩 1个本方武将攻击次数+
		this.WriteStepSkill(skill.player, skill.Name)
		skill.player.FindRandArray(true, 1)[0].AddBuff(skill.player, "连弩", 1, true, 3, skill.Params[0], 1)
	case NanMan: // 南蛮 战鼓，对对方所有武将造成点伤害
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindTeam(false) {
			v.ReduceHp(skill.player, skill.Params[0], 0, 0, 0, false, false, false, false, 2)
		}
	case TengJia: //藤甲 使本方1名武将防御力+，但受到火属性武将攻击时伤害+
		target := skill.player.FindBuffArray(true, skill.player, "藤甲", false, 1)
		if target != nil {
			this.WriteStepSkill(skill.player, skill.Name)
			target[0].AddBuff(skill.player, "藤甲", -1, false, 4, skill.Params[0], 1)
			//绑定藤甲发动技能
			newskillrun := &SkillStruct{"6", TengJiaRun, skill.Params, target[0]}
			target[0].Skill.PushBack(newskillrun)
		}
	case TengJiaRun: //藤甲执行
		if skill.player.Rattacker.Base.HasAttribute(3) { //如果攻击者是火属性，则触发火烧藤甲
			this.WriteStepSkill(skill.player.Rattacker, "火烧藤甲")
			skill.player.Rdamage += skill.Params[1]
		}
	case TiaoXin: //挑衅 每回合让对方武将获得坚守属性
		target := skill.player.FindTauntArray(false, false, 1)
		if target != nil {
			this.WriteStepSkill(skill.player, skill.Name)
			target[0].AddTaunt()
		}
	case ShiCheng: //师承 体力低于点获得诸葛亮的八卦技
		if skill.player.SkillTemp[ShiCheng] == 0 && skill.player.Base.Hp <= skill.Params[0] {
			skill.player.SkillTemp[ShiCheng] = 1
			this.WriteStepSkill(skill.player, skill.Name)
			skill.player.AddSkill(skill.player.Base.FindSkill(566, BaGua))
			skill.player.Main.AllStep += "4:" + strconv.Itoa(int(skill.player.PositionId)) + "," + strconv.Itoa(int(skill.player.PositionId)) + "," + BaGua + ",-1,0|"
		}
	case FuZheng: //辅政 由我方武将代替它的回合
		target := skill.player.FindWeExceptMe()
		if target != nil {
			this.WriteStepSkill(skill.player, skill.Name)
			target[0].StartAction()
		}
	case WuDi: //武帝 曹操受到的伤害时，其他每名武将承担点伤害
		damage := skill.player.Rdamage
		team := skill.player.FindWeExceptMe()
		if team != nil {
			this.WriteStepSkill(skill.player, skill.Name)
			for _, v := range team {
				if !strings.Contains(v.Base.Name, "曹") {
					continue
				}
				if damage > skill.Params[0] && damage != 0 {
					damage -= skill.Params[0]
					v.ReduceHp(skill.player.Rattacker, skill.Params[0], 0, 0, 0, false, false, false, true, 2)
				} else if damage != 0 {
					damage = 0
					v.ReduceHp(skill.player.Rattacker, skill.Params[0]-damage, 0, 0, 0, false, false, false, true, 2)
				}
			}
			skill.player.Rdamage = damage
		}
	case GuiXin: //归心 本方其他武将回合开始时，会为曹操回复1点体力
		for _, v := range skill.player.FindWeExceptMe() {
			//添加归心执行
			newskill := &SkillStruct{"2", GuiXinRun, skill.Params, v}
			v.SkillTempAi[GuiXin] = skill.player
			v.Skill.PushBack(newskill)
		}
	case GuiXinRun: //归心执行
		if skill.player.SkillTempAi[GuiXin].Alive == true { //如果曹操还活着
			this.WriteStepSkill(skill.player.SkillTempAi[GuiXin], GuiXin)
			//给他加血
			skill.player.SkillTempAi[GuiXin].AddHp(skill.player, skill.Params[0])
		}
	case GangLie: //刚烈 反戈，20%概率，反戈100%伤害，且自己不受伤害
		if skill.player.Main.Rands.Int31n(100) < int32(skill.Params[0]) {
			this.WriteStepSkill(skill.player, skill.Name)
			skill.player.Rreflectdamage = skill.Params[1]
			skill.player.Rdamage = 0
		}
	case QiangXi: //强袭 对对方所有武将造成等同于防御的伤害
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindTeam(false) {
			v.ReduceHp(skill.player, skill.player.Base.Defense, 0, 0, 0, false, false, false, false, 2)
		}
	case WenDi: //文帝 曹氏所有武将武技+，防御+
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindTeam(true) {
			if strings.Contains(v.Base.Name, "曹") {
				v.AddBuff(skill.player, "文帝", -1, false, 2, skill.Params[0], 1, 4, skill.Params[1], 1)
			}
		}
		//死亡时候调用结束技能
		newskillend := &SkillStruct{"11", WenDiEnd, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskillend)
	case WenDiEnd: //文帝
		for _, v := range skill.player.FindTeam(true) {
			v.RemoveBuff(skill.player, "文帝")
		}
	case LiaoHua: //廖化 激怒-行动力变为本方最快
		maxSpeed := skill.player.Base.AttackSpeed
		for _, v := range skill.player.FindTeam(true) {
			if v.Base.AttackSpeed > maxSpeed {
				maxSpeed = v.Base.AttackSpeed
			}
		}
		skill.player.AddBuff(skill.player, "极速", 1, false, 5, maxSpeed-skill.player.Base.AttackSpeed+1, 1)
	case HuiXue: //回血 1:我方一名武将 2:我方全部武将 回复血
		switch skill.Params[0] {
		case 1:
			skill.player.FindLessHp(true).AddHp(skill.player, skill.Params[1])
		case 2:
			for _, v := range skill.player.FindTeam(true) {
				v.AddHp(skill.player, skill.Params[1])
			}
		}
	case Buff: //增益 0: 0自己 1我方全体 2敌方全体
		//1:持续回合 2:是否叠加
		//value 1:选择增益类型 2:效果 2:好坏...
		var target []*Ai
		switch skill.Params[0] {
		case 0:
			target = make([]*Ai, 1)
			target[0] = skill.player
		case 1:
			target = skill.player.FindTeam(true)
		case 2:
			target = skill.player.FindTeam(false)
		}
		value := skill.Params[3:]
		//0->false 1->true
		last := int8(skill.Params[1])
		if skill.Params[1] == 0 {
			last = -1
		}
		canAdd := false
		if skill.Params[2] == 1 {
			canAdd = true
		}
		for _, v := range target {
			v.AddBuff(skill.player, skill.player.Base.Name, last, canAdd, value...)
		}
	case AddAttribute: //加力
		number := len(skill.Params) / 2
		for i := 0; i < number; i++ {
			skill.player.Main.AllStep += "5:" + strconv.Itoa(int(skill.player.PositionId)) + "," + strconv.Itoa(int(skill.Params[i])) + "," + strconv.Itoa(int(skill.Params[i+1])) + "|"
		}
	case NotTauntDamage: //非坚破 对非坚守武将伤害+
		if skill.player.Atarget.Taunt == false { //如果目标是非坚守
			skill.player.Adamage += skill.Params[0]
		}
	case ZhuRong: //祝融 对对方1名武将造成4点伤害，对其余武将造成1点伤害
		target := skill.player.FindRandArray(false, 1)[0]
		all := skill.player.FindTeam(true)
		all = append(all, skill.player.FindTeam(false)...)
		for _, v := range all {
			if v == target {
				v.ReduceHp(skill.player, skill.Params[0], 0, 0, 0, false, false, false, true, 2)
			} else {
				v.ReduceHp(skill.player, skill.Params[1], 0, 0, 0, false, false, false, false, 2)
			}
		}
	case XianXing: //现形 让敌方埋伏武将现形 0:随机一个 1:全部
		switch skill.Params[0] {
		case 0:
			target := skill.player.FindHide(false, 1)
			if target != nil {
				target[0].RemoveHide()
			}
		case 1:
			for _, v := range skill.player.FindHide(false, 255) {
				v.RemoveHide()
			}
		}
	case TongGong: //攻击伤害变为友方武将上一次的数值
		//为所有友方武将绑定技能，将它的攻击力记录下来
		for _, v := range skill.player.FindWeExceptMe() {
			newskill := &SkillStruct{"4", TongGongWrite, skill.Params, v}
			v.SkillTempAi[TongGong] = skill.player
			v.Skill.PushBack(newskill)
		}
		//为自己攻击前绑定触发同功
		newskill := &SkillStruct{"4", TongGongRun, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskill)
	case TongGongWrite: //给同功武将记录攻击力
		skill.player.SkillTempAi[TongGong].SkillTemp[TongGong] = int8(skill.player.Adamage)
	case TongGongRun: //发动同功
		skill.player.Adamage = uint8(skill.player.SkillTemp[TongGong])
	case ZhuGeZhan: //战鼓，选择一个友方武将，在自己回合结束时，立即执行此武将的回合
		target := skill.player.FindWeExceptMe()
		if target != nil {
			target[0].StartAction()
			target[0].ActionZero()
		}
	case YiJi: //被对方武将攻击时，体力-4，并使其目标变为对方一个武将
		skill.player.AddBuff(skill.player, "虚弱", -1, true, 1, skill.Params[0], 0)
		skill.player.Rattacker.AddChaos(1)
	case LiuBa: //遗计，对斩杀刘巴的对方武将造成点伤害，对其余对方武将点伤害
		for _, v := range skill.player.FindTeam(false) {
			if v == skill.player.Rattacker {
				v.ReduceHp(skill.player, skill.Params[0], 0, 0, 0, false, false, false, false, 2)
			} else {
				v.ReduceHp(skill.player, skill.Params[1], 0, 0, 0, false, false, false, false, 2)
			}
		}
	case MaZhong: //马忠 如果连续三次攻击同一武将，则必出暴击
		if skill.player.SkillTemp[MaZhong] == 0 {
			skill.player.SkillTemp[MaZhong] = 1
			skill.player.SkillTempAi[MaZhong] = skill.player.Atarget
		} else {
			if skill.player.SkillTempAi[MaZhong] == skill.player.Atarget {
				if skill.player.SkillTemp[MaZhong] == 2 { // 出暴击
					skill.player.Acritdamage = skill.player.CritDamage
					skill.player.SkillTempAi[MaZhong] = nil
				} else {
					skill.player.SkillTemp[MaZhong]++
				}
			} else {
				skill.player.SkillTemp[MaZhong] = 0
				skill.player.SkillTempAi[MaZhong] = skill.player.Atarget
			}
		}
	case WangPing: //攻击一个武将后，锁定此武将
		if skill.player.SkillTempAi[WangPing] == nil {
			skill.player.SkillTempAi[WangPing] = skill.player.Atarget
		} else {
			if skill.player.SkillTempAi[WangPing].Alive == true {
				skill.player.Atarget = skill.player.SkillTempAi[WangPing]
			} else {
				skill.player.SkillTempAi[WangPing] = skill.player.Atarget
			}
		}
	case LiHui: //战鼓，恢复一个友方武将等同于防御力一半数值的体力
		target := skill.player.FindRandArray(true, 1)
		if target != nil {
			target[0].AddHp(skill.player, target[0].Base.Defense/2)
		}
	case ChengXiang: //称象 战斗开始时，选本方一个曹氏武将，每回合开始时，交换武技和防御
		var target *Ai
		team := skill.player.FindWeExceptMe()
		r := skill.player.Main.Rands.Perm(len(team))
		//模拟随机访问数组
		for _, v := range r {
			if strings.Contains(team[v].Base.Name, "曹") {
				target = team[v]
				break
			}
		}
		skill.player.SkillTempAi[ChengXiang] = target
		//绑定称象发动技能
		newskillRun := &SkillStruct{"2", ChengXiangRun, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskillRun)
	case ChengXiangRun: //称象发动
		if skill.player.SkillTempAi[ChengXiang] != nil && skill.player.SkillTempAi[ChengXiang].Alive == true {
			this.WriteStepSkill(skill.player, ChengXiang)
			//存储防御力数值
			defense := skill.player.SkillTempAi[ChengXiang].Base.Defense
			skill.player.SkillTempAi[ChengXiang].Base.Defense = skill.player.SkillTempAi[ChengXiang].Base.Attack
			skill.player.SkillTempAi[ChengXiang].Base.Attack = defense
			//如果身上有加武技的buff，重新增加一次
			buffAdd := uint8(0)
			for _, v := range skill.player.SkillTempAi[ChengXiang].Buff {
				for key, val := range v.Effect {
					if key == 2 {
						//如果是武技增益则重新增加武技
						if val[1] == 1 {
							buffAdd += val[0]
						} else {
							//如果是减益
							if val[0] > buffAdd {
								buffAdd = 0
							} else {
								buffAdd -= val[0]
							}
						}
					}
				}
			}
			skill.player.SkillTempAi[ChengXiang].Base.Attack += buffAdd
			skill.player.Main.AllStep += "18:" + strconv.Itoa(int(skill.player.SkillTempAi[ChengXiang].PositionId)) + "," + strconv.Itoa(int(buffAdd)) + "|"
		}
	case HuBao: //虎豹 本方其他曹氏武将行动力+25
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindTeam(true) {
			if strings.Contains(v.Base.Name, "曹") {
				v.AddBuff(skill.player, "虎豹", -1, false, 5, skill.Params[0], 1)
			}
		}
		//死亡时候调用结束技能
		newskillend := &SkillStruct{"11", HuBaoEnd, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskillend)
	case HuBaoEnd: //虎豹结束
		for _, v := range skill.player.FindTeam(true) {
			v.RemoveBuff(skill.player, "虎豹")
		}
	case MingDi: //明帝 本方曹氏武将每进行1个回合，曹睿武技永久+，防御永久+
		for _, v := range skill.player.FindTeam(true) {
			if strings.Contains(v.Base.Name, "曹") {
				newskillRun := &SkillStruct{"3", MingDiRun, skill.Params, skill.player}
				v.Skill.PushBack(newskillRun)
			}
		}
	case MingDiRun: //明帝执行
		if skill.player.Alive {
			//如果明帝技能目标还活着
			this.WriteStepSkill(skill.player, MingDi)
			skill.player.AddBuff(skill.player, "明帝", -1, true, 2, skill.Params[0], 1, 4, skill.Params[1], 1)
		}
	case CaoZu: // 曹祖 所有曹氏武将武技+，防御+，行动力+，体力+，技能伤害+
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindTeam(true) {
			if strings.Contains(v.Base.Name, "曹") {
				v.AddBuff(skill.player, "曹祖", -1, false, 2, skill.Params[0], 1, 4, skill.Params[1], 1, 5, skill.Params[2], 1, 1, skill.Params[3], 1, 9, skill.Params[4], 1)
			}
		}
		//死亡时候调用结束技能
		newskillend := &SkillStruct{"11", CaoZuEnd, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskillend)
	case CaoZuEnd: //文帝
		for _, v := range skill.player.FindTeam(true) {
			v.RemoveBuff(skill.player, "武")
			v.RemoveBuff(skill.player, "防")
			v.RemoveBuff(skill.player, "行")
			v.RemoveBuff(skill.player, "体")
			v.RemoveBuff(skill.player, "技")
		}
	case CaoXiu: //曹休 如果第一次攻击造成伤害不多于,则再攻击一次
		newskillStart := &SkillStruct{"2", CaoXiuStart, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskillStart)
		newskillRun := &SkillStruct{"4", CaoXiuRun, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskillRun)
	case CaoXiuStart:
		skill.player.SkillTemp[CaoXiu] = 0
	case CaoXiuRun:
		if skill.player.SkillTemp[CaoXiu] == 0 && skill.player.Atarget.Rdamage <= skill.Params[0] {
			//如果第一次攻击伤害不足,再次攻击这个目标
			skill.player.SkillTemp[CaoXiu] = 1
			skill.player.Attack(skill.player.Atarget)
		}
	case CaoAng: //曹昂 本方曹氏武将伤害+
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindTeam(true) {
			if strings.Contains(v.Base.Name, "曹") {
				v.AddBuff(skill.player, "爆裂", -1, false, 2, skill.Params[0], 1)
			}
		}
		//死亡时候调用结束技能
		newskillend := &SkillStruct{"11", CaoAngEnd, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskillend)
	case CaoAngEnd: //曹昂结束
		for _, v := range skill.player.FindTeam(true) {
			v.RemoveBuff(skill.player, "曹")
		}
	case CaoZhi: //曹植 当其他曹氏武将阵亡时，对对方所有武将造成5点伤害
		//为所有曹氏武将绑定阵亡时触发的技能
		for _, v := range skill.player.FindWeExceptMe() {
			if strings.Contains(v.Base.Name, "曹") {
				newskillrun := &SkillStruct{"11", CaoZhiRun, skill.Params, skill.player}
				v.Skill.PushBack(newskillrun)
			}
		}
	case CaoZhiRun: //死亡时 触发曹植技能
		if skill.player.Alive {
			//如果曹植还活着
			for _, v := range skill.player.FindTeam(false) {
				v.ReduceHp(skill.player, skill.Params[0], 0, 0, 0, false, false, false, false, 2)
			}
		}
	case CaoSong: //其他曹氏武将，体力+，体力上限+
		for _, v := range skill.player.FindWeExceptMe() {
			if strings.Contains(v.Base.Name, "曹") {
				v.AddBuff(skill.player, "精力", -1, false, 1, skill.Params[0], 1)
				if 255-v.MaxHP < skill.Params[1] {
					v.MaxHP = 255
				} else {
					v.MaxHP += skill.Params[1]
				}
				if 255-v.LimitHP < skill.Params[1] {
					v.LimitHP = 255
				} else {
					v.LimitHP += skill.Params[1]
				}
			}
		}
	case CaoShuang: //曹爽 本方曹氏武将每受到一次，受到伤害的武将防御+
		//给每个曹氏武将绑定技能
		for _, v := range skill.player.FindTeam(true) {
			if strings.Contains(v.Base.Name, "曹") {
				newskillrun := &SkillStruct{"7", CaoShuangRun, skill.Params, v}
				v.Skill.PushBack(newskillrun)
				v.SkillTempAi[CaoShuang] = skill.player
			}
		}
	case CaoShuangRun: //曹爽执行
		if skill.player.SkillTempAi[CaoShuang].Alive {
			//如果曹爽活着，则触发
			skill.player.AddBuff(skill.player.SkillTempAi[CaoShuang], "铁甲", -1, true, 4, skill.Params[0], 1)
		}
	case CaoDe: //曹德 战鼓，为本方一个曹氏武将加【护盾】
		team := skill.player.FindTeam(true)
		r := skill.player.Main.Rands.Perm(len(team))
		//模拟随机访问数组
		for _, v := range r {
			if strings.Contains(team[v].Base.Name, "曹") {
				team[v].AddBuff(skill.player, "护盾", -1, true, 11, skill.Params[0], 1)
				break
			}
		}
	case LuoYi: //裸衣 战鼓，武技永久+，防御永久-
		this.WriteStepSkill(skill.player, skill.Name)
		skill.player.AddBuff(skill.player, "裸衣", -1, true, 2, skill.Params[0], 1, 4, skill.Params[1], 0)
	case RenJie: //忍戒 所有武将每使用或触发一次技能，武技+
		all := skill.player.FindTeam(true)
		all = append(all, skill.player.FindTeam(false)...)
		for _, v := range all {
			newskill := &SkillStruct{"16", RenJieRun, skill.Params, v}
			v.SkillTempAi[RenJie] = skill.player
			v.Skill.PushBack(newskill)
		}
	case QuHu: //驱虎 本方技能伤害+
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindTeam(true) {
			v.AddBuff(skill.player, QuHu, -1, false, 9, skill.Params[0], 1)
		}
	case MouZhu: //谋主 本方谋士（阴武将）获得“技能伤害+”
		this.WriteStepSkill(skill.player, skill.Name)
		for _, v := range skill.player.FindTeam(true) {
			if v.Base.HasAttribute(5) {
				v.AddBuff(skill.player, MouZhu, -1, false, 9, skill.Params[0], 1)
			}
		}
	case QiCe: //奇策 战鼓，随机对对方全部武将共造成12点伤害
		this.WriteStepSkill(skill.player, skill.Name)
		damage := skill.Params[0]
		targetArray := skill.player.FindRandArray(false, 100)
		damageArray := make([]uint8, len(targetArray))
		for damage != 0 { //只要伤害不为0，一直循环分配伤害
			for k, _ := range damageArray { //随机循环敌方单位，随机赋值伤害
				randDamage := uint8(skill.player.Main.Rands.Int31n(int32(damage))) + 1
				if randDamage > damage {
					randDamage = damage
				}
				damageArray[k] += randDamage
				damage -= randDamage
				if damage == 0 {
					break
				}
			}
		}
		for k, _ := range targetArray { //循环敌方单位，造成伤害
			if damageArray[k] != 0 {
				targetArray[k].ReduceHp(skill.player, damageArray[k], 0, 0, 0, false, false, false, false, 2)
			}
		}
	case LuanWu: //乱舞 每次攻击对对方随机造成x~y点伤害
		this.WriteStepSkill(skill.player, skill.Name)
		damage := uint8(skill.player.Main.Rands.Int31n(15)) + 1
		skill.player.Atarget.ReduceHp(skill.player, damage, 0, 0, 0, false, false, false, true, 2)
	case WeiMu: //不会受到技能造成的伤害；对方所有遗计效果失效
		this.WriteStepSkill(skill.player, skill.Name)
		//绑定被攻击技能
		newskill := &SkillStruct{"6", WeiMu1, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskill)
		//为所有敌方武将绑定帷幕2
		for _, v := range skill.player.FindTeam(false) {
			newskill := &SkillStruct{"17", WeiMu2, skill.Params, v}
			v.Skill.PushBack(newskill)
			v.SkillTempAi[WeiMu] = skill.player
		}
	case WeiMu1: //不会受到技能造成的伤害
		if skill.player.Rcanmiss == false && skill.player.Rcanreduce == false && skill.player.Rcanreflect == false { //如果既不能闪避又不能反击也不能反伤，认定为技能攻击
			this.WriteStepSkill(skill.player, WeiMu)
			skill.player.Rdamage = 0
		}
	case FuJi: //伏计 不受群体技能的伤害
		if skill.player.Rispoint == false { //如果是群体攻击，伤害为0
			this.WriteStepSkill(skill.player, skill.Name)
			skill.player.Rdamage = 0
		}
	case GangLi: //刚戾 为他回复体力的武将，可以攻击1次
		this.WriteStepSkill(skill.player, skill.Name)
		if target := skill.player.FindRandArray(false, 1); skill.player.Hvalue > 0 && target != nil {
			skill.player.Hadder.Attack(target[0])
		}
	case QingTengJia: //轻藤甲 每回合回血 受到伤害概率减伤 火加倍
		newskillRun := &SkillStruct{"2", QingTengJiaRun, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskillRun)
		newskillDefense := &SkillStruct{"6", QingTengJiaDefense, skill.Params, skill.player}
		skill.player.Skill.PushBack(newskillDefense)
	case QingTengJiaRun: //轻藤甲每回合
		skill.player.AddHp(skill.player, skill.Params[0])
	case QingTengJiaDefense: //轻藤甲受到攻击
		if skill.player.Main.Rands.Int31n(100) < int32(skill.Params[1]) {
			if skill.player.Rdamage > 1 {
				skill.player.Rdamage -= skill.Params[2]
			} else {
				skill.player.Rdamage = 0
			}
		}
		if skill.player.Rattacker.Base.HasAttribute(3) { //如果攻击者是火属性，则触发火烧藤甲
			this.WriteStepSkill(skill.player.Rattacker, "火烧藤甲")
			skill.player.Rdamage *= skill.Params[3]
		}
	case ShiDun: //石盾 受到伤害减免
		skill.player.Rdamage = uint8(uint16(skill.player.Rdamage) * uint16(100-skill.Params[0]) / 100)
	case PoJia: //破甲
		skill.player.Atarget.AddBuff(skill.player, "破甲", -1, true, 4, skill.Params[0], 0)
	}
}

func (this *AllSkill) WriteStepSkill(_ai *Ai, name string) {
	_ai.Main.AllStep += "3:" + strconv.Itoa(int(_ai.PositionId)) + "," + name + "|"
	//循环技能,只调用技能触发后释放的技能
	for e := _ai.Skill.Front(); e != nil; e = e.Next() {
		if e.Value.(*SkillStruct).Where == "16" {
			//指针thisSkill 代替 指针e.Value.(*SkillStruct)
			thisSkill := e.Value.(*SkillStruct)
			switch thisSkill.Name {
			case RenJieRun: // 忍戒发动
				if _ai.SkillTempAi[RenJie].Alive != false {
					_ai.Main.AllStep += "3:" + strconv.Itoa(int(_ai.SkillTempAi[RenJie].PositionId)) + "," + RenJie + "|"
					_ai.AddBuff(_ai.SkillTempAi[RenJie], RenJie, -1, true, 2, thisSkill.Params[0], 1)
				}
			}
		}
	}
}
