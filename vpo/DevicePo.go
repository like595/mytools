package vpo

// 设备信息
type DevicePo struct {
	Id               string //id
	CreateBy         string //创建人
	CreateTime       string //创建日期
	UpdateBy         string //更新人
	UpdateTime       string //更新日期
	SysOrgCode       string //所属部门
	DeviceTypeId     string //设备类型id
	SiteType         string //地点类型；1-隧道；2-路段；
	PileNo           string //桩号
	Direction        string //方向
	Lane             string //车道号
	Xx               string //x坐标
	Yy               string //y坐标
	Cycle            int    //周期
	CommType         string //通讯类型
	CommData         string //通讯参数
	DeviceData1      string //参数1
	DeviceData2      string //参数2
	DeviceData3      string //参数3
	DeviceData4      string //参数4
	DeviceData5      string //参数5
	DeviceData6      string //参数6
	DeviceData7      string //参数7
	DeviceData8      string //参数8
	DeviceData9      string //参数9
	DeviceData10     string //参数10
	InstallDate      string //安装日期
	//设备状态 1-正常；2-故障；3-告警；4-离线；
	DeviceStatus     int
	DeviceStatsuTime string //设备状态更新时间
	ServerId         string //所在服务器
	Name             string //设备名称
	SiteId           string //地点id
	ShiGongDanWei    string //施工单位
	DeviceType       string //设备类型
	DeviceClass      string //设备类别

}
