CREATE DATABASE IF NOT EXISTS userDB;

USE userDB;

create table if not exists familyroom (
    id INT primary key auto_increment not NULL,
    roomname varchar(128) not null unique
);

create table if not exists users (
    id INT primary key auto_increment not NULL,
    username varchar(255) not null unique,
    passhash BINARY(60) not null,
    firstname varchar(64) not null,
    lastname varchar(128) not null,
    photourl varchar(2083) not null,
    personrole VARCHAR(255),
    roomname varchar(128)
);

create table if not exists userlogin (
    id INT PRIMARY KEY auto_increment NOT NULL, 
    userid int not NULL,
    timesignin DATETIME NOT NULL,
    ipaddr VARCHAR(2083) NOT NULL
);