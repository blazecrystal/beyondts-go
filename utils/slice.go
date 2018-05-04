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

// attention: returned slice is not src
func RemoveFromSlice(src interface{}, indexToRemove int) interface{} {
    slice := reflect.ValueOf(src)
    length := slice.Len()
    if indexToRemove < length && indexToRemove > -1 {
        return reflect.AppendSlice(slice.Slice(0, indexToRemove), slice.Slice(indexToRemove + 1, length)).Interface()
    }
    return src
}