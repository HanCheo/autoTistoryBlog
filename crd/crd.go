package crd

import (
	"context"
	"fmt"
	_ "log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/hancheo/tistory/common"
	"github.com/joho/godotenv"
)

func stringAddSlash(s string) string {
	s = strings.ReplaceAll(s, `'`, `\'`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

//WriteStockDiscussion 검색할 종목명, 글 제목, 글 내용 : 종토방 글 자동으로 작성하기
func WriteStockDiscussion(stockType int, stockName, contentTitle, contentBody string) {
	godotenv.Load("../.env")

	var naver_ID, naver_PW string

	if stockType < 2 {
		naver_ID = os.Getenv("NAVER_ID_KOSPI")
		naver_PW = os.Getenv("NAVER_PW_KOSPI")
	} else {
		naver_ID = os.Getenv("NAVER_ID_KOSDAQ")
		naver_PW = os.Getenv("NAVER_PW_KOSDAQ")
	}
	contentTitle = stringAddSlash(contentTitle)
	contentBody = stringAddSlash(contentBody)

	// chrome 실행 옵션 설정
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false), //headless를 false로 하면 브라우저가 뜨고, true로 하면 브라우저가 뜨지않는 headless 모드로 실행됨. 기본값은 true.
	)

	contextVar, cancelFunc := chromedp.NewExecAllocator(context.Background(), opts...)
	contextVar, cancelFunc = chromedp.NewContext(contextVar)
	defer cancelFunc()

	// var strVar string
	ch := chromedp.WaitNewTarget(contextVar, func(info *target.Info) bool {
		return info.URL != ""
	})

	// 종목 검색 페이지 -> 종목 토론실 이동
	err := chromedp.Run(contextVar,
		chromedp.Navigate(`https://search.naver.com/search.naver?where=nexearch&sm=top_hty&fbm=1&ie=utf8&query=`+stockName), // 시작 URL
		chromedp.WaitVisible(`ul.lst2 li.ed a`),
		chromedp.Click(`ul.lst2 li.ed a`, chromedp.BySearch),
	)
	if err != nil {
		panic(err)
	}

	newCtx, cancel := chromedp.NewContext(contextVar, chromedp.WithTargetID(<-ch))
	defer cancel()

	// 게시판 -> 글쓰기 클릭 -> 로그인
	err = chromedp.Run(newCtx,
		chromedp.WaitVisible(`table[summary='게시판 옵션'] tbody > tr a[href^='/item/board_write_edit.nhn?']`),
		chromedp.Click(`table[summary='게시판 옵션'] tbody > tr a[href^='/item/board_write_edit.nhn?']`),

		chromedp.SetValue(`#id`, naver_ID, chromedp.ByID),
		chromedp.SetValue(`#pw`, naver_PW, chromedp.ByID),
		chromedp.WaitVisible(`input[value='로그인']`),
		chromedp.Click(`input[value='로그인']`, chromedp.BySearch),
	)
	if err != nil {
		common.CheckErr(err)
	}

	var res interface{}
	// 글작성 -> 완료버튼 클릭
	var iframes []*cdp.Node

	err = chromedp.Run(newCtx,
		chromedp.WaitReady("#write", chromedp.ByID),
		chromedp.Sleep(2*time.Second),
		chromedp.Nodes(`iframe[title="글쓰기 영역"]`, &iframes, chromedp.ByQuery),
	)

	if err != nil {
		fmt.Println(err)
	}

	err = chromedp.Run(newCtx,
		chromedp.WaitReady(`input#title`, chromedp.ByQuery, chromedp.FromNode(iframes[0])),
		chromedp.Evaluate(`document.querySelector('iframe[title="글쓰기 영역"]').contentWindow.document.body.querySelector('input#title').value="`+contentTitle+`"`, &res),
		chromedp.Evaluate(`document.querySelector('iframe[title="글쓰기 영역"]').contentWindow.document.body.querySelector('textarea#body').value="`+contentBody+`"`, &res),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`img#submit_btn`, chromedp.FromNode(iframes[0])),
		chromedp.Sleep(60*time.Second),
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(contentTitle)
	// fmt.Println(res)
}

func RunLocalHost() {
	godotenv.Load("../.env")

	url := "https://www.tistory.com/oauth/authorize?client_id=" + os.Getenv("CLIENT_ID") + "&redirect_uri=" + os.Getenv("REDIRECT_URL") + "&response_type=code"

	// chrome 실행 옵션 설정
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false), //headless를 false로 하면 브라우저가 뜨고, true로 하면 브라우저가 뜨지않는 headless 모드로 실행됨. 기본값은 true.
	)

	contextVar, cancelFunc := chromedp.NewExecAllocator(context.Background(), opts...)
	contextVar, cancelFunc = chromedp.NewContext(contextVar)
	defer cancelFunc()

	// 종목 검색 페이지 -> 종목 토론실 이동
	err := chromedp.Run(contextVar,
		chromedp.Navigate(url), // 시작 URL
		chromedp.WaitVisible(`button.confirm`, chromedp.ByQuery),
		chromedp.Click(`button.confirm`, chromedp.ByQuery),
	)
	if err != nil {
		panic(err)
	}
}

// document.querySelector("").click()
// document.querySelector("table[summary='게시판 옵션'] tbody > tr img[alt='글쓰기']").click()
