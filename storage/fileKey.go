package storage

import (
	"errors"
	"sync"

	"gitlab.com/SeaStorage/SeaStorage-TP/crypto"
)

type FileKey struct {
	Index string
	Used  int
	Key   string
}

type FileKeyMap struct {
	sync.RWMutex
	keys []*FileKey
}

func NewFileKey(key string) *FileKey {
	return &FileKey{Index: crypto.SHA512HexFromHex(key), Key: key, Used: 0}
}

func NewFileKeyMap() *FileKeyMap {
	return &FileKeyMap{keys: make([]*FileKey, 0)}
}

func (fkm *FileKeyMap) GetKey(index string) *FileKey {
	fkm.Lock()
	defer fkm.Unlock()
	for _, key := range fkm.keys {
		if key.Index == index {
			return key
		}
	}
	return nil
}

func (fkm *FileKeyMap) SearchKey(key string) *FileKey {
	index := crypto.SHA512HexFromHex(key)
	fkm.Lock()
	defer fkm.Unlock()
	for _, key := range fkm.keys {
		if key.Index == index {
			return key
		}
	}
	return nil
}

func (fkm *FileKeyMap) AddKey(key string, used bool) string {
	fkm.Lock()
	defer fkm.Unlock()
	index := crypto.SHA512HexFromHex(key)
	for _, fileKey := range fkm.keys {
		if fileKey.Index == index {
			if used {
				fileKey.Used++
			}
			return index
		}
	}
	fileKey := NewFileKey(key)
	if used {
		fileKey.Used++
	}
	fkm.keys = append(fkm.keys, fileKey)
	return index
}

func (fkm *FileKeyMap) UpdateKeyUsed(keyUsed map[string]int) {
	fkm.Lock()
	defer fkm.Unlock()
	for index, used := range keyUsed {
		var fileIndex int
		var fileKey *FileKey
		for i, key := range fkm.keys {
			if key.Index == index {
				fileIndex = i
				fileKey = key
				break
			}
		}
		if fileKey == nil {
			continue
		}
		fileKey.Used += used
		if fileKey.Used <= 0 {
			fkm.keys = append(fkm.keys[:fileIndex], fkm.keys[fileIndex+1:]...)
		}
	}
}

func (fkm *FileKeyMap) PublishKey(publicKey, key string) error {
	keyBytes := crypto.AESKeyEncryptedByPublicKey(key, publicKey)
	fileKey := fkm.SearchKey(crypto.BytesToHex(keyBytes))
	if fileKey != nil {
		if crypto.AESKeyVerify(publicKey, key, fileKey.Key) {
			fileKey.Key = key
			return nil
		}
	}
	return errors.New("invalid key or not exists")
}
