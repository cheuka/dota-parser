package getStats

import (
	"github.com/dotabuff/manta/dota"
	"log"
)


//统计每个英雄的造成的总输出： CreateTotalDamages
//有BUG：冰魂给队友套BUFF之后，队友平A造成的伤害，算在队友身上，而不是冰魂（和dotabuff不一致）
func calcCreateTotalDamages(replayData *ReplayData) {
	for _, aDamageLog := range replayData.allDamageLogs {
		if isHeroToOpponentHeroCombatLog(aDamageLog) {
			allHeroStats[aDamageLog.GetDamageSourceName()].CreateTotalDamages += aDamageLog.GetValue()
		}
	}
}

//统计每个英雄的造成的致死输出：CreateDeadlyDamages
//一条death记录对应了一条致死deadlyDamagelog记录（该记录的health=0，timestamp与death记录保持一致）
//death次数=致死damage次数=肉山盾复活英雄次数+英雄死亡次数（和KDA中D的总和相等）=肉山盾复活英雄次数+英雄击杀英雄次数（和KDA中A的总和相等）+非英雄单位（防御塔等）击杀英雄次数
//举例2562582896( totalKD=87,89): death次数(93)=肉山盾复活英雄次数（4）+英雄击杀英雄次数（87）+非英雄单位（防御塔）击杀英雄次数（2）
func calcCreateDeadlyDamages(replayData *ReplayData) {
	for _, deadlyDamagelog := range replayData.allDamageLogs {
		if deadlyDamagelog.GetHealth() == 0 && isToOpponentHeroCombatLog(deadlyDamagelog) {
			replayData.teamDeath[deadlyDamagelog.GetTargetTeam()]++
			for _, aDamagelog := range replayData.allDamageLogs {
				if isHeroToOpponentHeroCombatLog(aDamagelog) && isDamagelogCount(deadlyDamagelog, aDamagelog) {
					allHeroStats[aDamagelog.GetDamageSourceName()].CreateDeadlyDamages += aDamagelog.GetValue()
				}
			}
		}
	}
}

//判断aDamagelog是否应该计入deadlyDamagelog表示的这次击杀
//注意参数的顺序：deadlyDamagelog是致死伤害记录（包含了最全的GetAssistPlayers信息）
//暂定计入英雄死亡前17秒(冰魂A杖大持续时间)内受到的所有伤害[没有找到任何的官方文档]
func isDamagelogCount(deadlyDamagelog, aDamagelog *dota.CMsgDOTACombatLogEntry) bool {
	aDamagelogTimeStamp := aDamagelog.GetTimestamp()
	deadlyDamagelogTimeStamp := deadlyDamagelog.GetTimestamp()
	if aDamagelogTimeStamp <= deadlyDamagelogTimeStamp && aDamagelogTimeStamp >= deadlyDamagelogTimeStamp-17.0 && aDamagelog.GetTargetName() == deadlyDamagelog.GetTargetName() {
		return true
	}
	return false
}

//获取敌方死亡次数
func calcTeamDeath(replayData *ReplayData){
	for _, heroStates := range allHeroStats{
		teamNumber := findTeamNumberFromSteamId(heroStates.Steamid, replayData)
		for opponentTeamNumber, deathNumber := range replayData.teamDeath{
			//teamnumber和选手teamnumber不同的时猴，则为敌方死亡次数
			if opponentTeamNumber != uint32(teamNumber){
				heroStates.OpponentHeroDeaths = deathNumber
			}
		}
		heroStates.CreateDeadlyDamagesPerDeath = float32(heroStates.CreateDeadlyDamages) / float32(heroStates.OpponentHeroDeaths)
		heroStates.CreateDeadlyStiffControlPerDeath = float32(heroStates.CreateDeadlyStiffControl) / float32(heroStates.OpponentHeroDeaths)


		log.Printf("player: %v, hero : %v, opponentdeath : %v, damage perdeath : %v, control per death : %v", heroStates.PlayerName, heroStates.HeroName, heroStates.OpponentHeroDeaths, heroStates.CreateDeadlyDamagesPerDeath, heroStates.CreateDeadlyStiffControlPerDeath)
	}
}

func findTeamNumberFromSteamId(steamid uint64, replayData *ReplayData) int32{
	for _, aPlayInfo := range replayData.dotaGameInfo.GetPlayerInfo() {
		if(aPlayInfo.GetSteamid() == steamid){
			return aPlayInfo.GetGameTeam()
		}
	}
	return 0
}