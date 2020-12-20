package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Holder_Get_Non_Existing_ID_Should_Fail(t *testing.T) {
	holder := NewInMemoryHolder()

	game, err := holder.Get(ID(1))

	assert.Nil(t, game)
	assert.Equal(t, ErrGameDoesNotExist, err)
}

func Test_Holder_Insert_And_Get(t *testing.T) {
	holder := NewInMemoryHolder()

	given := Fake{}

	id, err := holder.Insert(given)

	assert.NoError(t, err)

	actual, err := holder.Get(id)

	assert.NoError(t, err)
	assert.Equal(t, given, actual)
}
