package vpo

// 消息体主体
type MainBody struct {
	//数据类型
	DataType string `json:"dataType"`
	//设备类别
	DeviceClass string `json:"deviceClass"`
	//设备类型
	DeviceType string `json:"deviceType"`
	//设备编号
	DeviceCode string `json:"deviceCode"`
	//数据
	Data interface{} `json:"data"`
}

// 消息体-设备状态
type DeviceStatus struct {
	//设备状态；0-未知；1-正常；2-故障；3-告警；4-离线；
	Status int `json:"status"`
	//状态内容
	StatusContent string `json:"statusContent"`
	//设备状态变更时间
	ChangeTime string `json:"changeTime"`
}

// 消息体-报警数据
type DeviceAlarm struct {
	//报警时间
	AlarmTime string `json:"alarmTime"`
	//报警描述
	AlarmContent string `json:"alarmContent"`
	//消警时间
	DisAlertTime string `json:"disAlertTime"`
	//消警描述
	DisAlertContent string `json:"disAlertContent"`
	//确认人
	ConfiremdBy string `json:"confiremdBy"`
	//确认时间
	ConfirmTime string `json:"confirmTime"`
	//确认描述
	ConfirmContent string `json:"confirmContent"`
	//采集值
	DeviceValue string `json:"deviceValue"`
	//越限制
	LimitValue string `json:"limitValue"`
}
