package main

import (
	"log"
	"new_stats/getStats"
)

func main() {
	allHeroStats, err := getStats.GetStats("C:/2545034458.dem")
	if err != nil {
		log.Fatalf("解析录像失败: %s", err)
	}

	log.Printf("英雄对敌方英雄造成的伤害统计：\n")
	for _, v := range allHeroStats {
		log.Printf("%s(Steamid=%d)——总伤害：%d，致死伤害：%d\n", v.HeroName, v.Steamid, v.CreateTotalDamages, v.CreateDeadlyDamages)
	}
}
