CREATE TABLE IF NOT EXISTS `studentss` (
     `id` BIGINT NOT NULL AUTO_INCREMENT,
    `student_code` VARCHAR(128) NOT NULL ,
    `first_name` VARCHAR(128) NOT NULL,
    `last_name` VARCHAR(512) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL ,
    `deleted_at` TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (`id`),
 INDEX `idx_created_at` (`created_at`),
 INDEX `idx_deleted_at` (`deleted_at`)
);


