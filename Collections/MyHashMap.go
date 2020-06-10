package Collections

type MyHashMap struct {
	val [1010]list
}

type list struct {
	values []node
}

type node struct {
	key   int
	value int
}

func NewHashMap() MyHashMap {
	return MyHashMap{val: [1010]list{}}
}

func (m *MyHashMap) Put(key int, value int) {
	r := key % 1009
	// 如果存在
	l := m.val[r].values
	for i := 0; i < len(l); i++ {
		if l[i].key == key {
			l[i].value = value // 如果存在就修改值
			return
		}
	}
	m.val[r].values = append(m.val[r].values, node{
		key:   key,
		value: value,
	})
}

func (m *MyHashMap) Get(key int) int {
	r := key % 1009
	l := m.val[r].values
	for i := 0; i < len(l); i++ {
		if l[i].key == key {
			return l[i].value
		}
	}
	return -1
}

func (m *MyHashMap) Remove(key int) {
	r := key % 1009
	l := m.val[r].values
	for i := 0; i < len(l); i++ {
		if l[i].key == key {
			m.val[r].values = append(l[:i], l[i+1:]...)
		}
	}
}
