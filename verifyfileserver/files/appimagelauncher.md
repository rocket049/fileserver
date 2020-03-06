# appimagelauncher - 帮助 appimage 程序把自己加入启动菜单
安装：`go get -v gitee.com/rocket049/appimagelauncher`

只需在`appimage`程序中调用`Create`函数，就可以在`~/.local/share/applications`中为当前的 `appimage` 程序创建`.desktop`启动器，把自己加入启动菜单。如果程序不是以`APPIMAGE`格式打包的，本程序会返回一个`error`。

函数 `func Create(desktopFile, iconFile string, force bool) error`的参数：

- `desktopFile`：`appimage` 的根目录中的 `.desktop` 文件名字。
- `iconFile`：`appimage` 的根目录中的图标文件名字。
- `force`：强制更新，如果 `force` 的值为 `false`，当启动器已经存在，并且比 `APPIMAGE` 更新时，不会重复创建启动器。

### 示例代码
```
import "gitee.com/rocket049/appimagelauncher"

appimagelauncher.Create("appimage-name.desktop", "icon-name.png", false)
```