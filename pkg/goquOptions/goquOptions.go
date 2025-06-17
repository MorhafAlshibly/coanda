package goquOptions

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type SelectDataset struct {
	Locked bool
}

func (s *SelectDataset) Apply(dataset *goqu.SelectDataset) *goqu.SelectDataset {
	if s.Locked {
		dataset = dataset.ForUpdate(exp.Wait)
	}
	return dataset
}
