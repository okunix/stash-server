package webutil

import (
	"net/http"
	"strconv"
)

func GetUintQueryParam(
	r *http.Request,
	key string,
	bitSize int,
	defaultValue ...uint64,
) (uint64, error) {
	uintStr := r.URL.Query().Get(key)
	res, err := strconv.ParseUint(uintStr, 10, bitSize)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
	}
	return res, err
}

func GetIntQueryParam(
	r *http.Request,
	key string,
	bitSize int,
	defaultValue ...int64,
) (int64, error) {
	uintStr := r.URL.Query().Get(key)
	res, err := strconv.ParseInt(uintStr, 10, bitSize)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
	}
	return res, err
}
