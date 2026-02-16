package main

import "sort"

// buildTree builds a proper process tree for filtered lists.
// Returns processes in depth-first order with TreePrefix set.
func buildTree(processes []*Process) []*Process {
	pidMap := make(map[int]*Process)
	for _, p := range processes {
		pidMap[p.PID] = p
	}

	childrenMap := make(map[int][]*Process)
	for _, p := range processes {
		if _, ok := pidMap[p.PPID]; ok {
			childrenMap[p.PPID] = append(childrenMap[p.PPID], p)
		}
	}

	// Find roots: processes whose parent isn't in our list
	var roots []*Process
	for _, p := range processes {
		if _, ok := pidMap[p.PPID]; !ok {
			roots = append(roots, p)
		}
	}
	sort.Slice(roots, func(i, j int) bool {
		return roots[i].ElapsedSec < roots[j].ElapsedSec
	})

	var result []*Process
	var flatten func(p *Process, depth int, isLast bool, prefix string)
	flatten = func(p *Process, depth int, isLast bool, prefix string) {
		if depth == 0 {
			p.TreePrefix = ""
		} else {
			connector := "\u251c\u2500 " // ├─
			if isLast {
				connector = "\u2514\u2500 " // └─
			}
			p.TreePrefix = prefix + connector
		}
		p.TreeDepth = depth
		result = append(result, p)

		children := childrenMap[p.PID]
		sort.Slice(children, func(i, j int) bool {
			return children[i].ElapsedSec < children[j].ElapsedSec
		})

		for i, child := range children {
			isLastChild := i == len(children)-1
			var childPrefix string
			if depth == 0 {
				childPrefix = ""
			} else if isLast {
				childPrefix = prefix + "   "
			} else {
				childPrefix = prefix + "\u2502  " // │
			}
			flatten(child, depth+1, isLastChild, childPrefix)
		}
	}

	for _, root := range roots {
		flatten(root, 0, true, "")
	}

	return result
}

// addParentNames adds parent info for --all tree mode.
// Instead of a full tree, just shows the parent name inline.
func addParentNames(processes []*Process) {
	pidMap := make(map[int]*Process)
	for _, p := range processes {
		pidMap[p.PID] = p
	}
	for _, p := range processes {
		if parent, ok := pidMap[p.PPID]; ok {
			name := parent.Parsed
			if len(name) > 15 {
				name = name[:15]
			}
			p.ParentName = name
		}
	}
}
