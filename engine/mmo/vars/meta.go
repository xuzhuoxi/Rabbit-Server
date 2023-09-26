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
	Type    VarType
	Default interface{}
}

var (
	playerVarMetas map[string]VarMeta
)

func GetPlayerVarType(key string) (t VarType, ok bool) {
	if m, b := playerVarMetas[key]; b {
		return m.Type, true
	}
	return
}

func GetPlayerVarDefault(key string) (v interface{}, ok bool) {
	if m, b := playerVarMetas[key]; b {
		return m.Default, true
	}
	return
}

func init() {
	playerVarMetas = make(map[string]VarMeta)
	playerVarMetas[PlayerPosX] = VarMeta{Type: VarForever, Default: int32(0)}
	playerVarMetas[PlayerPosY] = VarMeta{Type: VarForever, Default: int32(0)}
	playerVarMetas[PlayerPosZ] = VarMeta{Type: VarForever, Default: int32(0)}
	playerVarMetas[PlayerInputX] = VarMeta{Type: VarForever, Default: int32(0)}
	playerVarMetas[PlayerInputY] = VarMeta{Type: VarForever, Default: int32(0)}
	playerVarMetas[PlayerInputZ] = VarMeta{Type: VarForever, Default: int32(0)}
	playerVarMetas[PlayerInputJump] = VarMeta{Type: VarSoon, Default: false}

	playerVarMetas[PlayerFace] = VarMeta{Type: VarForever, Default: uint8(0)}
	playerVarMetas[PlayerAction] = VarMeta{Type: VarForever, Default: uint32(0)}

	playerVarMetas[PlayerHp] = VarMeta{Type: VarForever, Default: uint32(0)}
	playerVarMetas[PlayerBuff] = VarMeta{Type: VarForever, Default: uint32(0)}
	playerVarMetas[PlayerNick] = VarMeta{Type: VarForever, Default: ""}
	playerVarMetas[PlayerTeam] = VarMeta{Type: VarForever, Default: ""}
	playerVarMetas[PlayerTeamCorps] = VarMeta{Type: VarForever, Default: ""}
}
