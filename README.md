# goi
### 简介
goi是一个代码注入工具, 模仿goc的工作原理，将代码复制到临时目录对源文件进行代码修改后，执行编译。 
主要是提供给mock功能使用，也支持自定义插入目的代码。

### 功能
* 支持mock功能
* 支持录制功能
* 支持自定义代码插入
* 兼容goc

### 获取
``
go get git.garena.com/shopee/seller-server/seller-marketing/marketing-common/goi@latest
``

### 示例
默认配置，目前支持shopbff项目

``
goi build --buildflags="-mod=vendor" . -o xxx  --debug --spex-mocker --http-mocker
``

自定义配置，支持其他任何项目

``
goi build --buildflags="-mod=vendor" . -o bs  --debug -c xxx.json
``

参数详情，可以通过goi help查看

