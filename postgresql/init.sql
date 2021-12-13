create database user_data;
use user_data;
create table user_data_test
(
    id       int primary key,
    login    varchar(30) not null,
    password varchar(30) not null
);
