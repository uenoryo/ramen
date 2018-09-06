DROP TABLE IF EXISTS `reminder`;
CREATE TABLE `reminder` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `memo` varchar(255) NOT NULL,
  `remember_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2554 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
