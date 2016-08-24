package getStats

import "github.com/dotabuff/manta/dota"

//判断这条战斗记录，（作用）对象是不是对敌方英雄（英雄本体，非幻象、召唤物），需要满足以下全部条件才返回yes；否则返回no
//[1]TargetName的key在allHeroStats中
//[2]非幻象
//[3](作用)双方不属于同一阵营
func isToOpponentHeroCombatLog(aCombatLog *dota.CMsgDOTACombatLogEntry) bool {
	_, isTargetExist := allHeroStats[aCombatLog.GetTargetName()]
	if isTargetExist && !aCombatLog.GetIsTargetIllusion() && aCombatLog.GetAttackerTeam() != aCombatLog.GetTargetTeam() {
		return true
	}
	return false
}

//判断这条战斗记录，是不是己方英雄（通过其幻象、召唤物）作用于敌方英雄，需要满足以下全部条件才返回yes；否则返回no
//[1]DamageSourceName的key在allHeroStats中
//[2]isToOpponentHeroCombatLog返回yes
func isHeroToOpponentHeroCombatLog(aCombatLog *dota.CMsgDOTACombatLogEntry) bool {
	_, isDamageSourceExist := allHeroStats[aCombatLog.GetDamageSourceName()] //输出来自哪个英雄（英雄的幻象、召唤物造成的伤害，真正的源头还是英雄）
	if isDamageSourceExist && isToOpponentHeroCombatLog(aCombatLog) {
		return true
	}
	return false
}
