package main

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/Factom-Asset-Tokens/base58"
)

func main() {
	var key [32]byte

	found := make(map[string][]byte)
	var foundMtx sync.RWMutex

	var wg sync.WaitGroup
	for i := uint16(0); i < 0xFFFF; i++ {
		prefix := []byte{byte(i >> 8), byte(i)}
		adrStr := base58.CheckEncode(key[:], prefix...)
		prefixStr := adrStr[:2]
		switch prefixStr {
		case "FA", "FE", "Fe", "Fs", "EC", "Es":
			key := key
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i <= 0xFF; i++ {
					key[0] = byte(i)
					adrStr := base58.CheckEncode(key[:], prefix...)
					if adrStr[:2] != prefixStr {
						return
					}
				}
				foundMtx.RLock()
				lastPrefix, ok := found[prefixStr]
				foundMtx.RUnlock()

				if !ok {
					lastPrefix = []byte{0xff, 0xff}
				}

				if bytes.Compare(prefix, lastPrefix) >= 0 {
					return
				}

				foundMtx.Lock()
				found[prefixStr] = prefix
				foundMtx.Unlock()
			}()
		}
	}
	wg.Wait()

	for _, prefixStr := range []string{"EC", "Es", "FA", "FE", "Fe", "Fs"} {
		fmt.Printf("%s 0x%x\n", prefixStr, found[prefixStr])
	}

}
