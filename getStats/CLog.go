package getStats

import (
	"log"

	"github.com/dotabuff/manta/dota"
)

var debug = true

func SetDebug(isDebug bool) {
	debug = isDebug
}

func Clog(format string, v ...interface{}) {
	if debug {
		log.Printf(format, v...)
	}
}

//输出控制详情
func printControlDetail(replayData *ReplayData, aModifierLog *dota.CMsgDOTACombatLogEntry) {
	targetName := "未知目标"
	v, isTargetExist := allHeroStats[aModifierLog.GetTargetName()]
	if isTargetExist {
		targetName = v.HeroName
	}
	attackerName := "未知攻击者"
	v, isAttackerExist := allHeroStats[aModifierLog.GetAttackerName()]
	if isAttackerExist {
		attackerName = v.HeroName
	}
	damageSourceName := "未知攻击源头"
	v, isDamageSourceExist := allHeroStats[aModifierLog.GetDamageSourceName()]
	if isDamageSourceExist {
		damageSourceName = v.HeroName
	}

	eventName := "未知事件"
	logType := aModifierLog.GetType()
	switch logType {
	// case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DAMAGE:
	// 	replayData.allDamageLogs = append(replayData.allDamageLogs, m)
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_REMOVE:
		eventName = "remove"
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_ADD:
		eventName = "add"
	case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_REFRESH:
		eventName = "refresh"
	}

	Clog("%v : 事件:%v，技能:%v，目标:%v，攻击者:%v，攻击源头:%v，预计持续时间:%v，实际作用时间:%v，stun_duration:%v。", timeStampToString(aModifierLog.GetTimestamp()-replayData.gameStartTime), eventName, aModifierLog.GetInflictorName(), targetName, attackerName, damageSourceName, aModifierLog.GetModifierDuration(), aModifierLog.GetModifierElapsedDuration(), aModifierLog.GetStunDuration())
}
