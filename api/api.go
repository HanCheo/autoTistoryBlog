package api

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/hancheo/stockcode"
	"github.com/hancheo/tistory/common"
	"github.com/labstack/echo"
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

	if modifyDate == "" {
		date = common.GetDateFormat("yyyy", "mm", "dd")
	} else {
		date = "2021" + modifyDate
	}
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
		log.Fatal(errors.New("ROW가 0개 입니다"))
	}

	return returnData
}

// GetStockDataDaily 찾을 테이블, 매수 주체, 정렬 방식 (desc == 순매수 상위 "" == 순매도 상위)
func GetStockDataMonthly(table, buyType, order, modifyDate string, db *sql.DB) []ContentsStock {
	var sql, date string

	if modifyDate == "" {
		date = common.GetDateFormat("yyyy", "mm", "dd")
	} else {
		date = "2021" + modifyDate
	}
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

	if tmpDate == "" {
		date = common.GetDateFormat("yyyy", "mm", "dd")
	} else {
		date = "2021" + tmpDate
	}
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

	if tmpDate == "" {
		date = common.GetDateFormat("yyyy", "mm", "dd")
	} else {
		date = "2021" + tmpDate
	}
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
func GetContinueBuy(market, buyType, order string, db *sql.DB) {}

// DataToHTML 데이터를 html 테이블 형태로 반환
func DataToHTML(datas []ContentsStock) string {
	var str, stockName string
	var num1, num2, num3 int64
	var href string
	total := len(datas)

	str += "<div style='font-size: 13px; display:flex; justify-content: space-around; flex-direction: row; flex-wrap: wrap;'>"
	str += "<div style='width: 100%; font-size: 10px; text-align: left;'> (단위 = 금액: 백만원, 수량: 1주, 평균단가: 1원) </div>" // 테이블 주석
	str += "<table class='stockTable'>"                                                                            // 테이블 시작
	str += "<tr><th style='width:40px'>순위</th><th>종목명</th><th>금액</th><th>수량</th><th>평균단가</th></tr>"                //테이블 첫줄
	// 두번째 테이블이 비어있는지 확인용
	tableTrigger := false

	// 1위부터 각 데이터 행으로 추가
	for i := 0; i < total/2; i++ {
		stockName = datas[i].StockName
		num1, _ = strconv.ParseInt(datas[i].CostVolume, 10, 64)
		num2, _ = strconv.ParseInt(datas[i].CountVolume, 10, 64)
		num3, _ = strconv.ParseInt(datas[i].OneStock, 10, 64)
		// if datas[i][0] == "" {
		href = "https://search.naver.com/search.naver?where=nexearch&sm=top_hty&fbm=0&ie=utf8&query=" + stockName + "#_cs_root"
		// } else {
		// 	href = "https://finance.naver.com/item/main.nhn?code=" + datas[i][0]
		// }
		if num1 == 0 && num2 == 0 && num3 == 0 {
			str += "<tr class='noneDisplay'>"
			tableTrigger = true
		} else {
			str += "<tr>"
		}
		str += "<td>" + strconv.Itoa(i+1) + "</td><td><a href='" + href + "' target='_blank'>" + stockName + "</a></td><td>" + common.Comma(num1) + "</td><td>" + common.Comma(num2) + "</td><td>" + common.Comma(num3) + "</td>"
		str += "</tr>"
	}
	str += "</table>" //테이블 닫기

	num1, _ = strconv.ParseInt(datas[total/2].CostVolume, 10, 64)
	num2, _ = strconv.ParseInt(datas[total/2].CountVolume, 10, 64)
	num3, _ = strconv.ParseInt(datas[total/2].OneStock, 10, 64)

	if tableTrigger || (num1 == 0 && num2 == 0 && num3 == 0) {
		str += "</div>"
		return str
	}

	str += "<table class='stockTable'>"
	str += "<tr class='secondTableHead'><th style='width:40px'>순위</th><th>종목명</th><th>금액</th><th>수량</th><th>평균단가</th></tr>" //테이블 첫줄

	// 26, 50위부터 각 데이터 행으로 추가
	for i := total / 2; i < total; i++ {
		stockName = datas[i].StockName
		num1, _ = strconv.ParseInt(datas[i].CostVolume, 10, 64)
		num2, _ = strconv.ParseInt(datas[i].CountVolume, 10, 64)
		num3, _ = strconv.ParseInt(datas[i].OneStock, 10, 64)
		// if datas[i][0] == "" {
		href = "https://search.naver.com/search.naver?where=nexearch&sm=top_hty&fbm=0&ie=utf8&query=" + stockName + "#_cs_root"
		// } else {
		// 	// href = "https://finance.naver.com/item/main.nhn?code=" + datas[i][0]
		// }
		if num1 == 0 && num2 == 0 && num3 == 0 {
			str += "<tr class='noneDisplay'>"
		} else {
			str += "<tr>"
		}
		str += "<td>" + strconv.Itoa(i+1) + "</td><td><a href='" + href + "' target='_blank'>" + stockName + "</a></td><td>" + common.Comma(num1) + "</td><td>" + common.Comma(num2) + "</td><td>" + common.Comma(num3) + "</td>"
		str += "</tr>"
	}
	str += "</table></div>" //테이블 닫기

	return str
}

//InsertCode 주식 종목 코드 및 이름 데이터 입력
func InsertCode(c echo.Context) error {
	stockcode.InsertStockCode()
	return nil
}
