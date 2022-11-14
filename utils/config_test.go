package utils

import "testing"

var testPath1 = "classpath:aplication.toml"

func TestReadClassPath(t *testing.T) {
	t.Logf("got path:\n%s\nafter process:\n%s", testPath1, ReadClassPath(testPath1))
}
