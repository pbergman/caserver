package storage

import (
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/pbergman/caserver/util"
)

// DiskStorage is a Storage implementation that uses a
// directory filesystem as persistence storage.
type DiskStorage struct {
	records []*DiskRecord
	key     *[32]byte
	path    string
	lock    sync.RWMutex
}

func NewDiskStorage(path string, key *[32]byte) *DiskStorage {
	// noop function so we don`t check errors
	os.MkdirAll(path, 0700)
	return &DiskStorage{
		records: make([]*DiskRecord, 0),
		key:     key,
		path:    path,
	}
}

func (d *DiskStorage) walkNames(call func(string) bool) error {
	list, err := ioutil.ReadDir(d.path)
	if err != nil {
		return err
	}
	for i, c := 0, len(list); i < c; i++ {
		if list[i].Mode().IsRegular() {
			if false == call(list[i].Name()) {
				break
			}
		}
	}
	return nil
}

func (d *DiskStorage) Each(call func(Record) bool) error {
	errs := new(util.Errors)
	d.walkNames(func(name string) bool {
		if key := NewStorageKeyFromString(name); key != nil {
			if record, err := d.Open(key); err != nil {
				errs.Append(err)
			} else {
				return call(record)
			}
		}
		return true
	})
	if len(*errs) > 0 {
		return errs
	} else {
		return nil
	}
}

func (d *DiskStorage) Persist(r Record) (*StorageKey, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	if record, ok := r.(*DiskRecord); !ok {
		return nil, fmt.Errorf("invalid record type, expected *DiskRecord got %T", r)
	} else {

		if record.GetCertificate() != nil {
			if record.GetCertificate().Subject.CommonName == "" {
				return nil, errors.New("missing required CN field for certificate subject")
			}
		}

		if record.GetCertificateRequest() != nil {
			if record.GetCertificateRequest().Subject.CommonName == "" {
				return nil, errors.New("missing required CN field for certificate request subject")
			}
		}

		if pem := record.GetCertificate(); pem != nil {
			if pem.IsCA {
				record.mode |= MODE_IS_CA
			}
		}

		var file *os.File
		var err error

		if record.id == nil {
			file, err = ioutil.TempFile(d.path, "")
			if err != nil {
				return nil, err
			}
		} else {
			file, err = os.OpenFile(filepath.Join(d.path, record.id.String()), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				// should not happen, but file is removed
				// so we will re persist and remove ref.
				if os.IsNotExist(err) {
					record.id = nil
					return d.Persist(record)
				} else {
					return nil, err
				}
			}
		}

		hasher := sha1.New()
		writer := io.MultiWriter(file, hasher)

		defer file.Close()

		if raw, err := record.MarshalBinary(); err != nil {
			return nil, err
		} else {
			if _, err := writer.Write(raw); err != nil {
				return nil, err
			}
		}

		record.id = NewStorageKeyFromBytes(hasher.Sum(nil))
		return record.id, os.Rename(file.Name(), filepath.Join(d.path, record.id.String()))
	}
}

func (d *DiskStorage) Has(key *StorageKey) bool {
	d.lock.RLock()
	defer d.lock.RUnlock()
	stat, err := os.Stat(filepath.Join(d.path, key.String()))
	if err != nil {
		return false
	}
	return stat.Mode().IsRegular()
}

func (d *DiskStorage) Open(key *StorageKey) (Record, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	file, err := os.Open(filepath.Join(d.path, key.String()))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no record exist for key %s", key.String())
		} else {
			return nil, err
		}
	}
	if file != nil {
		defer file.Close()
		raw, buf := make([]byte, 0), make([]byte, 1024)
		for {
			n, err := file.Read(buf)
			raw = append(raw, buf[:n]...)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return nil, err
				}
			}
		}
		record := NewDiskRecord(d, key)
		if err := record.UnmarshalBinary(raw); err != nil {
			return nil, err
		} else {
			return record, nil
		}
	}
	return nil, nil
}

func (d *DiskStorage) Lookup(id string) (Record, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	var ret string
	if length := len(id); length <= 40 {
		d.walkNames(func(kid string) bool {
			if kid[:length] == id {
				ret = kid
				return false
			}
			return true
		})
	} else {
		ret = id[:40]
	}
	if ret == "" {
		return nil, nil
	} else {
		return d.Open(NewStorageKeyFromString(ret))
	}
}

func (d *DiskStorage) Search(cn string) (found Record) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	d.walkNames(func(kid string) bool {
		if record, _ := d.Open(NewStorageKeyFromString(kid)); record != nil {
			if pem := record.GetCertificate(); pem != nil {
				if pem.Subject.CommonName == cn {
					found = record
					return false
				}
			}
			if csr := record.GetCertificateRequest(); csr != nil {
				if csr.Subject.CommonName == cn {
					found = record
					return false
				}
			}
		}
		return true
	})
	return
}

func (d *DiskStorage) Remove(key *StorageKey) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	return os.Remove(filepath.Join(d.path, key.String()))
}

func (d *DiskStorage) GetCa() (list []*StorageKey) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	wg := new(sync.WaitGroup)
	d.walkNames(func(kid string) bool {
		if file, err := os.Open(filepath.Join(d.path, kid)); err == nil {
			wg.Add(1)
			go func(f *os.File) {
				b := make([]byte, 1)
				defer f.Close()
				defer wg.Done()
				if _, err := file.ReadAt(b, sha256.Size); err == nil {
					if MODE_IS_CA == (MODE_IS_CA & b[0]) {
						list = append(list, NewStorageKeyFromString(filepath.Base(file.Name())))
					}
				}
			}(file)
		}
		return true
	})
	wg.Wait()
	return
}

func (d *DiskStorage) NewRecord() Record {
	return NewDiskRecord(d, nil)
}
