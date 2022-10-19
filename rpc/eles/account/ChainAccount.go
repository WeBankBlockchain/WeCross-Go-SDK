package account

import "fmt"

type ChainAccountType string

const (
	FISCO_BOCS_2 ChainAccountType = "BCOS2.0"
	GM_BCOS_2    ChainAccountType = "GM_BCOS2.0"
	FABRIC_1     ChainAccountType = "Fabric1.4"
	FABRIC_2     ChainAccountType = "Fabric2.0"
)

type ChainAccount interface {
	GetAccountType() ChainAccountType
	ToString() string
}

type CommonAccount struct {
	KeyID       int              `json:"keyID"`
	AccountType ChainAccountType `json:"type"`
	IsDefault   bool             `json:"isDefault"`
}

func (ca *CommonAccount) ToString() string {
	str := fmt.Sprintf("CommonAccount{keyID=%d, type='%s', isDefault=%t}", ca.KeyID, ca.AccountType, ca.IsDefault)
	return str
}

func (ca *CommonAccount) GetAccountType() ChainAccountType {
	return ca.AccountType
}
