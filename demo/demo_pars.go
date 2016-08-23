package demo

import (
	"fmt"

	"github.com/dotabuff/manta"
	"os"
	"log"
	"github.com/dotabuff/manta/dota"
	"strings"
)

func ParseDemo(filePath string) error{
	//logFile, err := os.OpenFile("D:\\goLog.txt", os.O_RDWR | os.O_CREATE, 0)
	//defer logFile.Close()

	//logger := log.New(logFile, "", log.Ldate | log.Ltime);
	// Create a new parser instance from a file. Alternatively see NewParser([]byte)

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("unable to open file: %s", err)
	}
	defer f.Close()

	parser, err := manta.NewStreamParser(f)
	if err != nil {
		log.Fatalf("unable to create parser: %s", err)
	}

	// Register a callback, this time for the OnCUserMessageSayText2 event.
	parser.Callbacks.OnCUserMessageSayText2(func(m *dota.CUserMessageSayText2) error {
		log.Printf("%s said: %s\n", m.GetParam1(), m.GetParam2())
		return nil
	})


	parser.Callbacks.OnCMsgDOTACombatLogEntry(func(m *dota.CMsgDOTACombatLogEntry) error {
		logType := m.GetType()
		switch logType {
		//case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DAMAGE:
		//	logger.Printf("CMsgDOTACombatLogEntry %v, %v\n", m.GetType())
		//	attack_name, has := p.LookupStringByIndex("CombatLogNames", int32(m.GetAttackerName()))
		//	target_name, has2 := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetName()))
		//	if has && has2 {
		//		logger.Printf("%s  : %s attack %s at %d\n", timeStampToString(m.GetTimestamp() - gameStartTime), attack_name, target_name, m.GetValue())
		//	}
		//case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DEATH:
		//	//log.Printf("CMsgDOTACombatLogEntry %v\n", m.GetType())
		//	attack_name, has := p.LookupStringByIndex("CombatLogNames", int32(m.GetAttackerName()))
		//	target_name, has2 := p.LookupStringByIndex("CombatLogNames", int32(m.GetTargetName()))
		//	if has && has2 && *m.IsTargetHero {
		//		log.Printf("%s  : %s killed %s at %d\n", timeStampToString(m.GetTimestamp() - gameStartTime), attack_name, target_name, m.GetValue())
		//		log.Printf("assister %v \n", m.AssistPlayers)
		//		log.Printf("%v", m.String())
		//	}

		//case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_GAME_STATE:
		//	logger.Printf("game state is now %v, %v, %v\n", m.GetTimestamp(), m.GetValue(), gameStartTime)
		//	state := m.GetValue()
		//	if state == uint32(5) {
		//		gameStartTime = m.GetTimestamp()
		//	}

		}
		return nil
	})

	//parser.OnPacketEntity(func(pe *manta.PacketEntity, pet manta.EntityEventType) error {
	//	if pet == manta.EntityEventType_Create && strings.Compare(pe.ClassName, "CDOTA_NPC_Observer_Ward") == 0{
	//		log.Printf("EntityEvent : %v, %v\n\n\n", pe.ClassName, pet)
	//		for k,v := range pe.ClassBaseline.KV{
	//			log.Printf("ClassBaseline %v : %v\n", k, v)
	//		}
	//		for k,v := range pe.Properties.KV{
	//			log.Printf("Observer_Ward %v : %v\n", k, v)
	//		}
	//	}
	//	return nil
	//})
	// Start parsing the replay!
	damageCount, error := NewDamageCount(parser)
	if error == nil {
		damageCount.GetData()
	}
	parser.Start()
	log.Printf("Parse Complete!\n")
	return nil
}

func timeStampToString(stamp float32) string {
	str := "[%d : %d : %d]"
	minute := int32(stamp) / 60
	second := int32(stamp) % 60
	last := int32(stamp * 60) % 60
	str = fmt.Sprintf(str, minute, second, last)
	return str
}
