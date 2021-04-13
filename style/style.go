package style

import (
	"github.com/hancheo/tistory/common"
)

//BasicTheme 글 기본 테마 설정
func BasicTheme(contents, modifyDate, tp string) string {
	if tp == "monthly" {
		return style() + headerGuideTextMonthly(modifyDate) + contents + footerGuideText()
	} else {
		return style() + headerGuideText(modifyDate) + contents + footerGuideText()
	}

}

//Style 글 기본 스타일
func style() string {
	str := "<style>"
	str += ".toggleBtn {"
	str += "  border: 1px solid #f5aca6;"
	str += "  border-radius: 10px;"
	str += "  background-color: #fff;"
	str += "  color: #555;"
	str += "  cursor: pointer;"
	str += "  padding: 5px 10px;"
	str += "  -webkit-box-shadow: 0 0 10px 2px rgba(242, 191, 191, 1);"
	str += "  -moz-box-shadow: 0 0 10px 2px rgba(242, 191, 191, 1);"
	str += "  box-shadow: 0 0 10px 2px rgba(242, 191, 191, 1);"
	str += "}"
	str += "@media screen and (max-width: 1060px) {"
	str += "  .toggleBtn {"
	str += "    display: none;"
	str += "  }"
	str += "}"
	str += "@media screen and (min-width: 1061px) {"
	str += "  .toggleBtn {"
	str += "    display: block;"
	str += "  }"
	str += "  .area-aside {"
	str += "    position: relative;"
	str += "    right: 0px;"
	str += "    transition: right 1s;"
	str += "  }"
	str += "  .area-main {"
	str += "    transition: all 1s;"
	str += "  }"
	str += "}"
	str += "main.main {"
	str += "  overflow: hidden;"
	str += "}"
	str += "a {"
	str += "  text-decoration: none !important;"
	str += "}"
	str += ".stockTable tr,"
	str += ".stockTable tr {"
	str += "  border-top: 1px solid #dadce0;"
	str += "}"
	str += ".stockTable tr th,"
	str += ".stockTable tr td {"
	str += "  padding: 7px;"
	str += "  border-left: 1px solid #dadce0;"
	str += "}"
	str += ".stockTable {"
	str += "  max-width: 100% !important;"
	str += "  border: 1px solid #dadce0;"
	str += "  border-collapse: collapse;"
	str += "  flex: auto;"
	str += "}"
	str += "@media screen and (max-width: 649px) {"
	str += "  .secondTableHead {"
	str += "    display: none;"
	str += "  }"
	str += "  .table_wrap {"
	str += "    width: 100%;"
	str += "  }"
	str += "  .stockTable {"
	str += "    max-width: 100%;"
	str += "    width: 100%;"
	str += "  }"
	str += "  .stockTable:nth-child(odd) {"
	str += "    border-top: none;"
	str += "  }"
	str += "  .table_wrap:nth-child(odd) .stockTable {"
	str += "    border-top: none;"
	str += "  }"
	str += "}"
	str += "@media screen and (min-width: 650px) {"
	str += "  .stockTable {"
	str += "    width: 48%;"
	str += "  }"
	str += "  .table_wrap {"
	str += "    width: 48%;"
	str += "  }"
	str += "  .table_wrap .stockTable {"
	str += "    width: 100%;"
	str += "  }"
	str += "}"
	str += ".stockTable tr th:first-child,"
	str += ".stockTable tr td:first-child {"
	str += "  width: 40px;"
	str += "}"
	str += "</style>"

	return str
}

//HeaderGuideText 머리글
func headerGuideText(modifyDate string) string {
	var date string
	if modifyDate == "" {
		date = common.GetDateFormat("", "mm", "dd")
	} else {
		date = modifyDate
	}

	str := "<div style='text-align:center'> <img src='https://tistory3.daumcdn.net/tistory/3508113/skin/images/unnamed.jpg' /></div>"
	str += "<br><blockquote style='font-size: 13px; padding: 21px 25px !important' data-ke-style='style3'>" +
		"<h2 style='display: flex; margin-top:0; padding-top:0'>참고사항 <span style='font-size: 10px'>" + date + " 일자</span>" +
		"<button class='toggleBtn' style='margin-left:auto; font-size:11px' onclick='toggleCategory()'>확대/축소</button></h2>" +
		"<strong style='color: #ff1d1d'>이글은 주식장이 열리는날 7시 이전으로 자동생성 되는 글입니다. 구독해주시면 알림이 가요 !</strong><br>" +
		"1. 각 종목명을 클릭시 자동으로 네이버 검색으로 페이지가 넘어가져요 !<br>" +
		"2. 평균단가는 각 금액에서 100만원을 곱하고 수량으로 나눈것으로 대략적인 수치에요.<br>" +
		"3. 코스피-코스닥 순으로 나열해 되어있어요. <br>" +
		"4. 각 매수주체를 클릭하시면 내용이 변경되요 !<br> " +
		"5. 당일 거래종목수가 50종목이 안되는 경우도 있습니다 ! " +
		"</blockquote>" +
		"<div style='text-align: center; color:#ff1d1d; font-size: 10px;'># 본 글은 광고를 포함하고 있습니다. 광고 클릭에서 발생하는 수익금은 모두 블로그의 유지 및 관리, 그리고 콘텐츠 향상을 위해 쓰여집니다.</div>"
	return str
}

func headerGuideTextMonthly(modifyDate string) string {
	var date string
	if modifyDate == "" {
		date = common.GetDateFormat("", "mm", "dd")
	} else {
		date = modifyDate
	}

	str := "<div style='text-align:center'> <img src='https://tistory3.daumcdn.net/tistory/3508113/skin/images/unnamed.jpg' /></div>"
	str += "<br><blockquote style='font-size: 13px; padding: 21px 25px !important' data-ke-style='style3'>" +
		"<h2 style='display: flex; margin-top:0; padding-top:0'>참고사항 <span style='font-size: 10px'>" + date + " 일자</span>" +
		"<button class='toggleBtn' style='margin-left:auto; font-size:11px' onclick='toggleCategory()'>확대/축소</button></h2>" +
		"<strong style='color: #ff1d1d'>이글은 매월 말일 7시 이전으로 자동생성 되는 글입니다. 구독해주시면 알림이 가요 !</strong><br>" +
		"1. 각 종목명을 클릭시 자동으로 네이버 검색으로 페이지가 넘어가져요 !<br>" +
		"2. 평균단가는 각 금액에서 100만원을 곱하고 수량으로 나눈것으로 대략적인 수치에요.<br>" +
		"3. 코스피-코스닥 순으로 나열해 되어있어요. <br>" +
		"4. 각 매수주체를 클릭하시면 내용이 변경되요 !<br> " +
		"5. 월간 거래종목수가 100종목이 안되는 경우도 있습니다 ! " +
		"</blockquote>" +
		"<div style='text-align: center; color:#ff1d1d; font-size: 10px;'># 본 글은 광고를 포함하고 있습니다. 광고 클릭에서 발생하는 수익금은 모두 블로그의 유지 및 관리, 그리고 콘텐츠 향상을 위해 쓰여집니다.</div>"
	return str
}

//FooterGuideText 꼬리글
func footerGuideText() string {
	str := "<h3 style='color:#1f1f1f; text-align:center;'>투자에 좋은 도움이 되었으면 좋겠습니다 :)</h3>"
	return str
}

// NavUl nav
func NavUl() string {
	str := " <ul class='nav' id='myTab'>"
	str += " <li class='nav-item active'>"
	str += " <a onclick='changetab(0)' class='nav-link' id='personal-tab'>개인</a>"
	str += " </li>"
	str += " <li class='nav-item'>"
	str += " <a onclick='changetab(1)' class='nav-link' id='foreigner-tab'>외국인</a>"
	str += " </li>"
	str += " <li class='nav-item'>"
	str += " <a onclick='changetab(2)' class='nav-link' id='institutional-tab'>기관종합</a>"
	str += " </li>"
	str += " <li class='nav-item'>"
	str += " <a onclick='changetab(3)' class='nav-link' id='financial-tab'>금융투자</a>"
	str += " </li>"
	str += " <li class='nav-item'>"
	str += " <a onclick='changetab(4)' class='nav-link' id='insurance-tab'>보험</a>"
	str += " </li>"
	str += " <li class='nav-item'>"
	str += " <a onclick='changetab(5)' class='nav-link' id='investment-tab'>투신</a>"
	str += " </li>"
	str += " <li class='nav-item'>"
	str += " <a onclick='changetab(6)' class='nav-link' id='bank-tab'>은행</a>"
	str += " </li>"
	str += " <li class='nav-item'>"
	str += " <a onclick='changetab(7)' class='nav-link' id='otherFinance-tab'>기타금융</a>"
	str += " </li>"
	str += " <li class='nav-item'>"
	str += " <a onclick='changetab(8)' class='nav-link' id='pension-tab'>연기금등</a>"
	str += " </li>"
	str += " <li class='nav-item'>"
	str += " <a onclick='changetab(9)' class='nav-link' id='privateFund-tab'>사모펀드</a>"
	str += " </li>"
	str += " <li class='nav-item'>"
	str += " <a onclick='changetab(10)' class='nav-link' id='corporations-tab'>기타법인</a>"
	str += " </li>"
	str += "     </ul>"
	return str
}
