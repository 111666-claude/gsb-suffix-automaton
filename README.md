# suffix-automaton

后缀自动机（出现次数、不同子串数、最长公共子串），只用 Go 标准库。

```
go test ./...
go run ./cmd/suffix-automaton --text banana --count an
go run ./cmd/suffix-automaton --a abcd --b bcde
```

## 口径（README 为准）

- **出现次数要沿 link 树累加**：建好自动机后，把每个状态的计数按 `size` 从大到小加到 `link` 上。
  `--text banana --count an` 按口径是 `count(an)=2`，现在只在最后一个状态记了一次，给 `count(an)=1`。
- **最长公共子串要实现**：`--a abcd --b bcde` 按口径是 `lcs=bcd len=3`，现在直接报「最长公共子串还没实现」。
- 建机要求：clone 状态的 `size` 取 `states[previous].size + 1`、`link` 继承原状态并复制转移表；
  克隆后指向旧状态的转移要改指克隆。
- 不变量：不同子串数等于 `Σ size(v) - size(link(v))`（空串不算，`aaa` 给 3、`ab` 给 3）；
  `Count` 与「重叠计数」的语义一致（`aaa` 里 `aa` 出现 2 次）；同一输入重复建机结果相同。
- 规模：1 MB 文本建机要在 1 秒内完成，内存 O(状态数)（≤ 2n 个状态）；不许逐个子串暴力统计。

## 现在的行为

- `Build` 只做 `occ[last] = 1`：没有沿 link 树累加，除最后一次出现外都算不到。
- `LongestCommon` 直接返回 `ErrNoLCS`。
- `distinct` 与状态数是对的，所以 `aaa` 给 3、`ab` 给 3。

## 计划新增的能力

```
--text 文本      建自动机
--count 子串     统计出现次数
--a / --b        求两个串的最长公共子串
```

## 输出

```
states=8 distinct=15
count(an)=2
lcs=bcd len=3
```

## 目录

```
sam.go                 后缀自动机
cmd/suffix-automaton   命令行入口
sam_test.go            go test 用例
```
