package api

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hancheo/tistory/common"
)

//ContentsStock make contents struct
type ContentsStock struct {
	StockName   string
	CostVolume  string
	CountVolume string
	OneStock    string
}

//GetAccessToken 티스토리 접근 토큰 발급
func GetAccessToken(code string) string {
	tistoryAuth := "https://www.tistory.com/oauth/access_token?"
	tistoryAuth += "client_id=" + os.Getenv("CLIENT_ID") + ""
	tistoryAuth += "&client_secret=" + os.Getenv("CLIENT_SECRET") + ""
	tistoryAuth += "&redirect_uri=" + os.Getenv("REDIRECT_URL") + ""
	tistoryAuth += "&code=" + code
	tistoryAuth += "&grant_type=authorization_code"
	resp, err := http.Get(tistoryAuth)
	common.CheckErr(err)
	respBody, err := ioutil.ReadAll(resp.Body)
	return strings.Split(string(respBody), "=")[1]
}

// GetStockDataDaily 찾을 테이블, 매수 주체, 정렬 방식 (desc == 순매수 상위 "" == 순매도 상위)
func GetStockDataDaily(table, buyType, order, modifyDate string, db *sql.DB) []ContentsStock {
	var sql, date string

	date = common.GetDateFormat("yyyy", "", "") + modifyDate

	//검색
	sql = "select a.stockname, a." + buyType + ", b." + buyType + ", " +
		"case " +
		"when ABS(b." + buyType + ") > 0 " +
		"then (ABS(a." + buyType + ") * 1000000 / ABS(b." + buyType + "))" +
		"else 0 " +
		"end onestock " +
		"from daily.cost_craw a " +
		"left join daily.stock_craw b " +
		"on a.stockname = b.stockname and a.crawdate = b.crawdate " +
		"where a.crawdate = '" + date + "' and a.stocktype = " + table + " " +
		"order by a." + buyType + " " + order + " limit 50"

	rows, err := db.Query(sql)
	common.CheckErr(err)
	var row ContentsStock
	var returnData []ContentsStock
	count := 0
	for rows.Next() {
		count++
		rows.Scan(&row.StockName, &row.CostVolume, &row.CountVolume, &row.OneStock)
		returnData = append(returnData, row)
	}
	if count == 0 {
		log.Fatal(errors.New(date + "ROW가 0개 입니다"))
	}

	return returnData
}

// GetStockDataDaily 찾을 테이블, 매수 주체, 정렬 방식 (desc == 순매수 상위 "" == 순매도 상위)
func GetStockDataMonthly(table, buyType, order, modifyDate string, db *sql.DB) []ContentsStock {
	var sql, date string

	date = common.GetDateFormat("yyyy", "", "") + modifyDate

	//검색
	sql = "select a.stockname, a." + buyType + ", b." + buyType + ", " +
		"case " +
		"when ABS(b." + buyType + ") > 0 " +
		"then (ABS(a." + buyType + ") * 1000000 / ABS(b." + buyType + "))" +
		"else 0 " +
		"end onestock " +
		"from monthly.cost_craw a " +
		"left join monthly.stock_craw b " +
		"on a.stockname = b.stockname and a.crawdate = b.crawdate " +
		"where a.crawdate = '" + date + "' and a.stocktype = " + table + " " +
		"order by a." + buyType + " " + order + " limit 100"

	rows, err := db.Query(sql)
	common.CheckErr(err)
	var row ContentsStock
	var returnData []ContentsStock
	count := 0
	for rows.Next() {
		count++
		rows.Scan(&row.StockName, &row.CostVolume, &row.CountVolume, &row.OneStock)
		returnData = append(returnData, row)
	}
	if count == 0 {
		log.Fatal(errors.New("ROW가 0개 입니다"))
	}

	return returnData
}

// GetTopDataDaily 일간 Top 종목 가져오기
func GetTopDataDaily(table, buyType, order, tmpDate string, db *sql.DB) ContentsStock {
	var sql, date string

	date = common.GetDateFormat("yyyy", "", "") + tmpDate

	sql = "select a.stockname, a." + buyType + ", b." + buyType + ", " +
		"case " +
		"when ABS(b." + buyType + ") > 0 " +
		"then (ABS(a." + buyType + ") * 1000000 / ABS(b." + buyType + "))" +
		"else 0 " +
		"end onestock " +
		"from daily.cost_craw a " +
		"left join daily.stock_craw b " +
		"on a.stockname = b.stockname and a.crawdate = b.crawdate " +
		"where a.crawdate = '" + date + "' and a.stocktype = " + table + " " +
		"order by a." + buyType + " " + order + " limit 1"

	rows, err := db.Query(sql)
	common.CheckErr(err)

	var row ContentsStock

	count := 0
	for rows.Next() {
		count++
		rows.Scan(&row.StockName, &row.CostVolume, &row.CountVolume, &row.OneStock)
	}
	if count == 0 {
		log.Fatal(errors.New("ROW가 0개 입니다"))
	}

	return row
}

// GetTopDataMonthly 월간 Top 종목 가져오기
func GetTopDataMonthly(table, buyType, order, tmpDate string, db *sql.DB) ContentsStock {
	var sql, date string

	date = common.GetDateFormat("yyyy", "", "") + tmpDate

	sql = "select a.stockname, a." + buyType + ", b." + buyType + ", " +
		"case " +
		"when ABS(b." + buyType + ") > 0 " +
		"then (ABS(a." + buyType + ") * 1000000 / ABS(b." + buyType + "))" +
		"else 0 " +
		"end onestock " +
		"from monthly.cost_craw a " +
		"left join monthly.stock_craw b " +
		"on a.stockname = b.stockname and a.crawdate = b.crawdate " +
		"where a.crawdate = '" + date + "' and a.stocktype = " + table + " " +
		"order by a." + buyType + " " + order + " limit 1"

	rows, err := db.Query(sql)
	common.CheckErr(err)

	var row ContentsStock

	count := 0
	for rows.Next() {
		count++
		rows.Scan(&row.StockName, &row.CostVolume, &row.CountVolume, &row.OneStock)
	}
	if count == 0 {
		log.Fatal(errors.New("ROW가 0개 입니다"))
	}

	return row
}

//GetContinueBuy 연속 매수한 종목 가져오기 (예정)
func GetContinueBuy(market, buyType, order string, db *sql.DB) {
	// var sql string

	// sql = "select stockname, " + buyType + " from daily.buy_continue_day where ABS(" + buyType + ") > 3 order by " + buyType
	// rows, err := db.Query(sql)
	// common.CheckErr(err)

	// for rows.Next() {
	// 	rows.Scan()
	// }

}
