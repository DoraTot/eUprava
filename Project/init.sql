CREATE DATABASE IF NOT EXISTS e_uprava;

USE e_uprava;


CREATE TABLE IF NOT EXISTS attendance_record (
    id INT AUTO_INCREMENT PRIMARY KEY,
    child VARCHAR(100) NOT NULL,
    parent_auth0_id VARCHAR(100) NOT NULL,
    date DATE NOT NULL,
    missing BOOLEAN DEFAULT FALSE,
    justified BOOLEAN DEFAULT FALSE
    );

CREATE TABLE IF NOT EXISTS appointments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    child_name VARCHAR(255),
    parent_id VARCHAR(100) NOT NULL,
    doctor_id VARCHAR(100) NOT NULL,
    date_time DATETIME,
    notes TEXT
    );

CREATE TABLE IF NOT EXISTS medical_justifications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    child_name VARCHAR(255),
    doctor_id VARCHAR(100) NOT NULL,
    parent_id VARCHAR(100) NOT NULL,
    dated DATE,
    reason TEXT
    );