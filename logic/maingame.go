package logic

import (
	"math/rand"
)

type Maingame struct {
	MyTeam     []*Ai      //我方队伍指针
	EnemieTeam []*Ai      //敌方队伍指针
	Time       float32    //一秒的时间
	AllStep    string     //记录整个战斗过程
	Rands      *rand.Rand //随机数
}
