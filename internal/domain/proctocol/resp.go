package proctocol

// succ code msg 是公用的部分 这里直接统一使用 避免重复结构代码

type RespGeneral struct {
	Success   bool        `json:"success"`
	ErrorCode int32       `json:"errorCode"`
	ErrorMsg  string      `json:"errorMsg"`
	Data      interface{} `json:"data"`
}

type Test struct {
	Id int32 `json:"id"`
}

func (res *RespGeneral) SetGeneral(Success bool, ErrCode int32, Msg string) {
	res.Success = Success
	res.ErrorCode = ErrCode
	res.ErrorMsg = Msg
}

func (res *RespGeneral) SetData(Data interface{}) {
	res.Data = Data
}
