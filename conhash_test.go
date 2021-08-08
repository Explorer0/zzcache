package zzcache

import (
	"strconv"
	"testing"
)

func TestNewDistributeMap(t *testing.T) {
	hash := NewDistributeMap(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})


	hash.AddNode("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.GetNode(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// Adds 8, 18, 28
	hash.AddNode("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.GetNode(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	t.Logf("test consistent hash successfully!")
}
