package mmkv

type Manager interface {
	// OpenVault opens a vault with the given id.
	// If the vault does not exist, it will be created.
	// If id is empty, DefaultVaultID will be used.
	OpenVault(id string) (Vault, error)
}

type Vault interface {
	Keys() []string
	GetRaw(key string) ([]byte, bool)
	GetBytes(key string) ([]byte, error)
	GetString(key string) (string, error)
}
