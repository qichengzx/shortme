shortme
---------

### 基于Go的短网址 生成&解析 系统

[English](./README.md)

## 功能列表

- API提交生成短网址
- 根据HASH跳转到对应的网址

## TODO

- ~~校验API接口中URL参数的格式~~ ，仅验证网址是否包含http头。
- ~~redis结果查找失败后查找MySQL~~
- 配置信息放到配置文件中或环境变量中

## 数据库

创建数据表 :

```sql
DROP TABLE IF EXISTS `links`;
CREATE TABLE `links` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `hash` varchar(62) NOT NULL,
  `long_url` varchar(255) NOT NULL,
  `clicks` int(11) NOT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_hash` (`hash`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

```

## 运行

```
git clone git@github.com:qichengzx/shortme.git
cd shortme
go run main.go
```

打开 "[http://localhost:8000](http://localhost:8000)"

## 注意

此例中使用的go-hashids包的生成HASH的方法返回结果为数组，取第0个值作为HASH，需要修改

## License

GPL 3.0