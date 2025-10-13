CREATE DATABASE IF NOT EXISTS e_uprava;

USE e_uprava;


CREATE TABLE IF NOT EXISTS attendance_record (
    id INT AUTO_INCREMENT PRIMARY KEY,
    child VARCHAR(100) NOT NULL,
    parent_auth0_id VARCHAR(50) NOT NULL,
    date DATE NOT NULL,
    missing BOOLEAN DEFAULT FALSE,
    justified BOOLEAN DEFAULT FALSE
    );
