-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE `user` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
    `user_id` CHAR(22) CHARACTER SET UTF8MB4 NOT NULL DEFAULT '' COMMENT 'user shortuuid id',
    `first_name` CHAR(255) CHARACTER SET UTF8MB4 NOT NULL DEFAULT '' COMMENT 'user name',
    `last_name` CHAR(255) CHARACTER SET UTF8MB4 NOT NULL DEFAULT '' COMMENT 'user name',
    `email` CHAR(255) CHARACTER SET UTF8MB4 NOT NULL DEFAULT '' COMMENT 'user name',

    `db_create_time` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3) COMMENT 'database insertion time, please do not modify',
    `db_modify_time` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3) COMMENT 'database update time, please do not modify',
    `disabled_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'disabled time',
    `disabled` TINYINT DEFAULT '0' COMMENT 'soft delete',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_idx_user_id` (`user_id`),
    UNIQUE KEY `uniq_idx_email` (`email`)
) ENGINE = INNODB AUTO_INCREMENT = 1 DEFAULT CHARSET = UTF8MB4 COMMENT = 'user table';
