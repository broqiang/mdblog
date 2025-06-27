---
title: "将 Caps Lock 改建成智慧的 Control 和 Esc"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2019-04-22T16:51:54
updated_at: 2019-04-22T16:51:54
description: "在些代码的时候，有没有感觉左手小拇指有些疼，有没有发现大小写锁定键这么好的位置从来没有使用过，可以将这个优势按键替换成单击触发 ESC， 和其他按键组合就变成 Ctrl 键。"
tags: ["Linux"]
---

参考： [Smart Caps Lock: Remap to Control AND Escape](https://gist.github.com/tanyuan/55bca522bf50363ae4573d4bdcf06e2e)

在 Ubuntu 下非常简单，只需要三步

1. 安装 [xcape](https://github.com/alols/xcape)

```bash
sudo apt-get install xcape
```

1. 修改按键映射

```bash
# 修改大小写锁定为 ctrl
setxkbmap -option ctrl:nocaps

# 将短按左 ctrl 改为 ESC
xcape -e 'Control_L=Escape'
```

1. 将上面的命令加入到自动启动中，看个人习惯，如： ~/.bashrc 中
