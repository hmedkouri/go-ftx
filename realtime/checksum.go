package realtime

import (
	"hash/crc32"
	"strconv"
	"strings"
)

// CalcPartialOBChecksum calculates checksum of partial OB data received from WS
func CalcPartialOBChecksum(data *Orderbook) int64 {
	return CalcPartialChecksum(data.Bids, data.Asks)
}

func CalcPartialChecksum(bids [][2]float64, asks [][2]float64) int64 {
	var checksum strings.Builder
	var price, amount string
	for i := 0; i < 100; i++ {
		if len(bids)-1 >= i {
			price = checksumParseNumber(bids[i][0])
			amount = checksumParseNumber(bids[i][1])
			checksum.WriteString(price + ":" + amount + ":")
		}
		if len(asks)-1 >= i {
			price = checksumParseNumber(asks[i][0])
			amount = checksumParseNumber(asks[i][1])
			checksum.WriteString(price + ":" + amount + ":")
		}
	}
	checksumStr := strings.TrimSuffix(checksum.String(), ":")
	return int64(crc32.ChecksumIEEE([]byte(checksumStr)))
}

func checksumParseNumber(num float64) string {
	modifier := byte('f')
	if num < 0.0001 {
		modifier = 'e'
	}
	r := strconv.FormatFloat(num, modifier, -1, 64)
	if strings.IndexByte(r, '.') == -1 && modifier != 'e' {
		r += ".0"
	}
	return r
}