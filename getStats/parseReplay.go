package getStats

import (
	"fmt"
	"../dota2"
	"strings"

	"reflect"
	"sort"

	"github.com/dotabuff/manta"
	"github.com/dotabuff/manta/dota"
	"io"
	"log"
)

//package entity note 1.cdotaplayer.playerId对应 CDOTA_Unit_Hero_中的playerId, playId按照楼层排序
//2016/08/25 17:31:06 Properties, m_vecPlayerTeamData.0006.m_hSelectedHero : 15499785
//2016/08/25 17:31:06 Properties, m_vecPlayerTeamData.0006.m_iLevel : 1
//2016/08/25 17:31:06 Properties, m_vecPlayerTeamData.0006.m_iRespawnSeconds : -1
//2016/08/25 17:31:06 Properties, m_vecPlayerTeamData.0006.m_nSelectedHeroID : 84
//m_hSelectedHero对应 CDOTA_Unit_Hero_Ogre_Magi中的m_hInventoryParent or m_hModifierParent
// CDOTA_PlayerResource heroId 对应

//modifier_shadow_demon_disruption 毒狗的关 是否加进去还需判断目标是否同一个team
var SPECIAL_MODIFIERS = []string{"modifier_axe_berserkers_call"}

//playResourceEntity : m_vecPlayerTeamData
//0001
//2016/08/26 11:30:53 Properties, m_vecPlayerTeamData.0000.m_flTeamFightParticipation : 0.25
//2016/08/26 11:30:53 Properties, m_vecPlayerTeamData.0000.m_hSelectedHero : 10142214
//2016/08/26 11:30:53 Properties, m_vecPlayerTeamData.0000.m_iAssists : 2
//2016/08/26 11:30:53 Properties, m_vecPlayerTeamData.0000.m_iDeaths : 1
//2016/08/26 11:30:53 Properties, m_vecPlayerTeamData.0000.m_iLevel : 7
//2016/08/26 11:30:53 Properties, m_vecPlayerTeamData.0000.m_iRespawnSeconds : 27
//2016/08/26 11:30:53 Properties, m_vecPlayerTeamData.0000.m_nSelectedHeroID : 68
//2016/08/26 11:30:53 Properties, m_vecPlayerTeamData.0001.m_flTeamFightParticipation : 0.75

//2016/08/26 11:38:09 ClassBaseline, m_vecPlayerData.0000.m_iPlayerSteamID : 76561198046993283
//2016/08/26 11:38:09 ClassBaseline, m_vecPlayerData.0000.m_iPlayerTeam : 2
var playResourceEntity *manta.PacketEntity
var gameTime float32//比赛时间
func parseReplay(r io.Reader, replayData *ReplayData) error {
	parser, err := manta.NewStreamParser(r)
	if err != nil {
		return fmt.Errorf("初始化解析器失败: %s", err)
	}
	parser.Callbacks.OnCMsgDOTACombatLogEntry(func(m *dota.CMsgDOTACombatLogEntry) error {
		logType := m.GetType()
		switch logType {
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DAMAGE:
			replayData.allDamageLogs = append(replayData.allDamageLogs, m)
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_REMOVE:
			//printModifer(m, parser, replayData)
			replayData.allModifierLogs = append(replayData.allModifierLogs, m)
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_ADD:
			replayData.allModifierLogs = append(replayData.allModifierLogs, m)
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_GOLD:
			//加入金钱记录的功能
			replayData.allGoldLogs = append(replayData.allGoldLogs, m)
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_GAME_STATE:
			if m.GetValue() == uint32(5) {
				replayData.gameStartTime = m.GetTimestamp()
			} else if m.GetValue() == uint32(6) {
				replayData.gameEndTime = m.GetTimestamp()
			}
		}
		return nil
	})
	parser.Callbacks.OnCDemoFileInfo(func(m *dota.CDemoFileInfo) error {
		log.Printf("OnCDemoFileInfo")
		replayData.dotaGameInfo = m.GameInfo.Dota
		Clog(m.String())
		return nil
	})

	//记录GG
	//根据usermessage获取eneityIndex,根据Index获取playerId,根据playerId获取聊天人
	parser.Callbacks.OnCUserMessageSayText2(func(m *dota.CUserMessageSayText2) error {
		text := m.GetParam2()
		if strings.EqualFold("gg", strings.ToLower(text)) {
			if v,exist := parser.PacketEntities[int32(m.GetEntityindex())].FetchInt32("m_iPlayerID"); exist{
				replayData.ggCount[gameTime]=int32(v)
				Clog("said gg time : %v", replayData.ggCount[gameTime])
			}
		}
		Clog(m.String())
		return nil
	})

	parser.OnPacketEntity(func(entity *manta.PacketEntity, pet manta.EntityEventType) error {

		if strings.Contains(entity.ClassName, "CDOTA_PlayerResource") {
			//Clog("EntityEvent : %v, %v", entity.ClassName, pet)
			//Clog("ClassBaseline", entity.ClassBaseline)
			//Clog("Properties", entity.Properties)
			//Clog("\n\n")
			playResourceEntity = entity
		}

		if strings.EqualFold(entity.ClassName, "CDOTAGamerulesProxy") {
			if v, exist := entity.FetchFloat32("CDOTAGamerules.m_fGameTime"); exist {
				//log.Printf("EntityEvent : %v, %v, %v, %v, %v, %v, %v", pe.ClassName, pet, pe.ClassId, pe.Index, pe.Serial, timeStampToString(float32(parser.NetTick / 30) - gameStartTime), gameStartTime)
				//printProperties("properties : ", pe.Properties)
				gameTime = v;
			}
		}
		//for k, v := range entity.ClassBaseline.KV{
		//	if strings.Contains(k, "pick") || strings.Contains(k, "ban"){
		//		Clog("EntityEvent : %v, %v, %v, %v", entity.ClassName, pet, k, v)
		//	}
		//}
		recordHeroPosition(parser, entity, pet, replayData)

		return nil
	})
	log.Printf("开始解析")
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
	winTeam := replayData.dotaGameInfo.GetGameWinner()
	for _, aPlayInfo := range replayData.dotaGameInfo.GetPlayerInfo() {
		for combatLogName, aHeroStats := range allHeroStats {
			if strings.Contains(aPlayInfo.GetHeroName(), aHeroStats.HeroName) {
				aHeroStats.Steamid = aPlayInfo.GetSteamid()
				aHeroStats.MatchId = replayData.dotaGameInfo.GetMatchId()
				aHeroStats.PlayerName = aPlayInfo.GetPlayerName()
				aHeroStats.TeamNumber = aPlayInfo.GetGameTeam()
				if aPlayInfo.GetGameTeam() == winTeam{
					aHeroStats.IsWin = true
				}else{
					aHeroStats.IsWin = false
				}
				getHeroIdFromSteamId(combatLogName, replayData, aHeroStats, playResourceEntity)
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
		Clog("%v, %v : %v, %v\n", tag, k, ppt.KV[k], reflect.TypeOf(ppt.KV[k]))
	}
}

func printModifer(m *dota.CMsgDOTACombatLogEntry, p *manta.Parser, replayData *ReplayData) {
	if m.GetIsTargetHero() && m.GetAttackerName() != m.GetTargetName() && !m.GetTargetIsSelf() && !m.GetIsTargetIllusion() {
		Clog("%v , %v add %v from %v with %v", timeStampToString(m.GetTimestamp() - replayData.gameStartTime), lookForName(p, m.GetTargetName()), lookForName(p, m.GetInflictorName()), lookForName(p, m.GetAttackerName()), m.GetModifierDuration())
		Clog("%v, %v", m.GetStunDuration(), m.GetSilenceModifier())
	}
}

func lookForName(parser *manta.Parser, index uint32) string {
	str, has := parser.LookupStringByIndex("CombatLogNames", int32(index))
	if has {
		return str
	}
	return ""
}

//2016/08/26 11:30:53 Properties, m_vecPlayerTeamData.0000.m_hSelectedHero : 10142214
//2016/08/26 11:30:53 Properties, m_vecPlayerTeamData.0000.m_nSelectedHeroID : 68

//2016/08/26 11:38:09 ClassBaseline, m_vecPlayerData.0000.m_iPlayerSteamID : 76561198046993283
//2016/08/26 11:38:09 ClassBaseline, m_vecPlayerData.0000.m_iPlayerTeam : 2
//根据steamId 获取英雄ID
func getHeroIdFromSteamId(combatLogName uint32, replayData *ReplayData, aHeroStats *dota2.Stats, playResourceEntity *manta.PacketEntity) {
	steamId := aHeroStats.Steamid
	for index := 0; index < 10; index++ {
		indexStr := fmt.Sprintf("m_vecPlayerData.000%d.m_iPlayerSteamID", index)
		if v, ok := playResourceEntity.FetchUint64(indexStr); ok && v == steamId {
			if v, ok := playResourceEntity.FetchInt32(fmt.Sprintf("m_vecPlayerTeamData.000%d.m_nSelectedHeroID", index)); ok {
				aHeroStats.HeroId = uint32(v)
				aHeroStats.PlayerId = int32(index)
				if selectHero, ok := playResourceEntity.FetchUint32(fmt.Sprintf("m_vecPlayerTeamData.000%d.m_hSelectedHero", index)); ok {
					replayData.heroMap[combatLogName] = selectHero
					Clog("steamid : %v, heroId : %v, %v, %v", steamId, v, reflect.TypeOf(v), selectHero)
				}
			}

		}
	}
}

// 根据playerId来获取 hero stats实体
func getHeroStatesFromPlayerId(playerId int32) (*dota2.Stats, bool){
	for _, v := range allHeroStats{
		if v.PlayerId == playerId {
			return  v, true
		}
	}

	return nil,false
}
