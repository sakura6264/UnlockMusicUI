package mmkv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	mgr, err := NewManager("./testdata")
	assert.NoError(t, err)
	assert.NotNil(t, mgr)

	vault, err := mgr.OpenVault("")
	assert.NoError(t, err)
	assert.NotNil(t, vault)
}
