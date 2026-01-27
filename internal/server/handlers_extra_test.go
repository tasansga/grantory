package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tasansga/grantory/internal/storage"
)

func TestStoreFromLocalsVariants(t *testing.T) {
	t.Parallel()

	st := &storage.Store{}
	assert.Equal(t, st, storeFromLocals(st))
	assert.Equal(t, st, storeFromLocals(localStore{store: st}))
	assert.Equal(t, st, storeFromLocals(&localStore{store: st}))
	assert.Nil(t, storeFromLocals("unexpected"))
}
