CREATE DATABASE  IF NOT EXISTS `medium` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `medium`;
-- MySQL dump 10.13  Distrib 8.0.19, for macos10.15 (x86_64)
--
-- Host: 54.145.219.48    Database: medium
-- ------------------------------------------------------
-- Server version	8.0.19

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Temporary view structure for view `PostInfo`
--

DROP TABLE IF EXISTS `PostInfo`;
/*!50001 DROP VIEW IF EXISTS `PostInfo`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `PostInfo` AS SELECT 
 1 AS `postID`,
 1 AS `category`,
 1 AS `content`,
 1 AS `date`,
 1 AS `img_src`,
 1 AS `section`,
 1 AS `title`,
 1 AS `url`,
 1 AS `topic`,
 1 AS `upCount`*/;
SET character_set_client = @saved_cs_client;

--
-- Final view structure for view `PostInfo`
--

/*!50001 DROP VIEW IF EXISTS `PostInfo`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `PostInfo` AS select `Post`.`postID` AS `postID`,`Post`.`category` AS `category`,`Post`.`content` AS `content`,`Post`.`date` AS `date`,`Post`.`img_src` AS `img_src`,`Post`.`section` AS `section`,`Post`.`title` AS `title`,`Post`.`url` AS `url`,`PostTopic`.`topic` AS `topic`,`PostThumbUp`.`upCount` AS `upCount` from ((`Post` join `PostTopic` on((`Post`.`postID` = `PostTopic`.`postID`))) left join `PostThumbUp` on((`Post`.`postID` = `PostThumbUp`.`postID`))) union select `PostThumbUp`.`postID` AS `postID`,`PostThumbUp`.`upCount` AS `upCount`,`Post`.`category` AS `category`,`Post`.`content` AS `content`,`Post`.`date` AS `date`,`Post`.`img_src` AS `img_src`,`Post`.`section` AS `section`,`Post`.`title` AS `title`,`Post`.`url` AS `url`,`PostTopic`.`topic` AS `topic` from (`PostThumbUp` left join (`Post` join `PostTopic` on((`Post`.`postID` = `PostTopic`.`postID`))) on((`Post`.`postID` = `PostThumbUp`.`postID`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-03-27 18:18:12
