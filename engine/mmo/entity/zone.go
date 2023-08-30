// Package entity
// Created by xuzhuoxi
// on 2019-02-18.
// @author xuzhuoxi
package entity

//
//type ZoneConfig struct {
//}
//
//func NewIZoneEntity(zoneId string, zoneName string) basis.IZoneEntity {
//	return &ZoneEntity{ZoneId: zoneId, ZoneName: zoneName}
//}
//
//func NewZoneEntity(zoneId string, zoneName string) *ZoneEntity {
//	return &ZoneEntity{ZoneId: zoneId, ZoneName: zoneName}
//}
//
//type ZoneEntity struct {
//	ZoneId   string
//	ZoneName string
//	EntityChildSupport
//	ListEntityContainer
//
//	//RoomGroup *EntityListGroup
//	VariableSupport
//}
//
//func (o *ZoneEntity) UID() string {
//	return o.ZoneId
//}
//
//func (o *ZoneEntity) Name() string {
//	return o.ZoneName
//}
//
//func (o *ZoneEntity) EntityType() basis.EntityType {
//	return basis.EntityZone
//}
//
//func (o *ZoneEntity) InitEntity() {
//	o.EntityChildSupport = *NewEntityChildSupport()
//	o.ListEntityContainer = *NewListEntityContainer(0)
//	//e.RoomGroup = NewEntityListGroup(EntityRoom)
//	o.VariableSupport = *NewVariableSupport(o)
//}
//
////func (e *ZoneEntity) RoomList() []string {
////	return e.RoomGroup.Entities()
////}
////
////func (e *ZoneEntity) ContainRoom(roomId string) bool {
////	return e.RoomGroup.ContainEntity(roomId)
////}
////
////func (e *ZoneEntity) AddRoom(roomId string) error {
////	return e.RoomGroup.Accept(roomId)
////}
////
////func (e *ZoneEntity) RemoveRoom(roomId string) error {
////	return e.RoomGroup.Drop(roomId)
////}
