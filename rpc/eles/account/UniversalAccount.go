package account

import (
	"WeCross-Go-SDK/utils"
	"encoding/json"
	"strings"
)

type UniversalAccount struct {
	username      string
	password      string
	pubKey        string
	secKey        string
	uaID          string
	chainAccounts []ChainAccount
}

func (ua *UniversalAccount) ParseSelfFromJson(valueBytes []byte) {
	tempUa := ParseUniversalAccountFromJson(valueBytes)
	ua.username = tempUa.username
	ua.password = tempUa.password
	ua.pubKey = tempUa.pubKey
	ua.secKey = tempUa.secKey
	ua.uaID = tempUa.uaID
	ua.chainAccounts = tempUa.chainAccounts
}

// PlainUniversalAccount is used to prevent users from accessing to the private information (password and secKey) directly.
type PlainUniversalAccount struct {
	UserName      string         `json:"username"`
	PassWord      string         `json:"password"`
	PubKey        string         `json:"pubKey"`
	SecKey        string         `json:"secKey"`
	UaID          string         `json:"uaID"`
	ChainAccounts []ChainAccount `json:"chainAccounts"`
}

func NewUniversalAccount(username, password, pubKey, secKey, uaID string, accountsList []ChainAccount) *UniversalAccount {
	return &UniversalAccount{
		username:      username,
		password:      password,
		pubKey:        pubKey,
		secKey:        secKey,
		uaID:          uaID,
		chainAccounts: accountsList,
	}
}

// ToJson is just for testing
func (ua *UniversalAccount) toJson() []byte {
	plainUniversalAccount := &PlainUniversalAccount{
		UserName:      ua.username,
		PassWord:      ua.password,
		PubKey:        ua.pubKey,
		SecKey:        ua.secKey,
		UaID:          ua.uaID,
		ChainAccounts: ua.chainAccounts,
	}
	uaBytes, _ := json.Marshal(plainUniversalAccount)
	return uaBytes
}

func (ua *UniversalAccount) ToString() string {
	str := "{" + "\"username\":\"" + ua.username + "\"" + ", \"pubKey\":\"" + ua.pubKey + "\"" + ", \"uaID\":\"" + ua.uaID + "\""
	str += ", \"chainAccounts\": ["
	if ua.chainAccounts != nil {
		for i := 0; i < len(ua.chainAccounts); i++ {
			str += ua.chainAccounts[i].ToString() + ","
		}
		str = strings.TrimSuffix(str, ",")
	}
	str += "]}"
	return str
}

func ParseUniversalAccountFromJson(data []byte) *UniversalAccount {
	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}
	chainAccounts := make([]ChainAccount, 0)
	chainObjs, ok := m["chainAccounts"].([]interface{})
	if ok {
		for _, obj := range chainObjs {
			chainAccount := MappingChainAccount(obj)
			if chainAccount == nil {
				continue
			}
			chainAccounts = append(chainAccounts, chainAccount)
		}
	}
	return &UniversalAccount{
		username:      utils.AnyToString(m["username"]),
		password:      utils.AnyToString(m["password"]),
		pubKey:        utils.AnyToString(m["pubKey"]),
		secKey:        utils.AnyToString(m["secKey"]),
		uaID:          utils.AnyToString(m["uaID"]),
		chainAccounts: chainAccounts,
	}
}

func MappingChainAccount(obj interface{}) ChainAccount {
	objs, ok := obj.(map[string]interface{})
	if !ok {
		return nil
	}
	accountType, ok := objs["type"].(string)
	if !ok {
		return nil
	}
	valueBytes, err := json.Marshal(objs)
	if err != nil {
		return nil
	}
	switch ChainAccountType(accountType) {
	case FISCO_BOCS_2, GM_BCOS_2:
		bcosAccount := new(BCOSAccount)
		err = json.Unmarshal(valueBytes, bcosAccount)
		if err != nil {
			return nil
		}
		return bcosAccount
	case FABRIC_1, FABRIC_2:
		fabricAccount := new(FabricAccount)
		err = json.Unmarshal(valueBytes, fabricAccount)
		if err != nil {
			return nil
		}
		return fabricAccount
	default:
		return nil
	}
}
