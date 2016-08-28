package getStats

import (
	"fmt"

	"github.com/dotabuff/manta/dota"
)

//计算一次死亡的致死控制时间规则：
//[InflictorName]并不是技能名称，而是打击类型（类似于伤害类型）。【所有的add\remove记录都是以这个InflictorName为粒度的！！！】
//举例：白虎的箭、VS的锤、小强的尖刺外壳，都是InflictorName=110[客户端称为“眩晕”]；而小强的刺则为100[“穿刺”]（虽然眩晕和穿刺造成的结果是一样的）。
//[add记录]不是每一个控制技能都会有一条add记录！只有目标英雄第一次遭受InflictorName的时候，才会产生这类InflictorName的add记录。ModifierDuration表示这个InflictorName（从现在开始）将要持续多久。
//举例：小强的尖刺外壳晕了小黑，会有一条[眩晕]的add记录，在眩晕结束之前，VS补了一个锤子，这个锤子是不会产生add记录的。
//[remove记录]和add记录类似，当InflictorName持续时间到了的时候，会产生一条remove记录，ModifierElapsedDuration表示（最后一个技能）实际持续时间。
//【之前的bug之一】根据StunDuration>0判断某技能是否为控制技能时，只对add记录有效，因为StunDuration=ModifierDuration永远为正数。
//但是不能用来判断remove记录，因为StunDuration=ModifierElapsedDuration，是可以为0的(第二手补的控制直接把目标打死了)【这也是为什么最开始我无法打印出remove信息，但其实是有的】

//【计算方法】
//【bug】只根据add信息和remove信息，无法准确计算控制时间（当给予同一目标同一InflictorName的技能>=3个时，add和remove最多只能给出2个信息）
//对每一个InflictorName的add和remove，如果attackName，累加attackName的控制时间为remove的timestamp减去add的timestamp。（不用ModifierElapsedDuration是为了蓝胖的A杖晕）;
//如果attackName不同，累加第一个attackName的控制时间为min(ModifierDuration,remove_timestape-add_timestamp);
//累加第二个attackName的控制时间为ModifierElapsedDuration

//【展望】
//再找其他的信息进行计算。

func calcCreateDeadlyControl(replayData *ReplayData) {

	for _, deadlyDamagelog := range replayData.allDamageLogs {
		if deadlyDamagelog.GetHealth() == 0 && isToOpponentHeroCombatLog(deadlyDamagelog) {
			if _, isExist := allHeroStats[deadlyDamagelog.GetDamageSourceName()]; isExist {
				Clog("%v——<<<%s was killed by %s>>>\n", timeStampToString(deadlyDamagelog.GetTimestamp()-replayData.gameStartTime), allHeroStats[deadlyDamagelog.GetTargetName()].HeroName, allHeroStats[deadlyDamagelog.GetDamageSourceName()].HeroName)
			} else {
				Clog("%v——<<<%s was killed by NOT HERO>>>\n", timeStampToString(deadlyDamagelog.GetTimestamp()-replayData.gameStartTime), allHeroStats[deadlyDamagelog.GetTargetName()].HeroName)
			}
			//map保存modifier add的记录，在add和remove时如果inflictorName相同(比如晕眩)而attackName不同时，同时记录add和remove两条记录
			addModifierMap := make(map[uint32]*dota.CMsgDOTACombatLogEntry)
			for _, aModifierLog := range replayData.allModifierLogs {
				if isModifierlogCount1(replayData, deadlyDamagelog, aModifierLog, addModifierMap) {
					if aModifierLog.GetType() == dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_ADD {
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

func isModifierlogCount1(replayData *ReplayData, deadlyDamagelog, aModifierLog *dota.CMsgDOTACombatLogEntry, addModifierMap map[uint32]*dota.CMsgDOTACombatLogEntry) bool {
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
			if aModifierLog.GetType() == dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_ADD {

			} else {
				Clog("%v : %v removed %v from %v last %v", timeStampToString(aModifierLog.GetTimestamp()-replayData.gameStartTime), allHeroStats[aModifierLog.GetTargetName()].HeroName, aModifierLog.GetInflictorName(), allHeroStats[aModifierLog.GetAttackerName()].HeroName, aModifierLog.GetModifierElapsedDuration())

			}
			if addModifierMap[aModifierLog.GetInflictorName()] != nil {
				if addModifierMap[aModifierLog.GetInflictorName()].GetAttackerName() != aModifierLog.GetAttackerName() {
					Clog("%v : %v add %v from %v last %v", timeStampToString(addModifierMap[aModifierLog.GetInflictorName()].GetTimestamp()-replayData.gameStartTime), allHeroStats[addModifierMap[aModifierLog.GetInflictorName()].GetTargetName()].HeroName, addModifierMap[aModifierLog.GetInflictorName()].GetInflictorName(), allHeroStats[addModifierMap[aModifierLog.GetInflictorName()].GetAttackerName()].HeroName, addModifierMap[aModifierLog.GetInflictorName()].GetModifierDuration())
				}

			}
			return true
		}
	}
	return false
}

func isModifierlogCount(replayData *ReplayData, deadlyDamagelog, aModifierLog *dota.CMsgDOTACombatLogEntry, addModifierMap map[uint32]*dota.CMsgDOTACombatLogEntry) bool {
	modifierTimeStamp := aModifierLog.GetTimestamp()
	deadlyDamagelogTimeStamp := deadlyDamagelog.GetTimestamp()

	_, isAttackerExist := allHeroStats[aModifierLog.GetAttackerName()]
	if !isAttackerExist || aModifierLog.GetIsTargetIllusion() || aModifierLog.GetAttackerName() == aModifierLog.GetTargetName() {
		return false
	}

	// +1原因是控制时间有时候是在英雄死亡后结算
	if modifierTimeStamp <= deadlyDamagelogTimeStamp+float32(1) && modifierTimeStamp >= deadlyDamagelogTimeStamp-17.0 && aModifierLog.GetTargetName() == deadlyDamagelog.GetTargetName() {
		//remove记录根据add记录来判断是否为控制技能
		if _, isExist := addModifierMap[aModifierLog.GetInflictorName()]; isExist && aModifierLog.GetType() == dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_REMOVE {
			printControlDetail(replayData, aModifierLog)
			return true
		}
		//stun时间大于0或者是沉默技能或者在控制列表里
		if aModifierLog.GetStunDuration() > float32(0) || aModifierLog.GetSilenceModifier() || replayData.specialModifier[int32(aModifierLog.GetInflictorName())] != nil {
			//add记录判断是否为控制技能
			if aModifierLog.GetType() == dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_ADD {
				printControlDetail(replayData, aModifierLog)
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
