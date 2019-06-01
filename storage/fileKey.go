package storage

import (
	"errors"

	"gitlab.com/SeaStorage/SeaStorage-TP/crypto"
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

// NewFileKey is the construct for FileKey.
func NewFileKey(key string) *FileKey {
	return &FileKey{Index: crypto.SHA512HexFromHex(key), Key: key, Used: 0}
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

// SearchKey search the FileKey by key it self.
// If it exists, it will be return.
// Else, returns nil.
func (fkm *FileKeyMap) SearchKey(key string) *FileKey {
	index := crypto.SHA512HexFromHex(key)
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
	fileKey := NewFileKey(key)
	if used {
		fileKey.Used++
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
