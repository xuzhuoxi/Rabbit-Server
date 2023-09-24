// Package basis
// Created by xuzhuoxi
// on 2019-02-19.
// @author xuzhuoxi
package basis

import "math/rand"

var (
	ZeroXY  = XY{}
	ZeroXYZ = XYZ{}
)

func RandomXYZ() XYZ {
	rs := XYZ{
		X: rand.Int31n(200) - 100,
		Y: 0,
		Z: rand.Int31n(200) - 100,
	}
	return rs
}

type XY struct {
	X int32
	Y int32
}

type XYZ struct {
	X int32
	Y int32
	Z int32
}

func (xyz XYZ) XY() XY {
	return XY{X: xyz.X, Y: xyz.Y}
}

// NearXY 判断两点是否相近
// 用于转发附近消息
func NearXY(pos1 XY, pos2 XY, distance int32) bool {
	if distance == 0 {
		return true
	}
	x12 := pos1.X - pos2.X
	y12 := pos1.Y - pos2.Y
	return (x12*x12 + y12*y12) <= distance*distance
}

// NearXYZ 判断两点是否相近
// 用于转发附近消息
func NearXYZ(pos1 XYZ, pos2 XYZ, distance int32) bool {
	if distance == 0 {
		return true
	}
	x12 := pos1.X - pos2.X
	y12 := pos1.Y - pos2.Y
	z12 := pos1.Z - pos2.Z
	return (x12*x12 + y12*y12 + z12*z12) <= distance*distance
}
