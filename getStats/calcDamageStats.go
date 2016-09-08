package getStats

import (
	"github.com/dotabuff/manta/dota"
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
			isAloneKill := true
			//单杀的attackername
			damageSourceName := uint32(0)
			replayData.teamDeath[deadlyDamagelog.GetTargetTeam()]++
			isAloneCatch := true
			for _, aDamagelog := range replayData.allDamageLogs {
				if isHeroToOpponentHeroCombatLog(aDamagelog) && isDamagelogCount(deadlyDamagelog, aDamagelog) {
					allHeroStats[aDamagelog.GetDamageSourceName()].CreateDeadlyDamages += aDamagelog.GetValue()
					if damageSourceName == 0 {
						//如果是英雄伤害，记录attacker
						damageSourceName = aDamagelog.GetDamageSourceName()
					} else if damageSourceName != aDamagelog.GetDamageSourceName() {
						//如果后来有伤害记录和之前的attacker不一样，表示不是单杀
						isAloneKill = false
					}
					//如果判断标志true， 继续判断这一条记录，如果false 说明不是单抓，不再记录
					if isAloneCatch {
						isAloneCatch = isAloneCatched(aDamagelog, replayData)
					}

				}
			}

			if isAloneCatch && deadlyDamagelog.GetAttackerTeam() != 4 {
				Clog("%v : %v is alone catched", timeStampToString(deadlyDamagelog.GetTimestamp() - replayData.gameStartTime), allHeroStats[deadlyDamagelog.GetTargetName()].HeroName)
				allHeroStats[deadlyDamagelog.GetTargetName()].AloneBeCatchedNum++
			}

			//记录单杀次数， 判断条件：助攻人数小于等于1， 不是被野怪杀死
			if len(deadlyDamagelog.AssistPlayers) == 1 && deadlyDamagelog.GetAttackerTeam() != 4 {
				if _, exist := allHeroStats[damageSourceName]; exist && isAloneKill {
					Clog("%v killed %v alone at %v", allHeroStats[damageSourceName].HeroName, allHeroStats[deadlyDamagelog.GetTargetName()].HeroName, timeStampToString(deadlyDamagelog.GetTimestamp() - replayData.gameStartTime))
					allHeroStats[damageSourceName].AloneKilledNum++
				}
				allHeroStats[deadlyDamagelog.GetTargetName()].AloneBeKilledNum++
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
	if aDamagelogTimeStamp <= deadlyDamagelogTimeStamp && aDamagelogTimeStamp >= deadlyDamagelogTimeStamp - 17.0 && aDamagelog.GetTargetName() == deadlyDamagelog.GetTargetName() {
		return true
	}
	return false
}

//获取敌方死亡次数
func calcTeamDeath(replayData *ReplayData) {
	for _, heroStates := range allHeroStats {
		teamNumber := findTeamNumberFromSteamId(heroStates.Steamid, replayData)
		for opponentTeamNumber, deathNumber := range replayData.teamDeath {
			//teamnumber和选手teamnumber不同的时猴，则为敌方死亡次数
			if opponentTeamNumber != uint32(teamNumber) {
				heroStates.OpponentHeroDeaths = deathNumber
			}
		}
		heroStates.CreateDeadlyDamagesPerDeath = float32(heroStates.CreateDeadlyDamages) / float32(heroStates.OpponentHeroDeaths)
		heroStates.CreateDeadlyStiffControlPerDeath = float32(heroStates.CreateDeadlyStiffControl) / float32(heroStates.OpponentHeroDeaths)

		Clog("player: %v, hero : %v, opponentdeath : %v, damage perdeath : %v, control per death : %v", heroStates.PlayerName, heroStates.HeroName, heroStates.OpponentHeroDeaths, heroStates.CreateDeadlyDamagesPerDeath, heroStates.CreateDeadlyStiffControlPerDeath)
	}
}

func findTeamNumberFromSteamId(steamid uint64, replayData *ReplayData) int32 {
	for _, aPlayInfo := range replayData.dotaGameInfo.GetPlayerInfo() {
		if aPlayInfo.GetSteamid() == steamid {
			return aPlayInfo.GetGameTeam()
		}
	}
	return 0
}
