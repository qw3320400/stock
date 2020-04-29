package utils

import (
	"fmt"
	"strconv"
	"time"
)

func Log(msg string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), msg)
}

func LogErr(msg string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "{Error}", msg)
}

/*
在str的左或右添加0
:param str:待修改的字符串
:param length:总共的长度
:param direction:方向，True左，False右
:return:
*/
func AddZeroForString(content, length int64, direction bool) string {
	contentStr := strconv.FormatInt(content, 10)
	if len(contentStr) >= int(length) {
		return contentStr
	}
	fixStr := ""
	for i := 0; i < int(length)-len(contentStr); i++ {
		fixStr += "0"
	}
	if direction {
		return fixStr + contentStr
	} else {
		return contentStr + fixStr
	}
}

func StringIsIn(in string, arg ...string) bool {
	for _, v := range arg {
		if in == v {
			return true
		}
	}
	return false
}
