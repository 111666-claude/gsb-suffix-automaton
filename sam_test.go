package suffixauto

import "testing"

func TestCountMissAndHit(t *testing.T) {
	automaton := Build("banana")
	if automaton.Count("zz") != 0 {
		t.Fatal("不存在的子串应该是 0")
	}
	if automaton.Count("banana") == 0 {
		t.Fatal("整串应该至少出现一次")
	}
}

func TestDistinctCounts(t *testing.T) {
	if got := Build("aaa").Distinct(); got != 3 {
		t.Fatalf("aaa 的不同子串应该是 3：%d", got)
	}
	if got := Build("ab").Distinct(); got != 3 {
		t.Fatalf("ab 的不同子串应该是 3：%d", got)
	}
}

func TestEmptyText(t *testing.T) {
	automaton := Build("")
	if automaton.Count("a") != 0 || automaton.Distinct() != 0 {
		t.Fatal("空串没有子串")
	}
}

func TestSingleCharacter(t *testing.T) {
	automaton := Build("x")
	if automaton.States() < 2 {
		t.Fatalf("状态数不对：%d", automaton.States())
	}
	if automaton.Count("x") != 1 {
		t.Fatalf("单字符应该出现一次：%d", automaton.Count("x"))
	}
}

func TestStatesAreBounded(t *testing.T) {
	text := "mississippi"
	automaton := Build(text)
	if automaton.States() > 2*len(text) {
		t.Fatalf("状态数不该超过 2n：%d", automaton.States())
	}
}
