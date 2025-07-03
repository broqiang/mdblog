---
title: "KVM 安装"
author: "BroQiang"
created_at: 2018-06-02T11:32:42
updated_at: 2018-06-02T11:32:42
---

KVM 安装笔记

## 安装管理工具

### 检查 CPU 是否支持 HVM

安装 KVM ，必须要求 CPU 支持硬件虚拟化技术， 只要不是非常老的设备， 一般都是支持的。

```bash
$ cat /proc/cpuinfo | grep -E "vmx|svm" --color=auto

flags		: 省略更多... vmx ... 省略更多
```

查看是否能找出 vmx 或 svm 的 flag， 如果能够找出，就说明是支持硬件虚拟化。
其中 vmx 对应的是 Intel CPU， svm 对应的是 AMD CPU （这个没证实过，没有使用过 AMD）。

查看内核中是否包含 kvm 模块（一般默认都是存在的）

```bash
$ lsmod | grep kvm
kvm_intel             212992  0
kvm                   598016  1 kvm_intel
```

查看是否存在 kvm 设备（一般支持的情况下都会存在）

```bash
$ ls -l /dev/kvm
crw------- 1 root root 10, 232 6月   2 08:30 /dev/kvm
```

如果上面命令没有查询出结果， 可以查看下 bios 中查找关于 Virtualization Technology 的支持
选项， 并开启， 每个厂家的 bios 设置不太一样， 有些用的是 vt 的简写或者其他名称。

如果 bios 中没有这个选项， 确实也没有查询出 vmx 或 svm 的支持， 那就没法继续安装，
下面内容也就不用看了。

### 安装 qemu-kvm 和 管理工具

- Ubuntu:

```bash
sudo apt install qemu-kvm libvirt-bin bridge-utils virt-manager
```

### 启动 libvirtd 服务

```bash
# 启动服务
sudo systemctl start libvirtd
# 开机自动启动
sudo systemctl enable libvirtd

# 检查状态， 看到 active (running)
sudo systemctl status libvirtd
```

### 修改默认网络

软件安装完成后， 会自动创建一个虚拟桥， 可以通过修改默认的 default 文件来修改默认的网络配置。

编辑 `/etc/libvirt/qemu/networks/default.xml` 文件

```bash
<network>
  <name>default</name>
  <uuid>c3258a3c-ccf4-4776-8a26-bf537c6672c3</uuid>
  <forward mode='nat'/>
  <bridge name='virbr0' stp='on' delay='0'/>
  <mac address='52:54:00:4f:8b:86'/>
  <ip address='192.168.122.1' netmask='255.255.255.0'>
    <dhcp>
      <range start='192.168.122.2' end='192.168.122.254'/>
    </dhcp>
  </ip>
</network>
```

- bridge name='virbr0' 就是虚拟网卡（网桥）设备的名字, 可以通过 `ip addr` 查询到。

- ip address= 就是虚拟网卡（网桥）的地址， 也就是未来基于这个网桥的网关地址。

- 下面就是会启动一个 dhcp 服务， 默认可以分配的地址范围。

根据实际的需要去配置这个配置即可。

## 图形界面管理工具

可以通过 `sudo virt-manager` 打开图形界面管理工具， 或者从所有软件中找到
`Virtual Machine Manager` 这个工具

## 命令行管理工具

- virsh --help

- virt-install

…… 未完，待续
