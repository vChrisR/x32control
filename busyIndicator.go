package main

import (
	"github.com/therecipe/qt/core"
)

type BusyIndicator struct {
	core.QObject

	_ bool `property:"busy"`

	_ func() `constructor:"init"`
}

func (r *BusyIndicator) init() {
	r.SetBusy(true)
}
