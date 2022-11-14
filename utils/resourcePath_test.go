package utils

import "testing"

var resourcesPathTestCases = []struct {
	path    string
	isRight bool
}{
	{path: "testchain.org1.myresource", isRight: true},
	{path: "testchain.org1", isRight: false},
	{path: "testchain.org1.myresource.nonsense", isRight: false},
	{path: "testchain/org1/myresource", isRight: false},
	{path: "testchain-org1-myresource", isRight: false},
	{path: "testchain.org1!.myresource", isRight: false},
}

func TestCheckResourcePath(t *testing.T) {
	for index, testCase := range resourcesPathTestCases {
		if testCase.isRight != CheckResourcePath(testCase.path) {
			t.Fatalf("fail in case %d with path %s", index, testCase.path)
		}
	}
}
