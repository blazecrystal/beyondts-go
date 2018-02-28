package utils

import (
	"bytes"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"fmt"
	"encoding/base64"
	"encoding/hex"
)

var (
	ALPHA_NUM = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	HEX       = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
)

func Concat(val ...interface{}) string {
	buf := bytes.Buffer{}
	for _, v := range val {
		fmt.Fprint(&buf, v)
	}
	return buf.String()
}

func Concat2(sep string, val ...interface{}) string {
	buf := bytes.Buffer{}
	for i, v := range val {
		if i < len(val) - 1 {
			fmt.Fprint(&buf, v, sep)
		} else {
			fmt.Fprint(&buf, v)
		}
	}
	return buf.String()
}

func ConcatWithSept(sept string, str ...string) string {
	buf := bytes.Buffer{}
	for i, v := range str {
		buf.WriteString(v)
		if i < len(str)-1 {
			buf.WriteString(sept)
		}
	}
	return buf.String()
}

func ExistInSliece(slice []string, str string) bool {
	if slice == nil {
		return false
	}
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func IndexInSlice(slice []string, str string) int {
	for i, v := range slice {
		if v == str {
			return i
		}
	}
	return -1
}

func IsNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func Sort(src []string, desc bool) []string {
	if len(src) > 1 {
		for i := 0; i < len(src); i++ {
			for j := i + 1; j < len(src); j++ {
				if (!desc && src[i] > src[j]) || (desc && src[i] < src[j]) {
					src[i], src[j] = src[j], src[i]
				}
			}
		}
	}
	return src
}

//var lock sync.Mutex
func RandomString(length int) string {
	//lock.Lock()
	rand.Seed(time.Now().Unix() + rand.Int63())
	//rand1 := rand.New(rand.NewSource(time.Now().Unix() + rand.Int63n(4294967296)))
	//lock.Unlock()
	r := make([]rune, length)
	for i := 0; i < length; i++ {
		r[i] = ALPHA_NUM[rand.Intn(len(ALPHA_NUM))]
	}
	return string(r)
}

func SliceAtoi(strSlice []string) ([]int, error) {
	tmp := make([]int, len(strSlice))
	for i, v := range strSlice {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		tmp[i] = int(num)
	}
	return tmp, nil
}

func HexString2DecInt(hex string) int {
	rst := 0
	for i := 0; i < len(hex); i++ {
		num := IndexInSlice(HEX, strings.ToUpper(hex[i:i+1]))
		times := 1
		for j := i; j < len(hex)-1; j++ {
			times *= 16
		}
		rst += num * times
	}
	return rst
}

func DecInt2HexString(dec int) string {
	var tmp bytes.Buffer
	current := dec
	div, mod := dec, 0
	for div > 0 {
		mod, div = current%16, current/16
		tmp.WriteString(HEX[mod])
		current = div
	}
	hexBytes := tmp.Bytes()
	internal := len(hexBytes) / 2
	for i := 0; i < internal; i++ {
		hexBytes[i], hexBytes[len(hexBytes)-1-i] = hexBytes[len(hexBytes)-1-i], hexBytes[i]
	}
	return string(hexBytes)
}

func Bytes2String(data []byte, base64Encoding bool) string {
	if base64Encoding {
		return base64.StdEncoding.EncodeToString(data)
	} else {
		return strings.ToUpper(hex.EncodeToString(data))
	}
}

func String2Bytes(data string, base64Encoding bool) ([]byte, error) {
	var tmp []byte
	var err error
	if base64Encoding {
		tmp, err = base64.StdEncoding.DecodeString(data)
	} else {
		tmp, err = hex.DecodeString(data)
	}
	return tmp, err
}