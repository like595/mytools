package vtools

import (
	"fmt"
	"gopkg.in/ini.v1"
	"strconv"
	"strings"
)


type IniUtil struct {
	//配置信息
	iniFile *ini.File
}

func (this *IniUtil)Init(filePath string) {
	file, e := ini.Load(filePath)
	if e != nil {
		SugarLogger.Info("Fail to load " ,filePath, e.Error())

		return
	}
	this.iniFile = file
}

func (this *IniUtil)GetSection(sectionName string) *ini.Section {
	section, e := this.iniFile.GetSection(sectionName)
	if e != nil {
		SugarLogger.Info("未找到对应的配置信息:" + sectionName + e.Error())
		return nil
	}
	return section
}

func (this *IniUtil)GetSectionMap(sectionName string) map[string]string {
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

/**
获取字符串值
*/
func (this *IniUtil)GetVal(sectionName string, key string) string {
	var temp_val string
	section := this.GetSection(sectionName)
	if nil != section {
		temp_val = section.Key(key).Value()
	}
	return temp_val;
}

/**
获取字符串数组,通过,分割
*/
func (this *IniUtil)GetArr(sectionName string, key string) []string {
	temp_val := this.GetVal(sectionName, key)
	value := strings.Split(temp_val, ",")
	return value
}

/**
获取布尔值
*/
func (this *IniUtil)GetBool(sectionName string, key string) bool {
	temp_val := this.GetVal(sectionName, key)
	value, e := strconv.ParseBool(temp_val)
	if nil != e {
		SugarLogger.Error(e)
	}
	return value
}

/**
获取int
*/
func (this *IniUtil)GetInt(sectionName string, key string) int {
	temp_val := this.GetVal(sectionName, key)
	value, e := strconv.Atoi(temp_val)
	if nil != e {
		SugarLogger.Error(e)
	}
	return value
}

/**
获取int64
*/
func (this *IniUtil)GetInt64(sectionName string, key string) int64 {
	temp_val := this.GetVal(sectionName, key)
	value, e := strconv.ParseInt(temp_val, 0, 64);
	if nil != e {
		SugarLogger.Error(e)
	}
	return value
}

/**
获取float
*/
func (this *IniUtil)GetFloat(sectionName string, key string) float64 {
	temp_val := this.GetVal(sectionName, key)
	value, e := strconv.ParseFloat(fmt.Sprintf("%.2f", temp_val), 64)
	if nil != e {
		SugarLogger.Error(e)
	}
	return value
}