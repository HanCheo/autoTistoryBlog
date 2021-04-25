package insertstockdata

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/hancheo/tistory/common"
	"github.com/hancheo/tistory/conn"
)

//CsvInsertToDB stock info insert to db
func CsvInsertToDB(dateType string) {
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
		insertStockCostDataCsv(s, dateType, db)
	}
}

//종목 코드 및 이름 입력
func CodeNameCsv() {
	var insertDatas []nameCode
	file, err := os.Open("./files/csv/code_name.csv")
	fmt.Println(err)

	rdr := csv.NewReader(bufio.NewReader(file))
	// rdr.Comma = ';'
	// rdr.Comment = '#'
	rows, _ := rdr.ReadAll()
	fmt.Println("code_name CSV 가져오기 완료")

	for i, row := range rows {
		if i == 0 || row[0] == "" || row[1] == "" || row[0] == "'043500" {
			continue
		}
		row[1] = common.DecodeEUCKR(row[1])
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
		}
	}
	fmt.Printf("%d 건 데이터 입력완료", len(insertDatas))
}

//기관별, 외국인 매수/매도정보
func getCsv(dateType, fileName string) []stockData {
	var stockDatas []stockData
	var crawdate, st string
	crawdate = "2021" + fileName[:4]
	file, err := os.Open("./files/csv/" + dateType + "/" + fileName + ".csv")
	fmt.Println(err)
	fmt.Println(dateType + " " + fileName + "CSV 가져오기 완료")

	rdr := csv.NewReader(bufio.NewReader(file))
	// rdr.Comma = ';'
	// rdr.Comment = '#'
	rows, _ := rdr.ReadAll()
	if strings.Contains(fileName, "kosdaq") {
		st = "1"
	} else {
		st = "0"
	}
	for i, row := range rows {
		row[0] = common.DecodeEUCKR(row[0])

		for j, rowData := range row {
			if rowData == "" {
				row[j] = "0"
			}
		}

		if i <= 1 {
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
func insertStockCostDataCsv(fileName, dateType string, db *sql.DB) {
	var sql, table string
	var err error
	var count int
	datas := getCsv(dateType, fileName)

	db.Exec("set search_path='" + dateType + "'")

	if fileName[len(fileName)-4:] == "cost" {
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
