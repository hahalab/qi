QI(起)
-------

部署 Go 程序到阿里云函数计算的工具

原理
-------

原理和 https://github.com/apex/up 一致

自动部署 Go 程序到阿里的函数计算，并自动创建日志、权限、API Gateway 等服务

通过阿里云 API Gateway 调用函数计算，使用 python 做桥接 fork 用户的程序，可以完成部署任意二进制程序到函数计算

NOTICE
-------

因为阿里云 API Gateway 与 AWS 不同，route path 匹配规则非常简陋

比如在 AWS API 网关可以创建 '/{proxy+}' 与 '/ping' 两个路径，访问 '/pong' 会匹配到 '/{proxy+}'

而 阿里云 API 网关只可以做简单的直接匹配，如创建 '/' 与 '/ping' 两个路径，访问 '/pong' 会匹配失败，所以无法对后台的 Web 服务透传 path

可喜可贺，目前此坑暂弃，立此碑（此README）告诫有缘人

USAGE
-------

qi 已经完成对接 阿里云 API，自动创建函数、服务、Gateway、权限，日志，并打包、上传函数代码 等功能。
因上述原因暂停至此，但仍可使用部分功能。

构建
`make build`

帮助
`bin/qi --help`
