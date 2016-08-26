package getStats

import (
	"fmt"

	"github.com/dotabuff/manta/dota"
)

func calcCreateDeadlyControl(replayData *ReplayData) {

	for _, deadlyDamagelog := range replayData.allDamageLogs {
		if deadlyDamagelog.GetHealth() == 0 && isToOpponentHeroCombatLog(deadlyDamagelog) {
			if _, isExist := allHeroStats[deadlyDamagelog.GetDamageSourceName()]; isExist {
				Clog("%v——<<<%s was killed by %s>>>\n", timeStampToString(deadlyDamagelog.GetTimestamp() - replayData.gameStartTime), allHeroStats[deadlyDamagelog.GetTargetName()].HeroName, allHeroStats[deadlyDamagelog.GetDamageSourceName()].HeroName)
			} else {
				Clog("%v——<<<%s was killed by NOT HERO>>>\n", timeStampToString(deadlyDamagelog.GetTimestamp() - replayData.gameStartTime), allHeroStats[deadlyDamagelog.GetTargetName()].HeroName)
			}
			//map保存modifier add的记录，在add和remove时如果inflictorName相同(比如晕眩)而attackName不同时，同时记录add和remove两条记录
			addModifierMap := make(map[uint32]*dota.CMsgDOTACombatLogEntry)
			for _, aModifierLog := range replayData.allModifierLogs {
				if isModifierlogCount(replayData, deadlyDamagelog, aModifierLog, addModifierMap) {
					if (aModifierLog.GetType() == dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_ADD) {
						//add 将保存的状态加入到map中去
						addModifierMap[aModifierLog.GetInflictorName()] = aModifierLog
					} else {
						if addModifierMap[aModifierLog.GetInflictorName()] != nil {
							//如果inflictorName相同，attackName不同，分别记录add和remove时
							if addModifierMap[aModifierLog.GetInflictorName()].GetAttackerName() != aModifierLog.GetAttackerName() {
								//add的时候时间为getModifierDuration
								allHeroStats[addModifierMap[aModifierLog.GetInflictorName()].GetAttackerName()].CreateDeadlyStiffControl += addModifierMap[aModifierLog.GetInflictorName()].GetModifierDuration()
							}
							//remove 移除在map中保存的状态
							delete(addModifierMap, aModifierLog.GetInflictorName())
						}
						allHeroStats[aModifierLog.GetAttackerName()].CreateDeadlyStiffControl += aModifierLog.GetModifierElapsedDuration()
					}

				}
			}
		}
	}

}

func isModifierlogCount(replayData *ReplayData, deadlyDamagelog, aModifierLog *dota.CMsgDOTACombatLogEntry, addModifierMap map[uint32]*dota.CMsgDOTACombatLogEntry) bool {
	_, isAttackerExist := allHeroStats[aModifierLog.GetAttackerName()]
	if !isAttackerExist || aModifierLog.GetIsTargetIllusion() || aModifierLog.GetAttackerName() == aModifierLog.GetTargetName() {
		return false
	}

	modifierTimeStamp := aModifierLog.GetTimestamp()
	deadlyDamagelogTimeStamp := deadlyDamagelog.GetTimestamp()
	// +1原因是控制时间有时候是在英雄死亡后结算
	if modifierTimeStamp <= deadlyDamagelogTimeStamp + float32(1) && modifierTimeStamp >= deadlyDamagelogTimeStamp - 17.0 && aModifierLog.GetTargetName() == deadlyDamagelog.GetTargetName() {
		//stun时间大于0或者是沉默技能或者在控制列表里
		if aModifierLog.GetStunDuration() > float32(0) || aModifierLog.GetSilenceModifier() || replayData.specialModifier[int32(aModifierLog.GetInflictorName())] != nil {
			if aModifierLog.GetType() == dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_ADD {

			} else {
				Clog("%v : %v removed %v from %v last %v", timeStampToString(aModifierLog.GetTimestamp() - replayData.gameStartTime), allHeroStats[aModifierLog.GetTargetName()].HeroName, aModifierLog.GetInflictorName(), allHeroStats[aModifierLog.GetAttackerName()].HeroName, aModifierLog.GetModifierElapsedDuration())

			}
			if addModifierMap[aModifierLog.GetInflictorName()] != nil {
				if addModifierMap[aModifierLog.GetInflictorName()].GetAttackerName() != aModifierLog.GetAttackerName() {
					Clog("%v : %v add %v from %v last %v", timeStampToString(addModifierMap[aModifierLog.GetInflictorName()].GetTimestamp() - replayData.gameStartTime), allHeroStats[addModifierMap[aModifierLog.GetInflictorName()].GetTargetName()].HeroName, addModifierMap[aModifierLog.GetInflictorName()].GetInflictorName(), allHeroStats[addModifierMap[aModifierLog.GetInflictorName()].GetAttackerName()].HeroName, addModifierMap[aModifierLog.GetInflictorName()].GetModifierDuration())
				}

			}
			return true
		}
	}
	return false
}

//时间戳转化成游戏里的时间
func timeStampToString(stamp float32) string {
	str := "[%02d : %02d : %02d]"
	minute := int32(stamp) / 60
	second := int32(stamp) % 60
	last := int32(stamp * 60) % 60
	str = fmt.Sprintf(str, minute, second, last)
	return str
}
