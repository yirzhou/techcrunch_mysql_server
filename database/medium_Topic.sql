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
-- Table structure for table `Topic`
--

DROP TABLE IF EXISTS `Topic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Topic` (
  `topic` varchar(100) NOT NULL,
  PRIMARY KEY (`topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Topic`
--

LOCK TABLES `Topic` WRITE;
/*!40000 ALTER TABLE `Topic` DISABLE KEYS */;
INSERT INTO `Topic` VALUES ('3d-printing'),('500-startups'),('acer'),('aereo'),('after-school'),('airbnb'),('alibaba'),('alphabet'),('amazon'),('amazon-fire-phone'),('amazon-kindle'),('amazon-prime'),('amazon-web-services'),('american-jobs-act'),('andreessen-horowitz'),('android'),('apple'),('apple-music'),('apple-pay'),('apple-pencil'),('apple-tv'),('apple-watch'),('bastian-lehmann'),('beats'),('ben-horowitz'),('bill-gates'),('bitcoin'),('box'),('bumble'),('ces-2014'),('chromecast'),('cloudmagic-brings-fast-search-as-you-type-functionality-to-gmail-google-apps'),('coursera'),('craigslist'),('crowdfunding'),('dave-mcclure'),('dell'),('dropbox-inc'),('edward-snowden'),('ellen-pao'),('elon-musk'),('eric-schmidt'),('facebook'),('fidelity-investments'),('firefox'),('fitbit'),('flappy-bird'),('foursquare'),('galaxy-s6-edge'),('github'),('gmail'),('google'),('google-cardboard'),('google-chrome'),('google-driverless-car'),('google-drops-more-than-1-8-billion-on-newest-new-york-office'),('google-glass'),('google-maps'),('gopro'),('gurbaksh-chahal'),('htc'),('icloud'),('imgur'),('indiegogo'),('instacart'),('instagram'),('intel'),('internet-of-things'),('ios-7'),('ios-9'),('ipad'),('ipad-pro'),('iphone-5'),('iphone-5s'),('iphone-6'),('jan-koum'),('jeff-bezos'),('kakaotalk'),('kickstarter'),('lenovo'),('lg'),('linkedin'),('lyft'),('macbook'),('marc-andreessen'),('marissa-mayer'),('mark-zuckerberg'),('max-levchin'),('meerkat'),('microsoft'),('microsoft-azure'),('mozilla'),('net-neutrality'),('netflix'),('nexus-5'),('not_a_tag'),('not_a_topic1'),('not_a_topic2'),('nsa'),('oculus-rift'),('oculus-vr'),('oracle'),('palmer-luckey'),('paypal'),('periscope'),('peter-thiel'),('playstation-4'),('postmates'),('quantcast'),('radiumone'),('reddit'),('reid-hoffman'),('roku'),('salesforce'),('sam-altman'),('samsung'),('satya-nadella'),('sean-rad'),('search engine'),('secret'),('sergey-brin'),('shyp'),('silicon-valley'),('siri'),('skype'),('slack'),('snapchat'),('spacex'),('spotify'),('steve-ballmer'),('steve-jobs'),('surface-pro-3'),('tesla-model-s'),('tesla-motors'),('test'),('test1'),('tim-armstrong'),('tim-cook'),('tinder'),('toshiba'),('travis-kalanick'),('tumblr'),('twitch'),('twitter'),('uber'),('vimeo'),('waze'),('wearables'),('whatsapp'),('windows-10'),('windows-8-1'),('windows-phone'),('wwdc'),('xbox-one'),('xiaomi'),('y-combinator'),('yahoo'),('yahoo-mail'),('yik-yak'),('youtube'),('zenefits'),('zynga'),('zyngas-reported-7-10-billion-valuation-surpasses-that-of-ea');
/*!40000 ALTER TABLE `Topic` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-03-27 18:17:17
