package getStats

import (
	"math"
	"strings"

	"github.com/dotabuff/manta"
	"github.com/dotabuff/manta/dota"
)

//英雄位置的记录
//获取英雄位置对应的办法
// damageLog里的targetName获取 hero stats
// 根据playResource获取 m_vecPlayerTeamData.000%d.m_hSelectedHero 对应 targetName, 此处为replaysData.heroMap
// 英雄entity的对应关系是 entity的m_hModifierParent 对应 m_hSelectedHero 此处为replaysData.heroIndexMap
// 位置记录由entity的更新来获取，key为entity的index,  index可以由replaysData.heroIndexMap以m_hSelectedHero为key获取
// 位置记录为嵌套的map，第一层map的key为index, 第二层的key为时间。
// damageLog里英雄的位置， 根据targetName从第一个map里获取色狼selectHero, 在根据selectHero从第二个map里获取index，再根据index和timestamp的整数形式获取英雄的位置
// ps：不能跳过selectHero字段因为 entity的解析是在计算之前的，解析的时候还没有targetName和selectHero的对应关系
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
	//在英雄entity创建的时候创建map关系
	// 看是否有力量这个属性来判断是否是英雄，防止类似兽王的召唤物，暂时先用这个来解决bug，将来找到更加准确的字段再替换
	_, isHero := entity.FetchFloat32("m_flStrength")

	if len(replaydata.heroTackerMap) < 10 && pet == manta.EntityEventType_Create && strings.Contains(entity.ClassName, "CDOTA_Unit_Hero") && isHero {
		Clog("EntityEvent : %v, %v, %v, %v", entity.ClassName, pet, entity.Index, timeStampToString(float32(parser.NetTick)/30))
		//printProperties("baseLine : ", entity.ClassBaseline)
		//printProperties("properties : ", entity.Properties)
		//Clog("\n\n")
		if modifierParent, exist := entity.FetchUint32("m_hModifierParent"); exist {
			replaydata.heroIndexMap[modifierParent] = entity.Index
		}
		heroTackerMap := make(map[int32]*HeroPosition, 0)
		replaydata.heroTackerMap[entity.Index] = heroTackerMap
	}
	//if entity.Index == int32(404){
	//	cellX, _ := entity.FetchUint64("CBodyComponentBaseAnimatingOverlay.m_cellX")
	//	cellY, _ := entity.FetchUint64("CBodyComponentBaseAnimatingOverlay.m_cellY")
	//	vecX, _ := entity.FetchFloat32("CBodyComponentBaseAnimatingOverlay.m_vecX")
	//	vecY, _ := entity.FetchFloat32("CBodyComponentBaseAnimatingOverlay.m_vecY")
	//	Clog("entity.Index : %v, %v, %v, %v, %v, %v",  int32(parser.NetTick / 30), timeStampToString(float32(parser.NetTick / 30) - replaydata.gameStartTime), cellX, cellY, vecX, vecY)
	//}
	//在英雄entity的位置更新的时候记录下来，为避免数据量过大，以second为单位记录
	//NetTick为 tick rate， 和timestamp 1：30
	if _, exist := replaydata.heroTackerMap[entity.Index]; exist && pet == manta.EntityEventType_Update {
		//更新位置 tick数据量大，如果是同一秒，就更新，新的一秒的话就创建一个新的
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
				Cell_X:  int32(cellX),
				Cell_Y:  int32(cellY),
				Vec_X:   vecX,
				Vec_Y:   vecY,
				NetTick: parser.NetTick,
			}
		}
	}
}

//计算是否被单抓
//在 damage记录里 获取targetName 的友方，看是否存在距离小于1300码的
func isAloneCatched(log *dota.CMsgDOTACombatLogEntry, replayData *ReplayData) bool {
	isAlone := true
	targetName := log.GetTargetName()
	if heroStats, exists := allHeroStats[targetName]; exists {
		for teamMateName, teamMateHero := range allHeroStats {
			if teamMateName != targetName && heroStats.TeamNumber == teamMateHero.TeamNumber {
				isAlone = !isNear(targetName, teamMateName, int32(log.GetTimestampRaw()), replayData)
				if !isAlone {
					break
				}
			}
		}
	}
	return isAlone
}

//判断是否在指定距离之内
func isNear(targetName uint32, teamMate uint32, time int32, replayData *ReplayData) bool {
	near := false
	targetPosition, exist := findPosition(targetName, time, replayData)
	teamMatePosition, exist2 := findPosition(teamMate, time, replayData)
	//Clog("%v, isNear %v, %v, %v", timeStampToString(float32(time) - replayData.gameStartTime), allHeroStats[targetName].HeroName, allHeroStats[teamMate].HeroName, getDistance(targetPosition.Cell_X, targetPosition.Cell_Y, teamMatePosition.Cell_X, teamMatePosition.Cell_Y))
	//Clog("%v, isNear %v, %v, %v, %v", timeStampToString(float32(time) - replayData.gameStartTime), targetPosition.Cell_X, targetPosition.Cell_Y, teamMatePosition.Cell_X, teamMatePosition.Cell_Y)
	//Clog("%v, %v", exist, exist2)
	if exist && exist2 {
		//Clog("%v, isNearInDistance %v, %v, %v, %v", time, targetName, teamMate, targetPosition.Cell_X, teamMatePosition.Cell_X)
		near = isNearInDistance(targetPosition, teamMatePosition)
	}
	return near
}

//计算实际距离，此处以缩放128为标准，待之后得到更加准确的信息之后再更改
func isNearInDistance(targetPosition *HeroPosition, teamMatePosition *HeroPosition) bool {
	near := false
	distance := int32(math.Sqrt(math.Pow(float64(targetPosition.Cell_X-teamMatePosition.Cell_X), float64(2)) + math.Pow(float64(targetPosition.Cell_Y-teamMatePosition.Cell_Y), float64(2))))
	distance = distance * 128
	if distance < 1300 {
		near = true
	}
	//Clog("distance : %v", distance)
	//targetPosition.Cell_X - teamMatePosition.Cell_X
	//targetPosition.Cell_Y - teamMatePosition.Cell_Y
	return near
}

func getDistance(cell1X, cell1Y, cell2X, cell2Y int32) int32 {
	return int32(math.Sqrt(math.Pow(float64(cell1X-cell2X), float64(2)) + math.Pow(float64(cell1Y-cell2Y), float64(2))))
}

func findPosition(targetName uint32, time int32, replayData *ReplayData) (*HeroPosition, bool) {
	caculateTime := time
	for {
		targetPosition, exist := replayData.heroTackerMap[replayData.heroIndexMap[replayData.heroMap[targetName]]][caculateTime]
		if exist {
			return targetPosition, true
		} else {
			caculateTime--
		}

		if caculateTime <= 0 {
			Clog("findposition error: %v, %v, %v", targetName, allHeroStats[targetName].HeroName, time)
			break
		}
	}
	return nil, false
}
