package haproxyconfigparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeparateHostAndPortValidCase(t *testing.T) {
	host, port, err := SeparateHostAndPort("0.0.0.0:8080")
	assert.Equal(t, host, "0.0.0.0")
	assert.Equal(t, port, 8080)
	assert.Nil(t, err)
}

func TestSeparateHostAndPortInvalidCase(t *testing.T) {
	_, _, err := SeparateHostAndPort("invalid")
	assert.NotNil(t, err)

	_, _, err = SeparateHostAndPort("0.0.0.0:aaaaa")
	assert.NotNil(t, err)

	_, _, err = SeparateHostAndPort("0.0.0.0:8080:9080")
	assert.NotNil(t, err)
}

func TestUncomment(t *testing.T) {
	ret, enable := Uncomment("  this is normal line  ")
	assert.True(t, enable)
	assert.Equal(t, ret, "this is normal line")
}

func TestUncommentWithCommenting(t *testing.T) {
	ret, enable := Uncomment("  this is line with commenting # here is comment ")
	assert.True(t, enable)
	assert.Equal(t, ret, "this is line with commenting")
}

func TestUncommentWithCommentout(t *testing.T) {
	ret, enable := Uncomment("  #this is a comment-outed line")
	assert.False(t, enable)
	assert.Equal(t, ret, "this is a comment-outed line")
}

func TestSeparateConfigLine(t *testing.T) {
	ret, enable := SeparateConfigLine("acl reform   src 10.8.8.25   tab # glb-dev16")
	assert.True(t, enable)
	assert.Equal(t, len(ret), 5)
	assert.Equal(t, ret[0], "acl")
	assert.Equal(t, ret[1], "reform")
	assert.Equal(t, ret[2], "src")
	assert.Equal(t, ret[3], "10.8.8.25")
	assert.Equal(t, ret[4], "tab")
}

func TestSeparateConfigLineComment(t *testing.T) {
	ret, enable := SeparateConfigLine("  # this is comment line")
	assert.False(t, enable)
	assert.Equal(t, len(ret), 4)

	ret, enable = SeparateConfigLine("#server glb-api4  10.8.2.74:9080")
	assert.False(t, enable)
	assert.Equal(t, len(ret), 3)
	assert.Equal(t, ret[0], "server")

	ret, enable = SeparateConfigLine("server glb-api4#10.8.2.74:9080")
	assert.True(t, enable)
	assert.Equal(t, len(ret), 2)
	assert.Equal(t, ret[1], "glb-api4")
}
