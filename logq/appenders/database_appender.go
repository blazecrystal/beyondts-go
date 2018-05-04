package appenders

import (
    "database/sql"
    "strings"
    "github.com/blazecrystal/beyondts-go/utils"
    "time"
    "runtime"
    "os"
    "strconv"
    "fmt"
)

const (
    DEFAULT_DRIVER      = "mysql"
    DEFAULT_MAX_RID_LEN = 16
)

type DatabaseAppender struct {
    name, driver, url, sql string
    maxRidLen              int
    params                 [] string
    con                    *sql.DB
    stmt                   *sql.Stmt
}

func NewDatabaseAppender(name, driver, url, sqlstr string, maxRidLen int, params []string) (*DatabaseAppender, error) {
    if len(driver) < 1 {
        driver = DEFAULT_DRIVER
    }
    con, err := sql.Open(driver, url)
    if err != nil {
        return nil, err
    }
    stmt, err := con.Prepare(sqlstr)
    if err != nil {
        return nil, err
    }
    return &DatabaseAppender{name: name, driver: driver, url: url, sql: sqlstr, maxRidLen: maxRidLen, params: params, con: con, stmt: stmt}, nil
}

func (d *DatabaseAppender) GetName() string {
    return d.name
}

func (d *DatabaseAppender) GetType() string {
    return TYPE_DATABASE
}

func (d *DatabaseAppender) GetId() string {
    return TYPE_DATABASE + "#" + d.driver + "#" + d.url + "#" + d.sql + "#" + strings.Join(d.params, "#")
}

func (d *DatabaseAppender) Stop() {
    if d.stmt != nil {
        d.stmt.Close()
    }
    if d.con != nil {
        d.con.Close()
    }
}

func (d *DatabaseAppender) Write(loggerName, level, content string) {
    params := d.buildLog(loggerName, level, content)
    _, err := d.stmt.Exec(params...)
    if err != nil {
        fmt.Println(err)
    }
}

func (d *DatabaseAppender) buildLog(loggerName, level, content string) []interface{} {
    params := make([]interface{}, len(d.params))
    var rid string
    if d.maxRidLen > 0 {
        rid = utils.RandomString(d.maxRidLen)
    }
    var lf, sf, n string
    _, file, line, ok := runtime.Caller(2) // skip Caller & current stack
    if ok {
        file = utils.ToLocalFilePath(file)
        lf = file
        lastPathSeptIndex := strings.LastIndex(file, string(os.PathSeparator))
        sf = file[lastPathSeptIndex+1: ]
        n = strconv.Itoa(line)
    } else {
        lf = UNKNOWN
        sf = UNKNOWN
        n = UNKNOWN
    }
    for i, param := range d.params {
        params[i] = param
        if strings.Contains(params[i].(string), RANDOM_STRING_ID) {
            params[i] = strings.Replace(params[i].(string), RANDOM_STRING_ID, rid, -1)
        }
        if strings.Contains(params[i].(string), TIMESTAMP) {
            params[i] = strings.Replace(params[i].(string), TIMESTAMP, time.Now().Format("2006-01-02 15:04:05.000"), -1)
        }
        if strings.Contains(params[i].(string), LEVEL) {
            params[i] = strings.Replace(params[i].(string), LEVEL, level, -1)
        }
        if strings.Contains(params[i].(string), LONG_FILE_NAME) {
            params[i] = strings.Replace(params[i].(string), LONG_FILE_NAME, lf, -1)
        }
        if strings.Contains(params[i].(string), SHORT_FILE_NAME) {
            params[i] = strings.Replace(params[i].(string), SHORT_FILE_NAME, sf, -1)
        }
        if strings.Contains(params[i].(string), LINE_NUMBER) {
            params[i] = strings.Replace(params[i].(string), LINE_NUMBER, n, -1)
        }
        if strings.Contains(params[i].(string), LOGGER_NAME) {
            params[i] = strings.Replace(params[i].(string), LOGGER_NAME, loggerName, -1)
        }
        if strings.Contains(params[i].(string), CONTENT) {
            params[i] = strings.Replace(params[i].(string), CONTENT, content, -1)
        }
    }
    return params
}
