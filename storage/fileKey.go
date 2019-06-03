package storage

import (
	"errors"
	"github.com/yellowssi/SeaStorage-TP/crypto"
)

// FileKey store the information of key used to encrypt file.
type FileKey struct {
	Index string
	Used  int
	Key   string
}

// FileKeyMap provides file keys manage.
type FileKeyMap struct {
	Keys []*FileKey
}

// NewFileKeyMap is the construct for FileKeyMap.
func NewFileKeyMap() *FileKeyMap {
	return &FileKeyMap{Keys: make([]*FileKey, 0)}
}

// GetKey search the FileKey by index.
// If it exists, it will be return.
// Else, returns nil.
func (fkm *FileKeyMap) GetKey(index string) *FileKey {
	for _, key := range fkm.Keys {
		if key.Index == index {
			return key
		}
	}
	return nil
}

// AddKey add new information of key.
// If used, the used count of key will be 1.
// Else, it will be 0.
func (fkm *FileKeyMap) AddKey(key string, used bool) string {
	index := crypto.SHA512HexFromHex(key)
	for _, fileKey := range fkm.Keys {
		if fileKey.Index == index {
			if used {
				fileKey.Used++
			}
			return index
		}
	}
	var fileKey *FileKey
	if used {
		fileKey = &FileKey{Index: index, Key: key, Used: 1}
	} else {
		fileKey = &FileKey{Index: index, Key: key, Used: 0}
	}
	fkm.Keys = append(fkm.Keys, fileKey)
	return index
}

// UpdateKeyUsed update used count of keys by index and count.
func (fkm *FileKeyMap) UpdateKeyUsed(keyUsed map[string]int) {
	for index, used := range keyUsed {
		var fileIndex int
		var fileKey *FileKey
		for i, key := range fkm.Keys {
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
			fkm.Keys = append(fkm.Keys[:fileIndex], fkm.Keys[fileIndex+1:]...)
		}
	}
}

// PublishKey check key whether valid and publish it.
func (fkm *FileKeyMap) PublishKey(publicKey, keyIndex, key string) error {
	fileKey := fkm.GetKey(keyIndex)
	if fileKey != nil {
		fileKey.Key = key
		return nil
	}
	return errors.New("invalid key or not exists")
}
