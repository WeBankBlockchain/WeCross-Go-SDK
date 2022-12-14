package account

import (
	"reflect"
	"testing"
)

func TestUniversalAccountParse(t *testing.T) {
	bcosAccount := NewBCOSAccount(1, FISCO_BOCS_2, true, "XXX", "XXX", "address")
	bcosGMAccount := NewBCOSAccount(2, GM_BCOS_2, true, "XXX", "XXX", "address")
	fabric1Account := NewFabricAccount(3, FABRIC_1, true, "xxx", "xxx", "membershipID")
	fabric2Account := NewFabricAccount(4, FABRIC_2, true, "xxx", "xxx", "membershipID")
	universalAccount := NewUniversalAccount("hello", "world", "", "", "", []ChainAccount{bcosAccount, bcosGMAccount, fabric1Account, fabric2Account})
	t.Logf("String:\n%s", universalAccount.ToString())
	jsonBytes := universalAccount.toJson()
	t.Logf("Json:\n%s", jsonBytes)
	parsedUniversalAccount := ParseUniversalAccountFromJson(jsonBytes)
	if !reflect.DeepEqual(universalAccount, parsedUniversalAccount) {
		t.Fatalf("fail in toJson and parsing")
	}
}
