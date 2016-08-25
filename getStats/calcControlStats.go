package getStats

import (
	"fmt"
	"log"

	"github.com/dotabuff/manta/dota"
)

func calcCreateDeadlyControl(replayData *ReplayData) {
	for _, deadlyDamagelog := range replayData.allDamageLogs {
		if deadlyDamagelog.GetHealth() == 0 && isToOpponentHeroCombatLog(deadlyDamagelog) {
			if _, isExist := allHeroStats[deadlyDamagelog.GetDamageSourceName()]; isExist {
				log.Printf("%v——<<<%s was killed by %s>>>\n", timeStampToString(deadlyDamagelog.GetTimestamp()-replayData.gameStartTime), allHeroStats[deadlyDamagelog.GetTargetName()].HeroName, allHeroStats[deadlyDamagelog.GetDamageSourceName()].HeroName)
			} else {
				log.Printf("%v——<<<%s was killed by NOT HERO>>>\n", timeStampToString(deadlyDamagelog.GetTimestamp()-replayData.gameStartTime), allHeroStats[deadlyDamagelog.GetTargetName()].HeroName)
			}
			for _, aModifierLog := range replayData.allModifierLogs {
				if isModifierlogCount(replayData, deadlyDamagelog, aModifierLog) {
					allHeroStats[aModifierLog.GetAttackerName()].CreateDeadlyStiffControl += aModifierLog.GetModifierElapsedDuration()
				}
			}
		}
	}

}

func isModifierlogCount(replayData *ReplayData, deadlyDamagelog, aModifierLog *dota.CMsgDOTACombatLogEntry) bool {
	_, isAttackerExist := allHeroStats[aModifierLog.GetAttackerName()]
	if !isAttackerExist || aModifierLog.GetIsTargetIllusion() || aModifierLog.GetAttackerName() == aModifierLog.GetTargetName() {
		return false
	}
	modifierTimeStamp := aModifierLog.GetTimestamp()
	deadlyDamagelogTimeStamp := deadlyDamagelog.GetTimestamp()
	// +1原因是控制时间有时候是在英雄死亡后结算
	if modifierTimeStamp <= deadlyDamagelogTimeStamp+float32(1) && modifierTimeStamp >= deadlyDamagelogTimeStamp-17.0 && aModifierLog.GetTargetName() == deadlyDamagelog.GetTargetName() {
		//stun时间大于0或者是沉默技能或者在控制列表里
		if aModifierLog.GetStunDuration() > float32(0) || aModifierLog.GetSilenceModifier() || replayData.specialModifier[int32(aModifierLog.GetInflictorName())] != nil {
			if(aModifierLog.GetType() == dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_ADD){
				log.Printf("%v : %v add %v from %v last %v", timeStampToString(aModifierLog.GetTimestamp()-replayData.gameStartTime), allHeroStats[aModifierLog.GetTargetName()].HeroName, aModifierLog.GetInflictorName(), allHeroStats[aModifierLog.GetAttackerName()].HeroName, aModifierLog.GetModifierDuration())
			}else{
				log.Printf("%v : %v removed %v from %v last %v", timeStampToString(aModifierLog.GetTimestamp()-replayData.gameStartTime), allHeroStats[aModifierLog.GetTargetName()].HeroName, aModifierLog.GetInflictorName(), allHeroStats[aModifierLog.GetAttackerName()].HeroName, aModifierLog.GetModifierElapsedDuration())

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
	last := int32(stamp*60) % 60
	str = fmt.Sprintf(str, minute, second, last)
	return str
}
