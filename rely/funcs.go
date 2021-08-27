package rely

func Norm(x float64, ma float64, mi float64) float64 {
	return (x-mi)/(ma-mi)
}
