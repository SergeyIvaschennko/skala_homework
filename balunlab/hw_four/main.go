package main

import "fmt"

type OrderedMap struct {
	root *Node
	size int
}

type Node struct {
	key   int
	value int

	leftNode  *Node
	rightNode *Node
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		root: nil,
		size: 0,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.root == nil {
		m.size++
		m.root = &Node{key: key, value: value, leftNode: nil, rightNode: nil}
		return
	}

	current := m.root

	for {
		if current.key > key {
			if current.leftNode == nil {
				current.leftNode = &Node{key: key, value: value, leftNode: nil, rightNode: nil}
				m.size++
				break
			} else {
				current = current.leftNode
				continue
			}
		}

		if current.key < key {
			if current.rightNode == nil {
				current.rightNode = &Node{key: key, value: value, leftNode: nil, rightNode: nil}
				m.size++
				break
			} else {
				current = current.rightNode
				continue
			}
		}

		if current.key == key {
			return
		}
	}
}

func (m *OrderedMap) Erase(key int) {
	if m.root == nil {
		return
	}

	current := m.root
	var parent *Node
	var child *Node

	for {
		if current == nil {
			return
		}

		if current.key == key {
			break
		}

		parent = current

		if current.key > key {
			current = current.leftNode
			continue
		}

		if current.key < key {
			current = current.rightNode
			continue
		}
	}

	////ситуасион, когда нет детей
	if (current.leftNode == nil) && (current.rightNode == nil) {
		if parent == nil { // если один главный узел всего
			m.root = nil
			m.size--
			return
		}
		if parent.leftNode == current {
			parent.leftNode = nil
			m.size--
			return
		}
		if parent.rightNode == current {
			parent.rightNode = nil
			m.size--
			return
		}

		return
	}

	//ситуасион, когда есть ровно один ребенок
	if (current.leftNode != nil) != (current.rightNode != nil) {
		if current.leftNode != nil {
			child = current.leftNode
		}
		if current.rightNode != nil {
			child = current.rightNode
		}

		if parent == nil {
			m.root = child
			m.size--
			return
		}

		if parent.leftNode == current {
			parent.leftNode = child
			m.size--
			return
		}
		if parent.rightNode == current {
			parent.rightNode = child
			m.size--
			return
		}

		return
	}

	//cитуасион, когда два ребенка
	if (current.leftNode != nil) && (current.rightNode != nil) {
		targetParent := current
		target := current.leftNode

		for target.rightNode != nil {
			targetParent = target
			target = target.rightNode
		}

		current.key = target.key
		current.value = target.value

		if targetParent == current {
			targetParent.leftNode = target.leftNode
		} else {
			targetParent.rightNode = target.leftNode
		}

		m.size--
	}
}

func (m *OrderedMap) Contains(key int) bool {
	if m.root == nil {
		return false
	}

	current := m.root

	for {
		if current.key == key {
			return true
		}

		if current.key > key {
			if current.leftNode == nil {
				return false
			} else {
				current = current.leftNode
				continue
			}
		}

		if current.key < key {
			if current.rightNode == nil {
				return false
			} else {
				current = current.rightNode
				continue
			}
		}
	}
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	search(m.root, action)
}

func search(node *Node, action func(int, int)) {
	if node == nil {
		return
	}

	search(node.leftNode, action)
	action(node.key, node.value)
	search(node.rightNode, action)
}

func main() {
	m := NewOrderedMap()

	m.Insert(10, 10)
	m.Insert(5, 5)
	m.Insert(15, 15)
	m.Insert(2, 2)
	m.Insert(4, 4)
	m.Insert(12, 12)
	m.Insert(14, 14)

	printTree := func(title string) {
		fmt.Println("------", title, "------")
		fmt.Println("Size:", m.Size())
		m.ForEach(func(key, value int) {
			fmt.Printf("%d ", key)
		})
		fmt.Println("\n")
	}

	printTree("Исходное дерево")

	// Удаляем лист
	m.Erase(4)
	printTree("После удаления 4 (лист)")

	// Удаляем узел с одним ребёнком
	m.Erase(15)
	printTree("После удаления 15 (один ребёнок)")

	// Удаляем узел с двумя детьми
	m.Erase(5)
	printTree("После удаления 5 (два ребёнка)")

	// Удаляем корень
	m.Erase(10)
	printTree("После удаления 10 (корень)")
}
