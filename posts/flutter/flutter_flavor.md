---
title: "flutter flavor 配置"
author: "BroQiang"
created_at: 2024-06-10T01:53:53
updated_at: 2024-06-10T01:53:53
---

# flutter flavor 配置

[参考](https://dwirandyh.medium.com/create-build-flavor-in-flutter-application-ios-android-fb35a81a9fac)

## Flutter 中配置

### 创建 lib/flavors/flavor_config.dart 文件

> 文件内容只作为参考，按照实际的用到的去修改即可

```dart
class FlavorConfig {
  /// App 名字
  final String appName;

  /// Api 地址
  final String apiUrl;

  /// 当前环境
  final Flavor flavor;

  static FlavorConfig shared = FlavorConfig.create();

  factory FlavorConfig.create({
    String appName = "",
    String apiUrl = "",
    Flavor flavor = Flavor.dev,
  }) {
    return shared = FlavorConfig(appName: appName, apiUrl: apiUrl, flavor: flavor);
  }

  FlavorConfig({required this.appName, required this.apiUrl, this.flavor = Flavor.dev});
}
```

### 修改原来的 lib/main.dart 文件

为了防止文件名称混淆，这里直接将原来的 `lib/main.dart` 文件放在 `lib/flavors` 目录下，并且改名字为 `mainCommon.dart`

```dart
// 只修改了这里，其他内容还是该怎么样就怎么样
void mainCommon() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
      ),
      home: const MyHomePage(title: 'Flutter Demo Home Page'),
    );
  }
}

class MyHomePage extends StatelessWidget {
  const MyHomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      /// 需要使用的时候就可以直接这样调用
      body: Center(
        child: Column(
          children: [
            Text(FlavorConfig.shared.flavor.name),
            Text(FlavorConfig.shared.appName),
            Text(FlavorConfig.shared.apiUrl),
          ],
        ),
      ),
    );
  }
}
```

### 为每一个环境创建一个入口文件

- 创建 Dev `lib/flavors/main_dev.dart`

```dart
import 'flavor_config.dart';
import 'main_common.dart';

void main() {
  FlavorConfig.create(
    appName: 'app dev',
    apiUrl: 'http://127.0.0.1:8081',
    flavor: Flavor.dev,
  );

  mainCommon();
}
```

- Prod `lib/flavors/main_prod.dart`

```dart
import 'flavor_config.dart';
import 'main_common.dart';

void main() {
  FlavorConfig.create(
    appName: 'app prod',
    apiUrl: 'https://yourapiserver.com',
    flavor: Flavor.prod,
  );

  mainCommon();
}
```

## 配置 iOS 端 Flavor

用 xcode 打开文件， vscode 上是右键点击个目录下的 `ios` 目录， 然后选择 `Open in Xcode` 选项，
注意：`下面操作都是 xcode 中的`。

### 配置 Runner Configurations

选中 Runner -> PROJECT -> Configurations ， 在这个区域点击 + 号，将原本的 Debug、Release、Profile 复制一份。

![flutter_flavor_2024-06-12-09-16-16](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-09-16-16.png)

将刚刚复制的三个 Debug copy、Release copy、Profile copy 改名为 Debug-dev、Release-dev、Profile-dev。

再将原本的 Debug、Release、Profile 改名为 Debug-prod、Release-prod、Profile-prod。

![flutter_flavor_2024-06-12-09-35-32](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-09-35-32.png)

### 配置 Scheme

![flutter_flavor_2024-06-12-09-22-01](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-09-22-01.png)

![flutter_flavor_2024-06-12-09-25-24](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-09-25-24.png)

这里名字是叫的 dev， 名字要和 Flutter 里面配置的 Flavor 一致，防止使用时候出错， 然后再将原本存在的 Runner ， 改名为 prod

![flutter_flavor_2024-06-12-09-28-05](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-09-28-05.png)

然后将上面两个 Scheme 配置

![flutter_flavor_2024-06-12-09-30-36](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-09-30-36.png)

![flutter_flavor_2024-06-12-09-42-16](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-09-42-16.png)

prod 环境的是重命名的，应该会自动匹配上，也可以检查一下，是否正确

![flutter_flavor_2024-06-12-09-44-09](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-09-44-09.png)

### 配置 Bundle Identifier

打开 Runner -> Build Settings ， 搜索 Bundle Identifier， 将搜索出来的 `Product Bundle Identifier` 中 -dev 的值添加 `.dev`

![flutter_flavor_2024-06-12-09-51-33](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-09-51-33.png)

### 修改 App 名字

打开 Runner -> Build Settings 搜索 `Product Name`， 分别修改对应的 App 名字。

![flutter_flavor_2024-06-12-09-57-22](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-09-57-22.png)

这里的 `我的App` 就是安装后在手机桌面上显示的应用名称，按照实际名称修改，最好通过名称来区分不同环境， 比如在后面加个 `Dev`， 当然不改也不会有影响。

然后在 `info.plist` 中修改 Bundle display name 为 $(PRODUCT_NAME) 变量

![flutter_flavor_2024-06-12-10-04-05](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-10-04-05.png)

### 修改 App Icon

可以将准备好图标文件导入，如果没有，如果没有可以在 [https://www.appicon.co](https://www.appicon.co) 制作， 这里分别起名为 `AppIconDev`、`AppIconProd`， 名字随意，后面配置的时候能够对应上即可。

如果不需要区分图标， 可以直接替换原本的 AppIcon 即可，后面步骤就不需要了。

![flutter_flavor_2024-06-12-10-13-33](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-10-13-33.png)

图标导入成功后 在 `Runner` -> `Build Settings` 中搜索 `Primary App Icon` 将搜索出来的 `AppIconDev` 替换为 `AppIconProd`。

![flutter_flavor_2024-06-12-10-19-43](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-10-19-43.png)

### iOS 完成

到现在， iOS 的配置已经完成， 可以通过下面命令，分别测试下 dev 和 prod 环境。

```bash
# dev 环境
flutter run -t lib/flavors/main_dev.dart --flavor dev

# prod 环境
flutter run -t lib/flavors/main_prod.dart --flavor prod
```

## 配置 Android 端 Flavor

### build.gradle 中添加 Flavor 配置

打开 `android/app/build.gradle`， 在 `android{}` 中添加如下配置

```gradle
android {
  defaultConfig {
         ...
      }
  ...
  flavorDimensions "default"
  productFlavors {
      prod {
          dimension "default"
          resValue "string", "app_name", "YourAppName"
      }
      dev {
          dimension "default"
          applicationIdSuffix ".dev"
          resValue "string", "app_name", "YourAppNameDev"
          versionNameSuffix ".dev"
      }
  }
```

### 修改 App 名称

在 `android/app/src` 下创建两个文件（目录不存在，也创建） `dev/res/values/strings.xml` 和 `prod/res/values/strings.xml`

写入下面内容

```xml
<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string name="app_name">YourAppNameDev</string>
</resources>
```

修改 `android/app/src/main/AndroidManifest.xml` 中的 `android:label` 为 `@string/app_name`

```xml
<manifest xmlns:android="http://schemas.android.com/apk/res/android">
    <application
        android:label="@string/app_name"
        ......
```

### 修改 App Icon

将准备好的图标或者 [https://www.appicon.co](https://www.appicon.co) 制作的图标分别导入到 `android/app/src/dev/res`
和 `android/app/src/prod/res` 中， 如果不想要区分，直接放到 `android/app/src/main/res` 中即可。两个环境都可以使用一套图标。

![flutter_flavor_2024-06-12-10-30-45](https://image.broqiang.com/mdblog/flutter_flavor_2024-06-12-10-30-45.png)

### 完成

测试下打包

执行上面的 run 命令，这次选择设备选择 Android

## 编写 MakeFile

执行的命令有点长， 写个 MakeFile 来简化下

在项目根目录下创建文件 `Makefile`, 内容如下

```makefile
PHONY: run_dev run_dev_ios run_dev_android run_prod build_dev_apk build_prod_apk

FLAVOR_DEV=lib/flavors/main_dev.dart
FLAVOR_PROD=lib/flavors/main_dev.dart

run_dev:
	flutter run -t ${FLAVOR_DEV} --flavor dev

# -d 后面跟的事设备的 id， 通过 flutter devices 查询， 根据实际的加入
run_dev_ios:
	flutter run -t ${FLAVOR_DEV} --flavor dev -d 59b2a30e8b1edcc60809265090431cd1d9debeeb

# -d 后面跟的事设备的 id， 通过 flutter devices 查询， 根据实际的加入
run_dev_android:
	flutter run -t ${FLAVOR_DEV} --flavor dev -d rgvorccehepfpb5h

run_prod:
	flutter run -t ${FLAVOR_PROD} --flavor prod

build_dev_apk:
	flutter build apk -t ${FLAVOR_DEV} --flavor dev

build_prod_apk:
	flutter build apk -t ${FLAVOR_PROD} --flavor prod
```

使用的时候，直接在项目根目录下执行 make 命令即可

```bash
# 运行 dev 环境
make run_dev

# 打包 dev 环境 的 apk 包
make build_dev_apk
```

## 编写 VsCode launch.json

在项目根目录下创建文件 `.vscode/launch.json`， 内容如下

写入下面内容

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "tsappDev",
      "request": "launch",
      "type": "dart",
      "program": "lib/flavors/main_dev.dart",
      "args": ["--flavor", "dev"]
    },
    {
      "name": "tsappProd",
      "request": "launch",
      "type": "dart",
      "program": "lib/flavors/main_prod.dart",
      "args": ["--flavor", "prod"]
    }
  ]
}
```

至此 Flavor 配置完成， 可以愉快的在各种环境中玩耍。
