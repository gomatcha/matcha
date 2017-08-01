package radix

import (
	"fmt"
	"strconv"
	"strings"
)

type Radix struct {
	root *Node
}

func NewRadix() *Radix {
	return &Radix{
		root: &Node{},
	}
}

func (r *Radix) Range(f func(path []int64, node *Node)) {
	r.root._range([]int64{}, f)
}

func (r *Radix) Insert(path []int64) *Node {
	return r.root.insert(path)
}

func (r *Radix) At(path []int64) *Node {
	return r.root.at(path)
}

func (r *Radix) Delete(path []int64) {
	remove := r.root.delete(path)
	if remove {
		r.root.Value = nil
		r.root.exists = false
	}
}

func (r *Radix) String() string {
	return r.root.debugString()
}

type Node struct {
	children map[int64]*Node
	exists   bool
	Value    interface{}
}

func (n *Node) _range(p []int64, f func(path []int64, node *Node)) {
	if n.exists {
		f(p, n)
	}

	p = append(p, 0)
	for id, i := range n.children {
		p[len(p)-1] = id
		i._range(p, f)
	}
}

func (n *Node) insert(path []int64) *Node {
	if len(path) == 0 {
		n.exists = true
		return n
	}

	child, ok := n.children[path[0]]
	if !ok {
		child = &Node{}
		if n.children == nil {
			n.children = map[int64]*Node{}
		}
		n.children[path[0]] = child
	}
	return child.insert(path[1:])
}

func (n *Node) at(path []int64) *Node {
	if len(path) == 0 {
		if n.exists == false {
			return nil
		}
		return n
	}
	child, ok := n.children[path[0]]
	if !ok {
		return nil
	}
	return child.at(path[1:])
}

func (n *Node) delete(path []int64) bool {
	if len(path) == 0 {
		if len(n.children) == 0 { // If node has no children, remove self.
			n.exists = false
			n.Value = nil
			return true
		}
		// Otherwise mark as non-existant.
		n.exists = false
		n.Value = nil
		return false
	}

	child, ok := n.children[path[0]]
	if !ok { // If path doesn't exist, abort.
		return false
	}

	remove := child.delete(path[1:])
	if !remove { // If child doesn't want to be removed, abort.
		return false
	}

	// Remove child, and remove self if we don't exist.
	delete(n.children, path[0])
	n.Value = nil
	return !n.exists
}

func (n *Node) debugString() string {
	all := []string{}
	for k, i := range n.children {
		lines := strings.Split(i.debugString(), "\n")
		for idx, line := range lines {
			if idx == 0 {
				lines[idx] = padRight(strconv.Itoa(int(k)), " ", 5) + line
			} else {
				lines[idx] = "|    " + line
			}
		}
		all = append(all, lines...)
	}

	str := fmt.Sprintf("{%v %v}", n.exists, n.Value)
	if len(all) > 0 {
		str += "\n" + strings.Join(all, "\n")
	}
	return str
}

func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}
