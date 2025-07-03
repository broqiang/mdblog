---
title: "Go 语言学习路线指南"
author: "BroQiang"
created_at: 2018-05-28T01:09:25
updated_at: 2019-05-21T01:15:57
---

不知道是不是有同学打算开始学习 Golang，确不知道如何开始，至少我开始学习的时候就是这种感觉，为了这个，我查询了很多的帖子和问答。网上的 Golang 资料虽然不多，但是也不少，这个我的一个学习路线，从简单内容开始，可以作为参考。

## 第一步 Go 语言之旅

这个是一个官方的入门教程，或者说只是一个大概了解的教程，只介绍了一些简单的东西，并且没有太多的说明。不过这个教程支持在线执行代码，还是很不错的，这个时候你都不需要有本地的开发环境。不用想太多，现在就开始，把这个教程从头到尾看一遍，练习一遍，遇到不明白的地方也不要纠结，继续向后看就行了。

官方： [https://tour.golang.org](https://tour.golang.org)

中文网： [http://go-tour-zh.appspot.com](http://go-tour-zh.appspot.com/welcome/1)

## 第二步 开发环境

> 这里也可以忽略不看，因为每一个教程都会介绍怎么配置环境

#### 操作系统

个人推荐使用 Linux，可以使用 Ubuntu 或者 Fedora ，如果条件允许（不差钱） 也可以使用 Mac，当然使用 Windows 也是可以的，慢慢的就会知道 Windows 下做开发的纠结了。

#### 开发环境

Go 的安装非常的简单，没有太多的依赖，如果是 Linux 下安装基本上下载一个二进制包，解压配置上一个环境变量、GOROOT 既可以了，具体的可以查看官方的安装方法： [官网安装文档](https://golang.org/doc/install) 、 [中文安装文档](http://docscn.studygolang.com/doc/install)

#### 开发工具

可以选择一个自己喜欢的，个人建议要用个 IDE，我使用过 vim 、Sublime、Intellji idea ，最后发现还是 IDE 比较方便，尤其是代码追踪断点等。这个不纠结那种好，有人和我说 Sublime 和 vim 安装上插件也都可以，但是个人不推荐（我以前是 Sublime 重度用户，PHP 一直都在 Sublime 下开发）。

主流的文本编辑器及 IDE 的配置 [官方](https://github.com/golang/go/wiki/IDEsAndTextEditorPlugins) 都有介绍，选一个自己喜欢的就可以了。

## 第三步 看一套视频

有人可能喜欢看视频，有人可能喜欢看文档，这个根据个人爱好去选择，个人建议要看一套视频并且只看一套就够了，毕竟看视频的效率还是比较低的，看完视频一些基础的知识点就可以掌握，并且会知道一些专有名字的读法。我以前学 PHP 的时候就从来没看过视频，导致很多名词的发音都是错的，经常被人嘲笑……，当然如果英文非常的好的同学就不用纠结了。

网上 Golang 的视频不是很多，不过也有好多套，推荐 [无闻的 Go 编程基础](https://learnku.com/docs/go-fundamental-programming)，这个是被 [golangcaff.com](golangcaff.com) 的 [Summer
](https://golangcaff.com/users/1) 整理优化过的，看起来的效果会比一些其他网站好一些。

## 第四步 看一篇教程

教程也有很多，看个人的喜好吧，推荐看 [Go 入门指南](https://learnku.com/docs/the-way-to-go) ，这个也是由 [无闻](https://github.com/Unknwon) [翻译](https://github.com/Unknwon/the-way-to-go_ZH_CN) 的 [The Way to Go](https://sites.google.com/site/thewaytogo2012/) ，不过社区的版本对排版进行了优化并加入了后面没有翻译完的部分。

## 第五步 将标准库全部都过一遍

至少要叫常用的全都看一遍，如 strings / strconv / http 等，如果有能力可以将它们都记住，如果记忆力不太好（像我这样）至少也要知道有什么，用到的时候通过手册可以快速找到。

官方标准库： [https://golang.org/pkg/](https://golang.org/pkg/)

中文版的标准库： [https://studygolang.com/static/pkgdoc/main.html](https://studygolang.com/static/pkgdoc/main.html)

极力推荐 [https://github.com/astaxie/gopkg](https://github.com/astaxie/gopkg) ，可以在学习的时候看这个，有关于标准库的详细说明和示例，学习起来会容易一些，等全都明白了要使用的时候可以去查看上面的文档。

更新：

又发现了一个不错的学习标准库的资料： [《Go 语言标准库》The Golang Standard Library by Example](https://books.studygolang.com/The-Golang-Standard-Library-by-Example/) ，有点小遗憾就是不是很全，个别的包没有完成，不过 astaxie 的那个也不全，可以互相参考着看。

## 完成

到这个时候，你肯定已经入门了，剩下就开始写自己的东西吧，比如写一个博客，或者去学习一个框架，不知道怎么继续去深造的话就去招聘网站上看看自己喜欢的企业都要求什么，招聘要求会什么就去学什么。

## 2019-05-21 追加

一直没看这篇文档，也就忘了更新了， 今天更新一下

因为我原来是 PHP 程序员， 所以就从 web 开发入手的。

后来我又翻译了官方的： [Writing Web Applications](https://golang.org/doc/articles/wiki/) 这篇文档， 翻译： [Go 编写 Web 应用](https://broqiang.com/posts/writing-web-applications)

看了下 [httprouter](https://github.com/julienschmidt/httprouter) 的源码， 模仿它自己做了下路由的实现， [httprouter](https://broqiang.com/posts/httprouter-source-code-analysis) 。

看了下 [gin](https://github.com/gin-gonic/gin) 的源码， 并基于 gin 做了个博客 [broqiang.com](https://broqiang.com) ，源码 [https://github.com/BroQiang/mdblog](https://github.com/BroQiang/mdblog)

最近有空就会刷下 [https://leetcode-cn.com/](https://leetcode-cn.com/) ，用 go 做一遍实现，不过才刚刚刷了几十道题。

> 我也只能算是刚刚入门， 不是高手， 这是我学习的一个路线， 给新手一个借鉴。

本文来自 [https://broqiang.com](https://broqiang.com/posts/38) 没有版权限制，随意转载
