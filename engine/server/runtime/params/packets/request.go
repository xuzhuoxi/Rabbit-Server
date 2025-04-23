package packets

import "github.com/xuzhuoxi/infra-go/bytex"

type RequestPacket struct {
}

// ParsePacket
// block0 : eName utf8
// block1 : pid	utf8
// block2 : uid	utf8
// [n]其它信息
func (o *RequestPacket) ParsePacket(packetBytes []byte) (name string, pid string, uid string, data [][]byte) {
	if len(packetBytes) == 0 {
		return
	}
	index := 0
	buffToData := bytex.DefaultPoolBuffToData.GetInstance()
	defer bytex.DefaultPoolBuffToData.Recycle(buffToData)

	buffToData.WriteBytes(packetBytes)
	name = buffToData.ReadString()
	pid = buffToData.ReadString()
	uid = buffToData.ReadString()
	if buffToData.Len() > 0 {
		for buffToData.Len() > 0 {
			n, d := buffToData.ReadDataTo(packetBytes[index:]) //由于msgBytes前部分数据已经处理完成，可以利用这部分空间
			if nil == d {
				break
			}
			data = append(data, d)
			index += n
		}
	}
	return name, pid, uid, data
}
