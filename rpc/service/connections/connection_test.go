package connections

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/common"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/utils"
	"os"
	"testing"
)

var TEST_APPLICATION_CONFIG_DIR = "../../../fortests/tomldir"

func TestConnection_LoadToml(t *testing.T) {
	config, err := utils.GetToml(utils.ReadClassPath(common.APPLICATION_CONFIG_FILE, TEST_APPLICATION_CONFIG_DIR))
	if err != nil {
		t.Fatalf("fail in get toml: %v", err)
	}

	connection, err := NewConnection(config, TEST_APPLICATION_CONFIG_DIR)
	if err != nil {
		t.Fatalf("fail in new connection: %v", err)
	}

	t.Logf("got connection:\n %s", connection.ToString())

	caCert, errr := os.ReadFile(connection.GetCaCert())
	if errr != nil {
		t.Fatalf("fail in read ca cert: %v", errr)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)
	_, errr = tls.LoadX509KeyPair(connection.GetSslCert(), connection.GetSslKey())
	if err != nil {
		t.Fatalf("fail in load x509 key pair: %v", errr)
	}
}
