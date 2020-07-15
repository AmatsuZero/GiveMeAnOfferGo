package Core

import "github.com/pkg/xattr"

type DataReadingOption int

const (
	ReadingMappedIfSafe DataReadingOption = iota
	ReadingUncached
	ReadingMappedAlways
)

type DataWritingOption int

const (
	WritingAtomic DataWritingOption = iota
)

/// 给文件设置扩展属性
func setExtendedAttribute(name, path string, vale []byte, follow, overwrite bool) error {
	flags := 0
	if !follow {
		flags = xattr.XATTR_NOFOLLOW
	}
	if overwrite {
		flags |= 0
	} else {
		flags |= xattr.XATTR_CREATE
	}
	return xattr.SetWithFlags(path, name, vale, flags)
}

/// 获取扩展属性
func getExtendedAttributeAtPath(path string, traverseLink bool) ([]string, error) {

	return nil, nil
}
