/*
Navicat MySQL Data Transfer

Source Server         : 178
Source Server Version : 50556
Source Host           : 10.21.1.178:3306
Source Database       : ansible_manager

Target Server Type    : MYSQL
Target Server Version : 50556
File Encoding         : 65001

Date: 2017-12-29 09:24:19
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for ansible_host
-- ----------------------------
DROP TABLE IF EXISTS `ansible_host`;
CREATE TABLE `ansible_host` (
  `host_id` int(11) NOT NULL AUTO_INCREMENT,
  `project_id` int(11) NOT NULL,
  `host_alias` varchar(255) NOT NULL,
  `host_name` varchar(255) NOT NULL,
  `host_ip` varchar(255) NOT NULL,
  `host_user` varchar(255) DEFAULT NULL,
  `host_password` varchar(255) DEFAULT NULL,
  `host_key` text,
  PRIMARY KEY (`host_id`),
  KEY `project_id` (`project_id`),
  CONSTRAINT `ansible_host_ibfk_1` FOREIGN KEY (`project_id`) REFERENCES `ansible_project` (`project_id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for ansible_inventory
-- ----------------------------
DROP TABLE IF EXISTS `ansible_inventory`;
CREATE TABLE `ansible_inventory` (
  `inv_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `inv_name` varchar(255) NOT NULL,
  `inv_value` text,
  PRIMARY KEY (`inv_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for ansible_key
-- ----------------------------
DROP TABLE IF EXISTS `ansible_key`;
CREATE TABLE `ansible_key` (
  `key_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `key_name` varchar(255) NOT NULL,
  `key_value` text,
  PRIMARY KEY (`key_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for ansible_project
-- ----------------------------
DROP TABLE IF EXISTS `ansible_project`;
CREATE TABLE `ansible_project` (
  `project_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `project_name` varchar(255) NOT NULL,
  `project_created` datetime NOT NULL,
  PRIMARY KEY (`project_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for ansible_repository
-- ----------------------------
DROP TABLE IF EXISTS `ansible_repository`;
CREATE TABLE `ansible_repository` (
  `repo_id` int(11) NOT NULL AUTO_INCREMENT,
  `repo_name` varchar(255) NOT NULL,
  `repo_path` varchar(255) NOT NULL,
  `repo_desc` text,
  PRIMARY KEY (`repo_id`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for ansible_task
-- ----------------------------
DROP TABLE IF EXISTS `ansible_task`;
CREATE TABLE `ansible_task` (
  `task_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `tpl_id` int(11) NOT NULL,
  `playbook_tag` varchar(255) DEFAULT NULL,
  `task_name` varchar(255) NOT NULL,
  `task_status` varchar(255) DEFAULT NULL,
  `task_created` datetime NOT NULL,
  `task_start` datetime DEFAULT NULL,
  `task_end` datetime DEFAULT NULL,
  PRIMARY KEY (`task_id`),
  KEY `tpl_id` (`tpl_id`),
  CONSTRAINT `ansible_task_ibfk_1` FOREIGN KEY (`tpl_id`) REFERENCES `ansible_template` (`tpl_id`)
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for ansible_task_output
-- ----------------------------
DROP TABLE IF EXISTS `ansible_task_output`;
CREATE TABLE `ansible_task_output` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `task_id` int(11) NOT NULL,
  `time` datetime NOT NULL,
  `output` longtext NOT NULL,
  PRIMARY KEY (`id`),
  KEY `task_id` (`task_id`),
  CONSTRAINT `ansible_task_output_ibfk_1` FOREIGN KEY (`task_id`) REFERENCES `ansible_task` (`task_id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for ansible_template
-- ----------------------------
DROP TABLE IF EXISTS `ansible_template`;
CREATE TABLE `ansible_template` (
  `tpl_id` int(11) NOT NULL AUTO_INCREMENT,
  `project_id` int(11) NOT NULL,
  `repo_id` int(11) NOT NULL,
  `tpl_name` varchar(255) NOT NULL,
  `playbook` varchar(255) NOT NULL,
  `playbook_parse` text,
  PRIMARY KEY (`tpl_id`),
  KEY `project_id` (`project_id`),
  KEY `repo_id` (`repo_id`),
  CONSTRAINT `ansible_template_ibfk_1` FOREIGN KEY (`project_id`) REFERENCES `ansible_project` (`project_id`) ON DELETE CASCADE,
  CONSTRAINT `ansible_template_ibfk_2` FOREIGN KEY (`repo_id`) REFERENCES `ansible_repository` (`repo_id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for ansible_user
-- ----------------------------
DROP TABLE IF EXISTS `ansible_user`;
CREATE TABLE `ansible_user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_account` varchar(255) NOT NULL,
  `user_password` varchar(32) NOT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

INSERT INTO `ansible_user` (user_account,user_password) VALUES ('admin', 'e10adc3949ba59abbe56e057f20f883e');
-- ----------------------------
-- Table structure for ansible_vars
-- ----------------------------
DROP TABLE IF EXISTS `ansible_vars`;
CREATE TABLE `ansible_vars` (
  `vars_id` int(11) NOT NULL AUTO_INCREMENT,
  `repo_id` int(11) NOT NULL,
  `vars_type` varchar(20) NOT NULL,
  `vars_name` varchar(255) NOT NULL,
  `vars_path` varchar(255) NOT NULL,
  `vars_value` text NOT NULL,
  PRIMARY KEY (`vars_id`),
  KEY `repo_id` (`repo_id`),
  CONSTRAINT `ansible_vars_ibfk_1` FOREIGN KEY (`repo_id`) REFERENCES `ansible_repository` (`repo_id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8;
