//Created by xuzhuoxi
//on 2019-03-17.
//@author xuzhuoxi
package basis

import (
	"fmt"
	"testing"
)

func TestEntityType(t *testing.T) {
	fmt.Println(EntityNone, EntityRoom, EntityPlayer, EntityTeamCorps, EntityTeam, EntityChannel)
	fmt.Println("---")
	fmt.Println(EntityAll.Match(EntityRoom))
	fmt.Println(EntityRoom.Match(EntityAll))
	fmt.Println(EntityRoom.Match(EntityNone))
	fmt.Println("---")
	fmt.Println(EntityAll.Include(EntityRoom))
	fmt.Println(EntityRoom.Include(EntityAll))
	fmt.Println(EntityRoom.Include(EntityNone))
}
