package main

import (
	"log"
	"os"
	"strings"

	"github.com/dotabuff/manta"
	"github.com/dotabuff/manta/dota"
	"../dota2"
	"fmt"
)

//统计每个英雄的数据。ket=英雄在combatlog里面的index
var allHeroStats map[uint32]*dota2.Stats
//需要特殊记录的控制，无getstuntime和issilence，比如斧王的吼
var specialModifier map[int32]*string
var SPECIAL_MODIFIERS = []string{"modifier_axe_berserkers_call"}
var gameStartTime float32
func main() {
	gameStartTime = float32(0)
	f, err := os.Open("D:\\2545254388.dem")
	if err != nil {
		log.Fatalf("unable to open file: %s", err)
	}
	defer f.Close()
	parser, err := manta.NewStreamParser(f)
	if err != nil {
		log.Fatalf("unable to create parser: %s", err)
	}

	var allDamageLogs []*dota.CMsgDOTACombatLogEntry
	var allModifierLogs []*dota.CMsgDOTACombatLogEntry
	parser.Callbacks.OnCMsgDOTACombatLogEntry(func(m *dota.CMsgDOTACombatLogEntry) error {
		logType := m.GetType()
		switch logType {
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DAMAGE:
			allDamageLogs = append(allDamageLogs, m)
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_REMOVE:
			allModifierLogs = append(allModifierLogs, m)
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_GAME_STATE:
			//记录游戏开始时间
			if(m.GetValue() == uint32(5)){
				gameStartTime = m.GetTimestamp()
			}
		}

		return nil
	})

	parser.Start()
	log.Printf("Parse Complete!\n")

	//初始化allHeroStats
	if !initAllHeroStats(parser) {
		return
	}
	//获取统计结果至allHeroStats
	getHeroCreateDeadlyDamages(allDamageLogs, allModifierLogs)

	//打印结果
	log.Printf("英雄对敌方英雄造成的伤害统计：\n")
	for _, v := range allHeroStats {
		log.Printf("%s——总伤害：%d，致死伤害：%d\n", v.HeroName, v.CreateTotalDamages, v.CreateDeadlyDamages)
	}
	return
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

func isModifierlogCount(deadLog *dota.CMsgDOTACombatLogEntry, modifierLog *dota.CMsgDOTACombatLogEntry) bool {
	if !deadLog.GetIsTargetHero() || !deadLog.GetIsAttackerHero() {
		return false
	}

	modifierTimeStamp := modifierLog.GetTimestamp()
	deadlyDamagelogTimeStamp := deadLog.GetTimestamp()
	// +1原因是控制时间有时候是在英雄死亡后结算
	if modifierTimeStamp <= deadlyDamagelogTimeStamp + float32(1) && modifierTimeStamp >= deadlyDamagelogTimeStamp - 17.0 && modifierLog.GetTargetName() == deadLog.GetTargetName() {
		//目标是英雄，目标不是自己，目标不是幻象，stun时间大于0或者是沉默技能或者在控制列表里
		if modifierLog.GetIsTargetHero() && !modifierLog.GetTargetIsSelf() && !modifierLog.GetIsTargetIllusion() && (modifierLog.GetStunDuration() > float32(0) || modifierLog.GetSilenceModifier() || specialModifier[int32(modifierLog.GetInflictorName())] != nil) {
			log.Printf("%v %v was killed by %v\n", timeStampToString(deadLog.GetTimestamp() - gameStartTime), deadLog.GetTargetName(), deadLog.GetDamageSourceName())
			log.Printf("%v : %v removed %v from %v last %v", timeStampToString(modifierLog.GetTimestamp() - gameStartTime), allHeroStats[modifierLog.GetTargetName()].HeroName, modifierLog.GetInflictorName(), allHeroStats[modifierLog.GetAttackerName()].HeroName, modifierLog.GetModifierElapsedDuration())
			return true
		}
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
func getHeroCreateDeadlyDamages(allDamageLogs []*dota.CMsgDOTACombatLogEntry, allModifierLogs []*dota.CMsgDOTACombatLogEntry) {
	allHeroToHeroDamagelogs := getHeroCreateTotalDamages(allDamageLogs)
	for _, deadlyDamagelog := range allDamageLogs {
		//这里不能取allHeroToHeroDamagelogs是因为，可能有非英雄单位（防御塔）击杀掉英雄，然后其他英雄助攻打伤害的情况
		_, isTargetExist := allHeroStats[deadlyDamagelog.GetTargetName()]
		if isTargetExist && deadlyDamagelog.GetHealth() == 0 && !deadlyDamagelog.GetIsTargetIllusion() && deadlyDamagelog.GetAttackerTeam() != deadlyDamagelog.GetTargetTeam() {
			for _, aDamagelog := range allHeroToHeroDamagelogs {
				if isDamagelogCount(deadlyDamagelog, aDamagelog) {
					allHeroStats[aDamagelog.GetDamageSourceName()].CreateDeadlyDamages += aDamagelog.GetValue()
				}
			}
			getHeroCreateDeadlyControl(deadlyDamagelog, allModifierLogs)
		}
	}
}

func getHeroCreateDeadlyControl(deadLog *dota.CMsgDOTACombatLogEntry, modifierLogs []*dota.CMsgDOTACombatLogEntry) {
	for _, modifierLog := range modifierLogs {
		if isModifierlogCount(deadLog, modifierLog) {
			allHeroStats[modifierLog.GetAttackerName()].CreateDeadlyStiffControl += modifierLog.GetModifierElapsedDuration()
		}
	}
}

//debug版本：会输出每次击杀英雄的详细信息
func getHeroCreateDeadlyDamages_debug(allDamageLogs []*dota.CMsgDOTACombatLogEntry) {
	//记录所有英雄本体受到的伤害源属于敌方英雄（包括幻象、召唤物）的伤害记录
	allHeroToHeroDamagelogs := getHeroCreateTotalDamages(allDamageLogs)
	for _, deadlyDamagelog := range allDamageLogs {
		//分析每条DeadlyDamage记录
		if _, isTargetExist := allHeroStats[deadlyDamagelog.GetTargetName()]; isTargetExist && deadlyDamagelog.GetHealth() == 0 && !deadlyDamagelog.GetIsTargetIllusion() && deadlyDamagelog.GetAttackerTeam() != deadlyDamagelog.GetTargetTeam() {
			if _, isExist := allHeroStats[deadlyDamagelog.GetDamageSourceName()]; isExist {
				log.Printf("<<<%s was killed by %s>>>\n", allHeroStats[deadlyDamagelog.GetTargetName()].HeroName, allHeroStats[deadlyDamagelog.GetDamageSourceName()].HeroName)
			} else {
				log.Printf("<<<%s was killed by %s>>>\n", allHeroStats[deadlyDamagelog.GetTargetName()].HeroName, "not hero!")
			}

			allAssistHeroDamages := make(map[uint32]uint32) //测试用：输出每次英雄死亡的详细信息
			for _, aDamagelog := range allHeroToHeroDamagelogs {
				if isDamagelogCount(deadlyDamagelog, aDamagelog) {
					allHeroStats[aDamagelog.GetDamageSourceName()].CreateDeadlyDamages += aDamagelog.GetValue()
					//测试用：输出每次英雄死亡的详细信息
					if v, isExist := allAssistHeroDamages[aDamagelog.GetDamageSourceName()]; isExist {
						allAssistHeroDamages[aDamagelog.GetDamageSourceName()] = v + aDamagelog.GetValue()
					} else {
						allAssistHeroDamages[aDamagelog.GetDamageSourceName()] = aDamagelog.GetValue()
					}
				}
			}
			for i, v := range allAssistHeroDamages {
				log.Printf("%s:%d\n", allHeroStats[i].HeroName, v)
			}
		}
	}
}

//初始化：找出所有英雄在combatLog中的index
func initAllHeroStats(p *manta.Parser) bool {
	allHeroStats = make(map[uint32]*dota2.Stats)
	specialModifier = make(map[int32]*string)
	index := int32(0)
	for {
		name, has := p.LookupStringByIndex("CombatLogNames", index)
		//假设index在CombatLogNames中是没有间隔的，遍历CombatLogNames
		if !has {
			break
		}
		if strings.Contains(name, "npc_dota_hero_") {
			//因为后面要获取特殊控制技能，所以这里不再break
			allHeroStats[uint32(index)] = &dota2.Stats{HeroName: strings.TrimPrefix(name, "npc_dota_hero_")}
			//if len(allHeroStats) == 10 {
			//	break
			//}
		}
		//获取特殊控制技能
		for _, modifier := range SPECIAL_MODIFIERS {
			if strings.EqualFold(modifier, name){
				specialModifier[index] = &modifier
				break
			}
		}
		index = index + 1
	}

	if len(allHeroStats) != 10 {
		log.Printf("无法从combatLog中找到十个英雄的index\n")
		return false
	}
	return true
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
