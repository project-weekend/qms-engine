CREATE TABLE IF NOT EXISTS `projects` (
    `id`                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT                         COMMENT 'primary key',
    `name`              VARCHAR(50) NOT NULL DEFAULT ''                                 COMMENT 'unique project name',
    `description`       VARCHAR(250) NOT NULL DEFAULT ''                                COMMENT 'project description',
    `created_at`        TIMESTAMP DEFAULT CURRENT_TIMESTAMP                             COMMENT 'created time',
    `updated_at`        TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated time',
    `deleted_at`        TIMESTAMP NULL DEFAULT NULL                                     COMMENT 'deleted time',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_project_name` (`name`),
    INDEX idx_deleted_at (deleted_at)
);
