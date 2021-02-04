create database if not exists alertwebhook;
use alertwebhook;
create table if not exists  alertlog (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `alertname` varchar(512) DEFAULT NULL,
  `level` int(11) DEFAULT NULL,
  `name` varchar(512) DEFAULT NULL,
  `fingerprint` varchar(32) DEFAULT NULL,
  `count` int(11) DEFAULT NULL,
  `status` varchar(32) DEFAULT NULL,
  `update_time` timestamp NULL DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
)CHARSET=utf8;
