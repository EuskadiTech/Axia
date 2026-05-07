package doc_create

import (
	"axia4/tools"
	"axia4/types"
	"fmt"
	"strings"
	"time"
)

func addFieldText(doc *doc, f types.DocFieldText, font types.DocFont, flowHorizontal bool) error {

	textValue, exists := f.Captions["docFieldText"][doc.p.GetLang()]
	if !exists {
		return nil
	}

	textValue = strings.ReplaceAll(textValue, "{PAGE_END}", "{nb}")
	textValue = strings.ReplaceAll(textValue, "{PAGE_CUR}", fmt.Sprintf("%d", doc.p.PageNo()))
	textValue = strings.ReplaceAll(textValue, "{DATE_TODAY}", time.Now().Local().Format(tools.GetDatetimeFormat(font.DateFormat, false)))
	textValue = strings.ReplaceAll(textValue, "{DATETIME_NOW}", time.Now().Local().Format(tools.GetDatetimeFormat(font.DateFormat, true)))
	textValue = strings.ReplaceAll(textValue, "{TIME_NOW}", time.Now().Local().Format(tools.GetTimeFormat()))

	drawCellText(doc, font, f.SizeX, 0, flowHorizontal, 0, textValue)
	return nil
}
