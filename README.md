### 这个程序用于解决局域网中的几台电脑、手机之间临时性互相传送文件问题
编译安装：

使用`go`编译：
`go get github.com/rocket049/fileserver`

直接下载可执行程序：
[[https://github.com/rocket049/fileserver/releases](https://github.com/rocket049/fileserver/releases)
]([https://github.com/rocket049/fileserver/releases](https://github.com/rocket049/fileserver/releases)
)

下载解压后，把可执行程序移动到 `PATH` 中使用。

参数：

```
  -share string
    	Share files in this DIR (default ".")。分享这个目录中的文件。
  -upload string
    	Upload files to this DIR (default ".")。上传文件到这个目录。
```

运行后会显示访问 URL ，输入浏览器就可以下载上传。

同时也会显示一个二维码，可以用手机浏览器扫描访问。

建议：使用另一个小工具 `tray-controller` 控制本程序，方便随时打开和关闭。
 `tray-controller` 的主页：[https://github.com/rocket049/tray-controller](https://github.com/rocket049/tray-controller)

