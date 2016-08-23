package getStats

import (
	"log"
	"new_stats/dota2"
	"strings"

	"github.com/dotabuff/manta"
	"github.com/dotabuff/manta/dota"
)

//初始化：[1]找出所有英雄在combatLog中的index	[2]找出十名英雄使用者的SteamId
func initAllHeroStats(p *manta.Parser, dotaGameInfo *dota.CGameInfo_CDotaGameInfo) bool {
	allHeroStats = make(map[uint32]*dota2.Stats)
	index := int32(0)
	for {
		name, has := p.LookupStringByIndex("CombatLogNames", index)
		//假设index在CombatLogNames中是没有间隔的，遍历CombatLogNames
		if !has {
			break
		}
		if strings.Contains(name, "npc_dota_hero_") {
			allHeroStats[uint32(index)] = &dota2.Stats{HeroName: strings.TrimPrefix(name, "npc_dota_hero_")}
			if len(allHeroStats) == 10 {
				break
			}
		}
		index = index + 1
	}
	if len(allHeroStats) != 10 {
		log.Printf("无法从combatLog中找到十个英雄的index\n")
		return false
	}
	for _, aPlayInfo := range dotaGameInfo.GetPlayerInfo() {
		for _, aHeroStats := range allHeroStats {
			if strings.Contains(aPlayInfo.GetHeroName(), aHeroStats.HeroName) {
				aHeroStats.Steamid = aPlayInfo.GetSteamid()
				aHeroStats.MatchId = matchID
			}
		}
	}
	return true
}
