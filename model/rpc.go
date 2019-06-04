package model

import "fmt"

const (
	RPC_STATUS_OK = 200
	RPC_STATUS_BAD_REQUEST = 400
	RCP_STATUS_FAIL = 500
)

// transport rpc params
type RpcArgs struct {
	Metrics []MetricValue
}

type RpcReply struct {
	Message string
	Total   int
	Invalid int
	Latency int64
}

// code == 200 => success
// code == 400 => bad request
type SimpleRpcResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (this *SimpleRpcResponse) String() string {
	return fmt.Sprintf("<Code: %d>", this.Code)
}

type NullRpcRequest struct {
}

type InfluxRpcArgs struct {
	DB string
	Username string
	Password string
	Body [][]byte
}


type InfluxRpcResponse = SimpleRpcResponse
