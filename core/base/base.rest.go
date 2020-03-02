package base

// page select body
type Pager struct {
	// this page
	Page int64
	// page size
	PageSize int64
	// Total nums
	Total int64
	// body
	Body interface{}
}

func (pager *Pager) GetLimit() int64 {
	if pager.PageSize > 50 {
		pager.PageSize = 20
	}
	return pager.PageSize
}

func (pager *Pager) GetOffset() int64 {
	// page size max 50
	if pager.PageSize > 50 {
		pager.PageSize = 20
	}
	skipPage := pager.Page - 1
	if skipPage < 0 {
		skipPage = 0
	}
	return skipPage * pager.PageSize
}
