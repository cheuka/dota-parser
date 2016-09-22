package getStats

import (
	"github.com/dotabuff/manta/dota"
	"github.com/dotabuff/manta"
	"sort"
)

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

func fetchProperties(key string, entity *manta.PacketEntity) interface{}{
	value, exist := entity.Fetch(key)
	if exist{
		return value
	}

	return nil
}

func printIntKeyMaps(tag string, maps map[int32]*HeroPosition) {
	sorted_keys := make([]int, 0)
	for k, _ := range maps {
		sorted_keys = append(sorted_keys, int(k))
	}

	// sort 'string' key in increasing order
	sort.Ints(sorted_keys)

	lastK := 0
	inCountinudNum := 0
	for _, k := range sorted_keys {
		if k > 0 && k - lastK > 1{
			Clog("%v, in continued : %v, %v", tag, k, k - lastK)
			inCountinudNum++
		}
		lastK = k
	}
	Clog("%v, in continued number : %v", tag, inCountinudNum)
}

//记录打GG时间
func countGG(replaydata *ReplayData){
	sorted_keys := make([]float64, 0)
	for k, _ := range replaydata.ggCount {
		sorted_keys = append(sorted_keys, float64(k))
	}

	// sort 'string' key in increasing order
	sort.Float64s(sorted_keys)
	firstGG := false
	for _, k := range sorted_keys {
		playerId := replaydata.ggCount[float32(k)]
		if heroStats, exist := getHeroStatesFromPlayerId(playerId); exist && replaydata.gameEndTime - float32(k) < 15{
			if heroStats.IsWin{
				heroStats.IsReplyGG = 1
			} else {
				if !firstGG{
					heroStats.IsFirstGG = 1
					firstGG = true
				}
				heroStats.IsWriteGG = 1
			}
		}
	}
}