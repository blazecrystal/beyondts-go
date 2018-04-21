package utils

import (
    "reflect"
)

func CopySlice(src interface{}, srcOffset int, dst interface{}, dstOffset int, length int) {
    srcTmp := reflect.ValueOf(src)
    dstTmp := reflect.ValueOf(dst)
    for i := 0; i < length; i++ {
        dstTmp.Index(dstOffset + i).Set(srcTmp.Index(srcOffset + i))
    }
}

func ReverseByteSlice(src []byte) {
    length := len(src)
    for i := 0; i < length; i++ {
        if (i == length - i - 1) {
            break
        }
        src[i], src[length - i - 1] = src[length - i - 1], src[i]
    }
}


func ReverseInt8Slice(src []int8) {
    length := len(src)
    for i := 0; i < length; i++ {
        if (i >= length - i - 1) {
            break
        }
        src[i], src[length - i - 1] = src[length - i - 1], src[i]
    }
}