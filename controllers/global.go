package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"kingdoms/models"
	"strconv"
)

type GlobalController struct {
	beego.Controller
}

var (
	Bm cache.Cache //文件缓存
	Mc cache.Cache //内存缓存
)

func init() {
	//初始化文件缓存
	Bm, _ = cache.NewCache("file", `{"CachePath":"/cache","FileSuffix":".cache","DirectoryLevel":2,"EmbedExpiry":0}`)
	//初始化内存缓存
	Mc, _ = cache.NewCache("memory", `{"interval":600}`) //内存10分钟回收一次
}

/* ----------工具方法---------- */

/* 将[]string 转化为 []int */
func StringToIntArray(stringArray []string) []int {
	result := make([]int, len(stringArray))
	for k, _ := range stringArray {
		number, err := strconv.Atoi(stringArray[k])
		if err != nil { //转换错误
			return []int{}
		}
		result[k] = number
	}
	return result
}

/* md5简化用法 */
func MD5(text string) string {
	m := md5.New()
	m.Write([]byte(text))
	return hex.EncodeToString(m.Sum(nil))
}

/* 查找是否有缓存
 * 先从内存中查找，如果没有，再从文件中查找
 */
func FindCache(key string) interface{} {
	if beego.RunMode != "prod" { //不是部署模式不会查找缓存
		return nil
	}
	var value interface{}
	//从内存中查找
	if value = Mc.Get(key); value != nil && value != "" {

	} else if value = Bm.Get(key); value != nil && value != "" {
		//重新设置到内存中
		Mc.Put(key, value, 300)
	} else {
		return nil
	}
	return value
}

/* 设置缓存
 * 将内容设置到内存缓存保存5分钟，同时设置到文件缓存保存1小时
 */
func SetCache(key string, value interface{}) {
	if beego.RunMode != "prod" { //不是部署模式不会设置缓存
		return
	}
	//设置内存缓存
	Mc.Put(key, value, 300)
	//设置文件缓存
	Bm.Put(key, value, 3600)
}

/* 删除缓存 同时删除内存和文件的缓存 */
func DeleteCache(key string) {
	Mc.Delete(key)
	Bm.Delete(key)
}

/* 某IP是否操作过于频繁 */
func IsBusy(name string, ip string, time int64) bool {
	if cac := Mc.Get(name + ip); cac != nil && cac != "" { //过于频繁
		return true
	} else {
		Mc.Put(name+ip, 1, time) //冷却时间
		return false
	}
}

/* 重新赋值所有卡牌 */
func SetAllCards() {
	//赋值所有卡牌
	card := &models.Card{}
	allCards := card.FindAllInfo()
	cards = make(map[uint16]models.Card, len(allCards))
	for k, _ := range allCards {
		cards[allCards[k].Id] = *allCards[k]
	}
}
