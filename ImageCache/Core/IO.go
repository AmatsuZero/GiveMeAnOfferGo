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
	if !overwrite {
		flags = 0x0002 // windows 下没有 xattr.XATTR_CREATE，直接用数字代替
	}
	if follow {
		return xattr.SetWithFlags(path, name, vale, flags)
	} else {
		return xattr.LSetWithFlags(path, name, vale, flags)
	}
}

/// 获取扩展属性
func getExtendedAttributeNames(path string, follow bool) ([]string, error) {
	if follow {
		return xattr.List(path)
	}
	return xattr.LList(path)
}

func hasExtendedAttribute(name, path string, follow bool) (bool, error) {
	names, err := getExtendedAttributeNames(path, follow)
	if err != nil {
		return false, err
	}
	for _, n := range names {
		if n == name {
			return true, nil
		}
	}
	return false, nil
}

func getExtendedAttribute(name, path string, follow bool) ([]byte, error) {
	if follow {
		return xattr.Get(path, name)
	} else {
		return xattr.LGet(path, name)
	}
}

func removeExtendedAttribute(name, path string, follow bool) error {
	if follow {
		return xattr.Remove(path, name)
	} else {
		return xattr.LRemove(path, name)
	}
}
