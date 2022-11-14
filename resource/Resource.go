package resource

import (
	"WeCross-Go-SDK/common"
	"WeCross-Go-SDK/logger"
	"WeCross-Go-SDK/rpc"
	"WeCross-Go-SDK/rpc/eles/resources"
	"WeCross-Go-SDK/rpc/types"
	"WeCross-Go-SDK/rpc/types/response"
	"WeCross-Go-SDK/utils"
)

const (
	ResourceTag = "Resource"
)

type Resource struct {
	logger     *logger.Logger
	weCrossRPC rpc.WeCrossRPC
	path       string
}

func NewResource(weCrossRPC rpc.WeCrossRPC, path string) *Resource {
	logger := logger.NewLogger(ResourceTag)
	return &Resource{
		logger:     logger,
		weCrossRPC: weCrossRPC,
		path:       path,
	}
}

func (rsc *Resource) checkIPath() *common.WeCrossSDKError {
	if utils.CheckResourcePath(rsc.path) {
		return nil
	} else {
		return common.NewWeCrossSDKFromString(common.RESOURCE_ERROR, "Invalid iPath: "+rsc.path)
	}
}

func (rsc *Resource) checkWeCrossRPC() *common.WeCrossSDKError {
	if rsc.weCrossRPC == nil {
		return common.NewWeCrossSDKFromString(common.RESOURCE_ERROR, "WeCrossRPC not set")
	}
	return nil
}

func (rsc *Resource) Check() *common.WeCrossSDKError {
	if err := rsc.checkWeCrossRPC(); err != nil {
		return err
	}
	if err := rsc.checkIPath(); err != nil {
		return err
	}
	return nil
}

func (rsc *Resource) IsActive() bool {
	rsp, err := rsc.mustOkRequest(rsc.weCrossRPC.ListResources(false))
	if err != nil {
		rsc.logger.Errorf("got error %v", err)
		return false
	}
	err = checkResponse(rsp)
	if err != nil {
		rsc.logger.Errorf("got error %v", err)
		return false
	}
	gotResources, ok := rsp.Data.(*resources.Resources)
	if !ok {
		rsc.logger.Errorf("got response that is not ResourceResponse")
		return false
	}
	resourceDetails := gotResources.ResourceDetails
	isActiveFlag := false

	for _, rcdt := range resourceDetails {
		if rcdt.Path == rsc.path {
			isActiveFlag = true
			break
		}
	}

	return isActiveFlag
}

func (rsc *Resource) mustOkRequest(call *rpc.RemoteCall) (*types.Response, *common.WeCrossSDKError) {
	rsp, err := call.Send()
	if err != nil {
		rsc.logger.Errorf("Error in RemoteCall: %s", err.Content)
	}
	return rsp, err
}

func checkResponse(rsp *types.Response) *common.WeCrossSDKError {
	if rsp == nil {
		return common.NewWeCrossSDKFromString(common.RPC_ERROR, "response is nil")
	}
	if int(rsp.ErrorCode) != int(common.SUCCESS) || rsp.Data == nil {
		return common.NewWeCrossSDKFromString(common.RPC_ERROR, rsp.ToString())
	}
	return nil
}

func (rsc *Resource) Detail() (*resources.ResourceDetail, *common.WeCrossSDKError) {
	rsp, err := rsc.mustOkRequest(rsc.weCrossRPC.Detail(rsc.path))
	if err != nil {
		return nil, err
	}
	resourceDetail, ok := rsp.Data.(*resources.ResourceDetail)
	if !ok {
		return nil, common.NewWeCrossSDKFromString(common.INTERNAL_ERROR, "cannot parse the response as resource detail info")
	}
	return resourceDetail, nil
}

func (rsc *Resource) Call(method string, args ...string) ([]string, *common.WeCrossSDKError) {
	res := make([]string, 0)
	rsp, err := rsc.mustOkRequest(rsc.weCrossRPC.Call(rsc.path, method, args...))
	if err != nil {
		return res, err
	}
	err = checkResponse(rsp)
	if err != nil {
		return res, err
	}
	receipt, ok := rsp.Data.(*response.TXReceipt)
	if !ok {
		return res, common.NewWeCrossSDKFromString(common.INTERNAL_ERROR, "Resource.Call fail, cannot parse the response as transaction receipt")
	}
	if int(receipt.ErrorCode) != int(common.SUCCESS) {
		return res, common.NewWeCrossSDKFromString(common.CALL_CONTRACT_ERROR, "Resource.Call fail, receipt:"+receipt.ToString())
	}
	res = receipt.Result
	return res, nil
}

func (rsc *Resource) SendTransaction(method string, args ...string) ([]string, *common.WeCrossSDKError) {
	res := make([]string, 0)
	rsp, err := rsc.mustOkRequest(rsc.weCrossRPC.SendTransaction(rsc.path, method, args...))
	if err != nil {
		return res, err
	}
	err = checkResponse(rsp)
	if err != nil {
		return res, err
	}
	receipt, ok := rsp.Data.(*response.TXReceipt)
	if !ok {
		return res, common.NewWeCrossSDKFromString(common.INTERNAL_ERROR, "Resource.SendTransaction fail, cannot parse the response as transaction receipt")
	}
	if int(receipt.ErrorCode) != int(common.SUCCESS) {
		return res, common.NewWeCrossSDKFromString(common.CALL_CONTRACT_ERROR, "Resource.SendTransaction fail, receipt:"+receipt.ToString())
	}
	res = receipt.Result
	return res, nil
}

func (rsc *Resource) GetPath() string {
	return rsc.path
}

func (rsc *Resource) GetWeCrossRPC() rpc.WeCrossRPC {
	return rsc.weCrossRPC
}
