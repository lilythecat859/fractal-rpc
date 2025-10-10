package bench
import "testing"
func BenchmarkGetSignatures(b *testing.B) {
        for i := 0; i < b.N; i++ {
                _ = "dummy"
        }
}
