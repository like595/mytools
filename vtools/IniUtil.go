package vtools

import (
	"fmt"
	"gopkg.in/ini.v1"
	"strconv"
	"strings"
)

// 读取配置文件
type IniUtil struct {
	//配置信息
	iniFile *ini.File
}

// 初始化
// filePath文件路径
func (this *IniUtil) Init(filePath string) {
	if FileExists("../." + filePath) {
		filePath = "../." + filePath
	} else if FileExists("." + filePath) {
		filePath = "." + filePath
	}
	file, e := ini.Load(filePath)
	if e != nil {
		SugarLogger.Info("Fail to load ", filePath, e.Error())

		return
	}
	this.iniFile = file
}

func (this *IniUtil) GetSection(sectionName string) *ini.Section {
	section, e := this.iniFile.GetSection(sectionName)
	if e != nil {
		SugarLogger.Info("未找到对应的配置信息:" + sectionName + e.Error())
		return nil
	}
	return section
}

func (this *IniUtil) GetSectionMap(sectionName string) map[string]string {
	section, e := this.iniFile.GetSection(sectionName)
	if e != nil {
		SugarLogger.Info("未找到对应的配置信息:" + sectionName + e.Error())
		return nil
	}
	section_map := make(map[string]string, 0)
	for _, e := range section.Keys() {
		section_map[e.Name()] = e.Value()
	}
	return section_map
}

// 获取字符串值
func (this *IniUtil) GetVal(sectionName string, key string) string {
	var temp_val string
	section := this.GetSection(sectionName)
	if nil != section {
		temp_val = section.Key(key).Value()
	}
	return temp_val
}

// 获取字符串数组,通过,分割
func (this *IniUtil) GetArr(sectionName string, key string) []string {
	temp_val := this.GetVal(sectionName, key)
	value := strings.Split(temp_val, ",")
	return value
}

// 获取布尔值
func (this *IniUtil) GetBool(sectionName string, key string) bool {
	temp_val := this.GetVal(sectionName, key)
	value, e := strconv.ParseBool(temp_val)
	if nil != e {
		SugarLogger.Error(e)
	}
	return value
}

// 获取int
func (this *IniUtil) GetInt(sectionName string, key string) int {
	temp_val := this.GetVal(sectionName, key)
	value, e := strconv.Atoi(temp_val)
	if nil != e {
		SugarLogger.Error(e)
	}
	return value
}

// 获取int64
func (this *IniUtil) GetInt64(sectionName string, key string) int64 {
	temp_val := this.GetVal(sectionName, key)
	value, e := strconv.ParseInt(temp_val, 0, 64)
	if nil != e {
		SugarLogger.Error(e)
	}
	return value
}

// 获取float
func (this *IniUtil) GetFloat(sectionName string, key string) float64 {
	temp_val := this.GetVal(sectionName, key)
	value, e := strconv.ParseFloat(fmt.Sprintf("%.2f", temp_val), 64)
	if nil != e {
		SugarLogger.Error(e)
	}
	return value
}

// Crc16Ccitt 计算CRC-16/CCITT校验值（大端序）
func Crc16Ccitt(data []byte) []byte {
	poly := uint16(0x1021) // 多项式
	crc := uint16(0x0000)  // 初始值

	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc <<= 1
			}
			crc &= 0xFFFF // 保持16位长度
		}
	}

	// 返回大端序字节数组
	//return []byte{
	//	byte(crc >> 8),byte(crc & 0xFF)
	//}
	return []byte{byte(crc >> 8), byte(crc & 0xFF)}
}
