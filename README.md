# fortune-运势
使用Django搭建的运势项目 目标是 共享&自定义

## 运势池子
<details>
<summary>车万</summary>

- 作者：@妖
- 展示：建设中...

</details>
<details>
<summary>李清歌</summary>

- 作者：@aFox
- 展示：建设中...

</details>
<details>
<summary>原神</summary>

- 作者：[@石头](https://github.com/Katou-Kouseki)
- 展示：建设中...

</details>
<details>
<summary>阴阳师</summary>

- 作者：[@木理](https://github.com/Yiwen-Chan) 
- 展示：建设中...

</details>
<details>
<summary>碧蓝幻想</summary>

- 作者：@饿着吧，笨蛋。
- 展示：建设中...

</details>
<details>
<summary>公主连结</summary>

- 作者：[@Lostdegree](https://github.com/Lostdegree)
- 展示：建设中...

</details>
<details>
<summary>诺亚幻想</summary>

- 作者：@汐言
- 展示：建设中...

</details>

## 客户端下载
| 平台 | 依赖 | 插件地址 | 备注 |
| --- | --- | --- | --- |
| [先驱](https://www.xianqubot.com/) | .net framework 4.6.1 | [Fortune For XQ](https://github.com/Yiwen-Chan/fortune) |  |
| [先驱](https://www.xianqubot.com/) | [铃心自定义](http://qm.myepk.club/variable/) | [Fortune For EPK](https://github.com/Yiwen-Chan/fortune) |  |
| [Mirai](https://www.xianqubot.com/) | [Mirai-Native](https://github.com/iTXTech/mirai-native) & [铃心自定义](http://qm.myepk.club/variable/) | [Fortune For EPK](https://github.com/Yiwen-Chan/fortune) |  |
| [Mirai](https://www.xianqubot.com/) | [CQHTTP-Mirai](https://github.com/yyuueexxiinngg/cqhttp-mirai) | [Fortune For CQHTTP](https://github.com/Yiwen-Chan/fortune) |  |
| [MiraiGo](https://www.xianqubot.com/) | [Go-CQHTTP](https://github.com/Mrs4s/go-cqhttp) | [Fortune For CQHTTP](https://github.com/Yiwen-Chan/fortune) |  |

## 使用说明
本项目仅供学习交流，禁止商业化使用，侵删！
### 先驱
#### 环境要求：

先驱版本 >= 2020090301  [下载地址](http://api.xianqubot.com/index.php?newver=beta)

.net framework 4.6.2  [下载地址](https://dotnet.microsoft.com/download/dotnet-framework/net462)

#### 使用说明：

1.下载 fortune-运势.XQ.dll

2.将 dll文件 复制到 先驱\Plugin\ 路径下

3.重启先驱，启动插件 fortune-运势

4.修改自定义设置

#### 插件设置：

设置路径为：先驱\Config\fortune-运势\config.json

设置说明如下：

```
{
    '默认': {  //填群号，表示该群设置
        '触发': '运势',  //触发关键词，若为关则本群不会被触发
        '回复': '少女祈祷中......',  //收到关键词立刻回复内容
        '类型': '李清歌|碧蓝幻想|公主连结',  //池子见上方，多个池子用 " | " 隔开
        '限制': '全局'  //每日生成限制 【全局】为所有池子当天只生成一次 【池子】为当前池子当天生成一次 【关】当天无限制生成，仅用于测试
    },
}
```

## 致谢
### 特别感谢
- 此项目代码修改来自 [@fz6m](https://github.com/fz6m) 的 [项目-Vortune](https://github.com/fz6m/nonebot-plugin/tree/master/CQVortune) 
- 此项目背景模板来自 [@Lostdegree](https://github.com/Lostdegree) 的 [项目-Portune](https://github.com/Lostdegree/Portune)
