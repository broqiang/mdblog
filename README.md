# Bro Qiang 博客 [https://broqiang.com](https://broqiang.com)

正在开发中，计划用来显示 markdown 文档的静态博客

## 依赖

- [gin](https://github.com/gin-gonic/gin) 引擎及路由

- [BurntSushi/toml](https://github.com/BurntSushi/toml) 用来处理 toml 格式的配置文件

- [russross/blackfriday](https://github.com/russross/blackfriday/releases/tag/v1.5.2) 用来将 markdown 文本转换成 HTML 使用的是 1.5.2 版本。

- [microcosm-cc/bluemonday](https://github.com/microcosm-cc/bluemonday) HTML清理程序， 配合 blackfriday 来处理 markdown 的转换。

- [laravel-mix](https://github.com/JeffreyWay/laravel-mix) 相当于一个 webpack 的包装器，用来管理前端静态资源，以前用 Laravel 框架的时候觉得这个用来管理前端的内容很顺手，试了下，原来可以单独使用， 就用了。

- Jquery + Bootstrap 前端框架，用于页面展示。