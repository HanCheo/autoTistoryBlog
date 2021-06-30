# Tistory Auto Wrait

결과 : https://hanch.tistory.com/

티스토리 블로그에 글을 자동으로 작성하는 프로젝트입니다.<br>
글 내용은 각 매수 주체별 당일 매매 상위 50종목을 가져와 테이블 형태로 보여주는 글입니다.<br>

글 작성후 크롬드라이버를 이용하여 블로그 글 홍보를 진행하였고 <br>
그 결과 한달도 채 안되서 카카오, 구글애드센스 승인을 받았습니다.

## Used

 ![Postgresql](https://img.shields.io/badge/-Postgresql-336791?logo=Postgresql)
 ![Go](https://img.shields.io/badge/-Go-00ADD8?logo=GO)
 ![CSS](https://img.shields.io/badge/-CSS-1572B6?logo=CSS3)

## Dependencies
```go
import {
//Using local sever
	"github.com/labstack/echo" 
//Using PostgreSQL
    "github.com/lib/pq"
//Chrome Driver
	"github.com/chromedp/cdproto/cdp" 
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
//xml Response to JSON Convert
	"github.com/basgys/goxml2json"
//Using env File
	"github.com/joho/godotenv"
//UTF-8-bom to utf-8 encoding
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
//Using for Excel file
	"github.com/360EntSecGroup-Skylar/excelize"
}
```


 ## Preview
 <pre style="text-align: center;">
  <img src = "https://user-images.githubusercontent.com/38929712/114538158-05b57800-9c8e-11eb-8937-dd17db5fa11b.gif" style="width: 80%; min-width: 300px; max-width: 500px">
 </pre>


## Write Process
    1) HTS CSV 데이터 추출
    2) 사명이 변경된 종목 DB 업데이트
    3) 당일 매매정보 DB 입력
    4) 상위 50종목 조회
    5) 데이터 html table 형태로 변경
    6) 글 작성
    7) 각 주체별 상위 1등 종목 조회
    8) 크롬드라이버를 이용하여 네이버 종목토론방 정보 글 작성


## .env file
```go
-DB Contents-
DB_HOST_dev = "DB HOST"
DB_PORT_dev = "DB PORT"
DB_USERNAME_dev = "DB USER_NAME"
DB_PASSWORD_dev = "DB USER_PW"
DB_SCHEMAS_dev = "DB SCHEMA_NAME"
DB_NAME_dev = "DB NAME"

-Use Tistory Connect-
CLIENT_ID = "APP ID"
CLIENT_SECRET = "SECRET KEY"
REDIRECT_URL = "SEVER URL"


-Chrome Drover- // 한 계정당 글작성 갯수 제한으로 인해 2개의 아이디 사용 
NAVER_ID_KOSDAQ = "NAVER_ID"
NAVER_PW_KOSDAQ = "NAVER_PW"
NAVER_ID_KOSPI = "NAVER_ID"
NAVER_PW_KOSPI = "NAVER_PW"
```



# 개선점
1. 데이터를 HTS에서 수동으로 가져오는 문제
    - [ ] 증권사 API를 살펴볼 필요가 있음.
    - [ ] Mac에서는 안되므로 Window 전환 필요
2. 기존 CSV형태로 가져오게 되어있는데 excel로 데이터 전처리를 해야하는 문제.
    - [x] csv파일을 읽는 형태로 변환 필요 
