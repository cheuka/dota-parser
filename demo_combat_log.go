package demo

import (
	"github.com/dotabuff/manta/dota"
	"github.com/dotabuff/manta"
	"strings"
	"fmt"
	"log"
	"errors"
)

const KDA_COUNT_TIME = 60
const STRING_TABLE_COMBAT = "CombatLogNames"
const STRING_NPC_HERO = "npc_dota_hero_"

//class of demo damageCount
type DemoDamageCount struct {
	heroes     map[int32]*HeroDamageCounter
	heroIds    map[int32]int32
	heroesInit bool
	deathLog   []DeathLog
}

type HeroDamageCounter struct {
	steamId           *uint64
	heroName          *string
	combatLogHeroId   int32
	heroIndex         int32
	damageTotal       int32
	damageToKDA       int32
	damageSufferedLog []DamageLog
}

type DamageLog struct {
	attackName uint32
	targetName uint32
	damage     uint32
	damageSourceName uint32
	timestamp  float32
}

type DeathLog struct {
	targetName    uint32
	timestamp     float32
	assistPlayers []uint32
	damageLog []DamageLog
}

//add damage combat damage log to class
func (counter *DemoDamageCount)CountDamage(entry *dota.CMsgDOTACombatLogEntry) error {
	heroDamageCount := counter.heroes[int32(entry.GetTargetName())]
	if (heroDamageCount == nil) {
		return nil
	}

	damageLog := DamageLog{
		attackName: entry.GetAttackerName(),
		targetName: entry.GetTargetName(),
		damage: entry.GetValue(),
		timestamp: entry.GetTimestamp(),
		damageSourceName : entry.GetDamageSourceName(),
	}
	heroDamageCount.damageSufferedLog = append(heroDamageCount.damageSufferedLog, damageLog)

	attackHero := counter.heroes[int32(entry.GetAttackerName())]
	if attackHero == nil {
		attackHero = counter.heroes[int32(entry.GetDamageSourceName())]
	}

	if attackHero != nil {
		attackHero.damageTotal = attackHero.damageTotal + int32(entry.GetValue())
	}
	return nil
}

//add damage combat death hero log to class
func (counter *DemoDamageCount)CountDeath(entry *dota.CMsgDOTACombatLogEntry) error {
	heroDamageCount := counter.heroes[int32(entry.GetTargetName())]
	if (heroDamageCount == nil) {
		return nil
	}
	deathLog := DeathLog{
		targetName: entry.GetTargetName(),
		timestamp: entry.GetTimestamp(),
		assistPlayers: entry.AssistPlayers,
		damageLog: make([]DamageLog, 0),
	}
	counter.deathLog = append(counter.deathLog, deathLog)
	return nil
}

//for test
func (counter *DemoDamageCount)GetData() error {
	return nil
}

// begin analysis the data
func (counter *DemoDamageCount)Analysis() error {
	// get every hero death evnt
	for _, v := range counter.deathLog {
		targetName := v.targetName
		timeStamp := v.timestamp;
		assistPlayers := v.assistPlayers;
		//get death hero damage log
		deadHeroDamageCount := counter.heroes[int32(targetName)]
		assitDamageCount := make(map[int32]uint32, len(assistPlayers))

		for _, heroIndexes := range assistPlayers{
			assitDamageCount[counter.heroIds[int32(heroIndexes)]] = 0
		}

		if (deadHeroDamageCount != nil) {
			//get every damage event
			for _, sufferedDamage := range deadHeroDamageCount.damageSufferedLog {
				// time compare
				if timeStamp >= sufferedDamage.timestamp && timeStamp - sufferedDamage.timestamp < 20 {
					// attack hero must in the assist list
					countDamage := false
					for _, assistPlayer := range assistPlayers{
						if counter.heroIds[int32(assistPlayer)] == int32(sufferedDamage.attackName) {
							countDamage = true
							break
						}
					}

					attackHeroDamageCount := counter.heroes[int32(sufferedDamage.attackName)]
					if attackHeroDamageCount == nil{
						attackHeroDamageCount = counter.heroes[int32(sufferedDamage.damageSourceName)]
					}

					if (attackHeroDamageCount != nil && countDamage) {
						//append(v.damageLog, sufferedDamage)
						assitDamageCount[attackHeroDamageCount.combatLogHeroId] = assitDamageCount[attackHeroDamageCount.combatLogHeroId] + sufferedDamage.damage
						attackHeroDamageCount.damageToKDA = attackHeroDamageCount.damageToKDA + int32(sufferedDamage.damage)
					}
				} else if timeStamp < sufferedDamage.timestamp {
					break
				}
			}
			for k, v := range assitDamageCount{
				log.Printf("%s damage is %d", *counter.heroes[int32(k)].heroName, v)
			}
		}
	}
	counter.toDamage()
	return nil
}

// print damage
func (counter *DemoDamageCount)toDamage() error {
	for _, v := range counter.heroes {
		fmt.Printf("%v damage : %d %d\n", *v.heroName, v.damageToKDA, v.damageTotal)
	}

	return nil
}

// init
func NewDamageCount(p *manta.Parser) (*DemoDamageCount, error) {

	demo := &DemoDamageCount{
		heroes: make(map[int32]*HeroDamageCounter),
		heroIds: make(map[int32]int32),
		deathLog: make([]DeathLog, 0),

		heroesInit: false,
	}
	p.Callbacks.OnCMsgDOTACombatLogEntry(func(m *dota.CMsgDOTACombatLogEntry) error {
		logType := m.GetType()
		switch logType {
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DEATH:
			//fmt.Printf("DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DEATH %d, %d\n", m.GetAttackerName(), m.GetTargetName())
			if *m.IsTargetHero {
				demo.CountDeath(m)
			}
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DAMAGE:
			if !demo.heroesInit {
				initHero(p, demo)
				demo.heroesInit = true
			}

			if *m.IsTargetHero && !m.GetIsTargetIllusion(){
				demo.CountDamage(m)
			}
		}
		return nil
	})

	p.Callbacks.OnCDemoFileInfo(func(m *dota.CDemoFileInfo) error {

		demo.addPlayInfo(m.GameInfo.Dota.PlayerInfo)
		demo.Analysis()
		return nil
	})

	return demo, nil
}

func getHeroIds(p *manta.Parser) ([]int32, []string) {
	heroIds := make([]int32, 10)
	heroNames := make([]string, 10)
	index := 0
	i := int32(0)
	for i <= 50 {
		name, has := p.LookupStringByIndex(STRING_TABLE_COMBAT, i)
		if has && strings.Contains(name, STRING_NPC_HERO) {
			heroIds[index] = i
			heroNames[index] = name
			index = index + 1
			if index >= 10 {
				return heroIds, heroNames
			}
		}
		i = i + 1
	}
	return heroIds, heroNames
}

func initHero(p *manta.Parser, demo *DemoDamageCount) error {
	//fmt.Printf("initHero \n")
	heroIds, heroNames := getHeroIds(p)
	fmt.Printf("heroIds %v\n", heroIds)
	if len(heroIds) != 10 {
		return errors.New("heroId number error")
	}

	for n, _ := range heroNames {
		heroDamageCount := HeroDamageCounter{
			steamId: nil,
			combatLogHeroId: heroIds[n],
			heroIndex: 0,
			heroName: &heroNames[n],
			damageTotal: 0,
			damageToKDA: 0,
			damageSufferedLog: make([]DamageLog, 0),
		}
		log.Printf("initHero : %v, %d, %d", *heroDamageCount.heroName, n, heroDamageCount.combatLogHeroId)
		demo.heroes[heroIds[n]] = &heroDamageCount
	}

	return nil
}

func (demo *DemoDamageCount)addPlayInfo(playersInfo []*dota.CGameInfo_CDotaGameInfo_CPlayerInfo) {
	for n, v := range playersInfo {
		for _, value := range demo.heroes {
			if strings.Compare(*v.HeroName, *value.heroName) == 0{
				value.heroIndex = int32(n)
				value.steamId = v.Steamid
				demo.heroIds[int32(n)] = value.combatLogHeroId
				log.Printf("player info : %v, %v, %d, %d", *value.steamId, *value.heroName, n, value.combatLogHeroId)
				break
			}
		}
	}

}

