package paging

const (
	// LimitKeyName is the request key name.
	LimitKeyName = "limit"

	// PageKeyName is the request page key name.
	PageKeyName = "page"

	DefaultLimit    = 20
	DefaultMaxLimit = 5000
	DefaultPage     = 1
)

// Paginator page data
type Paginator struct {
	TotalCount int  `json:"record_count" snake:"record_count"`
	TotalPage  int  `json:"page_count" snake:"page_count"`
	Page       int  `json:"absolute_page" snake:"absolute_page"`
	Limit      int  `json:"page_size" snake:"page_size"`
	Offset     int  `json:"-" snake:"-"`
	SnakeCase  bool `json:"-" snake:"-"`
}

func (p *Paginator) SetTotalCount(count int) {
	p.TotalCount = count
	if p.Limit < 0 {
		// unlimit
		p.TotalPage = 1
	} else {
		p.TotalPage = p.TotalCount / p.Limit
	}
	if p.TotalCount%p.Limit > 0 || p.TotalPage == 0 {
		p.TotalPage++
	}
}
