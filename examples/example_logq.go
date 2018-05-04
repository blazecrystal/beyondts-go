package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"io"
	"os"
	"strconv"
	"time"
	"github.com/blazecrystal/beyondts-go/logq"
    "runtime"
)

func main() {
	os.Setenv("user.dir", "D:\\workspaces\\workspace-go\\beyondts")
	logq.Config("D:\\workspaces\\workspace-go\\beyondts\\src\\github.com\\blazecrystal\\beyondts-go\\examples\\logq.xml")
	//defer logq.Stop()
	count := 5000000
	start := time.Now()
	for ; count > 0; count-- {
		logq.GetLogger("test2/abc/dd").Info("aaaa----{}==={}=={}=", strconv.Itoa(count), "dddddd")
		//logq.GetLogger("test").Error(io.EOF, "aaaa----", strconv.Itoa(count))
	}
	fmt.Println(time.Since(start))
    mpr := &runtime.MemStats{}
    runtime.ReadMemStats(mpr)
    fmt.Println("total :", mpr.TotalAlloc, "; inuse :", mpr.Alloc, "; free :", mpr.Frees, "; sys :", mpr.Sys, "; other :", mpr.OtherSys)
    fmt.Println(mpr.BySize)
	time.Sleep(time.Millisecond * 3000)
	logq.Stop()
	fmt.Println("over!")
}
