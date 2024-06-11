package cache

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	B = 1 << (iota * 10)
	KB
	MB
	GB
	TB
	PB
)

func ParseSize(size string) (int64, string) {
	re, _ := regexp.Compile("[0-9]+")

	numStr := re.FindString(size)

	unit := strings.TrimSpace(re.ReplaceAllString(size, ""))

	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		log.Println("parse size fail : invalid number format")
		num = 100
		unit = "MB"
	}
	var byteNum int64 = 0
	switch unit {
	case "B":
		byteNum = num
	case "KB":
		byteNum = KB * num
	case "MB":
		byteNum = MB * num
	case "GB":
		byteNum = GB * num
	case "TB":
		byteNum = TB * num
	case "PB":
		byteNum = PB * num
	default:
		num = 0
		log.Println("parse size fail: invalid unit")
		num = 100
		unit = "MB"
		byteNum = num * MB
	}

	sizeStr := strconv.FormatInt(num, 10) + unit
	return byteNum, sizeStr
}
func GetValSize(val interface{}) int64 {
	bytes, _ := json.Marshal(val)
	size := int64(len(bytes))
	return size
}
