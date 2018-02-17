package examples

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"io"
	"os"
	"strconv"
	"time"
	"github.com/blazecrystal/beyondts-go/logq"
)

func main() {
	os.Setenv("user.dir", "D:\\workspaces\\workspace-go\\beyondts")
	logq.Go("D:\\workspaces\\workspace-go\\beyondts\\src\\beyondts\\examples\\logq.properties")
	defer logq.End()
	count := 1000
	start := time.Now()
	for ; count > 0; count-- {
		logq.GetLogger("test2.abc").Debug("aaaa----", strconv.Itoa(count))
		//logq.GetLogger("test2").Error(io.EOF, "aaaa----", strconv.Itoa(count))
	}
	fmt.Println(time.Since(start))
	time.Sleep(time.Millisecond * 3000)
	fmt.Println("over!")
}
