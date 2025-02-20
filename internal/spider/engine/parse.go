package engine

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"

	"github.com/PuerkitoBio/goquery"
)

var (
	classScheduleInfoReg = regexp.MustCompile(`([^\x{00}-\x{ff}]+)\s+([^\x{00}-\x{ff}]+)\s+([1-9,-]+)\[(\d-\d)]\s+([^\x{00}-\x{ff}]+[A-Z]*\d+)`)
	//classScheduleInfoReg = regexp.MustCompile(`(.+?)\s+?(.+?)\s?(\d?)\[(.+?)]\s?(.+?\d+)`)
)

type tdItem struct {
	Content string
	rowSpan int
	colSpan int
}

//func (it tdItem) String() string {
//	return fmt.Sprintf("(s:%s,r:%d,c:%d)", it.Content, it.rowSpan, it.colSpan)
//}

type htmlTableItem = [][]tdItem

func GetDocFromRes(res *http.Response) (*goquery.Document, error) {
	ct := res.Header.Get("Content-Type")
	//if ct == "application/vnd.ms-excel;charset=gbk" {
	//	ct = "application/vnd.ms-excel; charset=gbk"
	//}
	reader, err := charset.NewReader(res.Body, ct)
	//reader, err := charset.NewReaderLabel("gbk", res.Body)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(reader)
}

func ExpandTableFromHtml(doc *goquery.Document) ([][][]string, error) {
	var tables [][][]string
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		var htmlTable htmlTableItem
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			var rowData []tdItem
			row.Find("th, td").Each(func(k int, cell *goquery.Selection) {
				rowData = append(rowData, tdItem{
					Content: getText(cell),
					rowSpan: getRowspan(cell),
					colSpan: getColspan(cell),
				})
			})
			htmlTable = append(htmlTable, rowData)
		})

		t, _ := expandTable(htmlTable)
		tables = append(tables, t)
	})

	return tables, nil
}

func ExpandClassScheduleInfo(table [][]string) []ClassItem {
	var classSchedule []ClassItem
	for i := 1; i < len(table); i++ {
		for j := 2; j < len(table[i]); j++ {
			weekday := j - 1
			newItems := expandClassInfo(table[i][j], weekday)
			classSchedule = append(classSchedule, newItems...)
		}
	}
	return classSchedule
}

func expandClassInfo(s string, weekday int) []ClassItem {
	var classItems []ClassItem

	results := classScheduleInfoReg.FindAllStringSubmatch(s, -1)
	for _, result := range results {
		name := result[1]
		teacher := result[2]

		var weeks []int
		weeksStr := strings.Split(result[3], ",")
		for _, week := range weeksStr {
			if strings.Contains(week, "-") {
				parts := strings.Split(week, "-")
				start, _ := strconv.Atoi(parts[0])
				end, _ := strconv.Atoi(parts[1])

				for i := start; i <= end; i++ {
					weeks = append(weeks, i)
				}
			} else {
				num, _ := strconv.Atoi(week)
				weeks = append(weeks, num)
			}
		}

		var times []int
		parts := strings.Split(result[4], "-")
		for _, part := range parts {
			num, _ := strconv.Atoi(part)
			times = append(times, num)
		}

		classroom := result[5]
		item := ClassItem{
			Name:      name,
			Teacher:   teacher,
			Weekday:   weekday,
			Weeks:     weeks,
			Times:     times,
			Classroom: classroom,
		}
		classItems = append(classItems, item)
	}
	return classItems
}

func expandTable(data htmlTableItem) ([][]string, error) {
	var totalRow, totalCol int

	if len(data) > 0 {
		for _, item := range data[0] {
			totalCol += item.colSpan
		}
	}
	for i := 0; i < len(data); i++ {
		if len(data[i]) > 0 {
			totalRow += data[i][0].rowSpan
		}
	}
	table := make([][]string, totalRow)
	mark := make([][]bool, totalRow)
	for i := 0; i < totalRow; i++ {
		table[i] = make([]string, totalCol)
		mark[i] = make([]bool, totalCol)
	}

	for row, rowItems := range data {
		for col, item := range rowItems {
			r, c := findNext(mark, row, col)
			if r == -1 || c == -1 {
				return nil, fmt.Errorf("row %d col %d not found", row, col)
			}
			for i := 0; i < item.rowSpan; i++ {
				for j := 0; j < item.colSpan; j++ {
					if r+i < totalRow && c+j < totalCol {
						table[r+i][c+j] = item.Content
						mark[r+i][c+j] = true
					}
				}
			}
		}
	}
	return table, nil
}

func findNext(table [][]bool, r, c int) (int, int) {
	for i := r; i < len(table); i++ {
		for j := c; j < len(table[i]); j++ {
			if !table[i][j] {
				return i, j
			}
		}
	}
	return -1, -1
}

func getText(s *goquery.Selection) string {
	htmls, _ := s.Html()
	htmls = strings.ReplaceAll(htmls, "<br/>", " ")
	htmls = func(content string) string {
		buf := bytes.NewBuffer(nil)
		inTag := false
		for _, char := range content {
			switch char {
			case '<':
				inTag = true
			case '>':
				inTag = false
			default:
				if !inTag {
					buf.WriteRune(char)
				}
			}
		}
		return buf.String()
	}(htmls)
	return strings.TrimSpace(html.UnescapeString(htmls))
}

func getAttr(s *goquery.Selection, attr string) string {
	if val, exists := s.Attr(attr); exists {
		return val
	}
	return ""
}

func getRowspan(s *goquery.Selection) int {
	if s := getAttr(s, "rowspan"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			return 1
		}
		return v
	}
	return 1
}

func getColspan(s *goquery.Selection) int {
	if s := getAttr(s, "colspan"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			return 1
		}
		return v
	}
	return 1
}
