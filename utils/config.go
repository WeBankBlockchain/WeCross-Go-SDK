package utils

import (
	"WeCross-Go-SDK/common"
	"github.com/pelletier/go-toml"
	"os"
	"strings"
)

var classpath, _ = BackPwd()

func GetToml(fileName string) (*toml.Tree, *common.WeCrossSDKError) {
	config, err := toml.LoadFile(fileName)
	if err != nil {
		return nil, common.NewWeCrossSDKFromError(common.INTERNAL_ERROR, err)
	}
	return config, nil
}

// SetClassPath sets the path of classpath that helps to read toml files
func SetClassPath(inpath string) {
	classpath = inpath
}

// ReadClassPath will infer whther the inPath starts with a
func ReadClassPath(inPath string) (rawPath string) {
	if strings.HasPrefix(inPath, "classpath:") {
		rawPath = strings.Replace(inPath, "classpath:", classpath+string(os.PathSeparator), 1)
	} else {
		rawPath = inPath
	}
	return rawPath
}
