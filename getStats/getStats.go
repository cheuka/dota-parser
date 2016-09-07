package getStats

import (
	"fmt"

	"new_stats/dota2"

	"github.com/dotabuff/manta/dota"
	"strconv"
)

//ReplayData 用于存放解析录像时，从回调函数中获取到的数据
type ReplayData struct {
	allDamageLogs   []*dota.CMsgDOTACombatLogEntry
	allModifierLogs []*dota.CMsgDOTACombatLogEntry
	allGoldLogs []*dota.CMsgDOTACombatLogEntry
	dotaGameInfo    *dota.CGameInfo_CDotaGameInfo
	specialModifier map[int32]*string //需要特殊记录的控制，无getstuntime和issilence，比如斧王的吼
	gameStartTime   float32
	gameEndTime   float32
	teamDeath map[uint32]uint32 //队伍死亡次数， key:team id,value: 死亡次数
	heroMap map[uint32]uint32 //英雄combat log的target name对应的 英雄entity name, 由play resource的mSelectHero获得，再通过英雄entity的mModifierName找到对应的英雄
	heroIndexMap map[uint32]int32 //entitymap key 上面的resource的mSelectHero获得， value：entity的index 用来获取位置
	heroTackerMap map[int32]map[int32]*HeroPosition //英雄位置记录 第一个key为entity的index唯一， 第二个key为时间戳的整数形式，
}

//统计每个英雄的数据。ket=英雄在combatlog里面的index
var allHeroStats map[uint32]*dota2.Stats

//GetStats 解析一场比赛的录像，将得到的统计数据存放在allHeroStats中
func GetStats(filename string) (map[uint32]*dota2.Stats, error) {
	replayData := ReplayData{
		teamDeath : make(map[uint32]uint32, 0),
		heroMap : make(map[uint32]uint32, 0),
		heroIndexMap : make(map[uint32]int32, 0),
		heroTackerMap : make(map[int32]map[int32]*HeroPosition, 0),
	}
	//解析录像，获取数据
	err := parseReplay(filename, &replayData)
	if err != nil {
		return nil, fmt.Errorf("解析录像失败：%s", err)
	}
	//计算伤害指标至allHeroStats
	calcCreateTotalDamages(&replayData)
	calcCreateDeadlyDamages(&replayData)
	//计算控制指标至allHeroStats
	calcCreateDeadlyControl(&replayData)
	calcTeamDeath(&replayData)
	//计算金钱
	calcFarm(&replayData)
	printIntKeyMaps(strconv.Itoa(404), replayData.heroTackerMap[404])
	return allHeroStats, nil
}

