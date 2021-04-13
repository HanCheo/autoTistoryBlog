package common

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"net/http"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

// CheckErr err check
func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// CheckCode Server request check
func CheckCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with StatusCode :", res.StatusCode)
	}
}

// CleanString trim : 띄어쓰기 2개있는 것들 포함 모두 trim
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

// DecodeEUCKR str decode
func DecodeEUCKR(str string) string {
	var bufs bytes.Buffer

	wr := transform.NewWriter(&bufs, korean.EUCKR.NewDecoder())
	wr.Write([]byte(CleanString(str)))
	wr.Close()
	convVal := bufs.String()

	return convVal
}

// StrToNum 문자열을 파악하여 소수일경우 float64, 정수일경우 int 형으로 변환
func StrToNum(num string) (float64, int) {
	if num == "" {
		return 0, 0
	}

	re := regexp.MustCompile("[0-9.-]+")
	tmp := strings.Join(re.FindAllString(num, -1), "")
	if strings.Contains(tmp, ".") {
		nums, _ := strconv.ParseFloat(tmp, 32)
		convnum := math.Round(nums*100) / 100
		return convnum, 0
	}
	convnum, _ := strconv.Atoi(tmp)
	return 0, convnum
}

//Yymmdd : 시간을 주면 yymmdd 형태로 리턴
func Yymmdd(t time.Time) string {
	// 현재일
	y, m, d := t.Date()
	var year, month, day string

	year = strconv.Itoa(int(y - 2000))
	if int(m) < 10 {
		month = "0" + strconv.Itoa(int(m))
	} else {
		month = strconv.Itoa(int(m))
	}
	if int(d) < 10 {
		day = "0" + strconv.Itoa(int(d))
	} else {
		day = strconv.Itoa(int(d))
	}
	return year + month + day
}

// GetDateFormat set format now date
func GetDateFormat(y, m, d string) string {
	y2, m2, d2 := time.Now().Date()
	var year, month, day string

	year = strconv.Itoa(int(y2))
	month = strconv.Itoa(int(m2))
	day = strconv.Itoa(int(d2))

	//default == yyyy
	switch y {
	case "yy":
		year = year[2:]
		break
	case "":
		year = ""
	}

	switch m {
	case "mm":
		if int(m2) < 10 {
			month = "0" + month
		}
		break
	case "":
		month = ""
	}

	switch d {
	case "dd":
		if int(d2) < 10 {
			day = "0" + day
		}
		break
	case "":
		day = ""
	}

	return year + month + day
}

// InArray 값이 배열안에 있는지 확인
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {

	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

//Comma 숫자 1000단위 콤마
func Comma(v int64) string {
	sign := ""
	if v == math.MinInt64 {
		return "-9,223,372,036,854,775,808"
	}
	if v < 0 {
		sign = "-"
		v = 0 - v
	}
	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(v%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
}

//Openbrowser 브라우저창 오픈
func Openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

//GetStructKeys 구조체 키값만 추출
func GetStructKeys(m interface{}) []string {
	var structKeys []string
	t := reflect.TypeOf(m)
	for i := 0; i < t.NumField(); i++ {
		structKeys = append(structKeys, t.Field(i).Name)
	}
	return structKeys
}
