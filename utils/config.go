package utils

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/common"
	"github.com/pelletier/go-toml"
	"os"
	"strings"
)

func GetToml(fileName string) (*toml.Tree, *common.WeCrossSDKError) {
	config, err := toml.LoadFile(fileName)
	if err != nil {
		return nil, common.NewWeCrossSDKFromError(common.INTERNAL_ERROR, err)
	}
	return config, nil
}

// ReadClassPath will infer whther the inPath starts with a
func ReadClassPath(inPath string, classpath string) (rawPath string) {
	if strings.HasPrefix(inPath, "classpath:") {
		rawPath = strings.Replace(inPath, "classpath:", classpath+string(os.PathSeparator), 1)
	} else {
		rawPath = inPath
	}
	return rawPath
}
