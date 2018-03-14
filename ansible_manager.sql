SET FOREIGN_KEY_CHECKS=0;

DROP TABLE IF EXISTS `ansible_repository`;
CREATE TABLE `ansible_repository` (
  `repo_id` varchar(64) NOT NULL,
  `repo_name` varchar(255) NOT NULL,
  `repo_path` varchar(255) NOT NULL,
  `repo_group` text NOT NULL,
  `repo_tags`  text NOT NULL,
  `repo_vars`  text NOT NULL,
  `repo_notes` text NOT NULL,
  `repo_desc`  text,
  `created` datetime NOT NULL,
  PRIMARY KEY (`repo_id`)
) ENGINE=InnoDB CHARSET=utf8;


DROP TABLE IF EXISTS `ansible_host`;
CREATE TABLE `ansible_host` (
  `host_id` varchar(64) NOT NULL,
  `user_id` varchar(64) NOT NULL,
  `host_alias` varchar(255) NOT NULL,
  `host_name` varchar(255) NOT NULL,
  `host_ip` varchar(255) NOT NULL,
  `host_user` varchar(255) DEFAULT NULL,
  `host_password` varchar(255) DEFAULT NULL,
  `host_status` tinyint,
  `host_key` text,
  `created` datetime NOT NULL,
  PRIMARY KEY (`host_id`)
) ENGINE=InnoDB CHARSET=utf8;


DROP TABLE IF EXISTS `ansible_project`;
CREATE TABLE `ansible_project` (
  `project_id` varchar(64) NOT NULL,
  `user_id`    varchar(64) NOT NULL,
  `project_name` varchar(255) NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`project_id`)
) ENGINE=InnoDB CHARSET=utf8;


DROP TABLE IF EXISTS `ansible_project_host`;
CREATE TABLE `ansible_project_host` (
  `id`   int(11) NOT NULL AUTO_INCREMENT,
  `project_id` varchar(64) NOT NULL,
  `host_id` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`project_id`) REFERENCES `ansible_project` (`project_id`) ON DELETE CASCADE,
  FOREIGN KEY (`host_id`) REFERENCES `ansible_host` (`host_id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `ansible_task`;
CREATE TABLE `ansible_task` (
  `task_id`   varchar(64)  NOT NULL,
  `project_id`   varchar(64)  NOT NULL,
  `repo_id`   varchar(64)  NOT NULL,
  `task_name` varchar(255) NOT NULL,
  `task_group` text NOT NULL,
  `task_vars` text NOT NULL,
  `task_status` varchar(20) NOT NULL,
  `task_tag`  varchar(50),
  `task_start` datetime DEFAULT NULL,
  `task_end` datetime DEFAULT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`task_id`),
  FOREIGN KEY (`project_id`) REFERENCES `ansible_project` (`project_id`) ON DELETE CASCADE,
  FOREIGN KEY (`repo_id`) REFERENCES `ansible_repository` (`repo_id`) ON DELETE CASCADE
) ENGINE=InnoDB CHARSET=utf8;


DROP TABLE IF EXISTS `ansible_timer`;
CREATE TABLE `ansible_timer` (
  `timer_id`   varchar(64)  NOT NULL,
  `user_id`   varchar(64)  NOT NULL,
  `task_id`   varchar(64)  NOT NULL,
  `timer_name` varchar(255) NOT NULL,
  `timer_start` int(11) NOT NULL,
  `timer_interval` int(11) NOT NULL,
  `timer_repeat` int(11) NOT NULL,
  `timer_status` tinyint NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`timer_id`),
  FOREIGN KEY (`task_id`) REFERENCES `ansible_task` (`task_id`) ON DELETE CASCADE
) ENGINE=InnoDB CHARSET=utf8;

DROP TABLE IF EXISTS `ansible_user`;
CREATE TABLE `ansible_user` (
  `user_id`      varchar(64) NOT NULL,
  `user_account` varchar(255) NOT NULL,
  `user_password` varchar(32) NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

INSERT INTO `ansible_user` (user_id,user_account,user_password) VALUES ('2f8e409a-774c-440e-a281-3e21ef6467e0','admin', 'e10adc3949ba59abbe56e057f20f883e');
