package crypto

import (
    "github.com/blazecrystal/beyondts-go/utils"
)

const (
    LENGTH_IV = 32
    LENGTH_TJ = 64
    LENGTH_BYTES_OF_INT = 4
    LENGTH_DIGEST = 32
    LENGTH_BLOCK = 64
    TJ_START = 0xA748B290
    DEFAULT_IV = "a12BC!@&def568%$Gh9#ijKlmNoPqr43"
    DEFAULT_GAP = 256
)

type SM3 struct {
    src []byte
    bufOff, block int
    gap int16
    tj []int32
    buf, v []int8
    iv string
}

func SM3DigestString(src string, base64Encoding bool) string {
    return utils.Bytes2String(SM3Digest([]byte(src)), base64Encoding)
}

func SM3DigestString2(src string, base64Encoding bool, iv string, gap int16) string {
    return utils.Bytes2String(SM3Digest2([]byte(src), iv, gap), base64Encoding)
}

func SM3Digest(src []byte) []byte {
    sm3 := NewSM3()
    sm3.update(src)
    return sm3.doFinal()
}

func SM3Digest2(src []byte, iv string, gap int16) []byte {
    sm3 := New2SM3(iv, gap)
    sm3.update(src)
    return sm3.doFinal()
}

func NewSM3() *SM3 {
    return New2SM3(DEFAULT_IV, DEFAULT_GAP)
}

func New2SM3(iv string, gap int16) *SM3 {
    return &SM3{iv:iv, gap:gap, tj:genTj(gap), buf:make([]int8, LENGTH_BLOCK), v:genIV(iv)}
}

func Clone(sm3 SM3) *SM3 {
    copy := &SM3{iv:sm3.iv, gap:sm3.gap, tj:genTj(sm3.gap), buf:make([]int8, LENGTH_BLOCK), v:make([]int8, LENGTH_IV)}
    utils.CopySlice(sm3.buf, 0, copy.buf, 0, len(sm3.buf))
    copy.bufOff = sm3.bufOff
    utils.CopySlice(sm3.v, 0, copy.v, 0, len(sm3.v))
    return copy
}

func (sm3 *SM3) Reset() {
    sm3.bufOff = 0
    sm3.block = 0
    sm3.v = genIV(sm3.iv)
    for k, _ := range sm3.buf {
        sm3.buf[k] = 0
    }
}

func (sm3 *SM3) Sum(b []byte) []byte {
    if b == nil {
        sm3.update(b)
        return sm3.doFinal()
    }
    return sm3.doFinal()
}

func (sm3 *SM3) Write(b []byte) (nn int, err error)  {
    sm3.update(b)
    return len(b), nil
}

func (sm3 *SM3) Size() int {
    return LENGTH_DIGEST
}

func (sm3 *SM3) BlockSize() int {
    return LENGTH_BLOCK
}

func (sm3 *SM3) update(in []byte) {
    partLen := LENGTH_BLOCK - sm3.bufOff
    inputLen := len(in)
    dPos := 0
    if partLen < inputLen {
        copyByte2Int8Slice(in, dPos, sm3.buf, sm3.bufOff, partLen)
        inputLen -= partLen
        dPos += partLen
        sm3.doUpdate()
        for inputLen > LENGTH_BLOCK {
            copyByte2Int8Slice(in, dPos, sm3.buf, 0, LENGTH_BLOCK)
            inputLen -= LENGTH_BLOCK
            dPos += LENGTH_BLOCK
            sm3.doUpdate()
        }
    }
    copyByte2Int8Slice(in, dPos, sm3.buf, sm3.bufOff, inputLen)
    sm3.bufOff += inputLen
}

/*func (sm3 *SM3) UpdateByte(in byte) {
    buffer := []byte{in}
    sm3.Update(buffer)
}*/

func (sm3 *SM3) doFinal() []byte {
    out := make([]byte, LENGTH_DIGEST)
    tmp := sm3.doFinal1()
    utils.CopySlice(tmp, 0, out, 0, len(tmp))
    sm3.Reset()
    return out
}

func copyByte2Int8Slice(src []byte, srcOffset int, dst []int8, dstOffset int, length int) {
    for i := 0; i < length; i++ {
        dst[dstOffset + i] = int8(src[srcOffset + i])
    }
}

func copyInt82ByteSlice(src []int8, srcOffset int, dst []byte, dstOffset int, length int) {
    for i := 0; i < length; i++ {
        dst[dstOffset + i] = byte(src[srcOffset + i])
    }
}

func (sm3 *SM3) doFinal1() []byte {
    b := make([]int8, LENGTH_BLOCK)
    buffer := make([]int8, sm3.bufOff)
    utils.CopySlice(sm3.buf, 0, buffer, 0, len(buffer))
    tmp := padding(buffer, sm3.block)
    for i := 0; i < len(tmp); i += LENGTH_BLOCK {
        utils.CopySlice(tmp, i, b, 0, len(b))
        sm3.doHash(b)
    }
    return utils.Int82Bytes(sm3.v)
}

func (sm3 *SM3) doHash(b []int8) {
    tmp := sm3.cfByte(sm3.v, b)
    utils.CopySlice(tmp, 0, sm3.v, 0, len(sm3.v))
    sm3.block++
}

func (sm3 *SM3) doUpdate() {
    b := make([]int8, LENGTH_BLOCK)
    for i := 0; i < LENGTH_BLOCK; i += LENGTH_BLOCK {
        utils.CopySlice(sm3.buf, i, b, 0, len(b))
        sm3.doHash(b)
    }
    sm3.bufOff = 0
}

func genTj(gap int16) []int32 {
    tj := make([]int32, LENGTH_TJ)
    for i := 0; i < LENGTH_TJ; i++ {
        tj[i] = int32(TJ_START + i * int(gap))
    }
    return tj
}

func genIV(iv string) []int8 {
    ivBytes := make([]int8, LENGTH_IV)
    tmp := []byte(iv)
    for i := 0; i < LENGTH_IV; i++ {
        ivBytes[i] = int8(tmp[i % len(tmp)])
    }
    return ivBytes
}

func (sm3 *SM3) cfByte(v, b []int8) []int8 {
    tmpV := convertByte2Int(v);
    tmpB := convertByte2Int(b);
    return convertInt2Byte(sm3.cfInt32(tmpV, tmpB));
}

func (sm3 *SM3) cfInt32(v, x []int32) []int32 {
    var a, b, c, d, e, f, g, h, ss1, ss2, tt1, tt2 int32
    a = v[0];
    b = v[1];
    c = v[2];
    d = v[3];
    e = v[4];
    f = v[5];
    g = v[6];
    h = v[7];

    w, w1 := expand(x);

    for j := 0; j < 64; j++ {
        ss1 = (bitCycleLeft(a, 12) + e + bitCycleLeft(sm3.tj[j], j));
        ss1 = bitCycleLeft(ss1, 7);
        ss2 = ss1 ^ bitCycleLeft(a, 12);
        tt1 = ffj(a, b, c, int32(j)) + d + ss2 + w1[j];
        tt2 = ggj(e, f, g, int32(j)) + h + ss1 + w[j];
        d = c;
        c = bitCycleLeft(b, 9);
        b = a;
        a = tt1;
        h = g;
        g = bitCycleLeft(f, 19);
        f = e;
        e = p0(tt2);

    }

    out := make([]int32, 8)
    out[0] = a ^ v[0];
    out[1] = b ^ v[1];
    out[2] = c ^ v[2];
    out[3] = d ^ v[3];
    out[4] = e ^ v[4];
    out[5] = f ^ v[5];
    out[6] = g ^ v[6];
    out[7] = h ^ v[7];

    return out;
}

func expand(b []int32) ([]int32, []int32) {
    w := make([]int32 , 68)
    w1 := make([]int32 , 64)
    for i := 0; i < len(b); i++ {
        w[i] = b[i]
    }

    for i := 16; i < 68; i++ {
        w[i] = p1(w[i - 16] ^ w[i - 9] ^ bitCycleLeft(w[i - 3], 15)) ^ bitCycleLeft(w[i - 13], 7) ^ w[i - 6]
    }

    for i := 0; i < 64; i++ {
        w1[i] = w[i] ^ w[i + 4];
    }

    return w, w1
}

func ffj(x, y, z, j int32) int32 {
    if j >= 0 && j <= 15 {
        return x ^ y ^ z
    }
    return (x & y) | (x & z) | (y & z)
}

func ggj(x, y, z, j int32) int32 {
    if j >= 0 && j <= 15 {
        return x ^ y ^ z
    }
    return (x & y) | (^x & z)
}

func p0(x int32) int32 {
    rotateLeft(x, 9)
    y := bitCycleLeft(x, 9)
    rotateLeft(x, 17)
    z := bitCycleLeft(x, 17)
    return x ^ y ^ z
}

func rotateLeft(x int32, n uint) int32 {
    return (x << n) | (x >> (32 - n))
}

func p1(x int32) int32 {
    return x ^ bitCycleLeft(x, 15) ^ bitCycleLeft(x, 23)
}

func bitCycleLeft(n int32, bitLen int) int32 {
    bitLen %= 32;
    tmp := utils.Int322Int8(n);
    byteLen := bitLen / 8;
    length := bitLen % 8;
    if (byteLen > 0) {
        tmp = byteCycleLeft(tmp, byteLen);
    }

    if length > 0 {
        tmp = bitSmall8CycleLeft(tmp, uint(length));
    }

    return utils.Int82Int32(tmp)
}

func bitSmall8CycleLeft(in []int8, length uint) []int8 {
    tmp := make([]int8, len(in))
    for i := 0; i < len(tmp); i++ {
        t1 := (byte(in[i]) & 0x000000FF) << length
        t2 := (byte(in[(i + 1) % len(tmp)]) & 0x000000FF) >> (8 - length)
        t3 := byte(t1 | t2)
        tmp[i] = int8(t3)
    }
    return tmp
}

func byteCycleLeft(in []int8, byteLen int) []int8 {
    tmp := make([]int8, len(in))
    utils.CopySlice(in, byteLen, tmp, 0, len(in) - byteLen)
    utils.CopySlice(in, 0, tmp, len(in) - byteLen, byteLen)
    return tmp
}

func convertByte2Int(arr []int8) []int32 {
    out := make([]int32, len(arr) / LENGTH_BYTES_OF_INT)
    tmp := make([]byte, LENGTH_BYTES_OF_INT)
    for i := 0; i < len(arr); i += LENGTH_BYTES_OF_INT {
        //utils.CopySlice(arr, i, tmp, 0, LENGTH_BYTES_OF_INT)
        copyInt82ByteSlice(arr, i, tmp, 0, LENGTH_BYTES_OF_INT)
        out[i / LENGTH_BYTES_OF_INT] = utils.Bytes2Int32(tmp)
    }
    return out
}

func convertInt2Byte(arr []int32) []int8 {
    out := make([]int8, len(arr) * LENGTH_BYTES_OF_INT)
    for i := 0; i < len(arr); i++ {
        tmp := utils.Int322Int8(arr[i])
        for j := 0; j < LENGTH_BYTES_OF_INT; j++ {
            out[i * LENGTH_BYTES_OF_INT + j] = tmp[j]
        }
    }
    return out
}

func padding(in []int8, bLen int) []int8 {
    k := 448 - (8 * len(in) + 1) % 512
    if k < 0 {
        k = 960 - (8 * len(in) + 1) % 512
    }
    k++
    padd := make([]int8, k / 8)
    padd[0] = int8(-128)
    n := len(in) * 8 + bLen * 512
    out := make([]int8, len(in) + k / 8 + 8)
    pos := 0
    utils.CopySlice(in, 0, out, 0, len(in))
    pos += len(in)
    utils.CopySlice(padd, 0, out, pos, len(padd))
    pos += len(padd)
    tmp := utils.Int642Int8(int64(n))
    utils.ReverseInt8Slice(tmp)
    utils.CopySlice(tmp, 0, out, pos, len(tmp))
    return out
}