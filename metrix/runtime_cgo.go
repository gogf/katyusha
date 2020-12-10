// +build cgo
// +build !appengine

package metrix

import "runtime"

func numCgoCall() int64 {
	return runtime.NumCgoCall()
}
