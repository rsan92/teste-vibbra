-- Active: 1675298682566@@127.0.0.1@3306@test_vibbra
CREATE DATABASE IF NOT EXISTS test_vibbra;

USE test_vibbra;

DROP TABLE IF EXISTS usuarios;

CREATE TABLE usuarios(
    id int auto_increment primary key,
    name varchar(50) NOT NULL,
    email varchar(50) NOT NULL,
    login varchar(50) NOT NULL UNIQUE,
    password varchar(100) NOT NULL,
    lat FLOAT NOT NULL,
    ing FLOAT NOT NULL,
    address VARCHAR(200) NOT NULL,
    city varchar(50) NOT NULL,
    state varchar(2) NOT NULL,
    zip_code INT NOT NULL,
    register_date TIMESTAMP DEFAULT current_timestamp(),
    UNIQUE INDEX id_UNIQUE (id ASC) VISIBLE,
    UNIQUE INDEX login_UNIQUE (login ASC) VISIBLE
) ENGINE=INNODB;
