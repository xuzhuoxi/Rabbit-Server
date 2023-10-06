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
	playerVarMetas map[string]VarMeta
	roomVarMetas   map[string]VarMeta
)

func GetPlayerVarMeta(key string) (t VarMeta, ok bool) {
	if m, b := playerVarMetas[key]; b {
		return m, b
	}
	return
}

func GetRoomVarMeta(key string) (t VarMeta, ok bool) {
	if m, b := roomVarMetas[key]; b {
		return m, b
	}
	return
}

func RegisterPlayerVarMeta(meta VarMeta) {
	playerVarMetas[meta.Key] = meta
}

func RegisterRoomVarMeta(meta VarMeta) {
	roomVarMetas[meta.Key] = meta
}

func init() {
	playerVarMetas = make(map[string]VarMeta)
	playerVarMetas[PlayerPos] = VarMeta{Key: PlayerPos, Type: VarForever, Default: []int32{0, 0, 0}}
	playerVarMetas[PlayerInput] = VarMeta{Key: PlayerInput, Type: VarForever, Default: []int32{0, 0, 0}}
	playerVarMetas[PlayerInputJump] = VarMeta{Key: PlayerInputJump, Type: VarSoon, Default: false}

	playerVarMetas[PlayerFace] = VarMeta{Key: PlayerFace, Type: VarForever, Default: uint8(0)}
	playerVarMetas[PlayerActionState] = VarMeta{Key: PlayerActionState, Type: VarSoon, Default: uint32(0)}

	playerVarMetas[PlayerHp] = VarMeta{Key: PlayerHp, Type: VarForever, Default: uint32(0)}
	playerVarMetas[PlayerBuff] = VarMeta{Key: PlayerBuff, Type: VarForever, Default: uint32(0)}
	playerVarMetas[PlayerNick] = VarMeta{Key: PlayerNick, Type: VarForever, Default: ""}
	playerVarMetas[PlayerTeam] = VarMeta{Key: PlayerTeam, Type: VarForever, Default: ""}
	playerVarMetas[PlayerTeamCorps] = VarMeta{Key: PlayerTeamCorps, Type: VarForever, Default: ""}
}
