package utils

import "testing"

var testUrls = []struct {
	urlPrefix string
	expect    string
	valid     bool
}{
	{
		urlPrefix: "",
		expect:    "",
		valid:     true,
	},
	{
		urlPrefix: "/something/",
		expect:    "/something",
		valid:     true,
	},
	{
		urlPrefix: "something/",
		expect:    "/something",
		valid:     true,
	},
	{
		urlPrefix: "/some!thing",
		expect:    "",
		valid:     false,
	},
	{
		urlPrefix: "/somethingsomethingsomethingsomethingsomethingsomething",
		expect:    "",
		valid:     false,
	},
	{
		urlPrefix: "/some-thing",
		expect:    "/some-thing",
		valid:     true,
	},
	{
		urlPrefix: "/-something",
		expect:    "",
		valid:     false,
	},
	{
		urlPrefix: "/_something",
		expect:    "",
		valid:     false,
	},
	{
		urlPrefix: "/something_",
		expect:    "",
		valid:     false,
	},
	{
		urlPrefix: "/something-",
		expect:    "",
		valid:     false,
	},
	{
		urlPrefix: "/",
		expect:    "",
		valid:     false,
	},
	{
		urlPrefix: "/s",
		expect:    "/s",
		valid:     true,
	},
}

func TestFormatUrlPrefix(t *testing.T) {
	for index, tu := range testUrls {
		prefix, err := FormatUrlPrefix(tu.urlPrefix)
		if err != nil {
			if tu.valid {
				t.Fatalf("fail in test %d with url: %s, %v", index, tu.urlPrefix, err)
			} else {
				continue
			}
		}
		if prefix != tu.expect {
			t.Fatalf("fail in test %d with url: %s, expect: %s, got: %s", index, tu.urlPrefix, tu.expect, prefix)
		}
	}

}
