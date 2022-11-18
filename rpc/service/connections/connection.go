package connections

import (
	"fmt"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/common"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/utils"
	"github.com/pelletier/go-toml"
)

type Connection struct {
	server, sslKey, sslCert, caCert string
	sslSwitch                       int
	urlPrefix                       string
}

func NewConnection(config *toml.Tree) (*Connection, *common.WeCrossSDKError) {
	server, ok := config.Get("connection.server").(string)
	if !ok {
		return nil, common.NewWeCrossSDKFromString(common.FIELD_MISSING, "Something wrong with parsing [connection.server], please check configuration")
	}
	caCert, ok := config.Get("connection.caCert").(string)
	if !ok {
		return nil, common.NewWeCrossSDKFromString(common.FIELD_MISSING, "Something wrong with parsing [connection.caCert], please check configuration")
	}

	sslKey, ok := config.Get("connection.sslKey").(string)
	if !ok {
		return nil, common.NewWeCrossSDKFromString(common.FIELD_MISSING, "Something wrong with parsing [connection.sslKey], please check configuration")
	}

	sslCert, ok := config.Get("connection.sslCert").(string)
	if !ok {
		return nil, common.NewWeCrossSDKFromString(common.FIELD_MISSING, "Something wrong with parsing [connection.sslCert], please check configuration")
	}

	sslSwitch, ok := config.Get("connection.sslSwitch").(int64)
	if !ok {
		sslSwitch = 0 // default 0
		//return nil, common.NewWeCrossSDKFromString(common.FIELD_MISSING, "Something wrong with parsing [connection.sslSwitch], please check configuration")
	}

	urlPrefix, ok := config.Get("connection.urlPrefix").(string)
	if !ok {
		urlPrefix = "" // could be empty
		//return nil, common.NewWeCrossSDKFromString(common.FIELD_MISSING, "Something wrong with parsing [connection.urlPrefix], please check configuration")
	}
	formatedUrlPrefix, err := utils.FormatUrlPrefix(urlPrefix)
	if err != nil {
		return nil, err
	}
	connection := &Connection{
		server:    server,
		sslKey:    sslKey,
		sslCert:   sslCert,
		caCert:    caCert,
		sslSwitch: int(sslSwitch),
		urlPrefix: formatedUrlPrefix,
	}
	return connection, nil
}

func (conn *Connection) ToString() string {
	str := fmt.Sprintf("Connection{server='%s', sslKey='%s', sslCert='%s', caCert='%s', sslSwitch=%d, urlPrefix='%s'}", conn.server, conn.sslKey, conn.sslCert, conn.caCert, conn.sslSwitch, conn.urlPrefix)
	return str
}

func (conn *Connection) GetServer() string {
	return conn.server
}
func (conn *Connection) GetSslKey() string {
	return utils.ReadClassPath(conn.sslKey)
}
func (conn *Connection) GetSslCert() string {
	return utils.ReadClassPath(conn.sslCert)
}
func (conn *Connection) GetCaCert() string {
	return utils.ReadClassPath(conn.caCert)
}
func (conn *Connection) GetSslSwitch() int {
	return conn.sslSwitch
}
func (conn *Connection) GetUrlPrefix() string {
	return conn.urlPrefix
}
