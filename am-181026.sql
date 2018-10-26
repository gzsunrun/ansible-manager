ALTER TABLE `ansible_host`
ADD COLUMN `host_port`  int NULL DEFAULT 22 AFTER `host_ip`;

ALTER TABLE `ansible_host`
ADD COLUMN `host_tag`  int NULL AFTER `host_ip`;

ALTER TABLE `ansible_repository`
ADD COLUMN `repo_version`  varchar(255) NULL AFTER `repo_name`;