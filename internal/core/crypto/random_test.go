package crypto

import "testing"

func TestPasswordGen(t *testing.T) {
	t.Log(RandomPassword(10))
	t.Log(RandomAlphaNumericPassword(10))
}
