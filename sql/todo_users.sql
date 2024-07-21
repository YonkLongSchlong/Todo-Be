CREATE DATABASE IF NOT EXISTS `todo` 
USE `todo`;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` varchar(225) NOT NULL,
  `first_name` varchar(255) NOT NULL,
  `last_name` varchar(255) NOT NULL,
  `email` varchar(225) NOT NULL,
  `password` text NOT NULL,
  `avatar` text NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) 

