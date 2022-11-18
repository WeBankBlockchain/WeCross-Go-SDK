package utils

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/common"
	"regexp"
)

func FormatUrlPrefix(urlPrefix string) (string, *common.WeCrossSDKError) {
	pattern := "^/[A-Za-z0-9][\\w-]{0,17}$"
	pcre_pattern := "^/(?!_)(?!-)(?!.*?_$)(?!.*?-$)[\\w-]{1,18}$" // go package regexp not support ?!
	prefix := urlPrefix
	if len(prefix) == 0 {
		return "", nil
	}

	// something => /something
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}

	// /something/ => /something
	if prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}

	// ban /something- and /something_
	if len(prefix) == 0 {
		return "", common.NewWeCrossSDKFromString(common.FIELD_MISSING, "Error 'urlPrefix' in config, it should matches pattern: "+pcre_pattern)
	}
	if prefix[len(prefix)-1] == '-' || prefix[len(prefix)-1] == '_' {
		return "", common.NewWeCrossSDKFromString(common.FIELD_MISSING, "Error 'urlPrefix' in config, it should matches pattern: "+pcre_pattern)
	}

	// /something
	ok, err := regexp.Match(pattern, []byte(prefix))
	if err != nil {
		return "", common.NewWeCrossSDKFromError(common.FIELD_MISSING, err)
	}

	if !ok {
		return "", common.NewWeCrossSDKFromString(common.FIELD_MISSING, "Error 'urlPrefix' in config, it should matches pattern: "+pcre_pattern)
	} else {
		return prefix, nil
	}

}
