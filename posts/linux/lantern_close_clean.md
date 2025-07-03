---
title: "Linux 下 Lantern 关闭后不清除代理"
author: "BroQiang"
created_at: 2018-01-12T01:53:30
updated_at: 2018-01-12T01:53:30
---

个人觉得这个是一个 bug，Issues 中也有人提出了，不过官方一直没有解决，平时使用起来比较麻烦，只能自己写一个脚本来解决一下，虽然也不是很优雅，不过比起每次都手动去配置要方便一些。

## 在家目录下创建一个 bin 目录

```bash
mkdir ~/bin
```

## 然后在 bin 目录下新建一个 lantern_clear （名字随意起)，写入下面内容：

```bash
#!/bin/bash
#

if ps aux | grep 'lantern$' > /dev/null ; then
    killall -9 lantern
fi

gsettings set org.gnome.system.proxy mode none
```

## 给脚本添加执行权限

```bash
chmod +x ~/bin/lanternclear
```

然后执行脚本就可以将 lantern 结束了

终端直接执行或者 `Alt+F2` 执行 `lantern_clear`
