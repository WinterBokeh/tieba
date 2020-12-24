# tieba

  功能更完善的留言板，带板块功能。

  ## 进度

  - 基本完成用户模块
  - 一些安全方面的改善

  ## 暂时应用的内容：

  - 使用了xorm框架
  - 尝试用配置文件管理
  - 邮箱验证登录功能
  - redis存验证码
  - 密码非明文（md5加盐存储）
  - ~~使用cookie维护登录状态~~
  - ~~使用session维护登录状态~~
  - 使用json web token维护登录状态
  - 使用jwt存储用户信息，降低服务器数据库查询次数
  - 由于没有前端可以甩锅，postman要每次手动输token（很麻烦的方法）
  - 改用RESTful api风格
  - 添加[**后端接口文档**](https://github.com/WinterBokeh/tieba/blob/master/document.md)
  - 使用refreshToken刷新token，避免频繁登录
  - 用户模块的设计参考了各大网站
