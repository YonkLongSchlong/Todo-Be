CREATE DATABASE  IF NOT EXISTS `todo` 
USE `todo`;

DROP TABLE IF EXISTS `todos`;
CREATE TABLE `todos` (
  `id` varchar(100) NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `category` varchar(45) NOT NULL,
  `is_completed` tinyint NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `user_id` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_todos_users_idx` (`user_id`),
  CONSTRAINT `FK_todos_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) 
  ON DELETE RESTRICT 
  ON UPDATE RESTRICT
)


