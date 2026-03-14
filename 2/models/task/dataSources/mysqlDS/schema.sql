CREATE TABLE IF NOT EXISTS `tasks` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(128) NOT NULL,
    `description` VARCHAR(512) NOT NULL DEFAULT '',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NULL DEFAULT NULL,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_deleted_at` (`deleted_at`)
);
