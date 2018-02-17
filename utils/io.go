package utils

import (
	"bufio"
	"io"
)

func PipeTo(reader io.Reader, writer io.Writer, endFlag interface{}) error {
	byteCount := 0
	r := bufio.NewReader(reader)
	w := bufio.NewWriter(writer)
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil {
			switch endFlag.(type) {
			case error:
				if err == endFlag {
					w.Write(buf[:n])
					w.Flush()
					goto r
				} else {
					return err
				}
			}
		}
		switch endFlag.(type) {
		case int:
			if byteCount+n > endFlag.(int) {
				w.Write(buf[:endFlag.(int)-byteCount])
				w.Flush()
				goto r
			} else {
				byteCount += n
				w.Write(buf[:n])
				w.Flush()
				if n < 1 {
					goto r
				}
			}
		/*case string:
		  if string(buf[:n]) == endFlag.(string) {
		      goto r
		  } else {
		      byteCount += n
		      w.Write(buf[:n])
		      w.Flush()
		      if n < 1 {
		          goto r
		      }
		  }*/ // current not support string, string is complex, we may read this endFlag string in sperate 2 parts
		default:
			if n > 0 {
				w.Write(buf[:n])
				w.Flush()
			} else {
				goto r
			}
		}
	}
r:
	return nil
}
