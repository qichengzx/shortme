shortme
---------

### URL shortener service, powered by Go

[中文](./README_ZH.md)

## Features

- shortlen url
- redirect to url of hash

## Database

Create table :

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

## RUN

```
git clone git@github.com:qichengzx/shortme.git
cd shortme
go run main.go
```

open "[http://localhost:8000](http://localhost:8000)"

## License

GPL 3.0