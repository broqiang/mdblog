---
title: "将 vim 打造成 go 语言的 ide"
author: "BroQiang"
github_url: "https://broqiang.com"
created_at: 2019-04-23T00:01:16
updated_at: 2019-04-25T00:01:16
description: "通过对 vim 的配置，可以让 vim 像 ide 一样飞起来。"
tags: ["linux", "go", "vim", "ide"]
---

通过配置，使得 vim 可以看起来像 ide 一样来开发 Go 。

[github 地址](https://github.com/BroQiang/vim-go-ide)

说实在的，如果喜好折腾，喜好自己来配置环境，可以配置一下，如果不是很熟练 `vi/vim`，还是踏踏实实使用 IDE 吧，比如 VSCode（我也会用） 、 GoLand 还都是不错的，不过你要是看到了我的这篇文档，就证明还是想要折腾一下，下面就是我折腾后的结果。

![预览](https://github.com/BroQiang/vim-go-ide/raw/master/vim-go-ide.png)

> 本配置是在 Ubuntu 18.04 下完成的， vim 版本是 `VIM - Vi IMproved 8.0` ，开启 lantern （不考虑网络不通畅的情况）

## 配置 Go 环境

既然是 Go 的开发环境，第一步当前就是准备好 。

### 安装 Go

到 [golang.org](https://golang.org/dl/) 将安装包下载，并配置好环境， 推荐使用二进制版本，下载完成后直接解压缩就可以使用。如果无法访问 go 官网，可以考虑去 [golang.google.cn](https://golang.google.cn/dl) 去下载。

```bash
wget https://dl.google.com/go/go1.12.linux-amd64.tar.gz
sudo tar xzvf go1.12.linux-amd64.tar.gz -C /usr/local/
```

### 配置 PATH 及 GOPATH

创建 Go 的工作空间（目录）

```bash
# 这是默认的位置，也可以按照需求指定到其他目录
mkdir -p $HOME/go/{bin,pkg,src}
```

配置环境变量

用户 vim 创建一个配置文 `vim /etc/profile.d/go.sh` ，写入下面内容

```bash
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

## Vim 基本配置

打开 vimrc 文件 `vim ~/.vimrc` ，写入下面配置

```vim

"==============================================================================
" vim 内置配置
"==============================================================================

" 设置 vimrc 修改保存后立刻生效，不用在重新打开
" 建议配置完成后将这个关闭
autocmd BufWritePost $MYVIMRC source $MYVIMRC

" 关闭兼容模式
set nocompatible

set nu " 设置行号
set cursorline "突出显示当前行
" set cursorcolumn " 突出显示当前列
set showmatch " 显示括号匹配

" tab 缩进
set tabstop=4 " 设置Tab长度为4空格
set shiftwidth=4 " 设置自动缩进长度为4空格
set autoindent " 继承前一行的缩进方式，适用于多行注释

" 定义快捷键的前缀，即<Leader>
let mapleader=";"

" ==== 系统剪切板复制粘贴 ====
" v 模式下复制内容到系统剪切板
vmap <Leader>c "+yy
" n 模式下复制一行到系统剪切板
nmap <Leader>c "+yy
" n 模式下粘贴系统剪切板的内容
nmap <Leader>v "+p

" 开启实时搜索
set incsearch
" 搜索时大小写不敏感
set ignorecase
syntax enable
syntax on                    " 开启文件类型侦测
filetype plugin indent on    " 启用自动补全

" 退出插入模式指定类型的文件自动保存
au InsertLeave *.go,*.sh,*.php write
```

## 插件管理

插件的用途就是可以很方便的管理 vim 的各种插件，快速安装配置以及清除，网上现在的帖子多数都是使用的 [Vundle](https://github.com/VundleVim/Vundle.vim) 这个插件，不过个人觉得这个管理工具在插件安装多了的时候不是很流畅，更喜好使用 [vim-plug](https://github.com/junegunn/vim-plug) 这个插件，两个插件都是很清除的安装文档，这里是介绍 vim-plug 。

### 安装插件

在 Linux 下非常简单，直接通过 curl 下载即可（也可以手动下载，见官方文档）

```bash
curl -fLo ~/.vim/autoload/plug.vim --create-dirs \
    https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
```

### 配置插件

插件的配置也非常简单，只要将所有的插件配置在 `call plug#begin('~/.vim/plugged')` 和 `call plug#end()` 之间即可，常见的插件基本上都可以从 github 中找到，如果 github 找不到的话基本上 vim.org 的脚本都可以在 [vim-script](https://github.com/vim-scripts) 中找到备份

在刚刚的 `~/.vimrc` 下面继续添加插件相关的配置

```vim
" 插件开始的位置
call plug#begin('~/.vim/plugged')

" Shorthand notation; fetches https://github.com/junegunn/vim-easy-align
" 可以快速对齐的插件
Plug 'junegunn/vim-easy-align'

" 用来提供一个导航目录的侧边栏
Plug 'scrooloose/nerdtree'

" 可以使 nerdtree 的 tab 更加友好些
Plug 'jistr/vim-nerdtree-tabs'

" 可以在导航目录中看到 git 版本信息
Plug 'Xuyuanp/nerdtree-git-plugin'

" 查看当前代码文件中的变量和函数列表的插件，
" 可以切换和跳转到代码中对应的变量和函数的位置
" 大纲式导航, Go 需要 https://github.com/jstemmer/gotags 支持
Plug 'majutsushi/tagbar'

" 自动补全括号的插件，包括小括号，中括号，以及花括号
Plug 'jiangmiao/auto-pairs'

" Vim状态栏插件，包括显示行号，列号，文件类型，文件名，以及Git状态
Plug 'vim-airline/vim-airline'

" 有道词典在线翻译
Plug 'ianva/vim-youdao-translater'

" 代码自动完成，安装完插件还需要额外配置才可以使用
Plug 'Valloric/YouCompleteMe'

" 可以在文档中显示 git 信息
Plug 'airblade/vim-gitgutter'


" 下面两个插件要配合使用，可以自动生成代码块
Plug 'SirVer/ultisnips'
Plug 'honza/vim-snippets'

" 可以在 vim 中使用 tab 补全
"Plug 'vim-scripts/SuperTab'

" 可以在 vim 中自动完成
"Plug 'Shougo/neocomplete.vim'


" 配色方案
" colorscheme neodark
Plug 'KeitaNakamura/neodark.vim'
" colorscheme monokai
Plug 'crusoexia/vim-monokai'
" colorscheme github
Plug 'acarapetis/vim-colors-github'
" colorscheme one
Plug 'rakr/vim-one'

" go 主要插件
Plug 'fatih/vim-go', { 'tag': '*' }
" go 中的代码追踪，输入 gd 就可以自动跳转
Plug 'dgryski/vim-godef'

" markdown 插件
Plug 'iamcco/mathjax-support-for-mkdp'
Plug 'iamcco/markdown-preview.vim'

" 插件结束的位置，插件全部放在此行上面
call plug#end()
```

然后输入 `:w` 保存配置，在输入 `:PlugInstall` ，如下：

```vim
:w
:PlugInstall
```

插件会自动下载安装，看见上面显示 Finishing ... Done 的内容，插件安装成功

### 插件删除

如果想要删除插件，只要将不需要的插件注释或者删除，执行 `:PlugClean` 就可以自动清理了

## 插件配置

上面一起安装了很多个插件，有些插件要单独配置，记录到下面

### vim-go

这个是 go 语言支持插件，上面插件完成后还需要安装很多个 Go 的包才能正常工作，在 vim 中执行下面命令：

```vim
:GoInstallBinaries
```

出现 `vim-go: installing finished!` 安装成功，可以使用 Go 包的相关功能了

> 需要注意前面的 PATH 要配置正确，并且已经生效，如果配置正确没有生效，可以注销再登录查看

### YouCompleteMe

这个插件是用来自动完成的，不过需要手动做一些额外的配置

a. 安装以来关系

```bash
sudo apt install build-essential cmake python3-dev
```

b. 编译

```bash
cd ~/.vim/plugged/YouCompleteMe
# 编译，并加入 go 的支持
python3 install.py --go-completer
```

c. 配置和 `SirVer/ultisnips` 冲突的快捷键

```vim
let g:ycm_key_list_select_completion = ['<C-n>', '<space>']
let g:ycm_key_list_previous_completion = ['<C-p>', '<Up>']
let g:SuperTabDefaultCompletionType = '<C-n>'

" better key bindings for UltiSnipsExpandTrigger
let g:UltiSnipsExpandTrigger = "<tab>"
let g:UltiSnipsJumpForwardTrigger = "<tab>"
let g:UltiSnipsJumpBackwardTrigger = "<s-tab>"
```

### 其他插件

其他的插件配置不多，直接看配置文件的注释即可，下面将所有的配置全部贴出来

下面是全部的 `~/.vimrc` 中的配置, 也可以直接下载我配置好的 [vimrc](https://raw.githubusercontent.com/BroQiang/vim-go-ide/master/vimrc) 文件

```vim

"==============================================================================
" vim 内置配置
"==============================================================================

" 设置 vimrc 修改保存后立刻生效，不用在重新打开
" 建议配置完成后将这个关闭，否则配置多了之后会很卡
" autocmd BufWritePost $MYVIMRC source $MYVIMRC

" 关闭兼容模式
set nocompatible

set nu " 设置行号
set cursorline "突出显示当前行
" set cursorcolumn " 突出显示当前列
set showmatch " 显示括号匹配

" tab 缩进
set tabstop=4 " 设置Tab长度为4空格
set shiftwidth=4 " 设置自动缩进长度为4空格
set autoindent " 继承前一行的缩进方式，适用于多行注释

" 定义快捷键的前缀，即<Leader>
let mapleader=";"

" ==== 系统剪切板复制粘贴 ====
" v 模式下复制内容到系统剪切板
vmap <Leader>c "+yy
" n 模式下复制一行到系统剪切板
nmap <Leader>c "+yy
" n 模式下粘贴系统剪切板的内容
nmap <Leader>v "+p

" 开启实时搜索
set incsearch
" 搜索时大小写不敏感
set ignorecase
syntax enable
syntax on                    " 开启文件类型侦测
filetype plugin indent on    " 启用自动补全

" 退出插入模式指定类型的文件自动保存
au InsertLeave *.go,*.sh,*.php write

"==============================================================================
" 插件配置
"==============================================================================

" 插件开始的位置
call plug#begin('~/.vim/plugged')

" Shorthand notation; fetches https://github.com/junegunn/vim-easy-align
" 可以快速对齐的插件
Plug 'junegunn/vim-easy-align'

" 用来提供一个导航目录的侧边栏
Plug 'scrooloose/nerdtree'

" 可以使 nerdtree Tab 标签的名称更友好些
Plug 'jistr/vim-nerdtree-tabs'

" 可以在导航目录中看到 git 版本信息
Plug 'Xuyuanp/nerdtree-git-plugin'

" 查看当前代码文件中的变量和函数列表的插件，
" 可以切换和跳转到代码中对应的变量和函数的位置
" 大纲式导航, Go 需要 https://github.com/jstemmer/gotags 支持
Plug 'majutsushi/tagbar'

" 自动补全括号的插件，包括小括号，中括号，以及花括号
Plug 'jiangmiao/auto-pairs'

" Vim状态栏插件，包括显示行号，列号，文件类型，文件名，以及Git状态
Plug 'vim-airline/vim-airline'

" 有道词典在线翻译
Plug 'ianva/vim-youdao-translater'

" 代码自动完成，安装完插件还需要额外配置才可以使用
Plug 'Valloric/YouCompleteMe'

" 可以在文档中显示 git 信息
Plug 'airblade/vim-gitgutter'


" 下面两个插件要配合使用，可以自动生成代码块
Plug 'SirVer/ultisnips'
Plug 'honza/vim-snippets'

" 配色方案
" colorscheme neodark
Plug 'KeitaNakamura/neodark.vim'
" colorscheme monokai
Plug 'crusoexia/vim-monokai'
" colorscheme github
Plug 'acarapetis/vim-colors-github'
" colorscheme one
Plug 'rakr/vim-one'

" go 主要插件
Plug 'fatih/vim-go', { 'tag': '*' }
" go 中的代码追踪，输入 gd 就可以自动跳转
Plug 'dgryski/vim-godef'

" markdown 插件
Plug 'iamcco/mathjax-support-for-mkdp'
Plug 'iamcco/markdown-preview.vim'

" 插件结束的位置，插件全部放在此行上面
call plug#end()


"==============================================================================
" 主题配色
"==============================================================================

" 开启24bit的颜色，开启这个颜色会更漂亮一些
set termguicolors
" 配色方案, 可以从上面插件安装中的选择一个使用
colorscheme one " 主题
set background=dark " 主题背景 dark-深色; light-浅色


"==============================================================================
" vim-go 插件
"==============================================================================
let g:go_fmt_command = "goimports" " 格式化将默认的 gofmt 替换
let g:go_autodetect_gopath = 1
let g:go_list_type = "quickfix"

let g:go_version_warning = 1
let g:go_highlight_types = 1
let g:go_highlight_fields = 1
let g:go_highlight_functions = 1
let g:go_highlight_function_calls = 1
let g:go_highlight_operators = 1
let g:go_highlight_extra_types = 1
let g:go_highlight_methods = 1
let g:go_highlight_generate_tags = 1

let g:godef_split=2


"==============================================================================
" NERDTree 插件
"==============================================================================

" 打开和关闭NERDTree快捷键
map <F10> :NERDTreeToggle<CR>
" 显示行号
let NERDTreeShowLineNumbers=1
" 打开文件时是否显示目录
let NERDTreeAutoCenter=1
" 是否显示隐藏文件
let NERDTreeShowHidden=0
" 设置宽度
" let NERDTreeWinSize=31
" 忽略一下文件的显示
let NERDTreeIgnore=['\.pyc','\~$','\.swp']
" 打开 vim 文件及显示书签列表
let NERDTreeShowBookmarks=2

" 在终端启动vim时，共享NERDTree
let g:nerdtree_tabs_open_on_console_startup=1


"==============================================================================
"  majutsushi/tagbar 插件
"==============================================================================

" majutsushi/tagbar 插件打开关闭快捷键
nmap <F9> :TagbarToggle<CR>

let g:tagbar_type_go = {
	\ 'ctagstype' : 'go',
	\ 'kinds'     : [
		\ 'p:package',
		\ 'i:imports:1',
		\ 'c:constants',
		\ 'v:variables',
		\ 't:types',
		\ 'n:interfaces',
		\ 'w:fields',
		\ 'e:embedded',
		\ 'm:methods',
		\ 'r:constructor',
		\ 'f:functions'
	\ ],
	\ 'sro' : '.',
	\ 'kind2scope' : {
		\ 't' : 'ctype',
		\ 'n' : 'ntype'
	\ },
	\ 'scope2kind' : {
		\ 'ctype' : 't',
		\ 'ntype' : 'n'
	\ },
	\ 'ctagsbin'  : 'gotags',
	\ 'ctagsargs' : '-sort -silent'
\ }


"==============================================================================
"  nerdtree-git-plugin 插件
"==============================================================================
let g:NERDTreeIndicatorMapCustom = {
    \ "Modified"  : "✹",
    \ "Staged"    : "✚",
    \ "Untracked" : "✭",
    \ "Renamed"   : "➜",
    \ "Unmerged"  : "═",
    \ "Deleted"   : "✖",
    \ "Dirty"     : "✗",
    \ "Clean"     : "✔︎",
    \ 'Ignored'   : '☒',
    \ "Unknown"   : "?"
    \ }

let g:NERDTreeShowIgnoredStatus = 1



"==============================================================================
"  Valloric/YouCompleteMe 插件
"==============================================================================

" make YCM compatible with UltiSnips (using supertab)
let g:ycm_key_list_select_completion = ['<C-n>', '<space>']
let g:ycm_key_list_previous_completion = ['<C-p>', '<Up>']
let g:SuperTabDefaultCompletionType = '<C-n>'

" better key bindings for UltiSnipsExpandTrigger
let g:UltiSnipsExpandTrigger = "<tab>"
let g:UltiSnipsJumpForwardTrigger = "<tab>"
let g:UltiSnipsJumpBackwardTrigger = "<s-tab>"


"==============================================================================
"  其他插件配置
"==============================================================================

" markdwon 的快捷键
map <silent> <F5> <Plug>MarkdownPreview
map <silent> <F6> <Plug>StopMarkdownPreview

" tab 标签页切换快捷键
:nn <Leader>1 1gt
:nn <Leader>2 2gt
:nn <Leader>3 3gt
:nn <Leader>4 4gt
:nn <Leader>5 5gt
:nn <Leader>6 6gt
:nn <Leader>7 7gt
:nn <Leader>8 8gt
:nn <Leader>9 8gt
:nn <Leader>0 :tablast<CR>


"==============================================================================
" GVim 的配置
"==============================================================================
" 如果不使用 GVim ，可以不用配置下面的配置
if has('gui_running')
        colorscheme one
	" 设置启动时窗口的大小
	set lines=999 columns=999 linespace=4

	" 设置字体及大小
        set guifont=Roboto\ Mono\ 13

	set guioptions-=m " 隐藏菜单栏
	set guioptions-=T " 隐藏工具栏
	set guioptions-=L " 隐藏左侧滚动条
	set guioptions-=r " 隐藏右侧滚动条
	set guioptions-=b " 隐藏底部滚动条
        	" 在 gvim 下不会和 terminal 的 alt+数字的快捷键冲突，
	" 所以将 tab 切换配置一份 alt+数字的快捷键
	:nn <M-1> 1gt
	:nn <M-2> 2gt
	:nn <M-3> 3gt
	:nn <M-4> 4gt
	:nn <M-5> 5gt
	:nn <M-6> 6gt
	:nn <M-7> 7gt
	:nn <M-8> 8gt
        :nn <M-9> 9gt
        :nn <M-0> :tablast<CR>
endif
```
