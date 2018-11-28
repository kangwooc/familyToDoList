--  created table to store instances of those User structs
CREATE DATABASE IF NOT EXISTS userDB;

USE userDB;

create table if not exists users (
    id INT primary key auto_increment not NULL,
    username varchar(255) not null unique,
    passhash BINARY(60) not null,
    firstname varchar(64) not null,
    lastname varchar(128) not null,
    photourl varchar(2083) not null 
);

create table if not exists userlogin (
    id INT PRIMARY KEY auto_increment NOT NULL, 
    userid int not NULL,
    timesignin DATETIME NOT NULL,
    ipaddr VARCHAR(2083) NOT NULL
);