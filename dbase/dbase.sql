CREATE DATABASE tpoint;

USE tpoint;

CREATE TABLE `tb_point_log` (
  `id` varchar(36) NOT NULL,
  `member_id` int(11) NOT NULL,
  `point_type` tinyint(4) NOT NULL DEFAULT '0',
  `point_desc` varchar(255) NOT NULL,
  `point_before` int(11) DEFAULT '0',
  `point_amount` int(11) NOT NULL DEFAULT '0',
  `created_by` varchar(100) DEFAULT '0',
  `created_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `tb_member` (
  `member_id` int(11) NOT NULL AUTO_INCREMENT,
  `member_name` varchar(255) DEFAULT NULL,
  `member_email` varchar(255) DEFAULT NULL,
  `total_point` int(11) DEFAULT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `created_date` datetime(3) DEFAULT NULL,
  `updated_date` datetime DEFAULT NULL,
  `updated_by` varchar(255) DEFAULT NULL,
  `created_by` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`member_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

INSERT INTO `tb_member` (`member_id`, `member_name`, `member_email`, `total_point`, `is_active`, `created_date`, `updated_date`, `updated_by`, `created_by`)
VALUES
	(1,'Andi','andi@mail.com',1,1,'2020-07-26 18:17:38.354',NULL,NULL,'Admin');