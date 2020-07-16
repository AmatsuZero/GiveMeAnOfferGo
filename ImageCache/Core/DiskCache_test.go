package Core

import (
	"testing"
)

func TestMoveDir(t *testing.T) {
	config := DefaultImageCacheConfig()
	cache := newDiskCache(&config, "/Users/jiangzhenhua/Desktop/rhs")
	cache.MoveCacheDirectory("/Users/jiangzhenhua/Desktop/rhs", "/Users/jiangzhenhua/Desktop/lhs")
}
