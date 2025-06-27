---
title: "Go 不使用加号实现 A+B"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2018-06-17T20:32:41
updated_at: 2018-06-17T20:32:41
description: "Golang 不使用加号实现 A+B，此问题来自： https://www.lintcode.com/problem/a-b-problem/description"
tags: ["go"]
---

Go 不使用加号实现 A+B，此问题来自： [lintcode](https://www.lintcode.com/problem/a-b-problem/description)

## 问题描述：

说明
a 和 b 都是 32 位 整数么？

是的
我可以使用位运算符么？

当然可以
样例
如果 a=1 并且 b=2，返回 3。

挑战
显然你可以直接 return a + b，但是你是否可以挑战一下不这样做？（不使用++等算数运算符）

## 代码实现

```go
package main

import "fmt"

func main() {
	i := aplusb(5, 17)
	fmt.Println(i)
}

func aplusb(a int, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}

	x1 := a ^ b
	x2 := (a & b) << 1

	return aplusb(x1, x2)
}

```
