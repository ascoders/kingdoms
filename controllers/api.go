package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"kingdoms/models"
	"regexp"
	"strconv"
	"strings"
)

type ApiController struct {
	beego.Controller
}

func (this *ApiController) Prepare() {
	//如果是部署模式，防止其他站点使用此api
	if beego.RunMode == "prod" && this.Ctx.Input.Domain() != beego.AppConfig.String("webSite") && this.Ctx.Input.Domain() != beego.AppConfig.String("httpWebSite") {
		this.Ctx.Redirect(302, "http://"+beego.AppConfig.String("httpWebSite"))
		this.StopRun()
	}
	//为防止伪造preference，验证toekn
	if this.GetString("token") != beego.AppConfig.String("KingdomsToken") {
		this.StopRun()
	}
}

/* 测试战斗 */
func (this *ApiController) Test() {
	//查询全部卡牌
	card := models.Card{}
	cards := card.FindAll()
	this.Data["json"] = cards
	this.ServeJson()
}

/* 测试战斗提交表单 */
func (this *ApiController) TestPost() {
	//测试提交的表单
	valid := validation.Validation{}
	valid.Required(this.GetString("my"), "1")
	valid.Match(this.GetString("my"), regexp.MustCompile("^[,0-9]+$"), "2")
	valid.Required(this.GetString("enemie"), "3")
	valid.Match(this.GetString("enemie"), regexp.MustCompile("^[,0-9]+$"), "4")
	if valid.HasErrors() { //没有通过验证则退出
		return
	}
	//解析提交的敌我方数据
	myArray := strings.Split(this.GetString("my"), ",")
	enemieArray := strings.Split(this.GetString("enemie"), ",")
	mySlice := make([]string, len(myArray))
	enemieSlice := make([]string, len(enemieArray))
	for k, v := range myArray {
		mySlice[k] = v + ",0,0,0,0,0"
	}
	for k, v := range enemieArray {
		enemieSlice[k] = v + ",0,0,0,0,0"
	}
	//模拟战斗
	index := IndexController{}
	result := index.SimulateFight(mySlice, "0;0", enemieSlice, "0;0")
	this.Data["json"] = result
	this.ServeJson()
}

/* 每个卡牌属性页 */
func (this *ApiController) ShowCard() {
	//查询卡牌
	id, _ := strconv.Atoi(this.GetString("id"))
	this.Data["json"] = cards[uint16(id)]
	this.ServeJson()
}

/* 更新版本信息 */
func (this *ApiController) UpdateVersion() {
	version := &models.Version{}
	switch this.GetString("type") {
	case "card":
		version.UpdateVersion("c")
	case "level":
		version.UpdateVersion("l")
	case "image":
		updateImage()
		version.UpdateVersion("i")
	}
}

/* 新增武将post */
func (this *ApiController) AddCard() {
	//登陆提交的表单
	card := &models.Card{}
	//数据采集
	if err := this.ParseForm(card); err != nil {
		this.Data["json"] = -1
		this.ServeJson()
		this.StopRun()
	}
	//数据验证
	valid := validation.Validation{}
	b, _ := valid.Valid(card)
	if !b { //验证出错，停止
		this.Data["json"] = -2
		this.ServeJson()
		this.StopRun()
	}
	//卡牌id是否存在
	if card.IdExist(card.Id) {
		this.Data["json"] = -3
		this.ServeJson()
		this.StopRun()
	}
	//插入新卡牌
	card.Insert()
	SetAllCards()
	//输出
	this.Data["json"] = 1
	this.ServeJson()
}

/* 获取卡牌总数 */
func (this *ApiController) CardCount() {
	card := &models.Card{}
	this.Data["json"] = card.Count()
	this.ServeJson()
}

/* 显示卡牌列表 */
func (this *ApiController) CardList() {
	card := &models.Card{}
	page, _ := this.GetInt("page")
	result := card.Find(int(page-1)*10, 10)
	this.Data["json"] = result
	this.ServeJson()
}

/* 更新卡牌信息 */
func (this *ApiController) UpdateCard() {
	card := &models.Card{}
	//数据采集
	if err := this.ParseForm(card); err != nil {
		this.StopRun()
	}
	//数据验证
	valid := validation.Validation{}
	b, err := valid.Valid(card)
	if err != nil {
		this.StopRun()
	}
	if !b { //验证出错，停止
		this.StopRun()
	}
	//更新卡牌
	card.Update(card.Id)
	SetAllCards()
}

/* 删除某个卡牌 */
func (this *ApiController) DeleteCard() {
	card := &models.Card{}
	card.Delete(this.GetString("id"))
	SetAllCards()
}

/* 新增关卡post */
func (this *ApiController) AddLevel() {
	//登陆提交的表单
	level := &models.Level{}
	//数据采集
	if err := this.ParseForm(level); err != nil {
		this.Data["json"] = -1
		this.ServeJson()
		this.StopRun()
	}
	//数据验证
	valid := validation.Validation{}
	b, err := valid.Valid(level)
	if err != nil {
		// handle error
	}
	if !b { //验证出错，停止
		this.Data["json"] = -2
		this.ServeJson()
		this.StopRun()
	}
	//卡牌id是否存在
	if level.IdExist(level.Id) {
		this.Data["json"] = -3
		this.ServeJson()
		this.StopRun()
	}
	//插入新关卡
	level.Insert()
	//输出
	this.Data["json"] = 1
	this.ServeJson()
}

/* 获取关卡总数 */
func (this *ApiController) LevelCount() {
	level := &models.Level{}
	this.Data["json"] = level.Count()
	this.ServeJson()
}

/* 显示关卡列表 */
func (this *ApiController) LevelList() {
	level := &models.Level{}
	page, _ := this.GetInt("page")
	result := level.Find(int(page-1)*10, 10)
	this.Data["json"] = result
	this.ServeJson()
}

/* 更新关卡信息 */
func (this *ApiController) UpdateLevel() {
	level := &models.Level{}
	//数据采集
	if err := this.ParseForm(level); err != nil {
		this.StopRun()
	}
	//数据验证
	valid := validation.Validation{}
	b, err := valid.Valid(level)
	if err != nil {
		this.StopRun()
	}
	if !b { //验证出错，停止
		this.StopRun()
	}
	//更新卡牌
	level.Update(level.Id)
	SetAllCards()
}

/* 删除某个关卡 */
func (this *ApiController) DeleteLevel() {
	level := &models.Level{}
	level.Delete(this.GetString("id"))
}
