---
title: "grep 详细说明"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2019-05-28T14:52:43
updated_at: 2019-05-28T14:52:43
description: "详细记录了 grep 的使用"
tags: ["linux"]
---

grep 这个命令使用了很久了，一直没有太深入的去了解过， 基本只使用到了最基本的应用，今天详细的
梳理了下， 作为笔记记录下来。

## 格式

通过 `grep --help` 可以看到命令支持的格式及选项

```bash
Usage: grep [OPTION]... PATTERN [FILE]...
```

最常用的用法， 查询指定关键词所在的行：

```bash
$ grep root /etc/passwd
root:x:0:0:root:/root:/bin/bash
```

可以看到， 通过上面的命令就可以查询出 root 所在的行， 并将所在行全部显示。

## OPTION 常用选项

这里只列出了我觉得会常用到的选项， 更多的选项可以通过 help 查看

### -f, --file=FILE

PATTERN 是从文件中读取的

示例：

先在一个文件中（如 p.txt) 中写入下面内容：

```bash
root
nobody
```

然后通过 -f 参数来匹配

```bash
$ grep -f p.txt /etc/passwd
root:x:0:0:root:/root:/bin/bash
nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
```

可以看到， p.txt 中的每一行都会作为一个条件来匹配。

### -i, --ignore-case

忽略大小写， 这个是一个非常常用的选项。

示例：

```bash
$ grep -i port /etc/ssh/sshd_config
#Port 22
#GatewayPorts no
```

### -v, --invert-match

显示没有被匹配到的行

示例：

```bash
$ grep -v sbin /etc/passwd
root:x:0:0:root:/root:/bin/bash
sync:x:4:65534:sync:/bin:/bin/sync
whoopsie:x:112:117::/nonexistent:/bin/false
hplip:x:118:7:HPLIP system user,,,:/var/run/hplip:/bin/false
gnome-initial-setup:x:120:65534::/run/gnome-initial-setup/:/bin/false
gdm:x:121:125:Gnome Display Manager:/var/lib/gdm3:/bin/false
bro:x:1000:1000:BroQiang,,,:/home/bro:/bin/bash
```

可以看到， 只查询出了不包含 root 的行

### -o, --only-matching

只显示匹配的结果

```bash
$ grep -o root.*root /etc/passwd
root:x:0:0:root:/root
```

### 匹配的结果包含上下文所在行

这三个参数的作用是相同的， 分别表示匹配的上下文内容

- -B, --before-context=NUM print NUM lines of leading context

  可以查询出匹配结果所在行和它的上 NUM 行。

  ```bash
  $ grep -B 1 nobody /etc/passwd
  gnats:x:41:41::/var/lib/gnats:/usr/sbin/nologin
  nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin

  ```

- -A, --after-context=NUM print NUM lines of trailing context

  可以查询出匹配结果所在行和它的下 NUM 行。

  ```bash
  $ grep -A 1 nobody /etc/passwd
  nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
  systemd-network:x:100:102::/run/systemd/netif:/usr/sbin/nologin
  ```

- -C, --context=NUM print NUM lines of output context

  可以查询出匹配结果所在行和它的上下各 NUM 行。

  ```bash
  $ grep -C 1 nobody /etc/passwd
  gnats:x:41:41::/var/lib/gnats:/usr/sbin/nologin
  nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
  systemd-network:x:100:102::/run/systemd/netif:/usr/sbin/nologin
  ```

### -q, --quiet, --silent

这个是静默模式， 不显示查询的内容， 只有结果的状态， 一般脚本中会用到这个。

例如查询当前目录下是否存在否个文件， 不存在就创建（这只是个演示， 实际会比这个复杂）：

```bash
#!/bin/sh

ls | grep -q p.txt

if [ $? -eq 0 ]; then
    cat p.txt
fi
```

通过对 grep 添加 -q 参数， 就不会将 grep 出来的结果显示出来了。

## 基本正则表达式匹配

基本的匹配规则， 不是完全的正则表达式。

### 匹配的字符

- `.` 匹配任意单个字符

- `[]` 匹配指定范围的字符，如: [a-z]

- `[^]` 匹配范围以外的字符， 如： [^a-z] ， 匹配所有不是小写字符的行。

在范围匹配中， 还支持一些字符的特殊写法， 类似正则的 `\d` ， `\s` 等（ls 的时候也适用）。

- `[[:digit:]]` 所有数字

- `[[:lower:]]` 所有的小写字母

- `[[:upper:]]` 所有大写字母

- `[[:alpha:]]` 所有字母，包含大小写

- `[[:alnum:]]` 包含数字和字母

- `[[:punct:]]` 所有标点符号

- `[[:space:]]` 所有空白字符

示例：

r 和 t 之间只能包含两个小写字母

```bash
$ grep r[[:alpha:]][[:alpha:]]t /etc/passwd
root:x:0:0:root:/root:/bin/bash
```

### 匹配字符的次数

用在需要匹配的字符后面，表示重复的次数

- `*` 表示匹配字符的 0 次或多次

  ```bash
  $ grep ro\*t /etc/passwd
  root:x:0:0:root:/root:/bin/bash
  gnats:x:41:41:Gnats Bug-Reporting System (admin):/var/lib/gnats:/usr/sbin/nologin
  rtkit:x:109:114:RealtimeKit,,,:/proc:/usr/sbin/nologin
  ```

- `.*` 匹配任意长度的任意字符

  ```bash
  $ grep ro.*t /etc/passwd
  root:x:0:0:root:/root:/bin/bash
  ```

  贪婪模式匹配， 匹配最长的。

- `\?` 匹配前面字符 0 次或 1 次

  ```bash
  $ grep "roo\?t" /etc/passwd
  root:x:0:0:root:/root:/bin/bash
  ```

- `\+` 匹配前面字符 1 次或多次

  ```bash
  $ grep "ro\+t" /etc/passwd
  root:x:0:0:root:/root:/bin/bash
  ```

- `\{m\}` 匹配前面字符 m 次， m 是一个正整数

  ```bash
  $ grep "ro\{2\}t" /etc/passwd
  root:x:0:0:root:/root:/bin/bash
  ```

- `\{m,n\}` 匹配前面字符的范围， 最少 m 次， 最多 n 次

  ```bash
  $ grep "ro\{1,3\}" /etc/passwd
  root:x:0:0:root:/root:/bin/bash
  proxy:x:13:13:proxy:/bin:/usr/sbin/nologin
  rtkit:x:109:114:RealtimeKit,,,:/proc:/usr/sbin/nologin
  bro:x:1000:1000:BroQiang,,,:/home/bro:/bin/bash
  ```

### 搜索位置匹配

- `^PATTERN` 搜索的 PATTERN 只能出现在行首

  ```bash
  $ grep ^root /etc/passwd
  root:x:0:0:root:/root:/bin/bash
  ```

- `PATTERN$` 搜索的 PATTERN 只能出现在行尾

  ```bash
  $ grep "bash$" /etc/passwd
  root:x:0:0:root:/root:/bin/bash
  bro:x:1000:1000:BroQiang,,,:/home/bro:/bin/bash
  ```

- `^PATTERN$` 搜索的 PATTERN 匹配整行

  ```bash
  $ grep "^root:x:0:0:root:/root:/bin/bash$" /etc/passwd
  root:x:0:0:root:/root:/bin/bash
  ```

- `^$` 匹配空行， 不包含空白符的空行

- `^[[:space::]]*$` 可以有任意空白的空行

- `\<root\>` 只能出现在单词的词首， 词尾， `\<` `\>` 可以单独使用

  ```bash
  $ grep "\<root" /etc/passwd
  root:x:0:0:root:/root:/bin/bash
  ```

- `\(字符\)` 将一个或多个字符作为一个整理处理（分组）

  ```bash
  $ grep "\(ro\)\+" /etc/passwd
  root:x:0:0:root:/root:/bin/bash
  proxy:x:13:13:proxy:/bin:/usr/sbin/nologin
  rtkit:x:109:114:RealtimeKit,,,:/proc:/usr/sbin/nologin
  bro:x:1000:1000:BroQiang,,,:/home/bro:/bin/bash
  ```

## egrep 扩展正则表达式

可以直接通过 egrep 命令来搜索匹配，默认就是支持扩展正则表达式， 也可以通过 `grep -E` 参数，
使用 grep 命令来支持扩展正则表达式。

### egrep 匹配的字符

egrep 字符匹配和 grep 是类似的,基本相同。 直接看前面的 grep 的匹配字符即可。

### egrep 次数匹配

egrep 的次数匹配和 grep 也是类似的， 不过可以省略掉啰嗦的转译字符。

- `*` 前面的字符重复任意次数， 0 次或 多次

- `?` 前面的字符重复 0 次或 1 次

- `+` 前面的字符重复 1 次或 多次

- `{m}` 前面的字符重复的次数， m 是正整数

- `{m,n}` 之前前面字符最少重复 m 次， 最多重复 n 次

- 行首、行尾、词首、词尾 和 grep 是相同的， 见上方

- `()` 分组的用法也和 grep 相同， 只是不需要转译字符 `\` 了

- `|` 或的意思

  查询出 root 或 nobody 所在行。

  ```bash
  $ grep -E "root|nobody" /etc/passwd
  root:x:0:0:root:/root:/bin/bash
  nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
  ```

  查询出 ma 后面是 n 或 il 开头的行。

  ```bash
  $ grep -E "^ma(n|il)" /etc/passwd
  man:x:6:12:man:/var/cache/man:/usr/sbin/nologin
  mail:x:8:8:mail:/var/mail:/usr/sbin/nologin
  ```
