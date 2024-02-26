package vrabbitmqhelper

//"dataType":"data",
//"deviceClass":"CMS",
//"deviceType":"1",
//"deviceCode":"2",
//"data":{
var ShiJianMuBanHM = "2006-01-02 15:04:05.566"
/*
消息体的主体
 */
type MainMsg struct {
	DataType string `json:"dataType"`
	DeviceClass string `json:"deviceClass"`
	DeviceType string `json:"deviceType"`
	DeviceCode string `json:"deviceCode"`
	Data interface{} `json:"data"`
}

/*
消息体的内容
 */
type LogMsg struct {
	LeiXing string `json:"leiXing"`
	ShiJian string `json:"shiJian"`
	NeiRong string `json:"neiRong"`
}