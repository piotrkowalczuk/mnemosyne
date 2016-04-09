package mnemosyne

import "strconv"

func address(host string, port int) string {
	return host + ":" + strconv.FormatInt(int64(port), 10)
}
