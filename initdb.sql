create database if not exists logwebhook;
use logwebhook;
create table if not exists  alertlog (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `alertname` varchar(512) DEFAULT NULL,
  `count` int(11) DEFAULT NULL,
  `status` varchar(32) DEFAULT NULL,
  `update_time` timestamp NULL DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
)CHARSET=utf8;
