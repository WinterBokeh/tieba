# tieba

山寨贴吧

## 进度

- 基本完成用户模块 12.12

## 暂时应用的内容：
- 使用了xorm框架
- 尝试用配置文件管理
- 邮箱验证登录功能
- redis存验证码
- 密码非明文（md5加盐存储）
- 由于使用http协议，为保证安全性，使用服务器上的session存储登录状态