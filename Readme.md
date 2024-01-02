## 项目介绍

本项目为 [funcaptcha](https://github.com/acheong08/funcaptcha) 项目的Web实现，支持兼容PandoraNext格式的Arkose Token生成。

如果本项目对你有帮助，就给个小星星吧~

> [!WARNING]
>
> 本项目不保证生成的Arkose的可用性以及使用本项目不会导致封号等一系列问题，使用本项目造成的一切后果由使用者自行承担。

## 使用方法

1. 拉取项目至本地

2. 新建`harPool`文件夹

3. 根据Ninja项目中的[Har获取说明](https://github.com/gngpp/ninja/blob/main/doc/readme_zh.md#arkoselabs)下载Har文件至`harPool`文件夹下

4. 启动项目，`docker-compose up -d`

请求示例：

```
POST http://<ip>:23888/api/arkose/token
```
请求体格式：application/x-www-form-urlencoded
参数：`type`: `gpt-4`

响应示例：

```
{
    "token": "<Arkose Token>"
}
```