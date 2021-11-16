package utils

import (
	"fmt"
	"math"
)

func Paginate(currPage int, totalData int, limitData int) (page, string) {
	data := page{}
	if currPage <= 0 {
		currPage = 1
	}
	data.CurrentPage = currPage

	if totalData == 0 {
		return data, ""
	}

	if limitData <= 0 {
		limitData = 10
	}

	data.Pages = int(math.Ceil(float64(totalData) / float64(limitData)))
	if data.Pages <= 0 {
		return data, ""
	}

	if data.CurrentPage >= data.Pages {
		data.NextPage = 0
	} else {
		data.NextPage += data.CurrentPage + 1
	}

	if data.CurrentPage <= 1 {
		data.PrevPage = 0
	} else {
		data.PrevPage += data.CurrentPage - 1
	}
	data.Count = totalData
	return data, fmt.Sprintf(" OFFSET %d LIMIT %d ", (currPage-1)*limitData, limitData)
}

type page struct {
	Count       int         `json:"count"`
	CurrentPage int         `json:"current_page"`
	NextPage    int         `json:"next_page"`
	Pages       int         `json:"pages"`
	PrevPage    int         `json:"prev_page"`
	List        interface{} `json:"list"`
}
