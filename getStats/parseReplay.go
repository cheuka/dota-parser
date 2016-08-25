package getStats

import (
	"fmt"
	"new_stats/dota2"
	"os"
	"strings"

	"github.com/dotabuff/manta"
	"github.com/dotabuff/manta/dota"
	"log"
	"sort"
)
//package entity note 1.cdotaplayer.playerId对应 CDOTA_Unit_Hero_中的playerId, playId按照楼层排序
//2016/08/25 17:31:06 Properties, m_vecPlayerTeamData.0006.m_hSelectedHero : 15499785
//2016/08/25 17:31:06 Properties, m_vecPlayerTeamData.0006.m_iLevel : 1
//2016/08/25 17:31:06 Properties, m_vecPlayerTeamData.0006.m_iRespawnSeconds : -1
//2016/08/25 17:31:06 Properties, m_vecPlayerTeamData.0006.m_nSelectedHeroID : 84
//m_hSelectedHero对应 CDOTA_Unit_Hero_Ogre_Magi中的m_hInventoryParent or m_hModifierParent
// CDOTA_PlayerResource heroId 对应

var SPECIAL_MODIFIERS = []string{"modifier_axe_berserkers_call"}

func parseReplay(filename string, replayData *ReplayData) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("打开比赛录像失败: %s", err)
	}
	defer f.Close()
	parser, err := manta.NewStreamParser(f)
	if err != nil {
		return fmt.Errorf("初始化解析器失败: %s", err)
	}
	parser.Callbacks.OnCMsgDOTACombatLogEntry(func(m *dota.CMsgDOTACombatLogEntry) error {
		logType := m.GetType()
		switch logType {
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DAMAGE:
			replayData.allDamageLogs = append(replayData.allDamageLogs, m)
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_REMOVE:
			replayData.allModifierLogs = append(replayData.allModifierLogs, m)
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_GAME_STATE:
			if m.GetValue() == uint32(5) {
				replayData.gameStartTime = m.GetTimestamp()
			}
		}
		return nil
	})
	parser.Callbacks.OnCDemoFileInfo(func(m *dota.CDemoFileInfo) error {
		replayData.dotaGameInfo = m.GameInfo.Dota
		log.Printf(m.String())
		return nil
	})

	parser.OnPacketEntity(func(entity *manta.PacketEntity, pet manta.EntityEventType) error {

		if strings.Contains(entity.ClassName, "CDOTA_PlayerResource") {
			log.Printf("EntityEvent : %v, %v", entity.ClassName, pet)
			printProperties("ClassBaseline", entity.ClassBaseline)
			printProperties("Properties", entity.Properties)
			log.Printf("\n\n")
		}
		//for k, v := range entity.ClassBaseline.KV{
		//	if strings.Contains(k, "pick") || strings.Contains(k, "ban"){
		//		log.Printf("EntityEvent : %v, %v, %v, %v", entity.ClassName, pet, k, v)
		//	}
		//}
		return nil
	})
	parser.Start()                       //开始解析录像
	initAllHeroStats(parser, replayData) //初始化initAllHeroStats
	return nil                           //解析完成，返回数据
}

//初始化AllHeroStats：
//[1]找出所有英雄在combatLog中的index
//[2]找出十名英雄使用者的SteamId
func initAllHeroStats(parser *manta.Parser, replayData *ReplayData) error {
	allHeroStats = make(map[uint32]*dota2.Stats)
	replayData.specialModifier = make(map[int32]*string)
	index := int32(0)
	for {
		name, has := parser.LookupStringByIndex("CombatLogNames", index)
		//假设index在CombatLogNames中是没有间隔的，遍历CombatLogNames
		if !has {
			break
		}
		if strings.Contains(name, "npc_dota_hero_") {
			allHeroStats[uint32(index)] = &dota2.Stats{HeroName: strings.TrimPrefix(name, "npc_dota_hero_")}
		}
		//获取特殊控制技能
		for _, modifier := range SPECIAL_MODIFIERS {
			if strings.EqualFold(modifier, name) {
				replayData.specialModifier[index] = &modifier
				break
			}
		}
		index = index + 1
	}
	if len(allHeroStats) != 10 {
		return fmt.Errorf("无法从combatLog中找到十个英雄的index")
	}
	for _, aPlayInfo := range replayData.dotaGameInfo.GetPlayerInfo() {
		for _, aHeroStats := range allHeroStats {
			if strings.Contains(aPlayInfo.GetHeroName(), aHeroStats.HeroName) {
				aHeroStats.Steamid = aPlayInfo.GetSteamid()
				aHeroStats.MatchId = replayData.dotaGameInfo.GetMatchId()
			}
		}
	}
	return nil
}

func printProperties(tag string, ppt *manta.Properties) {
	sorted_keys := make([]string, 0)
	for k, _ := range ppt.KV {
		sorted_keys = append(sorted_keys, k)
	}

	// sort 'string' key in increasing order
	sort.Strings(sorted_keys)

	for _, k := range sorted_keys {
		log.Printf("%v, %v : %v\n", tag, k, ppt.KV[k])
	}
}
