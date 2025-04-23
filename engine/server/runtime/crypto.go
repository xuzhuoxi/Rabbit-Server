// Package runtime
// Create on 2025/4/22
// @author xuzhuoxi
package runtime

import (
	"github.com/xuzhuoxi/infra-go/cryptox"
	"sync"
)

type CryptoSupport struct {
	Cipher      cryptox.ICipher
	CipherMutex sync.RWMutex
}

func (o *CryptoSupport) SetPacketCipher(cipher cryptox.ICipher) {
	o.CipherMutex.Lock()
	defer o.CipherMutex.Unlock()
	o.Cipher = cipher
}

func (o *CryptoSupport) DecryptPacket(msgBytes []byte) ([]byte, error) {
	if nil == o.Cipher {
		return msgBytes, nil
	}
	o.CipherMutex.RLock()
	defer o.CipherMutex.RUnlock()
	rs, err := o.Cipher.Decrypt(msgBytes)
	if nil != err {
		return msgBytes, err
	}
	return rs, nil
}
