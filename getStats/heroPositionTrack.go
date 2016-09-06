package getStats

import (
	"github.com/dotabuff/manta"
	"strings"
	"github.com/dotabuff/manta/dota"
	"math"
)

type HeroPosition struct {
	Cell_X  int32
	Cell_Y  int32
	Vec_X   float32
	Vec_Y   float32
	NetTick uint32
}
//2016/09/05 22:27:35 baseLine : , CBodyComponentBaseAnimatingOverlay.m_cellX : 74
//2016/09/05 22:27:35 baseLine : , CBodyComponentBaseAnimatingOverlay.m_cellY : 74
//2016/09/05 22:27:35 baseLine : , CBodyComponentBaseAnimatingOverlay.m_cellZ : 132
//2016/09/05 22:27:35 baseLine : , CBodyComponentBaseAnimatingOverlay.m_vecX : 212
//2016/09/05 22:27:35 baseLine : , CBodyComponentBaseAnimatingOverlay.m_vecY : 211.96875
//2016/09/05 22:27:35 baseLine : , CBodyComponentBaseAnimatingOverlay.m_vecZ : 0
func recordHeroPosition(parser *manta.Parser, entity *manta.PacketEntity, pet manta.EntityEventType, replaydata *ReplayData) {
	if len(replaydata.heroTackerMap) < 10 && pet == manta.EntityEventType_Create && strings.Contains(entity.ClassName, "CDOTA_Unit_Hero") {
		Clog("EntityEvent : %v, %v, %v, %v", entity.ClassName, pet, entity.Index, timeStampToString(float32(parser.NetTick) / 30))
		//printProperties("baseLine : ", entity.ClassBaseline)
		//printProperties("properties : ", entity.Properties)
		//Clog("\n\n")
		if modifierParent, exist := entity.FetchUint32("m_hModifierParent"); exist {
			replaydata.heroIndexMap[modifierParent] = entity.Index
		}
		heroTackerMap := make(map[int32]*HeroPosition, 0)
		replaydata.heroTackerMap[entity.Index] = heroTackerMap
	}

	if _, exist := replaydata.heroTackerMap[entity.Index]; exist && pet == manta.EntityEventType_Update {

		timeStampInt := int32(parser.NetTick / 30)
		hPosition, positionExists := replaydata.heroTackerMap[entity.Index][timeStampInt]
		cellX, _ := entity.FetchUint64("CBodyComponentBaseAnimatingOverlay.m_cellX")
		cellY, _ := entity.FetchUint64("CBodyComponentBaseAnimatingOverlay.m_cellY")
		vecX, _ := entity.FetchFloat32("CBodyComponentBaseAnimatingOverlay.m_vecX")
		vecY, _ := entity.FetchFloat32("CBodyComponentBaseAnimatingOverlay.m_vecY")
		if positionExists {
			hPosition.Cell_X = int32(cellX)
			hPosition.Cell_Y = int32(cellY)
			hPosition.Vec_X = vecX
			hPosition.Vec_Y = vecY
			hPosition.NetTick = parser.NetTick
		} else {
			replaydata.heroTackerMap[entity.Index][timeStampInt] = &HeroPosition{
				Cell_X : int32(cellX),
				Cell_Y : int32(cellY),
				Vec_X :  vecX,
				Vec_Y :  vecY,
				NetTick : parser.NetTick,
			}
		}
	}
}

func isAloneCatched(log *dota.CMsgDOTACombatLogEntry, replayData *ReplayData) bool {
	isAlone := true
	targetName := log.GetTargetName()
	if heroStats, exists := allHeroStats[targetName]; exists {
		for teamMateName, teamMateHero := range allHeroStats {
			if teamMateName != targetName && heroStats.TeamNumber == teamMateHero.TeamNumber {
				isAlone = !isNear(targetName, teamMateName, int32(log.GetTimestamp()), replayData)
				if !isAlone {
					break
				}
			}
		}
	}
	return isAlone
}

func isNear(targetName uint32, teamMate uint32, time int32, replayData *ReplayData) bool {
	near := false
	//Clog("%v, isNear %v, %v, %v, %v", time, targetName, teamMate, replayData.heroMap[targetName], replayData.heroMap[teamMate])
	targetPosition, exist := replayData.heroTackerMap[replayData.heroIndexMap[replayData.heroMap[targetName]]][time]
	teamMatePosition, exist2 := replayData.heroTackerMap[replayData.heroIndexMap[replayData.heroMap[teamMate]]][time]
	//Clog("%v, %v", exist, exist2)
	if exist && exist2 {
		//Clog("%v, isNearInDistance %v, %v, %v, %v", time, targetName, teamMate, targetPosition.Cell_X, teamMatePosition.Cell_X)
		near = isNearInDistance(targetPosition, teamMatePosition)
	}
	return near
}

func isNearInDistance(targetPosition *HeroPosition, teamMatePosition *HeroPosition) bool {
	near := false
	distance := int32(math.Sqrt(math.Pow(float64(targetPosition.Cell_X - teamMatePosition.Cell_X), float64(2)) + math.Pow(float64(targetPosition.Cell_Y - teamMatePosition.Cell_Y), float64(2))))
	distance = distance * 160
	if distance < 1300{
		near = true
	}
	//Clog("distance : %v", distance)
	//targetPosition.Cell_X - teamMatePosition.Cell_X
	//targetPosition.Cell_Y - teamMatePosition.Cell_Y
	return near
}