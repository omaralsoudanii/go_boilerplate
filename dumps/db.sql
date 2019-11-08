	
CREATE TABLE `user` (
 `id` int(12) unsigned NOT NULL AUTO_INCREMENT,
 `first_name` varchar(255) CHARACTER SET utf8 NOT NULL,
 `last_name` varchar(255) CHARACTER SET utf8 NOT NULL,
 `username` varchar(255) CHARACTER SET utf8 NOT NULL,
 `email` varchar(500) CHARACTER SET utf8 NOT NULL,
 `password` varchar(500) CHARACTER SET utf8 NOT NULL,
 `avatar` varchar(500) CHARACTER SET utf8 DEFAULT NULL,
 `country_id` int(12) unsigned DEFAULT NULL,
 `city_id` int(12) unsigned DEFAULT NULL,
 `nationality_id` int(12) unsigned DEFAULT NULL,
 `gender` enum('male','female') DEFAULT NULL,
 `birth_date` datetime DEFAULT NULL,
 `created_at` datetime NOT NULL,
 `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 `flags` bit(8) NOT NULL DEFAULT b'0',
 PRIMARY KEY (`id`),
 UNIQUE KEY `username` (`username`),
 UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
