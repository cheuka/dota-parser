package getStats

import (
	"fmt"
	"log"
	"os"
	"strings"

	"new_stats/dota2"

	"github.com/dotabuff/manta"
	"github.com/dotabuff/manta/dota"
)

//统计每个英雄的数据。ket=英雄在combatlog里面的index
var allHeroStats map[uint32]*dota2.Stats
var matchID uint64

//ParseReplay 解析一场比赛的录像，将得到的统计数据存放在allHeroStats中
func ParseReplay(filename string) (map[uint32]*dota2.Stats, error) {

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("unable to open file: %s", err)
	}
	defer f.Close()

	parser, err := manta.NewStreamParser(f)
	if err != nil {
		log.Fatalf("unable to create parser: %s", err)
	}

	var allDamageLogs []*dota.CMsgDOTACombatLogEntry
	parser.Callbacks.OnCMsgDOTACombatLogEntry(func(m *dota.CMsgDOTACombatLogEntry) error {
		logType := m.GetType()
		switch logType {
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DAMAGE:
			allDamageLogs = append(allDamageLogs, m)
		}
		return nil
	})

	var dotaGameInfo *dota.CGameInfo_CDotaGameInfo
	parser.Callbacks.OnCDemoFileInfo(func(m *dota.CDemoFileInfo) error {
		dotaGameInfo = m.GameInfo.Dota
		return nil
	})

	parser.Start()
	//log.Printf("Parse Complete!\n")

	//初始化allHeroStats
	if !initAllHeroStats(parser, dotaGameInfo) {
		return nil, fmt.Errorf("初始化allHeroStats失败")
	}
	//获取统计结果至allHeroStats
	getHeroCreateDeadlyDamages(allDamageLogs)

	return allHeroStats, nil
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

//统计每个英雄的CreateTotalDamages（对敌方英雄[不含幻象、召唤物]造成的总输出）
//有BUG：冰魂给队友套BUFF之后，队友平A造成的伤害，算在队友身上，而不是冰魂（和dotabuff不一致）
//保存英雄对敌方英雄造成的伤害到[allHeroToHeroDamagelogs]
func getHeroCreateTotalDamages(allDamageLogs []*dota.CMsgDOTACombatLogEntry) []*dota.CMsgDOTACombatLogEntry {
	var allHeroToHeroDamagelogs []*dota.CMsgDOTACombatLogEntry
	for _, v := range allDamageLogs {
		_, isDamageSourceExist := allHeroStats[v.GetDamageSourceName()] //输出来自哪个英雄（英雄的幻象、召唤物造成的伤害，真正的源头还是英雄）
		_, isTargetExist := allHeroStats[v.GetTargetName()]             //对哪个英雄造成伤害（不包括对英雄幻象、召唤物的伤害）
		//不考虑对己方英雄造成的伤害（反补、臂章等等）
		if isDamageSourceExist && isTargetExist && !v.GetIsTargetIllusion() && v.GetAttackerTeam() != v.GetTargetTeam() {
			allHeroStats[v.GetDamageSourceName()].CreateTotalDamages += v.GetValue()
			allHeroToHeroDamagelogs = append(allHeroToHeroDamagelogs, v)
		}
	}
	return allHeroToHeroDamagelogs
}

//统计每个英雄的CreateDeadlyDamages（对敌方英雄[不含幻象、召唤物]造成的致死输出）
//一条death记录对应了一条致死deadlyDamagelog记录（该记录的health=0，timestamp与death记录保持一致）
//death次数=致死damage次数=肉山盾复活英雄次数+英雄死亡次数（和KDA中D的总和相等）=肉山盾复活英雄次数+英雄击杀英雄次数（和KDA中A的总和相等）+非英雄单位（防御塔等）击杀英雄次数
//举例2562582896( totalKD=87,89): death次数(93)=肉山盾复活英雄次数（4）+英雄击杀英雄次数（87）+非英雄单位（防御塔）击杀英雄次数（2）
func getHeroCreateDeadlyDamages(allDamageLogs []*dota.CMsgDOTACombatLogEntry) {
	allHeroToHeroDamagelogs := getHeroCreateTotalDamages(allDamageLogs)
	for _, deadlyDamagelog := range allDamageLogs { //这里不能取allHeroToHeroDamagelogs是因为，可能有非英雄单位（防御塔）击杀掉英雄，然后其他英雄助攻打伤害的情况
		_, isTargetExist := allHeroStats[deadlyDamagelog.GetTargetName()]
		if isTargetExist && deadlyDamagelog.GetHealth() == 0 && !deadlyDamagelog.GetIsTargetIllusion() && deadlyDamagelog.GetAttackerTeam() != deadlyDamagelog.GetTargetTeam() {
			for _, aDamagelog := range allHeroToHeroDamagelogs {
				if isDamagelogCount(deadlyDamagelog, aDamagelog) {
					allHeroStats[aDamagelog.GetDamageSourceName()].CreateDeadlyDamages += aDamagelog.GetValue()
				}
			}
		}
	}
}

//初始化：[1]找出所有英雄在combatLog中的index	[2]找出十名英雄使用者的SteamId
func initAllHeroStats(p *manta.Parser, dotaGameInfo *dota.CGameInfo_CDotaGameInfo) bool {
	allHeroStats = make(map[uint32]*dota2.Stats)
	index := int32(0)
	for {
		name, has := p.LookupStringByIndex("CombatLogNames", index)
		//假设index在CombatLogNames中是没有间隔的，遍历CombatLogNames
		if !has {
			break
		}
		if strings.Contains(name, "npc_dota_hero_") {
			allHeroStats[uint32(index)] = &dota2.Stats{HeroName: strings.TrimPrefix(name, "npc_dota_hero_")}
			if len(allHeroStats) == 10 {
				break
			}
		}
		index = index + 1
	}
	if len(allHeroStats) != 10 {
		log.Printf("无法从combatLog中找到十个英雄的index\n")
		return false
	}
	for _, aPlayInfo := range dotaGameInfo.GetPlayerInfo() {
		for _, aHeroStats := range allHeroStats {
			if strings.Contains(aPlayInfo.GetHeroName(), aHeroStats.HeroName) {
				aHeroStats.Steamid = aPlayInfo.GetSteamid()
				aHeroStats.MatchId = matchID
			}
		}
	}
	return true
}
