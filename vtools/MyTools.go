package vtools

import (
	"os"
)

// calculateCRC16 计算Modbus RTU的CRC校验
func CalculateCRC16(data []byte) []byte {
	// 初始CRC值
	const crcInitial = 0xFFFF
	// CRC-16/MODBUS多项式
	const crcPolynomial = 0xA001
	crc := crcInitial
	for _, b := range data {
		crc ^= int(b)
		for i := 0; i < 8; i++ {
			if crc&0x0001 > 0 {
				crc = (crc >> 1) ^ crcPolynomial
			} else {
				crc >>= 1
			}
		}
	}
	res := make([]byte, 2)
	res[0] = byte(crc & 0xFF)
	res[1] = byte(crc >> 8)
	return res
}

// 16位16进制，转成有符号int
func BytesToInt16(h byte, l byte) int16 {
	return int16(int(h)*0x100 + int(l))
}

// 有符号int转成16进制
func IntToBytes(data int) (byte, byte) {
	return byte(uint16(data) / 0x100), byte(uint16(data))
}

// 取byte数据的第几位
// data：待取数据；index：数据所在索引
func GetPoint(data byte, index int) int {
	data = data >> index
	return int(data) & 1
}

func BytesToString(byteArray []byte) string {
	res := ""
	for _, b := range byteArray {
		b1 := (byte)(b >> 4)
		if b1 >= 0x0 && b1 <= 0x9 {
			res += string(b1 + 0x30)
		} else {
			res += string(b1 + 0x37)
		}
		b2 := (byte)(b & 0x0F)
		if b2 >= 0x0 && b2 <= 0x9 {
			res += string(b2 + 0x30)
		} else {
			res += string(b2 + 0x37)
		}
		res += " "
	}
	return res
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// 处理其他可能的错误
	return false
}

// 判断数据是否改变，如果改变则拷贝数据到lastData。如果数据没有改变，则不拷贝。适用于Modbus协议等最后两位是校验位的数据
func IsDataChangeModbusRTU(plcData, lastData *[]byte) bool {
	isChange := false
	lenth := len((*plcData))
	if lenth != len((*lastData)) {
		return false
	}
	if (*plcData)[lenth-1] == (*lastData)[lenth-1] && (*plcData)[lenth-2] == (*lastData)[lenth-2] {
		return false
	}
	for i := 0; i < len((*plcData)) && i < len((*lastData)); i++ {
		if (*plcData)[i] != (*lastData)[i] {
			isChange = true
		}
	}
	//数据改变，拷贝数据
	if isChange {
		for i := 0; i < len((*plcData)) && i < len((*lastData)); i++ {
			(*lastData)[i] = (*plcData)[i]
		}
	}
	return isChange
}

// 判断数据是否改变，如果改变则拷贝数据到lastData。如果数据没有改变，则不拷贝。适用于Modbus协议等最后两位是校验位的数据
func IsDataChange(plcData, lastData *[]byte) bool {
	isChange := false

	for i := 0; i < len((*plcData)) && i < len((*lastData)); i++ {
		if (*plcData)[i] != (*lastData)[i] {
			isChange = true
		}
	}
	//数据改变，拷贝数据
	if isChange {
		for i := 0; i < len((*plcData)) && i < len((*lastData)); i++ {
			(*lastData)[i] = (*plcData)[i]
		}
	}
	return isChange
}
