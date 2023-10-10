// Package vars
// Create on 2023/9/26
// @author xuzhuoxi
package vars

type VarType int32

const (
	// VarNone
	// Undefined
	// 未定义
	VarNone VarType = iota
	// VarSoon
	// Soon
	// 瞬间
	VarSoon
	// VarTime
	// keep some time
	// 持续一段时间
	VarTime
	// VarForever
	// Forever until overwritten.
	// 永久直到被覆盖
	VarForever
)

type VarMeta struct {
	Key     string
	Type    VarType
	Default interface{}
}

var (
	roomVarMetas   map[string]VarMeta
	playerVarMetas map[string]VarMeta
	unitVarMetas   map[string]VarMeta
)

func GetRoomVarMeta(key string) (t VarMeta, ok bool) {
	if m, b := roomVarMetas[key]; b {
		return m, b
	}
	return
}

func GetPlayerVarMeta(key string) (t VarMeta, ok bool) {
	if m, b := playerVarMetas[key]; b {
		return m, b
	}
	return
}
func GetUnitVarMeta(key string) (t VarMeta, ok bool) {
	if m, b := unitVarMetas[key]; b {
		return m, b
	}
	return
}

func RegisterRoomVarMeta(meta VarMeta) {
	roomVarMetas[meta.Key] = meta
}

func RegisterPlayerVarMeta(meta VarMeta) {
	playerVarMetas[meta.Key] = meta
}

func RegisterUnitVarMeta(meta VarMeta) {
	unitVarMetas[meta.Key] = meta
}

func init() {
	initPlayerVarMetas()
	initUnitVarMetas()
}

func initPlayerVarMetas() {
	playerVarMetas = make(map[string]VarMeta)
	playerVarMetas[PlayerPos] = VarMeta{Key: PlayerPos, Type: VarForever, Default: []int32{0, 0, 0}}
	playerVarMetas[PlayerLook] = VarMeta{Key: PlayerLook, Type: VarForever, Default: int16(0)}
	playerVarMetas[PlayerInputMove] = VarMeta{Key: PlayerInputMove, Type: VarForever, Default: []int32{0, 0, 0}}
	playerVarMetas[PlayerInputJump] = VarMeta{Key: PlayerInputJump, Type: VarSoon, Default: false}

	playerVarMetas[PlayerActionState] = VarMeta{Key: PlayerActionState, Type: VarSoon, Default: uint32(0)}

	playerVarMetas[PlayerHp] = VarMeta{Key: PlayerHp, Type: VarForever, Default: uint32(0)}
	playerVarMetas[PlayerBuff] = VarMeta{Key: PlayerBuff, Type: VarForever, Default: uint32(0)}
	playerVarMetas[PlayerNick] = VarMeta{Key: PlayerNick, Type: VarForever, Default: ""}
	playerVarMetas[PlayerTeam] = VarMeta{Key: PlayerTeam, Type: VarForever, Default: ""}
	playerVarMetas[PlayerTeamCorps] = VarMeta{Key: PlayerTeamCorps, Type: VarForever, Default: ""}
}

func initUnitVarMetas() {
	unitVarMetas = make(map[string]VarMeta)
	unitVarMetas[UnitOwner] = VarMeta{Key: PlayerPos, Type: VarForever, Default: ""}
	unitVarMetas[UnitRoom] = VarMeta{Key: PlayerInputMove, Type: VarForever, Default: ""}
	unitVarMetas[UnitPos] = VarMeta{Key: PlayerInputMove, Type: VarForever, Default: []int32{0, 0, 0}}
	unitVarMetas[UnitLook] = VarMeta{Key: PlayerInputJump, Type: VarForever, Default: int16(0)}

	unitVarMetas[UnitInputMove] = VarMeta{Key: UnitInputMove, Type: VarForever, Default: []int32{0, 0, 0}}
	unitVarMetas[UnitInputJump] = VarMeta{Key: PlayerInputJump, Type: VarSoon, Default: false}
}
