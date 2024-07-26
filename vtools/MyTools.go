package vtools

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
