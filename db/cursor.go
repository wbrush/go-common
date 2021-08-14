package db

import (
	"encoding/base64"
	"strconv"
)

func EncodeIdToCursor(id int64) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(id, 10)))
}

func DecodeCursorToId(cursor string) (int64, error) {
	res, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseInt(string(res), 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
