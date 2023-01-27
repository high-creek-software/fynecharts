package fynecharts

import (
	"log"
	"testing"
)

func TestTickLabels(t *testing.T) {
	idxs, step, q, magnitude, err := generateTicks(3, 108, 7, containmentWithinData, defaultQ(), defaultWeights(), defaultLegibility)
	if err != nil {
		t.Error("got error generating ticks", err)
	}

	log.Println(idxs, step, q, magnitude)
}
