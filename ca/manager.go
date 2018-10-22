package ca

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509/pkix"
	"errors"

	"github.com/pbergman/caserver/config"
	"github.com/pbergman/caserver/storage"
)

func NewManager(config *config.Config, db storage.Storage) (*Manager, error) {

	manager := &Manager{db, NewFactory(config.PemNotAfter, config.CaNotAfter, nil), config, nil}

	if err := manager.Init(); err != nil {
		return nil, err
	} else {
		return manager, nil
	}
}

type Manager struct {
	storage storage.Storage
	factory FactoryInterface
	config  *config.Config
	ca      *storage.StorageKey
}

// Search will do a search based on the `CommonName` and return nil
// if none were found.
func (m *Manager) Search(cn string) storage.Record {
	return m.storage.Search(cn)
}

// Lookup will search based on a short hash
func (m *Manager) Lookup(id string) storage.Record {
	record, _ := m.storage.Lookup(id)
	return record
}

func (m *Manager) GetFactory() FactoryInterface {
	return m.factory
}

func (m *Manager) Get(s *storage.StorageKey) storage.Record {
	if r, e := m.storage.Open(s); e == nil && r != nil {
		return r
	} else {
		return nil
	}
}

func (m *Manager) Each(c func(storage.Record) bool) error {
	return m.storage.Each(c)
}

func (m *Manager) Save(r storage.Record) (*storage.StorageKey, error) {
	return m.storage.Persist(r)
}

func (m *Manager) Remove(key *storage.StorageKey) error {
	return m.storage.Remove(key)
}

func (m *Manager) GetCa() *storage.StorageKey {
	return m.ca
}

func (m *Manager) NewRecord() storage.Record {
	return m.storage.NewRecord()
}

func (m *Manager) SignCertificateRequest(csr, ca storage.Record) error {
	cert, err := m.factory.NewCertificate(csr.GetCertificateRequest(), ca.GetCertificate(), ca.GetPrivateKey())
	if err != nil {
		return err
	}
	csr.SetCertificate(cert)
	if _, err := m.storage.Persist(csr); err != nil {
		return err
	}
	return nil
}

func (m *Manager) NewCertificateRequest(hosts []string, subject pkix.Name, bits int) (storage.Record, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	csr, err := m.factory.NewCertificateRequest(key, subject, hosts)
	if err != nil {
		return nil, err
	}
	record := m.storage.NewRecord()
	record.SetPrivateKey(key)
	record.SetCertificateRequest(csr)
	if _, err := m.storage.Persist(record); err != nil {
		return nil, err
	} else {
		return record, nil
	}
}

// Init will do some check and setups for the manager, this manager only supports on
// active CA so if more than one is found it will return a error and it will create
// new certificates if none were found.
func (m *Manager) Init() error {
	if m.storage == nil {
		return errors.New("missing storage interface")
	}
	if m.config == nil {
		return errors.New("missing config")
	}
	list := m.storage.GetCa()
	switch len(list) {
	case 0:
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return err
		}
		cert, err := m.factory.NewCertificateAuthority(key, *m.config.CaSubject)
		if err != nil {
			return err
		}
		record := m.storage.NewRecord()
		record.SetPrivateKey(key)
		record.SetCertificate(cert)
		if kid, err := m.storage.Persist(record); err != nil {
			return err
		} else {
			m.ca = kid
		}
	case 1:
		m.ca = list[0]
	default:
		return errors.New("to many CA certificates found")
	}
	return nil
}
