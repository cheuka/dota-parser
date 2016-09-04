package getStats

import (
)

//gold reason
//0 发钱   原始工资625，赏金符， 炼金的被动
//1 死亡掉钱
//2 买活掉钱
//6 卖东西
//
//11 推塔  可靠
//12 击杀所获得的钱  击杀+助攻 可靠
//13 打野和补刀的钱 不可靠
//14 roshan 可靠
//15 杀小鸡 可靠

//工资 1分钟100块， 每秒1.666块
//计算GPM, RGPM,
//Reliable gold - Any bounty you get from hero kills, Roshan, couriers, Hand of Midas, Track gold and global gold from towers is added to your reliable gold pool.
//Unreliable gold - Everything else (starting gold, periodic gold, creep kills, neutrals, etc).
func calcFarm(replayData *ReplayData) {
	deadMap := make(map[float32]uint32)
	for _, logEntry := range replayData.allGoldLogs {
		//英雄死亡时，记录下时间戳，在计算fed gold时使用
		if logEntry.GetGoldReason() == uint32(1){
			deadMap[logEntry.GetTimestampRaw()] = logEntry.GetTargetName()
		}
	}

	for _, logEntry := range replayData.allGoldLogs {
		reason := logEntry.GetGoldReason()
		targetName := logEntry.GetTargetName()
		if allHeroStats[targetName] == nil {
			return
		}
		switch reason {
		case 0:
			allHeroStats[targetName].UnrRpm += logEntry.GetValue()
		case 1:
			//英雄死亡丢失金钱，为负数
			allHeroStats[targetName].DeadLoseGold -= logEntry.GetValue()
			//英雄死亡时，将敌方英雄击杀的金钱算到自己的fed gold中来
		case 2:
		case 6:
		case 11:
			allHeroStats[targetName].RGpm += logEntry.GetValue()
		case 12:
			//击杀英雄所获得金钱，在死亡英雄附近
			allHeroStats[targetName].RGpm += logEntry.GetValue()
			allHeroStats[targetName].KillHeroGold += logEntry.GetValue()
			//如果在deadmap里找的到死亡英雄，将改部分钱也计算到死亡英雄的fed gold里面
			if _, exist := deadMap[logEntry.GetTimestampRaw()]; exist {
				allHeroStats[deadMap[logEntry.GetTimestampRaw()]].FedEnemyGold += logEntry.GetValue()
			}
		case 13:
			//点金手算13 但是计入可靠金钱
			if logEntry.GetValue() == uint32(190) {
				allHeroStats[targetName].RGpm += logEntry.GetValue()
			} else {
				allHeroStats[targetName].UnrRpm += logEntry.GetValue()
			}
		case 14:
			allHeroStats[targetName].RGpm += logEntry.GetValue()
		case 15:
			allHeroStats[targetName].RGpm += logEntry.GetValue()
		}
	}
}
