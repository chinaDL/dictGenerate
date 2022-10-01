package dictGenerate

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	Whitespace         = "\t\n\r\v\f"
	AsciiLowercase     = "abcdefghijklmnopqrstuvwxyz"
	AsciiUppercase     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AsciiLetters       = AsciiLowercase + AsciiUppercase
	Digits             = "0123456789"
	Hexdigits          = Digits + "abcdef" + "ABCDEF"
	HexdigitsLowercase = Digits + "abcdef"
	HexdigitsUppercase = Digits + "ABCDEF"
	Octdigits          = "01234567"
	Punctuation        = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	Printable          = Digits + AsciiLetters + Punctuation + Whitespace
)

func GenerateDo(charset string, n int, fn func(string, context.CancelFunc)) {
	ctx, cancel := context.WithCancel(context.Background())
	rCh := make(chan string, 0)
	go Generate(charset, n, ctx, rCh)
	defer cancel()
For:
	for {
		select {
		case v, ok := <-rCh:
			if !ok {

				break For
			}
			fn(v, cancel)
		case <-ctx.Done():
			break For
		}

	}
}

func Generate(charset string, n int, context context.Context, dictCh chan string) {
	m := strings.Split(charset, "")
	mc := len(m)
	cIndexList := make([]string, n)
	count := int(math.Pow(float64(len(m)), float64(n)))
	fmt.Printf("字典总数: %d\n", count)

	for i := 0; i < n; i++ {
		cIndexList[i] = "0"
	}

	//startList := strings.Repeat(m[0], n)
	//fmt.Println(startList)
	//cIndex := n - 1
	r := ""
For:
	for c := 0; c < count; c++ {

		r = ""
		for _, ci := range cIndexList {
			index, _ := strconv.Atoi(ci)
			r += m[index%len(m)]
		}
		//fmt.Println(r)
		cIndexStr := strings.Join(cIndexList, ",")
		cIndex := anyToDecimal(cIndexStr, mc)
		cIndex++
		tStr := decimalToAny(cIndex, mc)
		//cIndexStr = fmt.Sprintf("%s%s", strings.Repeat("0", n-len(tStr)), tStr)
		cIndexList = listLeftPad(strings.Split(tStr, ","), "0", n)
		//cIndexList = strings.Split(cIndexStr, ",")
		if dictCh != nil {
			dictCh <- r
		}
		select {
		case <-context.Done():
			break For
		default:

		}
	}
	if dictCh != nil {
		close(dictCh)
	}
}

func listLeftPad(oList []string, padChar string, count int) []string {
	ret := make([]string, 0)
	c := count - len(oList)
	for i := 0; i < c; i++ {
		ret = append(ret, padChar)
	}
	for _, v := range oList {
		ret = append(ret, v)
	}
	return ret
}

func anyToDecimal(num string, n int) int {
	var newNum float64
	newNum = 0.0
	ns := strings.Split(num, ",")
	nNum := len(ns) - 1
	for _, value := range ns {
		tValue, _ := strconv.Atoi(value)
		tmp := float64(tValue)
		if tmp != -1 {
			newNum = newNum + tmp*math.Pow(float64(n), float64(nNum))
			nNum = nNum - 1
		} else {
			break
		}
	}
	return int(newNum)
}

func decimalToAny(num, n int) string {
	newNumStr := ""
	var remainder int
	var remainderString string
	for num != 0 {
		remainder = num % n

		remainderString = strconv.Itoa(remainder)

		newNumStr = "," + remainderString + newNumStr
		num = num / n
	}
	return newNumStr[1:]
}
