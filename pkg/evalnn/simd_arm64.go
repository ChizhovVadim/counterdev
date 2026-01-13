//go:build arm64

package evalnn

//go:noescape
func addNEON(dst, a, b []float32)

//go:noescape
func subNEON(dst, a, b []float32)

//go:noescape
func reluNEON(dst, src []float32)

//go:noescape
func dotProductNEON(a, b []float32) float32
