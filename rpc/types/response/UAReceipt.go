package response

import (
	"WeCross-Go-SDK/common"
	"WeCross-Go-SDK/rpc/eles/account"
	"WeCross-Go-SDK/utils"
	"encoding/json"
)

type UAReceipt struct {
	ErrorCode        common.ErrorCode          `json:"errorCode"`
	Message          string                    `json:"message"`
	Credential       string                    `json:"credential"`
	UniversalAccount *account.UniversalAccount `json:"universalAccount"`
}

func (ur *UAReceipt) ToString() string {
	str := "UAReceipt{" + "errorCode=" + ur.ErrorCode.ToString() + ", errorMessage='" + ur.Message + "'"
	if ur.Credential != "" {
		str += ", credential = '" + ur.Credential + "'"
	}
	str += "}"
	return str
}

func (ur *UAReceipt) ParseSelfFromJson(valueBytes []byte) {
	m := make(map[string]interface{})
	err := json.Unmarshal(valueBytes, &m)
	if err != nil {
		return
	}
	ur.ErrorCode = utils.AnyToErrorCode(m["errorCode"])
	ur.Message = utils.AnyToString(m["message"])
	ur.Credential = utils.AnyToString(m["credential"])
	if obj, ok := m["universalAccount"]; ok {
		uaBytes, err := json.Marshal(obj)
		if err != nil {
			ur.UniversalAccount = nil
			return
		}
		ur.UniversalAccount = account.ParseUniversalAccountFromJson(uaBytes)
	} else {
		ur.UniversalAccount = nil
	}
}
