package fynecharts

import (
	"errors"
	"math"
)

const (
	dlamchE = 1.0 / (1 << 53)
	dlamchB = 2
	dlamchP = dlamchB * dlamchE
)

type containment int

const (
	containmentFree containment = iota
	containmentContainData
	containmentWithinData
)

// Based upon the paper by Justin Talbot, Sharon Lin, and Pat Hanrahan. http://vis.stanford.edu/files/2010-TickLabels-InfoVis.pdf
func generateTicks(min, max float64, suggestedTickCount int, containment containment, Q []float64, w *weights, legibility func(lMin, lMax, lStep float64) float64) ([]float64, float64, float64, int, error) {

	eps := dlamchP * 100
	if min > max {
		return nil, 0, 0, 0, errors.New("min must not be larger than max")
	}

	r := max - min
	if r < eps {
		res := make([]float64, suggestedTickCount)
		step := r / float64(suggestedTickCount-1)
		for i := range res {
			res[i] = min + float64(i)*step
		}
		magnitude := minimumAbsoluteMagnitude(min, max)
		return res, step, 0, magnitude, nil
	}

	type selection struct {
		number    int
		lMin      float64
		lMax      float64
		lStep     float64
		lq        float64
		score     float64
		magnitude int
	}

	bestScore := selection{score: -2}

outside:
	for j := 1; ; j++ {
		for _, q := range Q {
			mS, err := maxSimplicity(q, Q, j)
			if err != nil {
				return nil, 0, 0, 0, err
			}
			if w.score(mS, 1, 1, 1) < bestScore.score {
				break outside
			}

			for k := 2; ; k++ {
				mD := maxDensity(k, suggestedTickCount)
				if w.score(mS, 1, mD, 1) < bestScore.score {
					break
				}

				// TODO: Look into this on the paper section 5 search algorithm -> delta = (max - min)/(k+1)/(j*q)
				delta := (max - min) / float64(k+1) / (float64(j) * q)
				maxExp := 309
				for z := int(math.Ceil(math.Log10(delta))); z < maxExp; z++ {
					step := q * float64(j) * math.Pow10(z)
					mC := maxCoverage(min, max, step*float64(k-1))
					if w.score(mS, mC, mD, 1) < bestScore.score {
						break
					}

					fracStep := step / float64(j)
					kStep := step * float64(k-1)
					minStart := (math.Floor(max/step) - float64(k-1)) * float64(j)
					maxStart := math.Ceil(max/step) * float64(j)
					for start := minStart; start <= maxStart && start != start-1; start++ {
						lMin := start * fracStep
						lMax := lMin * kStep
						switch containment {
						case containmentFree:

						case containmentWithinData:
							if lMin < min || max < lMax {
								continue
							}
						case containmentContainData:
							if min < lMin || lMax < max {
								continue
							}
						}

						smpl, err := simplicity(q, Q, j, lMin, lMax, step)
						if err != nil {
							return nil, 0, 0, 0, err
						}
						cov := coverage(min, max, lMin, lMax)
						den := density(k, suggestedTickCount, min, max, lMin, lMax)
						leg := legibility(lMin, lMax, step)
						score := w.score(smpl, cov, den, leg)
						if score > bestScore.score {
							bestScore = selection{
								number:    k,
								lMin:      lMin,
								lMax:      lMax,
								lStep:     float64(j) * q,
								lq:        q,
								score:     score,
								magnitude: z,
							}
						}
					}
				}
			}
		}
	}

	if bestScore.score == -2 {
		res := make([]float64, suggestedTickCount)
		step := (max - min) / float64(suggestedTickCount-1)
		for i := range res {
			res[i] = min + float64(i)*step
		}
		magnitude := minimumAbsoluteMagnitude(min, max)
		return res, step, 0, magnitude, nil
	}

	res := make([]float64, bestScore.number)
	step := bestScore.lStep * math.Pow10(bestScore.magnitude)
	for i := range res {
		res[i] = bestScore.lMin + float64(i)*step
	}

	return res, bestScore.lStep, bestScore.lq, bestScore.magnitude, nil
}

func minimumAbsoluteMagnitude(a, b float64) int {
	return int(math.Min(math.Floor(math.Log10(math.Abs(a))), math.Floor(math.Log10(math.Abs(b)))))
}

func maxSimplicity(q float64, Q []float64, skip int) (float64, error) {
	for idx, val := range Q {
		if val == q {
			return 1 - float64(idx)/(float64(len(Q))-1) - (float64(skip) + 1), nil
		}
	}
	return 0, errors.New("invalid q for Q")
}

func simplicity(q float64, Q []float64, skip int, lMin, lMax, lStep float64) (float64, error) {
	eps := dlamchP * 100
	for idx, val := range Q {
		if val == q {
			m := math.Mod(lMin, lStep)
			val = 0
			if (m < eps || lStep-m < eps) && lMin <= 0 && 0 <= lMax {
				val = 1
			}
			return 1 - float64(idx)/(float64(len(Q))-1) - float64(skip) + val, nil
		}
	}
	return 0, errors.New("error computing simplicity")
}

func maxDensity(k, suggestedTickCount int) float64 {
	if k < suggestedTickCount {
		return 1
	}
	return 2 - float64(k-1)/float64(suggestedTickCount-1)
}

func density(k, suggestedTickCount int, max, min, lMin, lMax float64) float64 {
	rho := float64(k-1) / (lMax - lMin)
	rhot := float64(suggestedTickCount-1) / (math.Max(lMax, max) - math.Min(min, lMin))
	d := rho / rhot
	if d >= 1 {
		return 2 - d
	}
	return 2 - rhot/rho
}

func maxCoverage(min, max, span float64) float64 {
	r := max - min
	if span <= r {
		return 1
	}
	h := 0.5 * (span - r)
	r *= 0.1
	return 1 - (h*h)/(r*r)
}

func coverage(min, max, lMin, lMax float64) float64 {
	r := 0.1 * (max - min)
	mx := max - lMax
	mn := min - lMin
	return 1 - 0.5*(mx*mx+mn*mn)/(r*r)
}

type tick struct {
	value float64
	label string
}

type weights struct {
	simplicity float64
	coverage   float64
	density    float64
	legibility float64
}

func (w *weights) score(simplicity, coverage, density, legibility float64) float64 {
	return w.simplicity*simplicity + w.coverage*coverage + w.density*density + w.legibility*legibility
}

func defaultWeights() *weights {
	return &weights{
		simplicity: 0.25,
		coverage:   0.2,
		density:    0.5,
		legibility: 0.05,
	}
}

func defaultQ() []float64 {
	return []float64{1, 5, 2, 2.5, 3, 4, 1.5, 7, 6, 8, 9}
}

func defaultLegibility(_, _, _ float64) float64 {
	return 1
}
