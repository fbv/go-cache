package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeek(t *testing.T) {
	c := New[string](NoExpiration[string]())
	c.Put("key", "value")
	v, ok := c.Peek("key")
	assert.True(t, ok)
	assert.Equal(t, "value", v)
	v, ok = c.Peek("no key")
	assert.False(t, ok)
	assert.Equal(t, "", v)
}

func TestRemove(t *testing.T) {
	c := New[string](NoExpiration[string]())
	c.Put("k1", "v1")
	c.Put("k2", "v2")
	c.Put("k3", "v3")
	c.Remove("k1", "k3")
	v, ok := c.Peek("k1")
	assert.False(t, ok)
	assert.Equal(t, "", v)
	v, ok = c.Peek("k2")
	assert.True(t, ok)
	assert.Equal(t, "v2", v)
	v, ok = c.Peek("k3")
	assert.False(t, ok)
	assert.Equal(t, "", v)
}
