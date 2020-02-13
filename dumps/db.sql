CREATE DATABASE go_boiler;

CREATE TABLE `user`
(
    `id`             int(12) unsigned                NOT NULL AUTO_INCREMENT,
    `first_name`     varchar(255) CHARACTER SET utf8 NOT NULL,
    `last_name`      varchar(255) CHARACTER SET utf8 NOT NULL,
    `username`       varchar(255) CHARACTER SET utf8 NOT NULL,
    `email`          varchar(500) CHARACTER SET utf8 NOT NULL,
    `password`       varchar(500) CHARACTER SET utf8 NOT NULL,
    `avatar`         varchar(500) CHARACTER SET utf8          DEFAULT NULL,
    `country_id`     int(12) unsigned                         DEFAULT NULL,
    `city_id`        int(12) unsigned                         DEFAULT NULL,
    `nationality_id` int(12) unsigned                         DEFAULT NULL,
    `gender`         enum ('male','female')                   DEFAULT NULL,
    `birth_date`     datetime                                 DEFAULT NULL,
    `created_at`     datetime                        NOT NULL,
    `updated_at`     datetime                        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `flags`          bit(8)                          NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `username` (`username`),
    UNIQUE KEY `email` (`email`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE `category`
(
    `id`    int(12) unsigned                NOT NULL AUTO_INCREMENT,
    `title` varchar(255) CHARACTER SET utf8 NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE `item`
(
    `id`          int(12) unsigned                NOT NULL AUTO_INCREMENT,
    `title`       varchar(255) CHARACTER SET utf8 NOT NULL,
    `description` varchar(255) CHARACTER SET utf8 NOT NULL,
    `price`       int(12) unsigned                NOT NULL,
    `category_id` int(12) unsigned                NOT NULL,
    `user_id`     int(12) unsigned                NOT NULL,
    `hash`        varchar(500) CHARACTER SET utf8          DEFAULT NULL,
    `created_at`  datetime                        NOT NULL,
    `updated_at`  datetime                        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `user_id` (`user_id`),
    KEY `category_id` (`category_id`),
    CONSTRAINT `item_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
    CONSTRAINT `item_ibfk_2` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE `item_images`
(
    `id`      int(12) unsigned                NOT NULL AUTO_INCREMENT,
    `item_id` int(12) unsigned                NOT NULL,
    `hash`    varchar(500) CHARACTER SET utf8 NOT NULL,
    `type`    varchar(255) CHARACTER SET utf8 NOT NULL,
    `size`    int(12) unsigned                NOT NULL,
    PRIMARY KEY (`id`),
    KEY `item_id` (`item_id`),
    CONSTRAINT `item_images_ibfk_1` FOREIGN KEY (`item_id`) REFERENCES `item` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;


