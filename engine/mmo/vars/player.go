package vars
const (
    // PlayerPos int32数组;坐标X、坐标Y、坐标Z
    PlayerPos string = "pos"
    // PlayerLook 朝向，int16
    PlayerLook string = "lk"
    // PlayerInputMove int32数组;输入X、输入Y、输入Z
    PlayerInputMove string = "im"
    // PlayerInputTarget int32数组;目标X、目标Y、目标Z
    PlayerInputTarget string = "it"
    // PlayerInputJump 输入状态Jump(bool)
    PlayerInputJump string = "ij"
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