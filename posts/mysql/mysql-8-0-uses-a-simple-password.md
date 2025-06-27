---
title: "mysql 8.0 使用简单密码"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2018-06-13T13:57:35
updated_at: 2018-06-13T13:57:35
description: "修改 mysql 8.0 的默认密码策略，允许设置简单密码"
tags: ["linux", "mysql"]
---

最近测试环境换上了 mysql 8.0.1 ，发现默认只允许使用非常复杂的密码： 大小写英文/数字/特殊字符。 这个对于个人开发环境来说就有点啰嗦了，毕竟要经常使用命令行来输入，每次都输入这么长的密码有点麻烦了。

## 最简单的配置方式

直接在 my.cnf 配置文件中 [mysqld] 部分加入下面参数，然后重启 mysqld 即可。

```bash
validate_password.policy = 0
validate_password.mixed_case_count = 0
validate_password.number_count = 0
validate_password.special_char_count = 0
validate_password.length = 0
```

## 解释说明

参考： [官方文档](https://dev.mysql.com/doc/refman/8.0/en/validate-password-options-variables.html)

### validate_password.policy

可以配置密码的复杂度，可以配置的级别：

- 0 or LOW
- 1 or MEDIUM
- 2 or STRONG

### validate_password.length

最终密码的长度，允许为 0 ，但是要注意这里有个坑，`validate_password.length` 的长度要大于 `validate_password.mixed_case_count` + `validate_password.number_count` + `validate_password.special_char_count` 的和。 例如默认这 3 个 参数的长度都是 1， 所以密码长度最小也只能是 4， 即使配置成了 1 或者 0 ，最终它也会自动变成 4 。 要是想使用 0 的长度，需要将另外三个参数也配置成 0 。
