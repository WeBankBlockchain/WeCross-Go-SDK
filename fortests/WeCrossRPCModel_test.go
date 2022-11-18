package fortests

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/logger"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/resource"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/eles/account"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/types/response"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	username = "org1-admin"
	password = "123456"
)

var candidateBcos2Account = &account.BCOSAccount{
	CommonAccount: &account.CommonAccount{
		KeyID:       0,
		AccountType: "BCOS2.0",
		IsDefault:   false,
	},
	PubKey: "-----BEGIN PUBLIC KEY-----\nMFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAE2Fr/usNx2w2UduCIdzeKOvU5nY/mQmZg\nkLXMJsLAVuYopXi7bfwZw5kKvJOyZM5I0hZZtrKbM8JVTp/Q0wZPug==\n-----END PUBLIC KEY-----\n",
	SecKey: "-----BEGIN PRIVATE KEY-----\nMIGEAgEAMBAGByqGSM49AgEGBSuBBAAKBG0wawIBAQQgAN6eGVBts+cAytEKmSpL\n4Fg7GNZtlx6XBj1+odqLNwmhRANCAATYWv+6w3HbDZR24Ih3N4o69Tmdj+ZCZmCQ\ntcwmwsBW5iileLtt/BnDmQq8k7JkzkjSFlm2spszwlVOn9DTBk+6\n-----END PRIVATE KEY-----\n",
	Ext:    "0x1c8a095cdb3296a05a29819c754acff66e88c5fb",
}

// TestBasicRPCInterface could be implemented only when the 'cross-all' scenario WeCross demo is on.
// More about it: https://wecross.readthedocs.io/zh_CN/latest/docs/tutorial/demo/demo_cross_all.html
func TestBasicRPCInterface(t *testing.T) {
	rand.Seed(time.Now().Unix())

	rpcService := service.NewWeCrossRPCService()
	rpcService.SetClassPath("./tomldir")

	err := rpcService.Init()
	if err != nil {
		t.Fatalf("fail in rpc service init, %v", err)
	}

	weCrossRPCModel := rpc.NewWeCrossRPCModel(rpcService)

	testRpcCase := func(method string, call *rpc.RemoteCall) {
		rsp, err := call.Send()
		if err != nil {
			t.Fatalf("fail in rpc method: %s, %v", method, err)
		}
		t.Logf("rpc method: %s, response: %s", method, rsp.ToString())
	}

	testUAChangeCase := func(method string, call *rpc.RemoteCall) {
		rsp, err := call.Send()
		if err != nil {
			t.Fatalf("fail in rpc method: %s, %v", method, err)
		}
		uaReceipt, ok := rsp.Data.(*response.UAReceipt)
		if !ok {
			t.Fatalf("fail in rpc method: %s, response data cannot be converted into UAReciept", method)
		}
		if uaReceipt.UniversalAccount == nil {
			t.Logf("rpc method: %s, response: %s", method, rsp.ToString())
		} else {
			// NOTE: when method is not 'ListAccount', the universal Account here in the response not contains the chain accounts infos, don't worry
			t.Logf("rpc method: %s, UniversalAccount info: %s", method, uaReceipt.UniversalAccount.ToFormatString())
		}

	}

	testListAccountsCase := func(method string, call *rpc.RemoteCall) {
		rsp, err := call.Send()
		if err != nil {
			t.Fatalf("fail in rpc method: %s, %v", method, err)
		}
		universalAccount, ok := rsp.Data.(*account.UniversalAccount)
		if !ok {
			t.Fatalf("fail in rpc method: %s, response data cannot be converted into Universal Account", method)
		}
		t.Logf("rpc method: %s, UniversalAccount info: %s", method, universalAccount.ToFormatString())
	}

	testRpcCase("Test", weCrossRPCModel.Test())

	testRpcCase("QueryPub", weCrossRPCModel.QueryPub())

	testRpcCase("QueryAuthCode", weCrossRPCModel.QueryAuthCode())

	call, err := weCrossRPCModel.Login(username, password)
	if err != nil {
		t.Fatalf("fail in login, %v", err)
	}
	testUAChangeCase("Login", call)

	testRpcCase("SupportedStubs", weCrossRPCModel.SupportedStubs())

	testRpcCase("ListResources", weCrossRPCModel.ListResources(false))

	testListAccountsCase("ListAccount", weCrossRPCModel.ListAccount())

	testRpcCase("SendTransaction", weCrossRPCModel.SendTransaction("payment.bcos-group1.HelloWorldGroup1", "set", time.Now().String()))

	testRpcCase("Call", weCrossRPCModel.Call("payment.bcos-group1.HelloWorldGroup1", "get"))

	testUAChangeCase("SetDefaultAccount", weCrossRPCModel.SetDefaultAccount("Fabric1.4", nil, 1))
	testListAccountsCase("ListAccount", weCrossRPCModel.ListAccount())

	testUAChangeCase("SetDefaultAccount", weCrossRPCModel.SetDefaultAccount("Fabric1.4", nil, 2))
	testListAccountsCase("ListAccount", weCrossRPCModel.ListAccount())

	testRpcCase("Invoke", weCrossRPCModel.Invoke("payment.fabric-mychannel.sacc", "set", "Leo", "awesome"))

	testRpcCase("Call", weCrossRPCModel.Call("payment.fabric-mychannel.sacc", "get", "Leo"))

	testUAChangeCase("Logout", weCrossRPCModel.Logout())

	newUserName := "user" + strconv.Itoa(rand.Int())[0:8]
	newPassWord := "654321"
	call, err = weCrossRPCModel.Register(newUserName, newPassWord)
	if err != nil {
		t.Fatalf("fail in register, %v", err)
	}
	testUAChangeCase("Register", call)

	call, err = weCrossRPCModel.Login(newUserName, newPassWord)
	if err != nil {
		t.Fatalf("fail in login, %v", err)
	}
	testUAChangeCase("Login", call)

	// NOTE: If you've run once this test, the added Account could be not added to the new universal account anymore.
	// If this happens, please revise the candidateBcos2Account or re-build wecross demo.
	testUAChangeCase("AddChainAccount", weCrossRPCModel.AddChainAccount("Bcos2.0", candidateBcos2Account))

	testListAccountsCase("ListAccount", weCrossRPCModel.ListAccount())

}

// TestBasicRPCInterface could be implemented only when the 'cross-all' scenario WeCross demo is on.
// More about it: https://wecross.readthedocs.io/zh_CN/latest/docs/tutorial/demo/demo_cross_all.html
func TestResourceInterface(t *testing.T) {
	rpcService := service.NewWeCrossRPCService()
	rpcService.SetClassPath("./tomldir")

	err := rpcService.Init()
	if err != nil {
		t.Fatalf("fail in rpc service init, %v", err)
	}

	weCrossRPCModel := rpc.NewWeCrossRPCModel(rpcService)

	testRpcCase := func(method string, call *rpc.RemoteCall) {
		rsp, err := call.Send()
		if err != nil {
			t.Fatalf("fail in rpc method: %s, %v", method, err)
		}
		t.Logf("rpc method: %s, response: %s", method, rsp.ToString())
	}

	testUAChangeCase := func(method string, call *rpc.RemoteCall) {
		rsp, err := call.Send()
		if err != nil {
			t.Fatalf("fail in rpc method: %s, %v", method, err)
		}
		uaReceipt, ok := rsp.Data.(*response.UAReceipt)
		if !ok {
			t.Fatalf("fail in rpc method: %s, response data cannot be converted into UAReciept", method)
		}
		if uaReceipt.UniversalAccount == nil {
			t.Logf("rpc method: %s, response: %s", method, rsp.ToString())
		} else {
			// NOTE: when method is not 'ListAccount', the universal Account here in the response not contains the chain accounts infos, don't worry
			t.Logf("rpc method: %s, UniversalAccount info: %s", method, uaReceipt.UniversalAccount.ToFormatString())
		}

	}

	// should first log in
	call, err := weCrossRPCModel.Login(username, password)
	if err != nil {
		t.Fatalf("fail in login, %v", err)
	}
	testUAChangeCase("Login", call)

	testRpcCase("ListResources", weCrossRPCModel.ListResources(false))

	// Here starts to test the resource API
	// ------------ Test FISCO BCOS Resource ------------------
	testResourceBcos := resource.NewResource(weCrossRPCModel, "payment.bcos-group1.HelloWorldGroup1")

	if !testResourceBcos.IsActive() {
		t.Fatalf("fail in resource IsActive")
	}

	detail, err := testResourceBcos.Detail()
	if err != nil {
		t.Fatalf("fail in resource Detail, %v", err)
	}
	t.Logf("The resource detail is:\n%s", detail.ToString())

	sendTransactionResult, err := testResourceBcos.SendTransaction("set", time.Now().String())
	if err != nil {
		t.Fatalf("fail in resource SendTransaction, %v", err)
	}
	t.Logf("The result of SendTransaction is:\n%v", sendTransactionResult)

	callResult, err := testResourceBcos.Call("get")
	if err != nil {
		t.Fatalf("fail in resource Call, %v", err)
	}
	t.Logf("The result of Call is:\n%v", callResult)

	// ------------ Test Hyperledger Fabric Resource ------------------
	testResourceFabric := resource.NewResource(weCrossRPCModel, "payment.fabric-mychannel.sacc")
	if !testResourceFabric.IsActive() {
		t.Fatalf("fail in resource IsActive")
	}

	detail, err = testResourceFabric.Detail()
	if err != nil {
		t.Fatalf("fail in resource Detail, %v", err)
	}
	t.Logf("The resource detail is:\n%s", detail.ToString())

	sendTransactionResult, err = testResourceFabric.SendTransaction("set", "Leo", "Awesome")
	if err != nil {
		t.Fatalf("fail in resource SendTransaction, %v", err)
	}
	t.Logf("The result of SendTransaction is:\n%v", sendTransactionResult)

	callResult, err = testResourceFabric.Call("get", "Leo")
	if err != nil {
		t.Fatalf("fail in resource Call, %v", err)
	}
	t.Logf("The result of Call is:\n%v", callResult)

}

// TestXATransactionRPCInterface could be implemented only when the 'xa-evidence' scenario WeCross demo is on.
// More about it: https://wecross.readthedocs.io/zh_CN/latest/docs/tutorial/demo/demo.html#id4
func TestXATransactionRPCInterface(t *testing.T) {
	stdoutLog := logger.NewStdLogSystem(os.Stdout, 0, logger.InfoLevel)
	logger.AddLogSystem(stdoutLog)

	rand.Seed(time.Now().Unix())

	rpcService := service.NewWeCrossRPCService()
	rpcService.SetClassPath("./tomldir")

	err := rpcService.Init()
	if err != nil {
		t.Fatalf("fail in rpc service init, %v", err)
	}

	weCrossRPCModel := rpc.NewWeCrossRPCModel(rpcService)

	testRpcCase := func(method string, call *rpc.RemoteCall) {
		rsp, err := call.Send()
		if err != nil {
			t.Fatalf("fail in rpc method: %s, %v", method, err)
		}
		t.Logf("rpc method: %s, response: %s", method, rsp.ToString())
	}

	testUAChangeCase := func(method string, call *rpc.RemoteCall) {
		rsp, err := call.Send()
		if err != nil {
			t.Fatalf("fail in rpc method: %s, %v", method, err)
		}
		uaReceipt, ok := rsp.Data.(*response.UAReceipt)
		if !ok {
			t.Fatalf("fail in rpc method: %s, response data cannot be converted into UAReciept", method)
		}
		if uaReceipt.UniversalAccount == nil {
			t.Logf("rpc method: %s, response: %s", method, rsp.ToString())
		} else {
			// NOTE: when method is not 'ListAccount', the universal Account here in the response not contains the chain accounts infos, don't worry
			t.Logf("rpc method: %s, UniversalAccount info: %s", method, uaReceipt.UniversalAccount.ToFormatString())
		}
	}

	// ---------------- test1: start->exec->commit ----------------------
	call, err := weCrossRPCModel.Login(username, password)
	if err != nil {
		t.Fatalf("fail in login, %v", err)
	}
	testUAChangeCase("Login", call)

	transactionID := strconv.Itoa(rand.Int())

	t.Logf("test1 transaction ID: %s", transactionID)

	testRpcCase("StartTransaction", weCrossRPCModel.StartXATransaction(transactionID, []string{"payment.bcos.evidence", "payment.fabric.evidence"}))

	// query the status of the XAtransaction
	testRpcCase("GetXATransaction", weCrossRPCModel.GetXATransaction(transactionID, []string{"payment.bcos.evidence", "payment.fabric.evidence"}))

	evidence1 := "evidence" + strconv.Itoa(rand.Int())[0:5]

	// NOTE: if use CallXA, should input the transactionID
	testRpcCase("Call", weCrossRPCModel.Call("payment.bcos.evidence", "queryEvidence", evidence1))
	testRpcCase("Call", weCrossRPCModel.Call("payment.fabric.evidence", "queryEvidence", evidence1))

	// NOTE: if use SendXATransaction, should input the transactionID
	testRpcCase("Invoke", weCrossRPCModel.Invoke("payment.bcos.evidence", "newEvidence", evidence1, "Leo is awesome."))
	testRpcCase("Invoke", weCrossRPCModel.Invoke("payment.fabric.evidence", "newEvidence", evidence1, "Leo is handsome."))

	testRpcCase("CommitTransaction", weCrossRPCModel.CommitXATransaction(transactionID, []string{"payment.bcos.evidence", "payment.fabric.evidence"}))

	testRpcCase("Call", weCrossRPCModel.Call("payment.bcos.evidence", "queryEvidence", evidence1))
	testRpcCase("Call", weCrossRPCModel.Call("payment.fabric.evidence", "queryEvidence", evidence1))

	// ---------------- take a break, and look through the historical XAtransactions -----------

	testRpcCase("ListXATransactions", weCrossRPCModel.ListXATransactions(10))

	// ---------------- test2: rollback ----------------------
	call, err = weCrossRPCModel.Login(username, password)
	if err != nil {
		t.Fatalf("fail in login, %v", err)
	}
	testUAChangeCase("Login", call)

	evidence2 := "evidence" + strconv.Itoa(rand.Int())[0:5]

	testRpcCase("Call", weCrossRPCModel.Call("payment.bcos.evidence", "queryEvidence", evidence2))
	testRpcCase("Call", weCrossRPCModel.Call("payment.fabric.evidence", "queryEvidence", evidence2))

	transactionID = strconv.Itoa(rand.Int())

	t.Logf("test2 transaction ID: %s", transactionID)

	testRpcCase("StartTransaction", weCrossRPCModel.StartXATransaction(transactionID, []string{"payment.bcos.evidence", "payment.fabric.evidence"}))

	testRpcCase("Invoke", weCrossRPCModel.Invoke("payment.bcos.evidence", "newEvidence", evidence2, "Leo is bad."))
	testRpcCase("Invoke", weCrossRPCModel.Invoke("payment.fabric.evidence", "newEvidence", evidence2, "Leo is sad."))

	testRpcCase("Call", weCrossRPCModel.Call("payment.bcos.evidence", "queryEvidence", evidence2))
	testRpcCase("Call", weCrossRPCModel.Call("payment.fabric.evidence", "queryEvidence", evidence2))

	// NOTE: this two normal transaction send should not be allowed as the XAtransaction is locked!
	testRpcCase("SendTransaction", weCrossRPCModel.SendTransaction("payment.bcos.evidence", "newEvidence", evidence2, "Krad is bad."))
	testRpcCase("SendTransaction", weCrossRPCModel.SendTransaction("payment.fabric.evidence", "newEvidence", evidence2, "Krad is sad."))

	testRpcCase("Call", weCrossRPCModel.Call("payment.bcos.evidence", "queryEvidence", evidence2))
	testRpcCase("Call", weCrossRPCModel.Call("payment.fabric.evidence", "queryEvidence", evidence2))

	// Now, roll back!
	testRpcCase("RollbackXATransaction", weCrossRPCModel.RollbackXATransaction(transactionID, []string{"payment.bcos.evidence", "payment.fabric.evidence"}))

	testRpcCase("Call", weCrossRPCModel.Call("payment.bcos.evidence", "queryEvidence", evidence2))
	testRpcCase("Call", weCrossRPCModel.Call("payment.fabric.evidence", "queryEvidence", evidence2))

}
