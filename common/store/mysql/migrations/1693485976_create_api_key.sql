-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE `api_key` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'id',
    `user_id` CHAR(22) CHARACTER SET UTF8MB4 NOT NULL DEFAULT '' COMMENT 'user shortuuid id',
    `api_key_id` CHAR(22) CHARACTER SET UTF8MB4 NOT NULL DEFAULT '' COMMENT 'api key shortuuid id',
    `type` TINYINT NOT NULL DEFAULT 0 COMMENT 'type',
    `expiration` DATETIME(3) DEFAULT NULL COMMENT 'expiration time',

    `db_create_time` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3) COMMENT 'database insertion time, please do not modify',
    `db_modify_time` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3) ON UPDATE CURRENT_TIMESTAMP (3) COMMENT 'database update time, please do not modify',
    `disabled_at` TIMESTAMP NULL DEFAULT NULL COMMENT 'disabled time',
    `disabled` TINYINT DEFAULT '0' COMMENT 'soft delete',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_idx_api_key_id` (`api_key_id`)
) ENGINE = INNODB AUTO_INCREMENT = 1 DEFAULT CHARSET = UTF8MB4 COMMENT = 'API key table';
