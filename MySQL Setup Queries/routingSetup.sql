-- MySQL Script generated by MySQL Workbench
-- do 05 apr 2018 10:49:47 CEST
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema dbrouting
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema dbrouting
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `dbrouting` DEFAULT CHARACTER SET utf8 ;
USE `dbrouting` ;

-- -----------------------------------------------------
-- Table `dbrouting`.`server`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `dbrouting`.`server` (
  `id` INT NOT NULL,
  `database_name` VARCHAR(45) NOT NULL,
  `database_port` VARCHAR(4) NOT NULL DEFAULT '3306',
  `server_port` VARCHAR(4) NOT NULL DEFAULT '8081',
  `api_port` VARCHAR(4) NOT NULL DEFAULT '8080',
  `address` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`));


-- -----------------------------------------------------
-- Table `dbrouting`.`category`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `dbrouting`.`category` (
  `server_0` INT NOT NULL,
  `server_1` INT NOT NULL,
  INDEX `server_1_idx` (`server_1` ASC),
  INDEX `server_0_idx` (`server_0` ASC),
  CONSTRAINT `server_0`
    FOREIGN KEY (`server_0`)
    REFERENCES `dbrouting`.`server` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `server_1`
    FOREIGN KEY (`server_1`)
    REFERENCES `dbrouting`.`server` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
