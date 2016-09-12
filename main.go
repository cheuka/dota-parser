package main

import (
	"compress/bzip2"
	"io"
	"io/ioutil"
	"log"
	"new_stats/dota2"
	"new_stats/getStats"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {

	demFileName := decompressBzip2ToDemFile("C:/2545299883.dem.bz2")
	textAGame(demFileName)
	//textAGame("C:/TI6/2545101126.dem")

	//writeToDB("root:123456@/dota2_new_stats?charset=utf8&parseTime=True&loc=Local", "C:/TI6/")
	//writeToDB("root:123456@/dota2_new_stats_for_cn?charset=utf8&parseTime=True&loc=Local", "D:/replays/")

}

func decompressBzip2ToDemFile(bz2FileName string) string {
	var demFileName string
	demFileName = strings.TrimSuffix(bz2FileName, ".bz2")

	bizFile, err := os.Open(bz2FileName)
	if err != nil {
		log.Fatalf("打开录像压缩文件失败: %s", err)
	}
	defer bizFile.Close()

	demFile, err := os.Create(demFileName)
	if err != nil {
		log.Fatalf("创建录像文件失败: %s", err)
	}

	bzip2Reader := bzip2.NewReader(bizFile)
	if err != nil {
		log.Fatalf("解压录像文件失败: %s", err)
	}

	io.Copy(demFile, bzip2Reader)
	return demFileName

}

func textAGame(fileName string) {
	allHeroStats, err := getStats.GetStats(fileName)
	if err != nil {
		log.Fatalf("解析录像失败: %s", err)
	}

	log.Printf("英雄对敌方英雄造成的伤害统计：\n")
	for _, v := range allHeroStats {
		log.Printf("%s(Steamid=%d)——总伤害：%d，致死伤害：%d，致死控制时间:%.2f, 单杀/被单杀/被单抓次数: %d/%d/%d, \n", v.HeroName, v.Steamid, v.CreateTotalDamages, v.CreateDeadlyDamages, v.CreateDeadlyStiffControl, v.AloneKilledNum, v.AloneBeKilledNum, v.AloneBeCatchedNum)
	}
}

func writeToDB(dbPath, replayDir string) {
	getStats.SetDebug(false)
	db, err := gorm.Open("mysql", dbPath)
	if err != nil {
		log.Printf("failed to connect database\n")
	}
	defer db.Close()

	db.AutoMigrate(&dota2.Stats{})

	dir, err := ioutil.ReadDir(replayDir)
	if err != nil {
		log.Printf("failed to open dir\n")
	}

	for i, aFile := range dir {
		aRepaly := replayDir + aFile.Name()
		matchID, _ := strconv.ParseUint(strings.TrimSuffix(aFile.Name(), ".dem"), 10, 64)
		log.Printf("正在解析第%d个录像：%d", i+1, matchID)
		allHeroStats, err := getStats.GetStats(aRepaly)
		if err != nil {
			log.Fatalf("解析录像失败: %s", err)
		}
		//写结果到数据库
		for _, aHeroStats := range allHeroStats {
			aHeroStats.MatchId = matchID
			db.Create(aHeroStats)
		}

	}

}
