---
title: "sort 包源码分析"
author: "BroQiang"
created_at: 2019-06-26T09:40:10
updated_at: 2019-06-26T09:40:10
---

最近自定义 sort 排序的时候， 发现了个问题， 就是根据 struct 的一个字段排序的时候， 排序完了
相对位置被改变了， 最后只能把多个字段都比较下， 才能保证位置。 就叫我产生了查看下源码的冲动，
看看到底是怎么回事。

这只是个人的一些观点， 如果有不对的地方欢迎指出
[issues](https://github.com/BroQiang/go-docs/issues) 或直接点击修改原文。

看完一遍代码， 感觉思想很重要，这个排序快把我熟悉或知道的排序算法全都用了一遍，
个人建议在看这个源码前要了解： 插入排序、希尔排序、堆排序、快速排序、选择排序，
因为在源码中都使用了， 否则看起来会有点困难。

## 开始

找到源码中的 sort 包， 查看下面包含了这些文件：

```bash
example_interface_test.go
example_keys_test.go
example_multi_test.go
example_search_test.go
example_test.go
example_wrapper_test.go
export_test.go
genzfunc.go
search.go
search_test.go
slice.go
sort.go
sort_test.go
zfuncversion.go
```

其中 `example_*` 的都是示例， 也都很必要全都看一遍， 可以更容易理解这个包怎么去用。

主要代码是再 `sort.go` 这个文件中， 在顶部定义了一个接口， 只要实现了这个接口（包括自定义的），
就可以使用 sort 包中的排序。

```go
type Interface interface {
    // Len is the number of elements in the collection.
    // 集合中元素的个数， 这里一般会很简单， 只要返回个长度就可以了
    Len() int
    // Less reports whether the element with
    // index i should sort before the element with index j.
    // Less 函数判断下标 i 的元素是否应该放在下标 j 的前面
    Less(i, j int) bool
    // Swap swaps the elements with indexes i and j.
    // 交换下标 i j 对应的元素
    Swap(i, j int)
}
```

## Sort 函数

> [排序算法稳定性](https://baike.baidu.com/item/%E6%8E%92%E5%BA%8F%E7%AE%97%E6%B3%95%E7%A8%B3%E5%AE%9A%E6%80%A7) 参考这里

它的入口有两个函数， `func Sort(data Interface)` 和 `func Stable(data Interface)` ,
Sort 是不稳定排序， Stable 是稳定排序， 不过从它的算法使用上来看， Sort 的速度会比 Stable
要快， 一般优先使用 Sort 。

Sort 和 Stable 传入的都是一个 Interface 类型的参数， 所以我们可以将准备排序数据集实现
Interface 就可以使用这两个方法来排序。先看下 Sort ， 后面再说 Stable 。

```go
func Sort(data Interface) {
    n := data.Len() // 实现的 Interface 接口的函数 Len 返回的长度
    // 调用快速排序函数， 但是它不仅仅用了快排的算法， 后面会说
    quickSort(data, 0, n, maxDepth(n))
}
```

上面函数中的 `maxDepth` 是快排递归的最大深度，返回的值为 2\*ceil(lg(n+1))，
他是一个快速排序切换堆排序的阀值， 在查看 quickSort 函数时候再说明。

```go
// 这个函数是根据数据集的长度来计算递归的深度
func maxDepth(n int) int {
    var depth int
    // 这里每次循环 i 右移一位， 相当于 i 每次循环都除以 2 , i /= 2
    // 个人觉得使用位运算性能会高一些吧。
    for i := n; i > 0; i >>= 1 {
        depth++
    }
    return depth * 2
}
```

> 这里有个小疑惑， 不明白为什么要乘以 2， 有人明白可以告诉我下

### quickSort

这是一个重量级的函数， Sort 函数中的主要实现都在这个函数中。

> 下面的分析中假设我们是从小到大排序的， 这样有些地方描述起来容易一些。

```go
func quickSort(data Interface, a, b, maxDepth int) {
    // 这个函数第一次进入的时候 a = 0， b = 数据集长度， 外面的 n
    // 暂时还不知道 a 是什么， 继续往下看
    // 看下面注释的说明， 如果数据集的长度大于 12， 就会进入这个循环，
    // 否则，就使用希尔排序， 先看下进入之后做了什么
    for b-a > 12 { // Use ShellSort for slices <= 12 elements
        // 如果递归到了最大深度， 就使用堆排序
        if maxDepth == 0 {
            // 调用堆排序函数， 一会再看这个函数
            heapSort(data, a, b)
            return
        }
        // 循环一次， 最大深度 -1， 相当于又深入（递归）了一层
        maxDepth--
        // 这个是求中位数的函数， 看到这里大概就明白了 a 和 b 是什么了
        // a 是数据集的左边， b 是数据集的右边， 熟悉快排的应该就明白了
        // 这个就是通过数据集， 左右边， 来求中位数， 这里返回了两个变量，
        // 只能暂停一下， 先看下 doPivot 的实现了
        // 看完 doPivot 再回到这里参考一下， 就可以知道了：
        // doPivot 它取一点为轴，把不大于中位数的元素放左边，大于轴的元素放右边，
        // 返回小于中位数部分数据的最后一个下标，以及大于轴部分数据的第一个下标。
        mlo, mhi := doPivot(data, a, b)
        // Avoiding recursion on the larger subproblem guarantees
        // a stack depth of at most lg(b-a).
        // 因为循环肯定比递归调用节省时间，但是两个子问题只能一个进行循环，另一个只能用递归。
        // 这里是把较小规模的子问题进行递归，较大规模子问题进行循环。
        if mlo-a < b-mhi {
            quickSort(data, a, mlo, maxDepth)
            a = mhi // i.e., quickSort(data, mhi, b)
        } else {
            quickSort(data, mhi, b, maxDepth)
            b = mlo // i.e., quickSort(data, a, mlo)
        }
    }

    // 如果元素的个数小于 12 个（无论是递归的还是首次进入）， 就先使用希尔排序
    // 然后再调用插入排序。
    if b-a > 1 {
        // Do ShellSort pass with gap 6
        // It could be written in this simplified form cause b-a <= 12
        for i := a + 6; i < b; i++ {
            if data.Less(i, i-6) {
                data.Swap(i, i-6)
            }
        }
        insertionSort(data, a, b)
    }
}
```

下面再去看看上面没有说的堆排序和最后出现的插入排序。

### doPivot 函数

这个函数看起来还有点长， 一点一点看吧， 快排的难点就在这里了， 看长度看这里的实现也比较复杂，
如果快速排序的算法熟悉的话， 这个函数可以很容易看明白， 如果看不明白就先去看下快速排序算法。

```go
func doPivot(data Interface, lo, hi int) (midlo, midhi int) {
    // 这里的 lo 和 hi 就是传入的 a 和 b， 代表左右边

    // 这里应该是取左右边的中间点， 下面注释说这样写是为了避免整数溢出
    // 不过这个写法挺巧妙的， 又学了一招
    m := int(uint(lo+hi) >> 1) // Written like this to avoid integer overflow.
    if hi-lo > 40 {
        // Tukey's ``Ninther,'' median of three medians of three.
        // 这里求中位数使用的是 Tukey's ninther ， 下面有专门解释的链接
        // 看完 Tukey's ninther 就应该知道这里做的是什么了吧，
        // 当两边间的元素超过 40 个的时候， 通过 9/3/3 来找到中位数

        // s 的位置是右边 - 左边 处理 8 这个位置， 暂时没明白为什么是 8，
        // 有人明白也可以告诉下我
        s := (hi - lo) / 8

        // medianOfThree 这个函数比较简单， 根据传入的位置， 进行比较，
        // 然后按照位置将元素交换， 如果纠结的话可以先到下面看这个函数的分析
        // 这一次执行完， lo 已经是三个数的中位数
        medianOfThree(data, lo, lo+s, lo+2*s)
        // 这里执行完， m 已经是这三个数的中位数
        medianOfThree(data, m, m-s, m+s)
        // 这里执行完， hi-1 已经是三个数的中位数（因为 hi 是传入的 n，数据集的长度，
        // 下标从 0 开始， 所以要 hi - 1
        medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
    }

    // 将三次中位数的结果再次求中位数， 相当于是用 3*3 个数来确定中位数
    medianOfThree(data, lo, m, hi-1)

    // Invariants are:
    //	data[lo] = pivot (set up by ChoosePivot)
    //	data[lo < i < a] < pivot
    //	data[a <= i < b] <= pivot
    //	data[b <= i < c] unexamined
    //	data[c <= i < hi-1] > pivot
    //	data[hi-1] >= pivot
    pivot := lo // 中位数定义为 lo
    a, c := lo+1, hi-1 // 将定义好的左边右移一位， 右边左移一位

    // 说实话， 这里用 a 和 c 这种变量， 还要向上去看代码才能知道是什么， 有点囧
    // 将左边和中位数进行比较， 一直到不满足条件为止
    for ; a < c && data.Less(a, pivot); a++ {
    }

    // 此时将 a 的位置赋值给 b（又是这样的变量……）
    b := a
    for {
        // 感觉这个和上面的 a<c 的循环做的是一个事，取反比较， 做了这一步应该是
        // 更严谨一些吧， 没有想到什么情况下能进入到这个循环
        for ; b < c && !data.Less(pivot, b); b++ { // data[b] <= pivot
        }
        // 用右边和中间数做比较， 不满足 Less 的时候停止
        for ; b < c && data.Less(pivot, c-1); c-- { // data[c-1] > pivot
        }

        // 比较小， 如果左边和右边重合或者已经再右边的右侧，就证明中间数左侧的数据
        // 全都是比右侧的小， 结束循环， 完成关于这个中位数的排序
        if b >= c {
            break
        }

        // 如果左侧的数据大于右侧， 就将数据交换， 完成排序， 再各自移动一位进行下一轮比较
        // data[b] > pivot; data[c-1] <= pivot
        data.Swap(b, c-1)
        b++
        c--
    }
    // 这里它说如果传入进来的右边 - 处理完了的右边界小于 3， 会出现重复， 保守一点，将比较
    // 的边界设置为 5 （暂时没明白什么意思， 继续往下看吧）
    // If hi-c<3 then there are duplicates (by property of median of nine).
    // Let's be a bit more conservative, and set border to 5.
    protect := hi-c < 5
    // protect 取反了， 就是这个值是大于 5 的
    // 并且传入的右边界 - 当前的右边 < 全部元素 / 4 （又没明白在做什么）
    if !protect && hi-c < (hi-lo)/4 {
        // Lets test some points for equality to pivot
        // 用一些特殊的点和中间数进行比较
        dups := 0
        // 用中位数和右边界的值比较， 如果中位数比右边界大, 就交换
        if !data.Less(pivot, hi-1) { // data[hi-1] = pivot
            data.Swap(c, hi-1)
            // 当前右边界向右移动一位（猜测是因为交换过来的没有进行过比较）
            c++
            dups++
        }
        // 如果移动后的左边界比中间数大（此时中间数有可能已经是上面交换完了的）
        // 就将当前的左边界向左移动一位
        if !data.Less(b-1, pivot) { // data[b-1] = pivot
            b--
            dups++
        }
        // m-lo = (hi-lo)/2 > 6
        // b-lo > (hi-lo)*3/4-1 > 8
        // ==> m < b ==> data[m] <= pivot
        // 用整个集合的中间数和求出的中间数进行比较, 如果比它大， 就交换，
        // 并且将当前的左边界再向左移动一位
        if !data.Less(m, pivot) { // data[m] = pivot
            data.Swap(m, b-1)
            b--
            dups++
        }
        // if at least 2 points are equal to pivot, assume skewed distribution
        // 如果上面的 if 进入了两次， 就证明现在是偏态分布（也就是左右不平衡的）
        protect = dups > 1
    }
    // 如果现在是不平衡的, 再次处理，将数据集平衡
    if protect {
        // Protect against a lot of duplicates
        // Add invariant:
        //	data[a <= i < b] unexamined
        //	data[b <= i < c] = pivot
        for {
            for ; a < b && !data.Less(b-1, pivot); b-- { // data[b] == pivot
            }
            for ; a < b && data.Less(a, pivot); a++ { // data[a] < pivot
            }
            if a >= b {
                break
            }
            // data[a] == pivot; data[b-1] < pivot
            data.Swap(a, b-1)
            a++
            b--
        }
    }
    // Swap pivot into middle
    data.Swap(pivot, b-1)
    // 最后返回的是处理完的左边界和右边界移动后的位置， 相当于都是中间数吧?
    // 因为它都是向中间移动的（这个暂时也是猜测的）
    // 这个函数大体看明白了， 就是计算中位数，然后各种移动，完成排序，
    // 不过还是有好多地方晕晕的， 如果下面的 b 和 c 换成有意义的变量，
    // 可能我就能确定它是什么了， 不是只能猜了。
    return b - 1, c
}
```

上面求中位数使用的是 Tukey's ninther 算法， 也叫 median of medians ，后面是链接，
感兴趣可以直接去看下：
[Tukey's ninther 算法](https://www.johndcook.com/blog/2009/06/23/tukey-median-ninther/)
[中文翻译](https://blog.csdn.net/mianshui1105/article/details/52691711)

### medianOfThree 函数

```go
func medianOfThree(data Interface, m1, m0, m2 int) {
    // 通过调用函数， 我们可以清楚， 这里面的 m1 m0 m2 分别对应的是数据集的索引（位置）
    // 这个函数比较简单， 就是将这个三个位置对应的值通过 Less 函数进行比较， 然后排序
    // Less 函数是 Interface 接口对应的方法的实现

    // sort 3 elements
    if data.Less(m1, m0) {
        data.Swap(m1, m0)
    }
    // data[m0] <= data[m1]
    if data.Less(m2, m1) {
        data.Swap(m2, m1)
        // data[m0] <= data[m2] && data[m1] < data[m2]
        if data.Less(m1, m0) {
            data.Swap(m1, m0)
        }
    }
    // 最终处理完的结果是， 第一个传入的数字在中间位置
    // now data[m0] <= data[m1] <= data[m2]
}
```

### 插入排序

先来看下这个代码的， 这个代码比较简单, 这貌似也没什么分析的， 就是一个最基础的插入排序。

```go
func insertionSort(data Interface, a, b int) {
    for i := a + 1; i < b; i++ {
        for j := i; j > a && data.Less(j, j-1); j-- {
            data.Swap(j, j-1)
        }
    }
}
```

### 堆排序

它主要用了两个函数 `heapSort` 和 `siftDown`， 又都不太长， 就把它们两个都放进来了，
先看下下面的 `heapSort` 。 代码注释中说建立一个最大堆， 这里就不纠结是最大还是最小，
因为最终的判断条件还是依赖 Less 函数的比较结果， 分析的时候就按照最大堆来描述。

```go
// siftDown implements the heap property on data[lo, hi).
// first is an offset into the array where the root of the heap lies.
// 这个函数是用来建堆, first 和堆排序本身没有关系， 因为这里的数组不一定是从 0 开始的，
// 所以需要有 first 来做偏移量， 比如 data[1], 就要是 data[first+1]
func siftDown(data Interface, lo, hi, first int) {
    // 这里看 lo 是堆的根节点
    root := lo
    for {
        // 左子节点的下标（因为最大堆是一个完全二叉树，所以可以确定出数组对应的下标）
        child := 2*root + 1
        if child >= hi {
            break // 如果左子节点的下标超出数组边界， 就停止
        }

        // child+1 是右子节点
        // 如果右子节点没有越界， 找出左右子节点中大的那一个
        // 这个稍微有点绕， 如果 child（左）， 比 child（右）大， child++
        // child 本身就变成了下一个元素
        if child+1 < hi && data.Less(first+child, first+child+1) {
            child++
        }

        // 再用 child 和根节点比较， 如果根节点大于子节点， 就可以退出此次插入
        if !data.Less(first+root, first+child) {
            return
        }

        // 如果根节点比子节点小， 就将子节点和根节点互换
        data.Swap(first+root, first+child)
        // 上面三步执行完， 就挑出了父节点，左右子节点间的交换， 保证根节点是最大的

        // 数据交换过， child 不一定是它的子节点中最大的了， 将child赋给 root，
        // 和它的子节点再比较， 直到满足最大堆的结构
        root = child
    }
}

func heapSort(data Interface, a, b int) {
    first := a // 左侧边界
    lo := 0 // 堆的根节点
    hi := b - a // 右侧边界 - 左侧边界，用来元素计数

    // Build heap with greatest element at top.
    // 这里要先建立一个最大堆(或最小堆， 根据 Less 函数的实现)
    for i := (hi - 1) / 2; i >= 0; i-- {
        // 这里将数组中的数据建立成一个堆的结果
        siftDown(data, i, hi, first)
    }

    // Pop elements, largest first, into end of data.
    // 注意这个循环是从后向前循环的
    for i := hi - 1; i >= 0; i-- {
        // 将第一个元素（堆顶的元素） 和 最后一个元素交换
        // 每一轮循环都将最大的元素放在了最后，下一轮最大元素都是前面一个
        // 这样就可以保证不会影响已经排序好的位置
        data.Swap(first, first+i)

        // 再次维护最大堆的结构
        siftDown(data, lo, i, first)
    }
}
```

到这里 Sort 这个入口就已经分析完了， 这个函数开始的排序是不稳定的， 如果想要排序的结果是稳定排序，
就要去分析下 `Stable` 这个函数

## Stable 函数

```go
func Stable(data Interface) {
    stable(data, data.Len())
}
```

这个函数就比较简单了， 直接只是调用了一下 stable 函数, 继续查看 stable 函数：

```go
func stable(data Interface, n int) {
    // 初始 blockSize 设置为 20
    blockSize := 20 // must be > 0
    a, b := 0, blockSize
    // 将切片按照每个 20 分成多个块， 然后对每个块进行插入排序
    for b <= n {
        insertionSort(data, a, b)
        a = b
        b += blockSize
    }
    // 这一个是对不到 20 的数据进行排序
    insertionSort(data, a, n)

    for blockSize < n {
        a, b = 0, 2*blockSize
        // 每次将两个 block 进行排序
        for b <= n {
            // 调用归并排序， 一会再看
            symMerge(data, a, a+blockSize, b)
            // 这里定义下标的偏移， 下一次循环就是下一组 block
            a = b
            b += 2 * blockSize
        }
        // 将剩余的元素排序
        if m := a + blockSize; m < n {
            symMerge(data, a, m, n)
        }

        // block 每次循环扩大两倍， 直到比元素的总个数大，就结束
        blockSize *= 2
    }
}
```

### 归并排序

```go
func symMerge(data Interface, a, m, b int) {
    // Avoid unnecessary recursions of symMerge
    // by direct insertion of data[a] into data[m:b]
    // if data[a:m] only contains one element.
    // 为了避免不必要的递归，当 data[a:m](第一个 block ) 或者 data[m:b](第 2 个 block )
    // 只有一个元素时，直接插入到另一个子数组中的对应位置。
    if m-a == 1 {
        // Use binary search to find the lowest index i
        // such that data[i] >= data[a] for m <= i < b.
        // Exit the search loop with i == b in case no such index exists.
        i := m
        j := b
        // 这里是找到一个中位数， 进行二分查找
        // 因为前面经过了插入排序， 所以可以保证每一个 block 中的数据都是有序的
        for i < j {
            h := int(uint(i+j) >> 1)
            if data.Less(h, a) {
                i = h + 1
            } else {
                j = h
            }
        }
        // Swap values until data[a] reaches the position before i.
        for k := a; k < i-1; k++ {
            data.Swap(k, k+1)
        }
        return
    }

    // Avoid unnecessary recursions of symMerge
    // by direct insertion of data[m] into data[a:m]
    // if data[m:b] only contains one element.
    // 这里和上面相同， 只是后一半插入前一半
    if b-m == 1 {
        // Use binary search to find the lowest index i
        // such that data[i] > data[m] for a <= i < m.
        // Exit the search loop with i == m in case no such index exists.
        i := a
        j := m
        for i < j {
            h := int(uint(i+j) >> 1)
            if !data.Less(m, h) {
                i = h + 1
            } else {
                j = h
            }
        }
        // Swap values until data[m] reaches the position i.
        for k := m; k > i; k-- {
            data.Swap(k, k-1)
        }
        return
    }

    // 看起来是要计算出一个中位数， 找出 a 到 b 之前的中间点
    mid := int(uint(a+b) >> 1)
    // 根据传入参数的是 symMerge(data, a, a+blockSize, b)
    // 这个 m 不就是中间点吗 ? 有点疑惑
    // mid + m ， 中间点加上中间点， 就是 = b ？ 继续往下看吧
    n := mid + m
    var start, r int
    // 这里做判断了， 如果 m > mid ， 就是说真正的中心点不是传入进来的
    // 应该是左边一半比右边一半的元素要多
    if m > mid {
        // 这 start 此时应该是 m 到 b 之前的位置
        start = n - b
        r = mid
    } else { // 正好数中间点， 或者 mid 在右边部分
        start = a // 左边元素的起点
        r = m // 左边元素结束的位置
    }
    p := n - 1 // 真正的最后一个元素

    // 能进入这个循环应该只能是上面 else 的情况下
    // 然后最终还是要叫 start 不小于 r， 和上面 if 中的情况类似
    for start < r {
        // 再求一下中位数
        c := int(uint(start+r) >> 1)
        if !data.Less(p-c, c) {
            start = c + 1
        } else {
            r = c
        }
    }

    // 烧脑了， 到这里分析不下去了， 过段时间再看了， 放一放， 看的多了就晕了
    // 如果有谁能看明白告诉下我也可以

    end := n - start
    if start < m && m < end {
        rotate(data, start, m, end)
    }

    if a < start && start < mid {
        symMerge(data, a, start, mid)
    }
    if mid < end && end < b {
        symMerge(data, mid, end, b)
    }
}
```

最后没有分析完， 先总结下吧， 就是当从 Stable 函数进入的时候， 就会使用归并排序进行排序，
中间还会使用插入排序预处理下数据， 因为只调用了插入和归并， 所以它是稳定排序。

## 应用

### 内部实现的排序

sort 包本身完成了 int float64 和 string 类型的数据排序， 使用起来也很简单， 分别调用：
`sort.Ints()` 、`sort.Strings` 和 `sort.Float64s` 即可。

它们的实现也很简单， 分别维护了一个 `IntSlice` 、 `Float64Slice` 和 `StringSlice` 的结构，
并且实现了 Interface 接口的 Len 、Less、 Swap 方法， 这里把 IntSlice 的实现复制出来看下，
其他的可以直接查看源码， 非常简单

```go
type IntSlice []int

func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i]
```

### 自定义排序实现

自己实现 sort 的接口也很简单， 直接实现了接口要求的三个方法即可。这里沾了一段我博客代码里的
[片段](https://github.com/BroQiang/mdblog/blob/master/app/mdfile/tags.go#L19)，
可以作为参考， 也是比较简单的。

```go
package mdfile

// Tags 标签
type Tags []Tag

// Tag 标签
type Tag struct {
    // 标签名称
    Title string

    // 标签下文章的数量
    Number int

    // 是否是选中的
    Active bool
}

// Len 实现 Sort 的接口
func (tags Tags) Len() int {
    return len(tags)
}

// Swap 实现的 Sort 接口
func (tags Tags) Swap(i, j int) {
    tags[i], tags[j] = tags[j], tags[i]
}

// Less 实现的 Sort 接口， 按照标签数量排序
func (tags Tags) Less(i, j int) bool {
    return tags[i].Number > tags[j].Number
}
```

### Reverse 函数

这个函数很巧妙， 比如 go 默认实现的 3 种基础类型， 它默认是从小到大， 如果想要从大到小排序，
实现也很方便， 用 Reverse 函数包装一个即可， 如：

```go
sort.Sort(sort.Reverse(sort.IntSlice([]int{1,2,3,4,5,6})))
```

这里因为要被 Reverse 包装， 索引只能使用原始的数据结构 IntSlice ，
其实 Ints 也只是调用了一下 `Sort(IntSlice(...))` 而已， 下面看看代码实现：

```go
type reverse struct {
	// This embedded Interface permits Reverse to use the methods of
	// another Interface implementation.
	Interface
}

// Less returns the opposite of the embedded implementation's Less method.
func (r reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

// Reverse returns the reverse order for data.
func Reverse(data Interface) Interface {
	return &reverse{data}
}
```

可以看到这个函数非常简单， 直接将我们自己的结构包装起来， 并且也只是实现了一个 Less 方法，
这个方法也非常简单，只是将我们原本的条件中的 i 和 j 互换了一下， 就完成了反序。

## 完结

看了几个小时才把代码看完， 还有一些地方是模糊的， 这个代码有点烧脑， 等有空再完善了，
不过这个 sort 写的真的很巧妙， 值得学习一下。
