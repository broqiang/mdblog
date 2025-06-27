---
title: "Ubuntu 17.10 安装 MShybrid 显卡"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2018-01-13T01:50:58
updated_at: 2018-01-13T01:50:58
description: "因为个人电脑雷神 st pro 是MShybrid的核显，记录下显卡驱动的安装过程。"
tags: ["linux"]
---

因为 雷神 st pro 是 MShybrid 的核显，ubuntu 对 MShybrid 的核显模式支持不好，所以进入到安装界面就会黑屏或卡死，使用下面方式进行安装。

## BIOS 中禁用集显

因为 MSHybrid 只能禁用集显，然后独显单独工作，其实正常这样也没什么问题，但是会增加笔记本的热度，导致风扇狂转。但是安装系统又只能先单独使用独显，等安装完成后再调回来。

进 bios，禁用核显，找到 MSHybrid 改为 DISCRETE

## 更新显卡驱动

直接使用软件和更新->附加驱动

- 选择 nvidia 的专有驱动进行更新

- Intel CPUs 的开源驱动也更新

如果在这里看不到的话，可以先将操作系统更新到最新，也非常建议先将操作系统更新到最新，否则更新完系统后下面的步骤还是要再操作一遍

```bash
sudo apt update
sudo apt upgrade -y
reboot
```

## 修改 /etc/default/grub

```bash
sudo vim /etc/default/grub

# 在 GRUB_CMDLINE_LINUX="" 上面添加一行
GRUB_GFXPAYLOAD_LINUX=1920*1080
# 这个 1920*1080 个人觉得和电脑的分辨率有关，这里 雷神 st pro 没有问题，其他品牌和型号没试过，可以参考自己的去配置
```

将上面的内容修改后保存，然后更新下引导

```bash
sudo update-grub
```

> 需要注意：这个 update-grub 在每次更新完内核都要执行一次，下次更新完后如果系统不能进入，可以将这个步骤从新走一遍即可

更新完成后重启系统，将显卡的 DISCRETE 再改回 MSHYbrid
