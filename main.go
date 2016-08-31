package main

import (
	"io/ioutil"
	"log"
	"new_stats/dota2"
	"new_stats/getStats"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

func main() {
	getStats.SetDebug(false)
	textAGame("D://2562582896.dem")
	//writeToDB()
}

func textAGame(fileName string) {
	allHeroStats, err := getStats.GetStats(fileName)
	if err != nil {
		log.Fatalf("解析录像失败: %s", err)
	}

	log.Printf("英雄对敌方英雄造成的伤害统计：\n")
	for _, v := range allHeroStats {
		log.Printf("%s(Steamid=%d)——总伤害：%d，致死伤害：%d，致死控制时间:%.2f\n", v.HeroName, v.Steamid, v.CreateTotalDamages, v.CreateDeadlyDamages, v.CreateDeadlyStiffControl)
	}
}

func writeToDB() {
	db, err := gorm.Open("mysql", "root:123456@/dota2_new_stats?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Printf("failed to connect database\n")
	}
	defer db.Close()

	db.AutoMigrate(&dota2.Stats{})

	dir, err := ioutil.ReadDir("C:/TI6_replays/")
	if err != nil {
		log.Printf("failed to open dir\n")
	}

	for i, aFile := range dir {
		aRepaly := "C:/TI6_replays/" + aFile.Name()
		matchID, _ := strconv.ParseUint(strings.TrimSuffix(aFile.Name(), ".dem"), 10, 64)
		log.Printf("正在解析第%d个录像：%d", i+1, matchID)
		allHeroStats, err := getStats.GetStats(aRepaly)
		if err != nil {
			//写结果到数据库
			for _, aHeroStats := range allHeroStats {
				db.Create(aHeroStats)
			}
		}
	}

}
