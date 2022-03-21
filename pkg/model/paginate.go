package model

type Pagination struct {
	Page    int32
	PerPage int32
}

func (p *Pagination) Offset() int32 {
	return (p.Page - 1) * p.PerPage
}
