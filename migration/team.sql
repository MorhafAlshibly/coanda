-- -----------------------------------------------------
-- Table `team`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `team` (
    `name` VARCHAR(255) NOT NULL,
    `owner` BIGINT UNSIGNED NOT NULL,
    `score` BIGINT NOT NULL DEFAULT 0,
    `data` JSON NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`name`)
);
CREATE UNIQUE INDEX `owner_UNIQUE` ON `team` (`owner` ASC) VISIBLE;
-- -----------------------------------------------------
-- Table `team_member`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `team_member` (
    `team` VARCHAR(255) NOT NULL,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `data` JSON NULL,
    `joined_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`),
    CONSTRAINT `fk_team_member_team` FOREIGN KEY (`team`) REFERENCES `team` (`name`) ON DELETE NO ACTION ON UPDATE NO ACTION
);
CREATE INDEX `fk_team_member_team_idx` ON `team_member` (`team` ASC) VISIBLE;
-- -----------------------------------------------------
-- Table `team_owner`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `team_owner` (
    `team` VARCHAR(255) NOT NULL,
    `user_id` BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`team`),
    CONSTRAINT `fk_team_owner_team1` FOREIGN KEY (`team`) REFERENCES `team` (`name`) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT `fk_team_owner_team_member1` FOREIGN KEY (`user_id`) REFERENCES `team_member` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
);
CREATE INDEX `fk_team_owner_team_member1_idx` ON `team_owner` (`user_id` ASC) VISIBLE;