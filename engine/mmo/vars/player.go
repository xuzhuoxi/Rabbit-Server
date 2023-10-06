package vars
const (
    // PlayerPos int32数组;坐标X、坐标Y、坐标Z
    PlayerPos string = "pos"
    // PlayerInput int32数组;输入X、输入Y、输入Z
    PlayerInput string = "ip"
    // PlayerInputJump 输入状态Jump(bool)
    PlayerInputJump string = "ij"
    // PlayerFace 面向(uint8)(8方向)
    PlayerFace string = "f"
    // PlayerMoveState 移动状态(uint32)
    PlayerMoveState string = "ms"
    // PlayerActionState 动作状态(uint32)
    PlayerActionState string = "as"
    // PlayerHp 耐久(uint32)
    PlayerHp string = "hp"
    // PlayerBuff Buff(uint32), 每一个位代表一个buff
    PlayerBuff string = "pbf"
    // PlayerNick 昵称(string)
    PlayerNick string = "pn"
    // PlayerTeam 队伍id(string)
    PlayerTeam string = "pt"
    // PlayerTeamCorps 军团Id(string)
    PlayerTeamCorps string = "pc"
)