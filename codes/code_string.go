package codes

import "strconv"

func (c Code) String() string {
	switch c {
	case Success:
		return "Success"
	case InternalError:
		return "InternalError"
	case FieldMissing:
		return "FieldMissing"
	case ResourceError:
		return "ResourceError"
	case IllegalSymbol:
		return "IllegalSymbol"
	case RemoteCallError:
		return "RemoteCallError"
	case RpcError:
		return "RpcError"
	case CallContractError:
		return "CallContractError"
	case LackAuthentication:
		return "LackAuthentication"
	case ResourceInactive:
		return "ResourceInactive"
	case InvalidContract:
		return "InvalidContract"
	default:
		return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
	}
}
