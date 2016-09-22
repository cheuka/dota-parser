// Code generated by protoc-gen-go.
// source: stats.proto
// DO NOT EDIT!

/*
Package dota2 is a generated protocol buffer package.

It is generated from these files:
	stats.proto

It has these top-level messages:
	Stats
*/
package dota2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Stats struct {
	MatchId                          uint64  `protobuf:"varint,1,opt,name=match_id,json=matchId" json:"match_id,omitempty"`
	AccountId                        uint32  `protobuf:"varint,2,opt,name=account_id,json=accountId" json:"account_id,omitempty"`
	Steamid                          uint64  `protobuf:"varint,15,opt,name=steamid" json:"steamid,omitempty"`
	PlayerName                       string  `protobuf:"bytes,3,opt,name=player_name,json=playerName" json:"player_name,omitempty"`
	HeroId                           uint32  `protobuf:"varint,4,opt,name=hero_id,json=heroId" json:"hero_id,omitempty"`
	HeroName                         string  `protobuf:"bytes,5,opt,name=hero_name,json=heroName" json:"hero_name,omitempty"`
	CreateTotalDamages               uint32  `protobuf:"varint,6,opt,name=create_total_damages,json=createTotalDamages" json:"create_total_damages,omitempty"`
	CreateDeadlyDamages              uint32  `protobuf:"varint,7,opt,name=create_deadly_damages,json=createDeadlyDamages" json:"create_deadly_damages,omitempty"`
	CreateTotalStiffControl          float32 `protobuf:"fixed32,8,opt,name=create_total_stiff_control,json=createTotalStiffControl" json:"create_total_stiff_control,omitempty"`
	CreateDeadlyStiffControl         float32 `protobuf:"fixed32,9,opt,name=create_deadly_stiff_control,json=createDeadlyStiffControl" json:"create_deadly_stiff_control,omitempty"`
	OpponentHeroDeaths               uint32  `protobuf:"varint,10,opt,name=opponent_hero_deaths,json=opponentHeroDeaths" json:"opponent_hero_deaths,omitempty"`
	CreateDeadlyDamagesPerDeath      float32 `protobuf:"fixed32,11,opt,name=create_deadly_damages_per_death,json=createDeadlyDamagesPerDeath" json:"create_deadly_damages_per_death,omitempty"`
	CreateDeadlyStiffControlPerDeath float32 `protobuf:"fixed32,12,opt,name=create_deadly_stiff_control_per_death,json=createDeadlyStiffControlPerDeath" json:"create_deadly_stiff_control_per_death,omitempty"`
	AloneKilledNum                   uint32  `protobuf:"varint,16,opt,name=aloneKilledNum" json:"aloneKilledNum,omitempty"`
	AloneBeKilledNum                 uint32  `protobuf:"varint,17,opt,name=aloneBeKilledNum" json:"aloneBeKilledNum,omitempty"`
	AloneBeCatchedNum                int32   `protobuf:"varint,18,opt,name=aloneBeCatchedNum" json:"aloneBeCatchedNum,omitempty"`
	RGpm                             uint32  `protobuf:"varint,19,opt,name=rGpm" json:"rGpm,omitempty"`
	UnrRpm                           uint32  `protobuf:"varint,20,opt,name=unrRpm" json:"unrRpm,omitempty"`
	KillHeroGold                     uint32  `protobuf:"varint,21,opt,name=killHeroGold" json:"killHeroGold,omitempty"`
	DeadLoseGold                     uint32  `protobuf:"varint,22,opt,name=deadLoseGold" json:"deadLoseGold,omitempty"`
	FedEnemyGold                     uint32  `protobuf:"varint,23,opt,name=fedEnemyGold" json:"fedEnemyGold,omitempty"`
	TeamNumber                       int32   `protobuf:"varint,24,opt,name=teamNumber" json:"teamNumber,omitempty"`
	IsWin                            bool    `protobuf:"varint,25,opt,name=isWin" json:"isWin,omitempty"`
	PlayerId                         int32   `protobuf:"varint,26,opt,name=player_id,json=playerId" json:"player_id,omitempty"`
	IsFirstGG                        int32   `protobuf:"varint,28,opt,name=isFirstGG" json:"isFirstGG,omitempty"`
	IsWriteGG                        int32   `protobuf:"varint,29,opt,name=isWriteGG" json:"isWriteGG,omitempty"`
	IsReplyGG                        int32   `protobuf:"varint,30,opt,name=isReplyGG" json:"isReplyGG,omitempty"`
	AloneCatchedNum                  int32   `protobuf:"varint,31,opt,name=aloneCatchedNum" json:"aloneCatchedNum,omitempty"`
	ConsumeDamage                    int32   `protobuf:"varint,32,opt,name=consumeDamage" json:"consumeDamage,omitempty"`
}

func (m *Stats) Reset()                    { *m = Stats{} }
func (m *Stats) String() string            { return proto.CompactTextString(m) }
func (*Stats) ProtoMessage()               {}
func (*Stats) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*Stats)(nil), "dota2.Stats")
}

func init() { proto.RegisterFile("stats.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 568 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x94, 0x5f, 0x6f, 0xd3, 0x3c,
	0x14, 0xc6, 0x95, 0xbd, 0x6b, 0x9b, 0x9e, 0x6d, 0xef, 0xb6, 0xb3, 0x6e, 0xf3, 0xfe, 0x47, 0x13,
	0xa0, 0x08, 0x21, 0x34, 0x8d, 0x4b, 0xc4, 0x0d, 0x2b, 0x84, 0x0a, 0x34, 0x50, 0x86, 0xb4, 0xcb,
	0xc8, 0x8b, 0x5d, 0x1a, 0x91, 0xc4, 0x91, 0xe3, 0x5e, 0xf4, 0x7b, 0xf3, 0x01, 0x90, 0x8f, 0xd3,
	0x91, 0x6c, 0x83, 0xbb, 0xfa, 0xf9, 0x3d, 0xcf, 0x39, 0x39, 0x47, 0x76, 0x61, 0xad, 0x36, 0xdc,
	0xd4, 0xaf, 0x2b, 0xad, 0x8c, 0xc2, 0x9e, 0x50, 0x86, 0x5f, 0x9e, 0xff, 0xf2, 0xa1, 0x77, 0x63,
	0x65, 0x3c, 0x00, 0xbf, 0xe0, 0x26, 0x9d, 0x25, 0x99, 0x60, 0x5e, 0xe0, 0x85, 0xab, 0xf1, 0x80,
	0xce, 0x13, 0x81, 0x27, 0x00, 0x3c, 0x4d, 0xd5, 0xbc, 0x34, 0x16, 0xae, 0x04, 0x5e, 0xb8, 0x11,
	0x0f, 0x1b, 0x65, 0x22, 0x90, 0xc1, 0xa0, 0x36, 0x92, 0x17, 0x99, 0x60, 0x9b, 0x2e, 0xd8, 0x1c,
	0xf1, 0x0c, 0xd6, 0xaa, 0x9c, 0x2f, 0xa4, 0x4e, 0x4a, 0x5e, 0x48, 0xf6, 0x5f, 0xe0, 0x85, 0xc3,
	0x18, 0x9c, 0x74, 0xcd, 0x0b, 0x89, 0xfb, 0x30, 0x98, 0x49, 0xad, 0x6c, 0xd9, 0x55, 0x2a, 0xdb,
	0xb7, 0xc7, 0x89, 0xc0, 0x23, 0x18, 0x12, 0xa0, 0x5c, 0x8f, 0x72, 0xbe, 0x15, 0x28, 0x75, 0x01,
	0xa3, 0x54, 0x4b, 0x6e, 0x64, 0x62, 0x94, 0xe1, 0x79, 0x22, 0x78, 0xc1, 0x7f, 0xc8, 0x9a, 0xf5,
	0xa9, 0x04, 0x3a, 0xf6, 0xdd, 0xa2, 0xb1, 0x23, 0x78, 0x09, 0xbb, 0x4d, 0x42, 0x48, 0x2e, 0xf2,
	0xc5, 0x7d, 0x64, 0x40, 0x91, 0x1d, 0x07, 0xc7, 0xc4, 0x96, 0x99, 0xb7, 0x70, 0xd8, 0xe9, 0x52,
	0x9b, 0x6c, 0x3a, 0x4d, 0x52, 0x55, 0x1a, 0xad, 0x72, 0xe6, 0x07, 0x5e, 0xb8, 0x12, 0xef, 0xb7,
	0x7a, 0xdd, 0x58, 0x7e, 0xe5, 0x30, 0xbe, 0x83, 0xa3, 0x6e, 0xc3, 0x6e, 0x7a, 0x48, 0x69, 0xd6,
	0x6e, 0xdb, 0x89, 0x5f, 0xc0, 0x48, 0x55, 0x95, 0x2a, 0x65, 0x69, 0x12, 0xda, 0x83, 0x90, 0xdc,
	0xcc, 0x6a, 0x06, 0x6e, 0xc2, 0x25, 0xfb, 0x24, 0xb5, 0x1a, 0x13, 0xc1, 0x31, 0x9c, 0x3d, 0x39,
	0x61, 0x52, 0x49, 0xed, 0xd2, 0x6c, 0x8d, 0x9a, 0x1e, 0x3d, 0x31, 0xeb, 0x37, 0xa9, 0xa9, 0x0c,
	0x7e, 0x85, 0xe7, 0xff, 0xf8, 0xec, 0x56, 0xad, 0x75, 0xaa, 0x15, 0xfc, 0x6d, 0x80, 0xfb, 0x82,
	0x2f, 0xe0, 0x7f, 0x9e, 0xab, 0x52, 0x7e, 0xce, 0xf2, 0x5c, 0x8a, 0xeb, 0x79, 0xc1, 0xb6, 0x68,
	0x84, 0x07, 0x2a, 0xbe, 0x84, 0x2d, 0x52, 0xde, 0xb7, 0x9c, 0xdb, 0xe4, 0x7c, 0xa4, 0xe3, 0x2b,
	0xd8, 0x6e, 0xb4, 0x2b, 0x7b, 0x41, 0x9d, 0x19, 0x03, 0x2f, 0xec, 0xc5, 0x8f, 0x01, 0x22, 0xac,
	0xea, 0xa8, 0x2a, 0xd8, 0x0e, 0x55, 0xa3, 0xdf, 0xb8, 0x07, 0xfd, 0x79, 0xa9, 0xe3, 0xaa, 0x60,
	0x23, 0x77, 0xeb, 0xdc, 0x09, 0xcf, 0x61, 0xfd, 0x67, 0x96, 0xe7, 0x76, 0xad, 0x91, 0xca, 0x05,
	0xdb, 0x25, 0xda, 0xd1, 0xac, 0xc7, 0xee, 0xe6, 0x8b, 0xaa, 0x25, 0x79, 0xf6, 0x9c, 0xa7, 0xad,
	0x59, 0xcf, 0x54, 0x8a, 0x0f, 0xa5, 0x2c, 0x16, 0xe4, 0xd9, 0x77, 0x9e, 0xb6, 0x86, 0xa7, 0x00,
	0xf6, 0x95, 0x5c, 0xcf, 0x8b, 0x3b, 0xa9, 0x19, 0xa3, 0xcf, 0x6f, 0x29, 0x38, 0x82, 0x5e, 0x56,
	0xdf, 0x66, 0x25, 0x3b, 0x08, 0xbc, 0xd0, 0x8f, 0xdd, 0xc1, 0xbe, 0x8b, 0xe6, 0x45, 0x65, 0x82,
	0x1d, 0x52, 0xc8, 0x77, 0xc2, 0x44, 0xe0, 0x31, 0x0c, 0xb3, 0xfa, 0x63, 0xa6, 0x6b, 0x13, 0x45,
	0xec, 0x98, 0xe0, 0x1f, 0xc1, 0xd1, 0x5b, 0x9d, 0x19, 0x19, 0x45, 0xec, 0x64, 0x49, 0x1b, 0xc1,
	0xd1, 0x58, 0x56, 0xf9, 0x22, 0x8a, 0xd8, 0xe9, 0x92, 0x36, 0x02, 0x86, 0xb0, 0x49, 0x9b, 0x6d,
	0x2d, 0xfc, 0x8c, 0x3c, 0x0f, 0x65, 0x7c, 0x06, 0x1b, 0xa9, 0x2a, 0xeb, 0x79, 0x21, 0xdd, 0xdd,
	0x62, 0x01, 0xf9, 0xba, 0xe2, 0x5d, 0x9f, 0xfe, 0x84, 0xde, 0xfc, 0x0e, 0x00, 0x00, 0xff, 0xff,
	0xfb, 0x5d, 0x30, 0xe7, 0x93, 0x04, 0x00, 0x00,
}
