package main

import (
	"compress/bzip2"
	"io"
	"log"
	"github.com/cheuka/dota-parser/getStats"
	"os"
	"strings"
	"encoding/json"
	"bufio"
)

func main() {
	getStats.SetDebug(false)
	//demFileName := decompressBzip2ToDemFile("C:/2545299883.dem.bz2")
	//textAGame(demFileName)

	//本地读取
	//bizFile, _ := os.Open("D://2545254388.dem")
	//defer bizFile.Close()
	//textAGame(bizFile)

	//stdin读取
	f := bufio.NewReader(os.Stdin)
	textAGame(f)
	//defer f.Close()
	//writeToDB("root:123456@/dota2_new_stats?charset=utf8&parseTime=True&loc=Local", "C:/TI6/")
	//writeToDB("root:123456@/dota2_new_stats_for_cn?charset=utf8&parseTime=True&loc=Local", "D:/replays/")

}

//func decompressBzip2ToInputStream(bz2FileName string){
//	bizFile, err := os.Open(bz2FileName)
//	if err != nil {
//		log.Fatalf("打开录像压缩文件失败: %s", err)
//	}
//	defer bizFile.Close()
//
//	bzip2Reader := bzip2.NewReader(bizFile)
//	if err != nil {
//		log.Fatalf("解压录像文件失败: %s", err)
//	}
//
//	io.Copy(demFile, bzip2Reader)
//}

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

func textAGame(r io.Reader) {
	allHeroStats, err := getStats.GetStats(r)
	if err != nil {
		log.Fatalf("解析录像失败: %s", err)
	}

	//log.Printf("英雄对敌方英雄造成的伤害统计：\n")
	//for _, v := range allHeroStats {
	//	//	log.Printf("%s(Steamid=%d)——总伤害：%d，致死伤害：%d，致死控制时间:%.2f, 单杀/被单杀/被单抓次数: %d/%d/%d, win : %v\n", v.HeroName, v.Steamid, v.CreateTotalDamages, v.CreateDeadlyDamages, v.CreateDeadlyStiffControl, v.AloneKilledNum, v.AloneBeKilledNum, v.AloneBeCatchedNum, v.IsWin)
	//	//}
	//	b, err := json.Marshal(v)
	//	if err == nil{
	//		fmt.Printf("%v \n", string(b))
	//	}
	//
	//}
	getStats.Clog("length of map : %d", len(allHeroStats))
	for _, v := range allHeroStats {
		getStats.Clog("%v, ggcount first %v, write %v, reply %v, healing %v, wardsBuy, %v, wardKill %v, runes %v, apm, %v", v.PlayerName, v.IsFirstGG, v.IsWriteGG, v.IsReplyGG, v.Healing, v.WardsBuy, v.WardsKill, v.RuneCount, v.Apm)
	}
	b, err := json.Marshal(allHeroStats)
	if err == nil {
		//fmt.Printf("%v", string(b))
		os.Stdout.Write(b)
	}
}

