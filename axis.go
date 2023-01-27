package fynecharts

type axis struct {
	min, max, dataRange float64
	normalizer          normalizer
}

func (a axis) normalize(x float64) float32 {
	return a.normalizer.normalize(a.min, a.max, x)
}

type normalizer interface {
	normalize(min, max, x float64) float32
}

var _ normalizer = linearNormalizer{}

type linearNormalizer struct {
}

func (ln linearNormalizer) normalize(min, max, x float64) float32 {
	return float32((x - min) / (max - min))
}
