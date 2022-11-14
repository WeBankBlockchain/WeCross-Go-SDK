package connections

import (
	"WeCross-Go-SDK/common"
	"WeCross-Go-SDK/utils"
	"crypto/tls"
	"crypto/x509"
	"os"
	"testing"
)

var TEST_APPLICATION_CONFIG_DIR = "../../../fortests/tomldir"

func TestConnection_LoadToml(t *testing.T) {
	utils.SetClassPath(TEST_APPLICATION_CONFIG_DIR)
	config, err := utils.GetToml(utils.ReadClassPath(common.APPLICATION_CONFIG_FILE))
	if err != nil {
		t.Fatalf("fail in get toml: %v", err)
	}

	connection, err := NewConnection(config)
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
