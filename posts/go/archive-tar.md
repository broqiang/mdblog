---
title: "archive/tar 实现打包压缩及解压"
author: "BroQiang"
created_at: 2018-06-28T20:07:39
updated_at: 2018-06-28T20:07:39
---

这个包比较简单，就是将文件进行打包和解包，要是熟悉 Linux 下的 tar 命令这个就很好理解了。 主要是通过 tar.Reader 读取 tar 包，通过 tar.Writer 写入 tar
包，在写入的过程中再设置一下头，详细的过程以示例的方式进行展示，可以查看代码里面的注释。

参考：

- [标准库 tar 中文文档](https://studygolang.com/static/pkgdoc/pkg/archive_tar.htm)

- [标准库 tar 官方文档](https://golang.org/pkg/archive/tar/)

## 单个文件操作

这个非常简单，就是读取一个文件，进行打包及解包操作即可。

### 单个文件打包

从 /etc/passwd 下复制了一个 passwd 文件到当前目录下，用来做压缩测试。什么文件都是可以的，自己随意写一个也行。这里的示例主要为了说明 tar ，没有处理路径，所以过程全部假设是在当前目录下执行。

```bash
cp /etc/passwd .
```

关于文件的打包直接查看示例代码，已经在示例代码中做了详细的注释。

示例代码（ [pack_single_file.go](example/pack_single_file.go) ）：

```go
package main

import (
	"os"
	"log"
	"archive/tar"
	"fmt"
	"io"
)

func main() {
	// 准备打包的源文件
	var srcFile = "passwd"
	// 打包后的文件
	var desFile = fmt.Sprintf("%s.tar",srcFile)

	// 需要注意文件的打开即关闭的顺序，因为 defer 是后入先出，所以关闭顺序很重要
	// 第一次写这个示例的时候就没注意，导致写完的 tar 包不完整

	// ###### 第 1 步，先准备好一个 tar.Writer 结构，然后再向里面写入内容。 ######
	// 创建一个文件，用来保存打包后的 passwd.tar 文件
	fw, err := os.Create(desFile)
	ErrPrintln(err)
	defer fw.Close()

	// 通过 fw 创建一个 tar.Writer
	tw := tar.NewWriter(fw)
	// 这里不要忘记关闭，如果不能成功关闭会造成 tar 包不完整
	// 所以这里在关闭的同时进行判断，可以清楚的知道是否成功关闭
	defer func() {
		if err := tw.Close(); err != nil {
			ErrPrintln(err)
		}
	}()

	// ###### 第 2 步，处理文件信息，也就是 tar.Header 相关的 ######
	// tar 包共有两部分内容：文件信息和文件数据
	// 通过 Stat 获取 FileInfo，然后通过 FileInfoHeader 得到 hdr tar.*Header
	fi, err := os.Stat(srcFile)
	ErrPrintln(err)
	hdr, err := tar.FileInfoHeader(fi, "")
	// 将 tar 的文件信息 hdr 写入到 tw
	err = tw.WriteHeader(hdr)
	ErrPrintln(err)

	// 将文件数据写入
	// 打开准备写入的文件
	fr, err := os.Open(srcFile)
	ErrPrintln(err)
	defer fr.Close()

	written, err := io.Copy(tw, fr)
	ErrPrintln(err)

	log.Printf("共写入了 %d 个字符的数据\n",written)
}

// 定义一个用来打印的函数，少写点代码，因为要处理很多次的 err
// 后面其他示例还会继续使用这个函数，就不单独再写，望看到此函数了解
func ErrPrintln(err error)  {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
```

#### 单个文件解包

这个也很简单，基本上将上面过程反过来，只需要处理 tar.Reader 即可，详细的描述见示例。

这里就用刚刚打包的 `passwd.tar` 文件做示例，如果怕结果看不出效果，可以将之前用的 passwd 源文件删除。

```bash
rm passwd
```

示例代码（ [unpack_single_file.go](example/unpack_single_file.go) ）：

```go
package main

import (
	"os"
	"archive/tar"
	"io"
	"log"
)

func main() {

	var srcFile = "passwd.tar"

	// 将 tar 包打开
	fr, err := os.Open(srcFile)
	ErrPrintln(err)
	defer fr.Close()

	// 通过 fr 创建一个 tar.*Reader 结构，然后将 tr 遍历，并将数据保存到磁盘中
	tr := tar.NewReader(fr)

	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next(){
		// 处理 err ！= nil 的情况
		ErrPrintln(err)
		// 获取文件信息
		fi := hdr.FileInfo()

		// 创建一个空文件，用来写入解包后的数据
		fw, err := os.Create(fi.Name())
		ErrPrintln(err)

		// 将 tr 写入到 fw
		n, err := io.Copy(fw, tr)
		ErrPrintln(err)
		log.Printf("解包： %s 到 %s ，共处理了 %d 个字符的数据。", srcFile,fi.Name(),n)

		// 设置文件权限，这样可以保证和原始文件权限相同，如果不设置，会根据当前系统的 umask 来设置。
		os.Chmod(fi.Name(),fi.Mode().Perm())

		// 注意，因为是在循环中，所以就没有使用 defer 关闭文件
		// 如果想使用 defer 的话，可以将文件写入的步骤单独封装在一个函数中即可
		fw.Close()
	}
}

func ErrPrintln(err error){
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
```

## 操作整个目录

我们实际中 tar 很少会去打包单个文件，一般都是打包整个目录，并且打包的时候通过 gzip 或者 bzip2 压缩。

如果要打包整个目录，可以通过递归的方式来实现。这里只演示了 gzip 方式压缩，这个实现非常简单，只需要在 fw 和 tw 之前加上一层压缩即可，详情见示例代码。

为了测试打包整个目录，复制了一个 log 目录到当前路径下。什么目录和文件都可以，只是因为这个里面内容比较多，就拿这个来做测试了。

```bash
# 出现没有权限的错误不用管它，复制过来多少是多少吧
cp -r /var/log/ .
```

详细的操作会在注释中说明，不过在之前单文件中出现过的步骤不再注释。

### 打包压缩

示例代码（ [targz.go](example/targz.go) ）：

```go
package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 修改日志格式，显示出错代码的所在行，方便调试，实际项目中一般不记录这个。

	var src = "apt"
	var dst = fmt.Sprintf("%s.tar.gz", src)

	// 将步骤写入了一个函数中，这样处理错误方便一些
	if err := Tar(src, dst); err != nil {
		log.Fatalln(err)
	}
}

func Tar(src, dst string) (err error) {
	// 创建文件
	fw, err := os.Create(dst)
	if err != nil {
		return
	}
	defer fw.Close()

	// 将 tar 包使用 gzip 压缩，其实添加压缩功能很简单，
	// 只需要在 fw 和 tw 之前加上一层压缩就行了，和 Linux 的管道的感觉类似
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// 创建 Tar.Writer 结构
	tw := tar.NewWriter(gw)
	// 如果需要启用 gzip 将上面代码注释，换成下面的

	defer tw.Close()

	// 下面就该开始处理数据了，这里的思路就是递归处理目录及目录下的所有文件和目录
	// 这里可以自己写个递归来处理，不过 Golang 提供了 filepath.Walk 函数，可以很方便的做这个事情
	// 直接将这个函数的处理结果返回就行，需要传给它一个源文件或目录，它就可以自己去处理
	// 我们就只需要去实现我们自己的 打包逻辑即可，不需要再去路径相关的事情
	return filepath.Walk(src, func(fileName string, fi os.FileInfo, err error) error {
		// 因为这个闭包会返回个 error ，所以先要处理一下这个
		if err != nil {
			return err
		}

		// 这里就不需要我们自己再 os.Stat 了，它已经做好了，我们直接使用 fi 即可
		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}
		// 这里需要处理下 hdr 中的 Name，因为默认文件的名字是不带路径的，
		// 打包之后所有文件就会堆在一起，这样就破坏了原本的目录结果
		// 例如： 将原本 hdr.Name 的 syslog 替换程 log/syslog
		// 这个其实也很简单，回调函数的 fileName 字段给我们返回来的就是完整路径的 log/syslog
		// strings.TrimPrefix 将 fileName 的最左侧的 / 去掉，
		// 熟悉 Linux 的都知道为什么要去掉这个
		hdr.Name = strings.TrimPrefix(fileName, string(filepath.Separator))

		// 写入文件信息
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		// 判断下文件是否是标准文件，如果不是就不处理了，
		// 如： 目录，这里就只记录了文件信息，不会执行下面的 copy
		if !fi.Mode().IsRegular() {
			return nil
		}

		// 打开文件
		fr, err := os.Open(fileName)
		defer fr.Close()
		if err != nil {
			return err
		}

		// copy 文件数据到 tw
		n, err := io.Copy(tw, fr)
		if err != nil {
			return err
		}

		// 记录下过程，这个可以不记录，这个看需要，这样可以看到打包的过程
		log.Printf("成功打包 %s ，共写入了 %d 字节的数据\n", fileName, n)

		return nil
	})
}
```

打包及压缩就搞定了，不过这个代码现在我还发现有个问题，就是不能处理软链接

#### 解包解压

这个过程基本就是把压缩的过程返回来，多了些创建目录的操作

```go
package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	var dst = "" // 不写就是解压到当前目录
	var src = "log.tar.gz"

	UnTar(dst, src)
}

func UnTar(dst, src string) (err error) {
	// 打开准备解压的 tar 包
	fr, err := os.Open(src)
	if err != nil {
		return
	}
	defer fr.Close()

	// 将打开的文件先解压
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return
	}
	defer gr.Close()

	// 通过 gr 创建 tar.Reader
	tr := tar.NewReader(gr)

	// 现在已经获得了 tar.Reader 结构了，只需要循环里面的数据写入文件就可以了
	for {
		hdr, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case hdr == nil:
			continue
		}

		// 处理下保存路径，将要保存的目录加上 header 中的 Name
		// 这个变量保存的有可能是目录，有可能是文件，所以就叫 FileDir 了……
		dstFileDir := filepath.Join(dst, hdr.Name)

		// 根据 header 的 Typeflag 字段，判断文件的类型
		switch hdr.Typeflag {
		case tar.TypeDir: // 如果是目录时候，创建目录
			// 判断下目录是否存在，不存在就创建
			if b := ExistDir(dstFileDir); !b {
				// 使用 MkdirAll 不使用 Mkdir ，就类似 Linux 终端下的 mkdir -p，
				// 可以递归创建每一级目录
				if err := os.MkdirAll(dstFileDir, 0775); err != nil {
					return err
				}
			}
		case tar.TypeReg: // 如果是文件就写入到磁盘
			// 创建一个可以读写的文件，权限就使用 header 中记录的权限
			// 因为操作系统的 FileMode 是 int32 类型的，hdr 中的是 int64，所以转换下
			file, err := os.OpenFile(dstFileDir, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			n, err := io.Copy(file, tr)
			if err != nil {
				return err
			}
			// 将解压结果输出显示
			fmt.Printf("成功解压： %s , 共处理了 %d 个字符\n", dstFileDir, n)

			// 不要忘记关闭打开的文件，因为它是在 for 循环中，不能使用 defer
			// 如果想使用 defer 就放在一个单独的函数中
			file.Close()
		}
	}

	return nil
}

// 判断目录是否存在
func ExistDir(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}
```

到这里解压就完成了，只是一个实验代码，还有很多不完善的地方，欢迎提出宝贵的意见。
