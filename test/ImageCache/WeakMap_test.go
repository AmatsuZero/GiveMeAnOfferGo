package ImageCache

import (
	"github.com/AmatsuZero/GiveMeAnOfferGo/ImageCache/Core"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

type Point3D struct {
	x, y, z int
}

func TestWeakMapS2W(t *testing.T) {
	wm := Core.NewWeakMap(Core.WeakMapStrongToWeak)
	pt := &Point3D{1, 2, 3}
	key := 1
	wm.Set(key, pt)
	rhs, ok := wm.Get(key)
	assert.True(t, ok)
	assert.Equal(t, pt, rhs)
	pt = nil
	// ouch
	for i := 1; i < 10; i++ {
		runtime.Gosched()
		runtime.GC()
	}
	rhs, _ = wm.Get(key)
	assert.Nil(t, rhs)
}
