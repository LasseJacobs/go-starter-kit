package model

type Pagination struct {
	Page    uint64
	PerPage uint64
	//Count   uint64
}

func (p *Pagination) Offset() uint64 {
	return (p.Page - 1) * p.PerPage
}
