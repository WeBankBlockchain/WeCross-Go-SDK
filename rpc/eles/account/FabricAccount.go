package account

import (
	"fmt"
	"strconv"
)

type FabricAccount struct {
	*CommonAccount
	PubKey string `json:"pubKey"`
	SecKey string `json:"secKey"`
	Ext    string `json:"ext"`
}

func NewFabricAccount(keyID int, accountType ChainAccountType, isDefault bool, pubKey, secKey, ext string) *FabricAccount {
	commonAccount := &CommonAccount{
		KeyID:       keyID,
		AccountType: accountType,
		IsDefault:   isDefault,
	}
	return &FabricAccount{commonAccount, pubKey, secKey, ext,}
}

func (fa *FabricAccount) ToString() string {
	str := "{" + "\"keyID\":\"" + strconv.Itoa(fa.KeyID) + "\"" + ", \"type\":\"" + string(fa.AccountType) + "\"" + ", \"cert\":\"" + fa.PubKey + "\"" + ", \"MembershipID\":\"" + fa.Ext + "\"" + ", \"isDefault\":\"" + fmt.Sprintf("%t", fa.IsDefault) + "\"}"
	return str
}

func (fa *FabricAccount) ToFormatString() string {
	str := "\t" + string(fa.AccountType) + " Account:\n" + "\tkeyID    : " + strconv.Itoa(fa.KeyID) + "\n" + "\ttype     : " + string(fa.AccountType) + "\n" + "\tMembershipID : " + fa.Ext + "\n" + "\tisDefault: " + fmt.Sprintf("%t", fa.IsDefault) + "\n\t----------\n"
	return str
}

func (fa *FabricAccount) ToDetailString() string {
	str := "\t" + string(fa.AccountType) + " Account:\n" + "\tkeyID    : " + strconv.Itoa(fa.KeyID) + "\n" + "\ttype     : " + string(fa.AccountType) + "\n" + "\tcert     : " + fa.PubKey + "\n" + "\tMembershipID : " + fa.Ext + "\n" + "\tisDefault: " + fmt.Sprintf("%t", fa.IsDefault) + "\n\t----------\n"
	return str
}
