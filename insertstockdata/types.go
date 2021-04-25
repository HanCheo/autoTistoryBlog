package insertstockdata

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
