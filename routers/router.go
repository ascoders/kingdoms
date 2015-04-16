package routers

import (
	"github.com/astaxie/beego"
	"kingdoms/controllers"
)

func init() {
	/* 三国军争游戏入口 */
	//--------提供游戏服务 post
	//注册
	beego.Router("/register", &controllers.IndexController{}, "post:Register")
	//登陆
	beego.Router("/login", &controllers.IndexController{}, "post:Login")
	//处理用户特殊操作的枚举方法
	beego.Router("/special", &controllers.IndexController{}, "post:Special")
	//开始游戏
	beego.Router("/startgame", &controllers.IndexController{}, "post:StartGame")
	//开始战斗
	beego.Router("/startfight", &controllers.IndexController{}, "post:StartFight")
	//下载数据库
	beego.Router("/download", &controllers.IndexController{}, "post:Download")
	//下载卡牌
	beego.Router("/downloadimage", &controllers.IndexController{}, "post:DownloadImage")
	//付款成功接收post
	beego.Router("/paysuccess", &controllers.IndexController{}, "post:PaySuccess")

	//--------对外api
	//战斗测试
	beego.Router("/api/test", &controllers.ApiController{}, "post:Test")
	//战斗测试-post
	beego.Router("/api/testpost", &controllers.ApiController{}, "post:TestPost")
	//显示卡牌页面
	beego.Router("/api/showcard", &controllers.ApiController{}, "post:ShowCard")
	//更新version
	beego.Router("/api/updateversion", &controllers.ApiController{}, "post:UpdateVersion")
	//添加卡牌
	beego.Router("/api/addcard", &controllers.ApiController{}, "post:AddCard")
	//获取卡牌总数
	beego.Router("/api/cardcount", &controllers.ApiController{}, "post:CardCount")
	//显示卡牌列表
	beego.Router("/api/cardlist", &controllers.ApiController{}, "post:CardList")
	//更新卡牌
	beego.Router("/api/updatecard", &controllers.ApiController{}, "post:UpdateCard")
	//删除卡牌
	beego.Router("/api/deletecard", &controllers.ApiController{}, "post:DeleteCard")
	//添加关卡
	beego.Router("/api/addlevel", &controllers.ApiController{}, "post:AddLevel")
	//获取关卡总数
	beego.Router("/api/levelcount", &controllers.ApiController{}, "post:LevelCount")
	//显示关卡列表
	beego.Router("/api/levellist", &controllers.ApiController{}, "post:LevelList")
	//更新关卡
	beego.Router("/api/updatelevel", &controllers.ApiController{}, "post:UpdateLevel")
	//删除关卡
	beego.Router("/api/deletelevel", &controllers.ApiController{}, "post:DeleteLevel")
}
