package utils

import (
	"WeCross-Go-SDK/common"
	"github.com/pelletier/go-toml"
)

func GetToml(fileName string) (*toml.Tree, *common.WeCrossSDKError) {
	config, err := toml.LoadFile(fileName)
	if err != nil {
		return nil, common.NewWeCrossSDKFromError(common.INTERNAL_ERROR, err)
	}
	return config, nil
}
