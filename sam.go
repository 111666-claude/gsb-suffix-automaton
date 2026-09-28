// Package suffixauto 实现后缀自动机。
// 缺陷：出现次数只记在最后一个状态上（没有沿 link 树累加），最长公共子串还没实现。
package suffixauto

import "errors"

// ErrNoLCS 表示还没实现。
var ErrNoLCS = errors.New("suffixauto: 最长公共子串还没实现")

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

// Build 建自动机。
func Build(text string) *Automaton {
	automaton := &Automaton{states: []state{{next: map[byte]int32{}, link: -1}}, last: 0}
	for index := 0; index < len(text); index++ {
		automaton.extend(text[index])
	}
	automaton.occ = make([]int64, len(automaton.states))
	automaton.occ[automaton.last] = 1 // 缺陷：只有最后一个状态记了一次
	return automaton
}

func (a *Automaton) extend(char byte) {
	current := int32(len(a.states))
	a.states = append(a.states, state{next: map[byte]int32{}, link: -1,
		size: a.states[a.last].size + 1})
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
	target := a.states[previous].next[char]
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
	for previous >= 0 && a.states[previous].next[char] == target {
		a.states[previous].next[char] = clone
		previous = a.states[previous].link
	}
	a.states[target].link = clone
	a.states[current].link = clone
	a.last = current
}

// Count 统计子串出现次数。缺陷：occ 没沿 link 树累加。
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

// Distinct 统计不同子串个数。
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

// LongestCommon 缺陷：还没实现。
func LongestCommon(first, second string) (string, error) {
	return "", ErrNoLCS
}
