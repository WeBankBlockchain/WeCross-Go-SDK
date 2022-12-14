package account

import (
	"fmt"
	"strconv"
)

type BCOSAccount struct {
	*CommonAccount
	PubKey string `json:"pubKey"`
	SecKey string `json:"secKey"`
	Ext    string `json:"ext"`
}

func NewBCOSAccount(keyID int, accountType ChainAccountType, isDefault bool, pubKey, secKey, ext string) *BCOSAccount {
	commonAccount := &CommonAccount{
		KeyID:       keyID,
		AccountType: accountType,
		IsDefault:   isDefault,
	}
	return &BCOSAccount{commonAccount, pubKey, secKey, ext}
}

func (ba *BCOSAccount) ToString() string {
	str := "{" + "\"keyID\":\"" + strconv.Itoa(ba.KeyID) + "\"" + ", \"type\":\"" + string(ba.AccountType) + "\"" + ", \"pubKey\":\"" + ba.PubKey + "\"" + ", \"address\":\"" + ba.Ext + "\"" + ", \"isDefault\":\"" + fmt.Sprintf("%t", ba.IsDefault) + "\"}"
	return str
}

func (ba *BCOSAccount) ToFormatString() string {
	str := "\t" + string(ba.AccountType) + " Account:\n" + "\tkeyID    : " + strconv.Itoa(ba.KeyID) + "\n" + "\ttype     : " + string(ba.AccountType) + "\n" + "\taddress  : " + ba.Ext + "\n" + "\tisDefault: " + fmt.Sprintf("%t", ba.IsDefault) + "\n\t----------\n"
	return str
}

func (ba *BCOSAccount) ToDetailString() string {
	str := "\t" + string(ba.AccountType) + " Account:\n" + "\tkeyID    : " + strconv.Itoa(ba.KeyID) + "\n" + "\ttype     : " + string(ba.AccountType) + "\n" + "\tpubKey   : " + ba.PubKey + "\n" + "\taddress  : " + ba.Ext + "\n" + "\tisDefault: " + fmt.Sprintf("%t", ba.IsDefault) + "\n\t----------\n"
	return str
}
