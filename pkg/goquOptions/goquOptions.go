package goquOptions

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type SelectDataset struct {
	Locked bool
}

func Apply(opts *SelectDataset, dataset *goqu.SelectDataset) *goqu.SelectDataset {
	if opts == nil {
		return dataset
	}
	if opts.Locked {
		dataset = dataset.ForUpdate(exp.Wait)
	}
	return dataset
}
