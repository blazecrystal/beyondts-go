package utils

import (
    "strings"
    "encoding/base64"
    "encoding/hex"
    "bytes"
)

const(
    UNIT = 8
    MASK = 0x000000FF
    PER_LENGTH_INT = 4
    PER_LENGTH_INT64 = 8
)

func Bytes2Int(src []byte) int {
    dst := 0
    tmp := 0
    for i := 0; i < len(src); i++ {
        dst <<= UNIT
        tmp = int(src[i] & MASK)
        dst |= tmp
    }
    return dst
}


func Int82Int(src []int8) int {
    dst := 0
    tmp := 0
    for i := 0; i < len(src); i++ {
        dst <<= UNIT
        tmp = int(src[i]) & MASK
        dst |= tmp
    }
    return dst
}
func Int82Int32(src []int8) int32 {
    return int32(Int82Int(src))
}

func Bytes2Int32(src []byte) int32 {
    return int32(Bytes2Int(src))
}

func Int2Bytes(src int) []byte {
    dst := make([]byte, PER_LENGTH_INT)
    for i := 0; i < len(dst); i++ {
        dst[i] = byte((src >> (uint(len(dst) - 1 - i) * UNIT)) & MASK)
    }
    return dst
}

func Int322Int8(src int32) []int8 {
    dst := make([]int8, PER_LENGTH_INT)
    for i := 0; i < len(dst); i++ {
        dst[i] = int8(src >> (uint(len(dst) - 1 - i) * UNIT))
    }
    return dst
}

func Int642Bytes(src int64) []byte {
    dst := make([]byte, PER_LENGTH_INT64)
    for i := 0; i < len(dst); i++ {
        dst[i] = byte((src >> uint((len(dst) - 1 - i) * UNIT)) & MASK);
    }
    return dst;
}


func Int642Int8(src int64) []int8 {
    dst := make([]int8, PER_LENGTH_INT64)
    for i := 0; i < len(dst); i++ {
        dst[i] = int8((src >> uint((len(dst) - 1 - i) * UNIT)) & MASK);
    }
    return dst;
}

func Int82Bytes(src []int8) []byte {
    rst := make([]byte, len(src))
    for k, v := range src {
        rst[k] = byte(v)
    }
    return rst
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