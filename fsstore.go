package imapsql

import (
	"os"
	"path/filepath"
)

// FSStore struct represents directory on FS used to store message bodies.
//
// Always use field names on initialization because new fields may be added
// without a major version change.
type FSStore struct {
	Root string
}

func (s *FSStore) Open(key string) (ExtStoreObj, error) {
	f, err := os.Open(filepath.Join(s.Root, key))
	if err != nil {
		return nil, ExternalError{
			Key:         key,
			Err:         err,
			NonExistent: os.IsNotExist(err),
		}
	}
	return f, nil
}

func (s *FSStore) Create(key string, blobSize int64) (ExtStoreObj, error) {
	f, err := os.Create(filepath.Join(s.Root, key))
	if err != nil {
		return nil, ExternalError{
			Key:         key,
			Err:         err,
			NonExistent: false,
		}
	}
	if blobSize != -1 {
		if err := f.Truncate(blobSize); err != nil {
			return nil, ExternalError{
				Key: key,
				Err: err,
			}
		}
	}
	return f, nil
}

func (s *FSStore) Delete(keys []string) error {
	for _, key := range keys {
		if err := os.Remove(filepath.Join(s.Root, key)); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return ExternalError{
				Key: key,
				Err: err,
			}
		}
	}
	return nil
}
