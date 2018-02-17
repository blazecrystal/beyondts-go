package appenders

import (
	"bytes"
	"database/sql"
	"strconv"
	"strings"
	"github.com/blazecrystal/beyondts-go/utils"
)

// in this appender, layout should be a correct insert sql. eg. insert into logs (logtime, level, location, logger, content) values (%d%, %p%, %sl%, %ln%, %c%)
type DatabaseAppender struct {
	name, table, columns, params, driver, url, user, pwd string
	maxIdLen                                             int
	log                                                  chan map[string]string
}

func NewDatabaseAppender(attrs map[string]string) *DatabaseAppender {
	tmp, err := strconv.Atoi(attrs["maxIdLen"])
	if err != nil {
		tmp = 0
	}
	return &DatabaseAppender{name: attrs["name"], table: attrs["table"], columns: attrs["columns"], params: attrs["params"], driver: attrs["driver"], url: attrs["url"], user: attrs["user"], pwd: attrs["pwd"], maxIdLen: tmp, log: make(chan map[string]string)}
}

func (a *DatabaseAppender) GetType() string {
	return Type_Database
}

func (a *DatabaseAppender) WriteLog(fills map[string]string) {
	a.log <- fills
}

func (a *DatabaseAppender) Stop() {
	end := make(map[string]string, 1)
	end["END"] = string(CMD_END)
	a.log <- end
}

func (a *DatabaseAppender) Run(flag chan string) {
	con, err := sql.Open(a.driver, a.url)
	defer con.Close()
	if err != nil {
		debugMsg("can't open a database connection (driver : ", a.driver, ", url : ", a.url, ", this appender will be skipped\n", err)
		return
	}
	sql := prepareSql(a.table, a.columns)
	stmt, err := con.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		debugMsg("can't prepare statement of insertion (sql : ", sql, "), this appender will be skipped\n", err)
		return
	}
	flag <- utils.Concat("+", a.name)
	for {
		tmp := <-a.log
		if bytes.Equal([]byte(tmp["END"]), CMD_END) {
			break
		}
		_, err := stmt.Exec(buildParams(a, tmp)...)
		if err != nil {
			debugMsg("failed to insert logs into database (sql : ", sql, "), this appender will continue with next log\n", err)
			continue
		}
	}
	flag <- utils.Concat("-", a.name)
}

func prepareSql(table, columns string) string {
	cols := strings.Split(columns, " ")
	sql := utils.Concat("insert into ", table, " (")
	var columnList, paramList string
	for i, col := range cols {
		columnList = utils.Concat(columnList, strings.TrimSpace(col))
		paramList = utils.Concat(paramList, "?")
		if i < len(cols)-1 {
			columnList = utils.Concat(columnList, ", ")
			paramList = utils.Concat(paramList, ", ")
		}
	}
	sql = utils.Concat(sql, columnList, ") values (", paramList, ")")
	return sql
}

func buildParams(a *DatabaseAppender, fills map[string]string) []interface{} {
	tmp := strings.Split(a.params, " ")
	paramList := make([]interface{}, len(tmp))
	for i, paramName := range tmp {
		if strings.EqualFold(paramName, RANDOM_STRING_ID) {
			paramList[i] = utils.RandomString(a.maxIdLen)
		} else {
			paramList[i] = fills[strings.TrimSpace(paramName)]
		}
	}
	return paramList
}

func (a *DatabaseAppender) Update(attrs map[string]string) {
	tmp, err := strconv.Atoi(attrs["maxIdLen"])
	if err != nil {
		tmp = 0
	}
	a.name, a.table, a.columns, a.params, a.driver, a.url, a.user, a.pwd, a.maxIdLen = attrs["name"], attrs["table"], attrs["columns"], attrs["params"], attrs["driver"], attrs["url"], attrs["user"], attrs["pwd"], tmp
}
