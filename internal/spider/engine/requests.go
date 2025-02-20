package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"kd_api/internal/spider/goreq"
)

func (c *KdClient) loginReq() (keyUrl string, err error) {
	res, err := c.Get(loginApi, defaultHeaders, nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	doc, err := GetDocFromRes(res)
	if err != nil {
		return "", err
	}

	keyUrl, exists := doc.Find("#kingo_encypt").Attr("src")
	if !exists {
		return "", errors.New("loginReq: keyUrl not exists")
	}

	return keyUrl, nil
}

func (c *KdClient) fetchKeyReq(keyUrl string) (string, string, error) {
	res, err := c.Get(keyUrl, nil, nil)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	doc, err := GetDocFromRes(res)
	if err != nil {
		return "", "", err
	}
	result := desKeyReg.FindStringSubmatch(doc.Text())
	var desKey, sessionId string
	if len(result) > 0 {
		desKey = result[1]
	}
	sessionId = c.GetCookie(res.Request.URL, "JSESSIONID")
	return desKey, sessionId, nil
}

func (c *KdClient) logonReq(params string) (string, error) {
	headers := defaultHeaders
	headers["X-Requested-With"] = "XMLHttpRequest"
	res, err := c.Post(logonApi+"?"+params, headers, nil, nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := goreq.ResToJson(res)
	if err != nil {
		return "", err
	}
	if v, ok := data["status"].(string); !ok || v != "200" {
		return "", fmt.Errorf("logon: (%#v)", data)
	}
	if v, ok := data["result"].(string); ok {
		return indexUrl + v, nil
	} else {
		return "", fmt.Errorf("logon: (%#v)", data)
	}
}

func (c *KdClient) entryIndex(url string) error {
	res, err := c.Get(url, defaultHeaders, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (c *KdClient) studentInfoReq() (*StudentInfo, error) {
	res, err := c.Get(studentInfoApi, defaultHeaders, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := GetDocFromRes(res)
	if err != nil {
		return nil, err
	}

	infoMap := make(map[string]string)
	content := doc.Text()
	matches := varMapReg.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) == 3 {
			k := strings.TrimSpace(match[1])
			v := strings.TrimSpace(match[2])
			v = strings.Replace(v, `"`, "", -1)
			v = strings.Replace(v, `'`, "", -1)
			infoMap[k] = v
		}
	}
	var schoolCode, StudentId uint32
	var StudentCode uint64
	u, err := strconv.ParseUint(infoMap["G_SCHOOL_CODE"], 10, 32)
	schoolCode = uint32(u)
	if err != nil {
		return nil, err
	}
	u, err = strconv.ParseUint(infoMap["G_LOGIN_ID"], 10, 32)
	StudentId = uint32(u)
	if err != nil {
		return nil, err
	}
	StudentCode, err = strconv.ParseUint(infoMap["G_USER_CODE"], 10, 64)
	if err != nil {
		return nil, err
	}
	return &StudentInfo{
		SchoolCode:  schoolCode,
		SchoolName:  infoMap["G_SCHOOL_NAME"],
		StudentId:   StudentId,
		StudentCode: StudentCode,
		Name:        infoMap["G_USER_NAME"],
	}, nil
}

func (c *KdClient) fetchTeachingEvaluateTablesReq(tableId, fre int) error {
	query := fmt.Sprintf("tableId=%d&fre=%d", tableId, fre)
	api := teachingEvaluateTableApi + "?" + query
	data := map[string]string{
		"xn":               "2024",
		"xq":               "0",
		"pjlc":             "2024003",
		"sfzbpj":           "1",
		"sfwjpj":           "1",
		"pjzt_m":           "20",
		"pjfsbz":           "0",
		"qyxjkc":           "",
		"records":          "",
		"menucode":         "S902",
		"hidKey":           "",
		"xwtj":             "1",
		"menucode_current": "S902",
	}
	formData := url.Values{}
	for k, v := range data {
		formData.Set(k, v)
	}

	headers := defaultHeaders
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	res, err := c.Post(api, headers, nil, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	doc, err := GetDocFromRes(res)
	if err != nil {
		return err
	}

	tables, err := ExpandTableFromHtml(doc)
	if err != nil {
		return err
	}

	formLen := len(tables[0]) - 1
	var params []map[string]interface{}
	for i := 0; i < formLen; i++ {
		s := doc.Find(fmt.Sprintf("#tr%d_wjdc a", i))
		t, exists := s.Attr("onclick")
		if exists {
			param := t[13 : len(t)-6]
			param = strings.ReplaceAll(param, `\`, "")

			paramMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(param), &paramMap)
			if err != nil {
				continue
			}
			params = append(params, paramMap)
		}
	}
	fmt.Println(params)

	c.submitTeachingEvaluateFormReq(params[0])
	//for _, v := range params {
	//
	//}

	return nil
}

func (c *KdClient) submitTeachingEvaluateFormReq(params map[string]interface{}) error {
	type SubmitParams struct {
		P           map[string]interface{}
		Scores      [20]int
		Suggestion  string
		StudentCode uint64
	}
	submitParams := SubmitParams{
		P: params,
		Scores: [20]int{
			100, 100, 100, 100, 100,
			100, 100, 100, 100, 100,
			100, 100, 100, 100, 100,
			100, 100, 100, 100, 100,
		},
		Suggestion:  "0",
		StudentCode: 201700008496,
	}

	scoresStr := ""
	totalScore := 0
	for i, v := range submitParams.Scores {
		scoresStr += fmt.Sprintf("%d@00%d@ ;", v, i)
		totalScore += v
	}
	scoresStr = "100@0002@ ;100@0004@ ;100@0005@ ;100@0008@ ;100@0010@ ;100@0012@ ;100@0014@ ;100@0016@ ;100@0018@ ;100@0020@ ;100@0022@ ;100@0024@ ;"
	totalScore = 1200

	data := map[string]interface{}{
		"wspjZbpjWjdcForm.pjlb_m":         submitParams.P["pjlb_m"],                           // from param.pjlb_m
		"wspjZbpjWjdcForm.sfzjjs":         submitParams.P["sfzjjs"],                           // from param
		"wspjZbpjWjdcForm.commitZB":       scoresStr,                                          // 分数@0011@;...;100@0020@;
		"wspjZbpjWjdcForm.commitWJText":   fmt.Sprintf("0006@#@%s;", submitParams.Suggestion), // 0006@#@0;  => 0006@#@建立里面的内容;
		"wspjZbpjWjdcForm.commitWJSelect": "",
		"wspjZbpjWjdcForm.xn":             submitParams.P["xn"],
		"wspjZbpjWjdcForm.xq":             submitParams.P["xq"],
		"wspjZbpjWjdcForm.jsid":           submitParams.P["jsid"],   // from param
		"wspjZbpjWjdcForm.kcdm":           submitParams.P["kcdm"],   // from param
		"wspjZbpjWjdcForm.skbjdm":         submitParams.P["skbjdm"], // ...
		"wspjZbpjWjdcForm.pjlc":           submitParams.P["pjlc"],   // ...
		"wspjZbpjWjdcForm.userCode":       submitParams.StudentCode, // 学生的code
		"wspjZbpjWjdcForm.pjzt_m":         submitParams.P["pjzt_m"], // ...
		"wspjZbpjWjdcForm.zbmb_m":         submitParams.P["zbmb_m"], // ?
		"wspjZbpjWjdcForm.wjmb_m":         "001",                    // ?
		"bfzfs_xx":                        "",
		"bfzfs_sx":                        "",
		"totalcj":                         totalScore,               // 总分
		"zbSize":                          12,                       //
		"zbmb":                            submitParams.P["zbmb_m"], // ?
		"sel_scorecj0":                    100,
		"sel_scorecj1":                    100,
		"sel_scorecj2":                    100,
		"sel_scorecj3":                    100,
		"sel_scorecj4":                    100,
		"sel_scorecj5":                    100,
		"sel_scorecj6":                    100,
		"sel_scorecj7":                    100,
		"sel_scorecj8":                    100,
		"sel_scorecj9":                    100,
		"sel_scorecj10":                   100,
		"sel_scorecj11":                   100,
		"wjmb":                            "001",                   // ?
		"wjSize":                          1,                       // ?
		"area0":                           submitParams.Suggestion, // 建议
		"menucode_current":                "S902",                  // ...
	}
	formData := url.Values{}
	for k, v := range data {
		encoded := url.QueryEscape(fmt.Sprintf("%v", v))
		formData.Set(k, encoded)
	}

	ss := ""
	for k, v := range data {
		data[k] = url.QueryEscape(fmt.Sprintf("%v", v))
		v = url.QueryEscape(fmt.Sprintf("%v", v))
		ss += fmt.Sprintf("%v=%v&", k, v)
	}
	fmt.Println(data)

	headers := defaultHeaders
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	d := "wspjZbpjWjdcForm.pjlb_m=08&wspjZbpjWjdcForm.sfzjjs=1&wspjZbpjWjdcForm.commitZB=100%400002%40+%3B100%400004%40+%3B100%400005%40+%3B100%400008%40+%3B100%400010%40+%3B100%400012%40+%3B100%400014%40+%3B100%400016%40+%3B100%400018%40+%3B100%400020%40+%3B100%400022%40+%3B100%400024%40+%3B&wspjZbpjWjdcForm.commitWJText=0006%40%23%400%3B&wspjZbpjWjdcForm.commitWJSelect=&wspjZbpjWjdcForm.xn=2024&wspjZbpjWjdcForm.xq=0&wspjZbpjWjdcForm.jsid=101017&wspjZbpjWjdcForm.kcdm=2017566&wspjZbpjWjdcForm.skbjdm=2017566&wspjZbpjWjdcForm.pjlc=2024003&wspjZbpjWjdcForm.userCode=201700008496&wspjZbpjWjdcForm.pjzt_m=20&wspjZbpjWjdcForm.zbmb_m=022&wspjZbpjWjdcForm.wjmb_m=001&bfzfs_xx=&bfzfs_sx=&totalcj=1200&zbSize=12&zbmb=022&sel_scorecj0=100&sel_scorecj1=100&sel_scorecj2=100&sel_scorecj3=100&sel_scorecj4=100&sel_scorecj5=100&sel_scorecj6=100&sel_scorecj7=100&sel_scorecj8=100&sel_scorecj9=100&sel_scorecj10=100&sel_scorecj11=100&wjmb=001&wjSize=1&area0=0&menucode_current=S902"
	//d := "wspjZbpjWjdcForm.pjlb_m=01&wspjZbpjWjdcForm.sfzjjs=1&wspjZbpjWjdcForm.commitZB=100%400012%40+%3B100%400013%40+%3B100%400011%40+%3B100%400014%40+%3B100%400015%40+%3B100%400016%40+%3B100%400017%40+%3B100%400018%40+%3B100%400019%40+%3B100%400020%40+%3B&wspjZbpjWjdcForm.commitWJText=0006%40%23%400%3B&wspjZbpjWjdcForm.commitWJSelect=&wspjZbpjWjdcForm.xn=2024&wspjZbpjWjdcForm.xq=0&wspjZbpjWjdcForm.jsid=101878&wspjZbpjWjdcForm.kcdm=20173507&wspjZbpjWjdcForm.skbjdm=20173507-011&wspjZbpjWjdcForm.pjlc=2024001&wspjZbpjWjdcForm.userCode=201700008496&wspjZbpjWjdcForm.pjzt_m=20&wspjZbpjWjdcForm.zbmb_m=015&wspjZbpjWjdcForm.wjmb_m=001&bfzfs_xx=&bfzfs_sx=&totalcj=1000&zbSize=10&zbmb=015&sel_scorecj0=100&sel_scorecj1=100&sel_scorecj2=100&sel_scorecj3=100&sel_scorecj4=100&sel_scorecj5=100&sel_scorecj6=100&sel_scorecj7=100&sel_scorecj8=100&sel_scorecj9=100&wjmb=001&wjSize=1&area0=0&menucode_current=S902"
	d = ss
	fmt.Println(d)
	res, err := c.Post(teachingEvaluateFormSubmitApi, headers, nil, bytes.NewBuffer([]byte(d)))
	//res, err := c.Post(teachingEvaluateFormSubmitApi, headers, nil, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	doc, err := GetDocFromRes(res)
	if err != nil {
		return err
	}
	fmt.Println(doc.Text())

	return nil
}
