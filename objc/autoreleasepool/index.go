package autoreleasepool

import "C"

type Page struct {
	// 对当前AutoreleasePoolPage 完整性的校验
	magic C.magic_t
	// 指向下一个即将产生的autoreleased对象的存放位置（当next == begin()时，表示AutoreleasePoolPage为空；当next == end()时，表示AutoreleasePoolPage已满
	next interface{}
	// 当前线程，表明与线程有对应关系
	pthread C.pthread_t
	// 指向父节点，第一个节点的 parent 值为 nil；
	parent *Page
	// 指向子节点，最后一个节点的 child 值为 nil；
	child *Page
	// 代表深度，第一个page的depth为0，往后每递增一个page，depth会加1；
	depth int32
	// 代表深度，第一个page的depth为0，往后每递增一个page，depth会加1；
	hiwat int32
}
