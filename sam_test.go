package suffixauto

import (
	"testing"
	"time"
)

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

func bruteForceCount(text, sub string) int64 {
	if len(sub) == 0 || len(sub) > len(text) {
		return 0
	}
	total := int64(0)
	for start := 0; start+len(sub) <= len(text); start++ {
		if text[start:start+len(sub)] == sub {
			total++
		}
	}
	return total
}

func bruteForceDistinct(text string) int64 {
	seen := map[string]struct{}{}
	for start := 0; start < len(text); start++ {
		for end := start + 1; end <= len(text); end++ {
			seen[text[start:end]] = struct{}{}
		}
	}
	return int64(len(seen))
}

func TestCountBanana(t *testing.T) {
	automaton := Build("banana")
	cases := map[string]int64{"an": 2, "na": 2, "a": 3, "n": 2, "banana": 1, "ana": 2}
	for sub, want := range cases {
		if got := automaton.Count(sub); got != want {
			t.Fatalf("count(%s)=%d，想要 %d", sub, got, want)
		}
	}
}

func TestOverlappingCount(t *testing.T) {
	automaton := Build("aaa")
	if got := automaton.Count("aa"); got != 2 {
		t.Fatalf("aaa 里 aa 应该出现 2 次：%d", got)
	}
	if got := automaton.Count("a"); got != 3 {
		t.Fatalf("aaa 里 a 应该出现 3 次：%d", got)
	}
	if got := automaton.Count("aaa"); got != 1 {
		t.Fatalf("aaa 里 aaa 应该出现 1 次：%d", got)
	}
}

func TestCloneState(t *testing.T) {
	text := "abb"
	automaton := Build(text)
	if automaton.States() <= len(text)+1 {
		t.Fatalf("%q 建机应该产生 clone 状态（状态数 > n+1）：%d", text, automaton.States())
	}
	for _, sub := range []string{"a", "b", "ab", "bb", "abb"} {
		if got, want := automaton.Count(sub), bruteForceCount(text, sub); got != want {
			t.Fatalf("clone 后 count(%s)=%d，想要 %d", sub, got, want)
		}
	}
	if got, want := automaton.Distinct(), bruteForceDistinct(text); got != want {
		t.Fatalf("clone 后 distinct=%d，想要 %d", got, want)
	}
}

func TestLongestCommon(t *testing.T) {
	lcs, err := LongestCommon("abcd", "bcde")
	if err != nil || lcs != "bcd" {
		t.Fatalf("abcd/bcde 的 lcs 应该是 bcd：%q err=%v", lcs, err)
	}
	lcs, err = LongestCommon("banana", "ananas")
	if err != nil || lcs != "anana" || len(lcs) != 5 {
		t.Fatalf("banana/ananas 的 lcs 应该是 anana：%q err=%v", lcs, err)
	}
	lcs, err = LongestCommon("abc", "abc")
	if err != nil || lcs != "abc" {
		t.Fatalf("相同串的 lcs 应该是整串：%q err=%v", lcs, err)
	}
}

func TestLongestCommonEdge(t *testing.T) {
	for _, pair := range [][2]string{{"abc", "xyz"}, {"", "abc"}, {"abc", ""}, {"", ""}} {
		lcs, err := LongestCommon(pair[0], pair[1])
		if err != nil || lcs != "" {
			t.Fatalf("%q/%q 应该没有公共子串：%q err=%v", pair[0], pair[1], lcs, err)
		}
	}
}

func TestDeterministicBuild(t *testing.T) {
	text := "mississippi"
	first, second := Build(text), Build(text)
	if first.States() != second.States() || first.Distinct() != second.Distinct() {
		t.Fatal("同一输入重复建机结果应该相同")
	}
	for _, sub := range []string{"s", "ss", "issi", "ppi", "mississippi"} {
		if first.Count(sub) != second.Count(sub) {
			t.Fatalf("count(%s) 两次建机不一致", sub)
		}
	}
}

func TestRandomAgainstBruteForce(t *testing.T) {
	rng := uint64(20260929)
	next := func(mod int) int {
		rng = rng*6364136223846793005 + 1442695040888963407
		return int(rng>>33) % mod
	}
	for round := 0; round < 200; round++ {
		length := 1 + next(40)
		data := make([]byte, length)
		for index := range data {
			data[index] = byte('a' + next(3))
		}
		text := string(data)
		automaton := Build(text)
		for query := 0; query < 5; query++ {
			var sub string
			if next(2) == 0 && length > 1 {
				start := next(length)
				end := start + 1 + next(length-start)
				sub = text[start:end]
			} else {
				subLength := 1 + next(5)
				buf := make([]byte, subLength)
				for index := range buf {
					buf[index] = byte('a' + next(3))
				}
				sub = string(buf)
			}
			if got, want := automaton.Count(sub), bruteForceCount(text, sub); got != want {
				t.Fatalf("text=%q count(%s)=%d，暴力=%d", text, sub, got, want)
			}
		}
		if got, want := automaton.Distinct(), bruteForceDistinct(text); got != want {
			t.Fatalf("text=%q distinct=%d，暴力=%d", text, got, want)
		}
	}
}

func TestBuildOneMBWithinOneSecond(t *testing.T) {
	if testing.Short() {
		t.Skip("短模式跳过性能测试")
	}
	rng := uint64(7)
	data := make([]byte, 1<<20)
	for index := range data {
		rng = rng*6364136223846793005 + 1442695040888963407
		data[index] = byte('a' + (rng>>33)&1)
	}
	text := string(data)
	start := time.Now()
	automaton := Build(text)
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("1 MB 建机超过 1 秒：%v", elapsed)
	}
	if automaton.States() > 2*len(text) {
		t.Fatalf("状态数不该超过 2n：%d", automaton.States())
	}
}
