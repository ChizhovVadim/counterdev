//go:build !arm64

package evalnn

func addNEON(dst, a, b []float32)           {}
func subNEON(dst, a, b []float32)           {}
func reluNEON(dst, src []float32)           {}
func dotProductNEON(a, b []float32) float32 { return 0 }
