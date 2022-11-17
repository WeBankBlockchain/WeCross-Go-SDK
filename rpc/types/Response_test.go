package types

import (
	"WeCross-Go-SDK/rpc/types/response"
	"bytes"
	"net/http"
	"testing"
)

type testCase struct {
	caseName     string
	responseJson string
	responseType response.ResponseType
}

var testCases = []testCase{
	{
		caseName:     "UAResponse",
		responseJson: uaResponseJson,
		responseType: response.UAResponse,
	},
	{
		caseName:     "XAResponse",
		responseJson: xaResponseJson,
		responseType: response.XAResponse,
	},
	{
		caseName:     "ResourceResponse",
		responseJson: resourceJson,
		responseType: response.ResourceResponse,
	},
}

var (
	uaResponseJson = `
{
  "errorCode": 0,
  "message": "unununknown",
  "data": {
    "errorCode": 0,
    "message": "unknown",
    "credential": "xxx",
    "universalAccount": {
      "username": "hello",
      "password": "world",
      "chainAccounts": [
        {
          "keyID": 1,
          "type": "BCOS2.0",
          "isDefault": true,
          "pubKey": "XXX",
          "secKey": "XXX",
          "ext": "address"
        },
        {
          "keyID": 2,
          "type": "Fabric1.4",
          "isDefault": true,
          "pubKey": "xxx",
          "secKey": "xxx",
          "ext": "membershipID"
        },
        {
          "keyID": 3,
          "type": "Fabric2.0",
          "isDefault": true,
          "pubKey": "xxx",
          "secKey": "xxx",
          "ext": "membershipID"
        }
      ]
    }
  }
}`
	xaResponseJson = `{
  "errorCode": 0,
  "data": {
    "status": 0,
    "chainErrorMessages": [
      {
        "path": "test/test1",
        "message": "test fail"
      },
      {
        "path": "test/test2",
        "message": "test fail"
      }
    ]
  },
  "xarawResponse": {
    "status": 0,
    "chainErrorMessages": [
      {
        "path": "test/test1",
        "message": "test fail"
      },
      {
        "path": "test/test2",
        "message": "test fail"
      }
    ]
  }
}`
	resourceJson = `
{
  "errorCode": 0,
  "data": {
    "resourceDetails": [
      {
        "path": "test.test.test",
        "distance": 2,
        "stubType": "testType",
        "properties": {
          "property2": 0,
          "property1": "unknown"
        }
      }
    ]
  },
  "resources": {
    "resourceDetails": [
      {
        "path": "test.test.test",
        "distance": 2,
        "stubType": "testType",
        "properties": {
          "property2": 0,
          "property1": "unknown"
        }
      }
    ]
  }
}`
)

type mockIOReadCloser struct {
	reader *bytes.Reader
}

func newMockMockIOReadCloser(input string) *mockIOReadCloser {
	buffer := bytes.NewReader([]byte(input))

	return &mockIOReadCloser{reader: buffer}
}

func (m mockIOReadCloser) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m mockIOReadCloser) Close() error {
	return nil
}

func generateHttpResponse(body string) *http.Response {
	httpResponse := &http.Response{
		Body: newMockMockIOReadCloser(body),
	}
	return httpResponse
}

func TestParseUAReceipt(t *testing.T) {
	rsp := ParseResponse(generateHttpResponse(testCases[0].responseJson), testCases[0].responseType)
	uaReceipt, _ := rsp.Data.(*response.UAReceipt)
	t.Logf("got Universal Account Info:\n%s", uaReceipt.UniversalAccount.ToString())
}

func TestParseResponse(t *testing.T) {
	for i, tc := range testCases {
		response := ParseResponse(generateHttpResponse(tc.responseJson), tc.responseType)
		t.Logf("case %d, name %s result:\n%s", i, tc.caseName, response.ToString())
	}
}
