// +build go1.5

package metrix

import "runtime"

func gcCPUFraction(memStats *runtime.MemStats) float64 {
	return memStats.GCCPUFraction
}
