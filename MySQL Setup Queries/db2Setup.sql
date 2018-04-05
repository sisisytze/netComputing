-- MySQL Script generated by MySQL Workbench
-- Wed 04 Apr 2018 12:32:46 PM CEST
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema db2
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema db2
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `db2` DEFAULT CHARACTER SET utf8 ;
USE `db2` ;

-- -----------------------------------------------------
-- Table `db2`.`sensortype`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `db2`.`sensortype` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `radius` INT NOT NULL,
  PRIMARY KEY (`id`));


-- -----------------------------------------------------
-- Table `db2`.`sensor`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `db2`.`sensor` (
  `uuid` VARCHAR(36) NOT NULL,
  `mac_address` VARCHAR(45) NOT NULL,
  `sensor_type_id` INT NOT NULL,
  `longitude` FLOAT NOT NULL,
  `latitude` FLOAT NOT NULL,
  PRIMARY KEY (`uuid`),
  INDEX `sensor_type_idx` (`sensor_type_id` ASC),
  CONSTRAINT `sensor_type`
    FOREIGN KEY (`sensor_type_id`)
    REFERENCES `db2`.`sensortype` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table `db2`.`measurement`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `db2`.`measurement` (
  `uuid` VARCHAR(36) NOT NULL,
  `sensor_uuid` VARCHAR(36) NOT NULL,
  `value` FLOAT NOT NULL,
  `at` DATETIME NOT NULL,
  PRIMARY KEY (`uuid`),
  INDEX `sensor_idx` (`sensor_uuid` ASC),
  CONSTRAINT `sensor_measure`
    FOREIGN KEY (`sensor_uuid`)
    REFERENCES `db2`.`sensor` (`uuid`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table `db2`.`sync_sensor`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `db2`.`sync_sensor` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `sensor_uuid` VARCHAR(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `sensor_idx` (`sensor_uuid` ASC),
  CONSTRAINT `sensor_sync`
    FOREIGN KEY (`sensor_uuid`)
    REFERENCES `db2`.`sensor` (`uuid`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


-- -----------------------------------------------------
-- Table `db2`.`sync_measurement`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `db2`.`sync_measurement` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `measurement_uuid` VARCHAR(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `measurement_idx` (`measurement_uuid` ASC),
  CONSTRAINT `measurement`
    FOREIGN KEY (`measurement_uuid`)
    REFERENCES `db2`.`measurement` (`uuid`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
