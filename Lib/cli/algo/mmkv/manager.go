package mmkv

import (
	"fmt"
	"os"
	"path"
)

const (
	DefaultVaultID = "mmkv.default"
)

type manager struct {
	dir    string
	vaults map[string]Vault
}

// NewManager creates a new MMKV Manager.
func NewManager(dir string) (Manager, error) {
	// check dir exists
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to stat dir: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory")
	}

	return &manager{
		dir:    dir,
		vaults: make(map[string]Vault),
	}, nil
}

func (m *manager) OpenVault(id string) (Vault, error) {
	if id == "" {
		id = DefaultVaultID
	}

	if v, ok := m.vaults[id]; ok {
		return v, nil
	}

	vault, err := m.openVault(id)
	if err != nil {
		return nil, fmt.Errorf("failed to open vault: %w", err)
	}
	m.vaults[id] = vault

	return vault, nil
}

func (m *manager) openVault(id string) (Vault, error) {
	metaFile, err := os.Open(path.Join(m.dir, id+".crc"))
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}
	defer metaFile.Close()

	vaultFile, err := os.Open(path.Join(m.dir, id))
	if err != nil {
		return nil, fmt.Errorf("failed to open vault file: %w", err)
	}
	defer vaultFile.Close()

	meta, err := loadMetadata(metaFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}

	v, err := loadVault(vaultFile, meta)
	if err != nil {
		return nil, fmt.Errorf("failed to load vault: %w", err)
	}

	return v, nil
}
