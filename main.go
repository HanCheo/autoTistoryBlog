package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	xj "github.com/basgys/goxml2json"

	"github.com/hancheo/tistory/api"
	"github.com/hancheo/tistory/common"
	"github.com/hancheo/tistory/conn"
	"github.com/hancheo/tistory/crd"
	"github.com/hancheo/tistory/insertstockdata"
	"github.com/hancheo/tistory/style"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

type stockData struct {
	stockName     string // 종목명
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

type writeData struct {
	topData api.ContentsStock
	buyType string
}

type newPost struct {
	Tistory tistory `json:",omitempty"`
}
type tistory struct {
	Status string `json:",omitempty"`
	PostId string `json:",omitempty"`
	Url    string `json:",omitempty"`
}

func main() {
	godotenv.Load(".env")
	e := echo.New()
	e.GET("/", authLogin)
	e.GET("/api", newContent)
	e.GET("/modify", modify)
	e.GET("/stockDiscussion", stockDiscussionApi)
	e.GET("/api2", newContentMonthly)

	e.Logger.Fatal(e.Start(":1323"))

}

// "/" 로 접근시 자동으로 티스토리 인증센터로 redierect
func authLogin(c echo.Context) error {
	common.Openbrowser("https://www.tistory.com/oauth/authorize?client_id=" + os.Getenv("CLIENT_ID") + "&redirect_uri=" + os.Getenv("REDIRECT_URL") + "&response_type=code")
	return nil
}

// 3일이상 연속 순매수 종목 작성
func continueBuyDaily(c echo.Context) error {
	var market, buyType, order string
	code := api.GetAccessToken(c.QueryParam("code"))
	db, err := conn.DbConn()
	common.CheckErr(err)

	api.GetContinueBuy(market, buyType, order, db)

	tistoryContents := "https://www.tistory.com/apis/post/write"

	res, err := http.PostForm(tistoryContents, url.Values{
		"access_token": {code},
		"postId":       {"189"},
		"blogName":     {"hanch"},
		"title":        {"개인,외국인,기관별 3일 이상 연속 순매수 종목"},
		"content":      {"종목"},
		"category":     {"465733"},
		"visibility":   {"0"},
		"tag":          {"연속매수, 기관별 주식, 주식매매현황, 기관주식, 외국인주식, 개인주식, 기관매수종목, 기관매도종목, 개인, 외국인, 기관종합, 금융투자, 보험, 투신, 은행, 기타금융, 연기금등, 사모펀드, 기타법인"},
		"password":     {"12345"},
	})
	fmt.Println("글 작성 완료")

	//get url in response
	var e newPost
	body, _ := ioutil.ReadAll(res.Body)
	xml := strings.NewReader(string(body))
	js, _ := xj.Convert(xml)
	json.Unmarshal([]byte(js.String()), &e)

	//Url
	fmt.Println(e.Tistory.Url)
	stockDiscussion(e.Tistory.Url)

	return nil
}

//코스피 코스닥 순서로 진행
//글 수정
func modify(c echo.Context) error {
	insertstockdata.CsvInsertToDB("daily")
	code := api.GetAccessToken(c.QueryParam("code"))
	db, err := conn.DbConn()
	common.CheckErr(err)

	tistoryContents := "https://www.tistory.com/apis/post/modify"
	var m stockData
	var keys []string
	keys = common.GetStructKeys(m)

	str2 := []string{"코스피", "코스닥"}
	str3 := []string{"개인", "외국인", "기관종합", "금융투자", "보험", "투신", "은행", "기타금융", "연기금등", "사모펀드", "기타법인"}

	month := common.GetDateFormat("", "m", "")
	// monthDay := common.GetDateFormat("", "mm", "dd")
	monthDay := "0407"

	var tableHTML string = ""
	tableHTML += style.NavUl()
	for idx, k := range keys {
		if idx == 0 {
			continue
		}
		if k == "personal" {
			tableHTML += "<div class='tab active'>"
		} else {
			tableHTML += "<div class='tab'>"
		}
		for i, s := range str2 {
			data := api.GetStockDataDaily(strconv.Itoa(i), k, "desc", monthDay, db)
			data2 := api.GetStockDataDaily(strconv.Itoa(i), k, "", monthDay, db)
			fmt.Println(s + " " + str3[idx-1] + " 상위 매수/매도 50 종목 가져오기 완료")
			tableHTML += "<h3>" + str3[idx-1] + "_" + s + "<h3>"
			tableHTML += "<h1 style='font-weight: bold; text-aling:center; color: #ef5369;'>순매수 종목</h1>"
			tableHTML += style.DataToHTML(data)
			tableHTML += "<h1 style='font-weight: bold; text-aling:center; color: #006dd7'>순매도 종목</h1>"
			tableHTML += style.DataToHTML(data2)

		}
		tableHTML += "</div>"
	}

	_, err1 := http.PostForm(tistoryContents, url.Values{
		"access_token": {code},
		"postId":       {"203"},
		"blogName":     {"hanch"},
		"title":        {"2021년 " + monthDay[:2] + "월 " + monthDay[2:] + "일 개인,외국인,기관별 순매수/매도 상위 50 종목"},
		"content":      {style.BasicTheme(tableHTML, "", "daily")},
		"category":     {"465734"}, // monthly 465733
		"visibility":   {"3"},
		"tag":          {month + "월 주식, 기관별 주식, 주식매매현황, 기관주식, 외국인주식, 개인주식, 기관매수종목, 기관매도종목, 개인, 외국인, 기관종합, 금융투자, 보험, 투신, 은행, 기타금융, 연기금등, 사모펀드, 기타법인"},
	})
	common.CheckErr(err1)

	return c.String(http.StatusOK, code)
}

//일간 글 작성 (매수 상위 50종목)
func newContent(c echo.Context) error {
	// dkosdaqcost, dkosdaqstock, dkospicost, dkospistock
	insertstockdata.CodeNameCsv()
	insertstockdata.CsvInsertToDB("daily")
	fmt.Println("일간 데이터 입력 완료")
	code := api.GetAccessToken(c.QueryParam("code"))
	db, err := conn.DbConn()
	common.CheckErr(err)

	tistoryContents := "https://www.tistory.com/apis/post/write"
	var m stockData
	var keys []string
	keys = common.GetStructKeys(m)

	str2 := []string{"코스피", "코스닥"}
	str3 := []string{"개인", "외국인", "기관종합", "금융투자", "보험", "투신", "은행", "기타금융", "연기금등", "사모펀드", "기타법인"}

	month := common.GetDateFormat("", "m", "")
	monthDay := common.GetDateFormat("", "mm", "dd")

	var tableHTML string = ""
	tableHTML += style.NavUl()
	for idx, k := range keys {
		if idx == 0 {
			continue
		}
		if k == "personal" {
			tableHTML += "<div class='tab active'>"
		} else {
			tableHTML += "<div class='tab'>"
		}
		for i, s := range str2 {
			data := api.GetStockDataDaily(strconv.Itoa(i), k, "desc", "", db)
			data2 := api.GetStockDataDaily(strconv.Itoa(i), k, "", "", db)
			fmt.Println(s + " " + str3[idx-1] + " 상위 매수/매도 50 종목 가져오기 완료")
			tableHTML += "<h3>" + str3[idx-1] + "_" + s + "<h3>"
			tableHTML += "<h1 style='font-weight: bold; text-aling:center; color: #ef5369;'>순매수 종목</h1>"
			tableHTML += style.DataToHTML(data)
			tableHTML += "<h1 style='font-weight: bold; text-aling:center; color: #006dd7'>순매도 종목</h1>"
			tableHTML += style.DataToHTML(data2)
		}
		tableHTML += "</div>"
	}

	resp2, err1 := http.PostForm(tistoryContents, url.Values{
		"access_token": {code},
		"blogName":     {"hanch"},
		"title":        {"2021년 " + monthDay[:2] + "월 " + monthDay[2:] + "일 개인,외국인,기관별 순매수/매도 상위 50 종목"},
		"content":      {style.BasicTheme(tableHTML, "", "daily")},
		"category":     {"465734"}, // monthly 465733
		"visibility":   {"3"},
		"tag":          {month + "월 주식, 기관별 주식, 주식매매현황, 기관주식, 외국인주식, 개인주식, 기관매수종목, 기관매도종목, 개인, 외국인, 기관종합, 금융투자, 보험, 투신, 은행, 기타금융, 연기금등, 사모펀드, 기타법인"},
	})

	common.CheckErr(err1)
	defer resp2.Body.Close()

	fmt.Println("글 작성 완료")

	//get url in response
	var e newPost
	body, _ := ioutil.ReadAll(resp2.Body)
	xml := strings.NewReader(string(body))
	js, _ := xj.Convert(xml)
	json.Unmarshal([]byte(js.String()), &e)

	//Url
	fmt.Println(e.Tistory.Url)
	stockDiscussion(e.Tistory.Url)

	return c.String(http.StatusOK, code)
}

//월간 글 작성 (매수 상위 100종목)
func newContentMonthly(c echo.Context) error {
	// dkosdaqcost, dkosdaqstock, dkospicost, dkospistock
	insertstockdata.CsvInsertToDB("monthly")
	fmt.Println("월간 데이터 입력 완료")
	code := api.GetAccessToken(c.QueryParam("code"))
	db, err := conn.DbConn()
	common.CheckErr(err)

	tistoryContents := "https://www.tistory.com/apis/post/write"
	var m stockData
	var keys []string
	keys = common.GetStructKeys(m)

	str2 := []string{"코스피", "코스닥"}
	str3 := []string{"개인", "외국인", "기관종합", "금융투자", "보험", "투신", "은행", "기타금융", "연기금등", "사모펀드", "기타법인"}

	month := common.GetDateFormat("", "m", "")
	monthDay := common.GetDateFormat("", "mm", "dd")

	var tableHTML string = ""
	tableHTML += style.NavUl()
	for idx, k := range keys {
		if idx == 0 {
			continue
		}
		if k == "personal" {
			tableHTML += "<div class='tab active'>"
		} else {
			tableHTML += "<div class='tab'>"
		}
		for i, s := range str2 {
			data := api.GetStockDataMonthly(strconv.Itoa(i), k, "desc", "", db)
			data2 := api.GetStockDataMonthly(strconv.Itoa(i), k, "", "", db)
			fmt.Println(s + " " + str3[idx-1] + " 상위 매수/매도 100 종목 가져오기 완료")
			tableHTML += "<h3>" + str3[idx-1] + "_" + s + "<h3>"
			tableHTML += "<h1 style='font-weight: bold; text-aling:center; color: #ef5369;'>순매수 종목</h1>"
			tableHTML += style.DataToHTML(data)
			tableHTML += "<h1 style='font-weight: bold; text-aling:center; color: #006dd7'>순매도 종목</h1>"
			tableHTML += style.DataToHTML(data2)

		}
		tableHTML += "</div>"
	}

	resp2, err1 := http.PostForm(tistoryContents, url.Values{
		"access_token": {code},
		"blogName":     {"hanch"},
		"title":        {"2021년 " + monthDay[:2] + "월 개인,외국인,기관별 순매수/매도 상위 100 종목"},
		"content":      {style.BasicTheme(tableHTML, "", "monthly")},
		"category":     {"465733"},
		"visibility":   {"3"},
		"tag":          {month + "월 주식, 기관별 주식, 주식매매현황, 기관주식, 외국인주식, 개인주식, 기관매수종목, 기관매도종목, 개인, 외국인, 기관종합, 금융투자, 보험, 투신, 은행, 기타금융, 연기금등, 사모펀드, 기타법인"},
	})

	common.CheckErr(err1)
	defer resp2.Body.Close()

	fmt.Println("글 작성 완료")

	//get url in response
	var e newPost
	body, _ := ioutil.ReadAll(resp2.Body)
	xml := strings.NewReader(string(body))
	js, _ := xj.Convert(xml)
	json.Unmarshal([]byte(js.String()), &e)

	fmt.Println(e.Tistory.Url)
	// stockDiscussion(e.Tistory.Url)

	return c.String(http.StatusOK, code)
}

func stockDiscussionApi(c echo.Context) error {
	stockDiscussion("")
	return nil
}

//네이버 종목토론방 글 작성 (각 매수주체별 매매 상위 1등 종목에 한함.)
func stockDiscussion(url string) error {
	// stockName   string
	db, err := conn.DbConn()
	common.CheckErr(err)

	str2 := []string{"코스피", "코스닥"}
	str3 := []string{"개인", "외국인", "기관종합", "금융투자", "보험", "투신", "은행", "기타금융", "연기금", "사모펀드", "기타법인"}
	var keys []string
	var m stockData
	var tmp [4][]writeData

	keys = common.GetStructKeys(m)
	for idx := range str2 {
		for i, s := range keys {
			if i == 0 || idx == 0 {
				continue
			}
			var tmpType writeData
			tmpType.topData = api.GetTopDataDaily(strconv.Itoa(idx), s, "desc", "", db)
			tmpType.buyType = str3[i-1]

			if idx == 0 {
				tmp[0] = append(tmp[0], tmpType) // 코스피 s 매수 상위 1
				fmt.Println("Append 0 : " + str3[i-1] + " 코스피 매수")
				tmpType.topData = api.GetTopDataDaily(strconv.Itoa(idx), s, "", "", db)
				tmp[1] = append(tmp[1], tmpType) // 코스피 s 매도 상위 1
				fmt.Println("Append 1 : " + str3[i-1] + "코스피 매도")
			} else {
				tmp[2] = append(tmp[2], tmpType) // 코스닥 s 매수 상위 1
				fmt.Println("Append 2 : " + str3[i-1] + "코스닥 매수")
				tmpType.topData = api.GetTopDataDaily(strconv.Itoa(idx), s, "", "", db)
				tmp[3] = append(tmp[3], tmpType) // 코스닥 s 매도 상위 1
				fmt.Println("Append 3 : " + str3[i-1] + "코스닥 매도")
			}
		}
	}

	var title string
	var content string
	var sN string

	for i, s := range tmp {
		for j := 0; j < len(s); j++ {
			if s[j].buyType == "" {
				continue
			} else {

				sN = s[j].topData.StockName

				if i%2 == 0 {
					content = sN + ` 금일 종목 주체별 매수 정보 !`
				} else {
					content = sN + ` 금일 종목 주체별 매도 정보 !`
				}
				title = sN + " " + s[j].buyType

				tmp1, _ := strconv.ParseInt(s[j].topData.CostVolume, 10, 64)
				tmp2, _ := strconv.ParseInt(s[j].topData.CountVolume, 10, 64)
				tmp3, _ := strconv.ParseInt(s[j].topData.OneStock, 10, 64)

				content += `\n\n` + s[j].buyType + ` : \n   금액(100만 단위) : ` + common.Comma(tmp1) + `\n   수량(1 단위) : ` + common.Comma(tmp2) + `\n   평균단가 : ` + common.Comma(tmp3)
			}
			// 같은 종목이 있으면 내용 추가
			for k := j + 1; k < len(s); k++ {
				if s[k].topData.StockName == sN {
					title += `, ` + s[k].buyType

					tmp1, _ := strconv.ParseInt(s[k].topData.CostVolume, 10, 64)
					tmp2, _ := strconv.ParseInt(s[k].topData.CountVolume, 10, 64)
					tmp3, _ := strconv.ParseInt(s[k].topData.OneStock, 10, 64)

					content += `\n\n` + s[k].buyType + ` : \n   금액(100만 단위) : ` + common.Comma(tmp1) + `\n   수량(1 단위) : ` + common.Comma(tmp2) + `\n   평균단가 : ` + common.Comma(tmp3)

					s[k].buyType = ""
				}
			}
			if i < 2 {
				if i == 0 {
					title += " 코스피 종목중 순매수 1위 !!"
				} else {
					title += " 코스피 종목중 순매도 1위 !!"
				}
			} else {
				if i == 2 {
					title += " 코스닥 종목중 순매수 1위 !!"
				} else {
					title += " 코스닥 종목중 순매도 1위 !!"
				}
			}
			content += `\n\n금일 각 매수주체별 상위 50종목 매매 정보 블로그\n`
			content += `hanch.tistory.com`
			crd.WriteStockDiscussion(i, sN, title, content)
			fmt.Println(title)
		}
	}

	return nil
}

//네이버 종목토론방 월간 글 작성 (각 매수주체별 매매 상위 1등 종목에 한함.)
func stockDiscussionMonthly(url string) error {
	// stockName   string
	db, err := conn.DbConn()
	common.CheckErr(err)

	str2 := []string{"코스피", "코스닥"}
	str3 := []string{"개인", "외국인", "기관종합", "금융투자", "보험", "투신", "은행", "기타금융", "연기금", "사모펀드", "기타법인"}
	var keys []string
	var m stockData
	var tmp [4][]writeData

	keys = common.GetStructKeys(m)
	for idx := range str2 {
		for i, s := range keys {
			if i == 0 {
				continue
			}
			var tmpType writeData
			tmpType.topData = api.GetTopDataMonthly(strconv.Itoa(idx), s, "desc", "", db)
			tmpType.buyType = str3[i-1]

			if idx == 0 {
				tmp[0] = append(tmp[0], tmpType) // 코스피 s 매수 상위 1
				fmt.Println("Append 0 : 코스피 매수")
				tmpType.topData = api.GetTopDataMonthly(strconv.Itoa(idx), s, "", "", db)
				tmpType.buyType = str3[i-1]
				tmp[1] = append(tmp[1], tmpType) // 코스피 s 매도 상위 1
				fmt.Println("Append 1 : 코스피 매도")
			} else {
				tmp[2] = append(tmp[2], tmpType) // 코스닥 s 매수 상위 1
				fmt.Println("Append 2 : 코스닥 매수")
				tmpType.topData = api.GetTopDataMonthly(strconv.Itoa(idx), s, "", "", db)
				tmpType.buyType = str3[i-1]
				tmp[3] = append(tmp[3], tmpType) // 코스닥 s 매도 상위 1
				fmt.Println("Append 3 : 코스닥 매도")
			}
			fmt.Println(i)
		}
	}

	var title string
	var content string
	var sN string

	for i, s := range tmp {
		for j := 0; j < len(s); j++ {
			if s[j].buyType == "" || j < 2 {
				continue
			} else {

				sN = s[j].topData.StockName

				if i%2 == 0 {
					content = sN + ` 월간 종목 주체별 매수 정보 !`
				} else {
					content = sN + ` 월간 종목 주체별 매도 정보 !`
				}
				title = sN + " " + s[j].buyType

				tmp1, _ := strconv.ParseInt(s[j].topData.CostVolume, 10, 64)
				tmp2, _ := strconv.ParseInt(s[j].topData.CountVolume, 10, 64)
				tmp3, _ := strconv.ParseInt(s[j].topData.OneStock, 10, 64)

				content += `\n\n` + s[j].buyType + ` : \n   금액(100만 단위) : ` + common.Comma(tmp1) + `\n   수량(1 단위) : ` + common.Comma(tmp2) + `\n   평균단가 : ` + common.Comma(tmp3)
			}
			// 같은 종목이 있으면 내용 추가
			for k := j + 1; k < len(s); k++ {
				if s[k].topData.StockName == sN {
					title += `, ` + s[k].buyType

					tmp1, _ := strconv.ParseInt(s[k].topData.CostVolume, 10, 64)
					tmp2, _ := strconv.ParseInt(s[k].topData.CountVolume, 10, 64)
					tmp3, _ := strconv.ParseInt(s[k].topData.OneStock, 10, 64)

					content += `\n\n` + s[k].buyType + ` : \n   금액(100만 단위) : ` + common.Comma(tmp1) + `\n   수량(1 단위) : ` + common.Comma(tmp2) + `\n   평균단가 : ` + common.Comma(tmp3)

					s[k].buyType = ""
				}
			}
			if i < 2 {
				if i == 0 {
					title += " 월간 코스피 종목중 순매수 1위 !!"
				} else {
					title += " 월간 코스피 종목중 순매도 1위 !!"
				}
			} else {
				if i == 2 {
					title += " 월간 코스닥 종목중 순매수 1위 !!"
				} else {
					title += " 월간 코스닥 종목중 순매도 1위 !!"
				}
			}
			content += `\n\n월간 각 매수주체별 상위 100종목 매매 정보 블로그\n`
			content += url
			crd.WriteStockDiscussion(i, sN, title, content)
		}
	}

	return nil
}
