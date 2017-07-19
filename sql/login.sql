/*
Navicat MariaDB Data Transfer

Source Server         : localhost
Source Server Version : 100121
Source Host           : localhost:3306
Source Database       : login

Target Server Type    : MariaDB
Target Server Version : 100121
File Encoding         : 65001

Date: 2017-07-19 16:52:01
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for accounts
-- ----------------------------
DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(16) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `password` char(60) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `status` tinyint(3) unsigned NOT NULL DEFAULT '41',
  `auth_key` char(32) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

-- ----------------------------
-- Records of accounts
-- ----------------------------
INSERT INTO `accounts` VALUES ('1', 'ubis', '$2a$10$.xhBdq6h30zS533vogmFpeUedBbbEbYcS6a4ldg/0fwXuj.CzBSHe', '32', 'xEQrKqsUuVDxmTdVCMMjFygOC5TxxGvX');

-- ----------------------------
-- Table structure for accounts_subpassword
-- ----------------------------
DROP TABLE IF EXISTS `accounts_subpassword`;
CREATE TABLE `accounts_subpassword` (
  `account` int(11) NOT NULL,
  `password` char(60) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `answer` char(60) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `question` tinyint(3) unsigned NOT NULL,
  `expires` datetime NOT NULL,
  PRIMARY KEY (`account`),
  CONSTRAINT `accounts_subpassword_ibfk_1` FOREIGN KEY (`account`) REFERENCES `accounts` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
