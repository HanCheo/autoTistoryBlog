package insertstockdata

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/hancheo/tistory/common"
	"github.com/hancheo/tistory/conn"
)

type stockData struct {
	crawDate      string // 날짜
	stockName     string // 종목명
	stockType     string // 0 : kospi, 1 kosdaq
	personal      string // 개인
	foreigner     string // 외국인
	institutional string // 기관계
	financial     string // 금융투자
	insurance     string // 보험
	investment    string // 투신
	bank          string // 은행
	otherFinance  string // 기타금융
	pension       string // 연기금등
	privateFund   string // 사모펀드
	corporations  string // 기타법인
}

type gStockData struct {
	stockname   string
	endcost     string
	yesterday   string
	fluctuate   string
	transation  string
	costvolume  string
	gtransation string
	gcostvolume string
}

type nameCode struct {
	stockname string
	code      string
}

//ExcelInsertToDB stock info insert to db
func ExcelInsertToDB(dateType string) {
	db, err := conn.DbConn()
	common.CheckErr(err)
	defer db.Close()
	var date string
	if dateType == "daily" {
		date = common.GetDateFormat("", "mm", "dd")
	} else {
		date = common.GetDateFormat("yy", "mm", "")
	}
	var str []string = []string{date + "kospicost", date + "kosdaqcost", date + "kospistock", date + "kosdaqstock"}

	for _, s := range str {
		insertStockCostData(s, dateType, db)
	}
}

//종목 코드 및 이름 입력
func CodeNameExcel() {
	var insertDatas []nameCode
	file, err := excelize.OpenFile("./files/excel/code_name.xlsm")
	common.CheckErr(err)

	fmt.Println("code_name 엑셀 가져오기 완료")

	rows, err := file.GetRows(file.GetSheetName(0))
	// var codeArr []string
	for i, row := range rows {
		if i == 0 || row[0] == "" || row[1] == "" || row[0] == "'043500" {
			continue
		}
		stockcode := strings.ReplaceAll(row[0], "'", "")
		exist, _ := common.InArray(nameCode{code: stockcode, stockname: row[1]}, insertDatas)

		if exist {
			fmt.Println(stockcode, row[1])
			continue
		}
		insertDatas = append(insertDatas, nameCode{code: stockcode, stockname: row[1]})
	}
	fmt.Println("배열 입력 완료")
	db, err := conn.DbConn()
	common.CheckErr(err)

	for _, data := range insertDatas {
		sql := "insert into daily.code_name(code, stockname) values "
		sql += "('" + data.code + "', upper('" + data.stockname + "')),"
		sql = sql[:len(sql)-1]
		sql += " on conflict (code) Do update set stockname = upper('" + data.stockname + "');"
		_, err = db.Exec(sql)
		if err != nil {
			fmt.Println(data.stockname, err)
			// common.CheckErr(err)
		}
	}
	fmt.Printf("%d 건 데이터 입력완료", len(insertDatas))
}

//기관별, 외국인 매수/매도정보
func getExcel(dateType, stockType string) []stockData {
	var stockDatas []stockData
	var crawdate, st string
	crawdate = common.GetDateFormat("yyyy", "mm", "dd")
	// crawdate = "2021" + stockType[:4]
	file, err := excelize.OpenFile("./files/excel/" + dateType + "/" + stockType + ".xlsm")
	fmt.Println(err)
	fmt.Println(dateType + " " + stockType + "엑셀 가져오기 완료")

	if strings.Contains(stockType, "kosdaq") {
		st = "1"
	} else {
		st = "0"
	}

	rows, err := file.GetRows(file.GetSheetName(0))
	for i, row := range rows {
		if i == 0 {
			continue
		}
		tmpCode := stockData{
			crawDate:      crawdate,
			stockName:     row[0],
			stockType:     st,
			personal:      row[1],
			foreigner:     row[2],
			institutional: row[3],
			financial:     row[4],
			insurance:     row[5],
			investment:    row[6],
			bank:          row[7],
			otherFinance:  row[8],
			pension:       row[9],
			privateFund:   row[10],
			corporations:  row[11],
		}
		stockDatas = append(stockDatas, tmpCode)
	}
	return stockDatas
}

//기관별, 외국인 매수/매도정보 입력
func insertStockCostData(stockType, dateType string, db *sql.DB) {
	var sql, table string
	var err error
	var count int
	datas := getExcel(dateType, stockType)

	db.Exec("set search_path='" + dateType + "'")

	if stockType[len(stockType)-4:] == "cost" {
		table = "cost"
	} else {
		table = "stock"
	}
	sql = ("select count(*) from " + table + "_craw where stocktype = " + datas[0].stockType + " and crawdate = '" + datas[0].crawDate + "';")
	row := db.QueryRow(sql)
	err = row.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}

	if count > 0 {
		res, err := db.Exec("Delete from " + table + "_craw where crawdate = '" + datas[0].crawDate + "';")
		fmt.Println(datas[0].crawDate+"날짜 데이터가 존재합니다. row를 삭제후 재입력 합니다. \r\n삭제수 :", res)
		if err != nil {
			fmt.Println(err)
		}
	}

	sql = "insert into " + table + "_craw (crawdate, stockname, stocktype,personal, foreigner, institutional, financial, insurance, investment, bank, otherfinance, pension, privatefund, corporations) values "

	for _, data := range datas {
		sql += "("
		sql += "'" + data.crawDate + "',"
		sql += "upper('" + data.stockName + "'),"
		sql += "'" + data.stockType + "',"
		sql += data.personal + ", "
		sql += data.foreigner + ", "
		sql += data.institutional + ", "
		sql += data.financial + ", "
		sql += data.insurance + ", "
		sql += data.investment + ", "
		sql += data.bank + ", "
		sql += data.otherFinance + ", "
		sql += data.pension + ", "
		sql += data.privateFund + ", "
		sql += data.corporations
		sql += "),"
	}
	sql = sql[:len(sql)-1]
	sql += ";"
	res, err := db.Exec(sql)
	fmt.Println(res)
	if err != nil {
		fmt.Println(err)
	}
}
