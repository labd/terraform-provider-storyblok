package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	p := New()

	assert.NotNil(t, p)

}
