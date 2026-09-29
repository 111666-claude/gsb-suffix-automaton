// 命令行入口：建自动机、统计子串出现次数、求最长公共子串。
package main

import (
	"flag"
	"fmt"
	"os"

	suffixauto "example.com/suffixauto"
)

func main() {
	text := flag.String("text", "", "建自动机的文本")
	count := flag.String("count", "", "统计出现次数的子串")
	first := flag.String("a", "", "最长公共子串的第一个串")
	second := flag.String("b", "", "最长公共子串的第二个串")
	flag.Parse()

	if *first != "" || *second != "" {
		lcs, err := suffixauto.LongestCommon(*first, *second)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("lcs=%s len=%d\n", lcs, len(lcs))
		return
	}

	automaton := suffixauto.Build(*text)
	fmt.Printf("states=%d distinct=%d\n", automaton.States(), automaton.Distinct())
	if *count != "" {
		fmt.Printf("count(%s)=%d\n", *count, automaton.Count(*count))
	}
}
