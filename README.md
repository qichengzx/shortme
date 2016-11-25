shortme
---------

### 基于Go的短网址 生成&解析 系统

## 功能列表

- API提交生成短网址
- API根据HASH解析出对应的长网址
- 根据HASH跳转到对应的网址

## TODO

- 校验API接口中URL参数的格式
- ~~redis结果查找失败后查找MySQL~~
- 配置信息放到配置文件中或环境变量中

## 注意

此例中使用的go-hashids包的生成HASH的方法返回结果为数组，取第0个值作为HASH，需要修改

## License

GPL 3.0