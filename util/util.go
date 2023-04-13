package util

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
)

// 正则的封装
func QuickRegexp(raw string, patten string) ([][]string, error) {
	reg := regexp.MustCompile(patten)
	res := reg.FindAllStringSubmatch(raw, -1)
	if len(res) == 0 {
		return nil, errors.New("no match")
	}
	return res, nil
}

func JSONResponse(code int, msg string, data interface{}) (int, map[string]interface{}) {
	return code, gin.H{"code": code, "msg": msg, "data": data}
}
