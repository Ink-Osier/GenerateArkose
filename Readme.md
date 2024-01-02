## 项目介绍

本项目为 [funcaptcha](https://github.com/acheong08/funcaptcha) 项目的Web实现，支持兼容PandoraNext格式的Arkose Token生成。

## 使用方法

1. 拉取项目至本地

2. 根据Ninja项目中的[Har获取说明](https://github.com/gngpp/ninja/blob/main/doc/readme_zh.md#arkoselabs)下载Har文件至`harPool`文件夹下

3. 启动项目，`docker-compose up -d`

请求示例：

```
POST http://<ip>:23888/api/arkose/token

{
    "type": "gpt-4"
}
```

响应示例：

```
{
    "token": "<Arkose Token>"
}
```