CREATE DATABASE IF NOT EXISTS test_vibbra;

USE test_vibbra;

DROP TABLE IF EXISTS usuarios;

CREATE TABLE usuarios(
    id int auto_increment primary key,
    login varchar(50) NOT NULL UNIQUE,
    password varchar(100) NOT NULL,
    register_date TIMESTAMP DEFAULT current_timestamp(),
    UNIQUE INDEX id_UNIQUE (id ASC) VISIBLE,
    UNIQUE INDEX usuario_UNIQUE (login ASC) VISIBLE
) ENGINE=INNODB;