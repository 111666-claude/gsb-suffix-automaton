// Package suffixauto 实现后缀自动机（出现次数、不同子串数、最长公共子串）。
package suffixauto

type state struct {
	next map[byte]int32
	link int32
	size int32
}

// Automaton 是后缀自动机。
type Automaton struct {
	states []state
	last   int32
	occ    []int64
}

// Build 建自动机，并沿 link 树累加出现次数。
func Build(text string) *Automaton {
	automaton := &Automaton{last: 0}
	automaton.states = make([]state, 1, 2*len(text)+1)
	automaton.states[0] = state{next: map[byte]int32{}, link: -1}
	automaton.occ = make([]int64, 1, 2*len(text)+1)
	for index := 0; index < len(text); index++ {
		automaton.extend(text[index])
	}
	automaton.propagate()
	return automaton
}

func (a *Automaton) extend(char byte) {
	current := int32(len(a.states))
	a.states = append(a.states, state{next: map[byte]int32{}, link: -1,
		size: a.states[a.last].size + 1})
	a.occ = append(a.occ, 1)
	previous := a.last
	for previous >= 0 {
		if _, ok := a.states[previous].next[char]; ok {
			break
		}
		a.states[previous].next[char] = current
		previous = a.states[previous].link
	}
	if previous < 0 {
		a.states[current].link = 0
		a.last = current
		return
	}
	target, _ := a.states[previous].next[char]
	if a.states[previous].size+1 == a.states[target].size {
		a.states[current].link = target
		a.last = current
		return
	}
	clone := int32(len(a.states))
	copied := state{next: map[byte]int32{}, link: a.states[target].link,
		size: a.states[previous].size + 1}
	for key, value := range a.states[target].next {
		copied.next[key] = value
	}
	a.states = append(a.states, copied)
	a.occ = append(a.occ, 0) // clone 初始计数为 0
	for previous >= 0 && a.states[previous].next[char] == target {
		a.states[previous].next[char] = clone
		previous = a.states[previous].link
	}
	a.states[target].link = clone
	a.states[current].link = clone
	a.last = current
}

// propagate 按 size 从大到小把每个状态的计数累加到 link 上。
func (a *Automaton) propagate() {
	maxSize := int32(0)
	for index := range a.states {
		if a.states[index].size > maxSize {
			maxSize = a.states[index].size
		}
	}
	bucket := make([]int32, maxSize+1)
	for index := range a.states {
		bucket[a.states[index].size]++
	}
	for index := 1; index < len(bucket); index++ {
		bucket[index] += bucket[index-1]
	}
	order := make([]int32, len(a.states))
	for index := len(a.states) - 1; index >= 0; index-- {
		size := a.states[index].size
		bucket[size]--
		order[bucket[size]] = int32(index)
	}
	for index := len(order) - 1; index > 0; index-- {
		parent := a.states[order[index]].link
		if parent >= 0 {
			a.occ[parent] += a.occ[order[index]]
		}
	}
}

// Count 统计子串出现次数（重叠出现分别计数）。
func (a *Automaton) Count(sub string) int64 {
	current := int32(0)
	for index := 0; index < len(sub); index++ {
		next, ok := a.states[current].next[sub[index]]
		if !ok {
			return 0
		}
		current = next
	}
	return a.occ[current]
}

// Distinct 统计不同子串个数（空串不算）。
func (a *Automaton) Distinct() int64 {
	total := int64(0)
	for index := 1; index < len(a.states); index++ {
		total += int64(a.states[index].size - a.states[a.states[index].link].size)
	}
	return total
}

// States 是状态数。
func (a *Automaton) States() int {
	return len(a.states)
}

// LongestCommon 求两个串的最长公共子串：在 first 的自动机上扫 second，
// 失配就沿 link 跳。没有公共子串时返回空串。
func LongestCommon(first, second string) (string, error) {
	automaton := Build(first)
	current, length := int32(0), int32(0)
	best, bestEnd := int32(0), 0
	for index := 0; index < len(second); index++ {
		char := second[index]
		for current != 0 {
			if _, ok := automaton.states[current].next[char]; ok {
				break
			}
			current = automaton.states[current].link
			length = automaton.states[current].size
		}
		if next, ok := automaton.states[current].next[char]; ok {
			current = next
			length++
		} else {
			current = 0
			length = 0
		}
		if length > best {
			best = length
			bestEnd = index + 1
		}
	}
	return second[bestEnd-int(best) : bestEnd], nil
}
