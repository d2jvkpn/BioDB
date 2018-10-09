-- author: d2jvkpn
-- version: 0.1
-- release: 2018-10-01
-- project: https://github.com/d2jvkpn/BioDB
-- license: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)

-- sudo mysql -u hello -p;
create database BioDB;

insert into mysql.user (User, Password, Host)
    values ("hello", password(""), "localhost");
-- set password for "hello"@"localhost" = PASSWORD("");

grant all privileges on BioDB.* to 'hello'@'localhost';
flush privileges;

-- mysql -u hello BioDB;

use BioDB;

create table description (
    table_name   varchar(63)    not null,
    description  varchar(255)   not null,
    primary key  (table_name)    
)

create table GO (
    id           int              auto_increment,
    GO_id        nchar(10)        not null,
    prot_id      varchar(255),
    class        nchar(1),
    genes        varchar(255)     not null,
    taxon_id     int              not null,
    check        (taxon_id > 0),
    primary key  (id)
);
    --product      varchar(255),


create table Taxon (
    taxon_id         int           not null,
    scientific_name  varchar(255)  not null,
    primary key      (taxon_id)
);
    -- lineage          varchar(255),
    --common_name      varchar(255),


create table Pathway_code (
    orgcode      varchar(7)   not null,
    organism     varchar(255)  not null,
    lineage      varchar(255)  not null,
    primary key  (code)
);

create table Pathway_LX3 (
    C_id         char(6)       not null,
    C_name       varchar(255)  not null,
    B_id         char(6)       not null,
    B_name       varchar(255)  not null,
    A_id         char(6)       not null,
    A_name       varchar(255)  not null,
    primary key  (C_id)
);


create table Pathway_infor (
    orgcode          varchar(7)     not null,
    pathway_id       varchar(63)    not null,
    pathway_name     varchar(255),
    id               varchar(63)    not null,
    name             varchar(63),
    desciption       varchar(255),
    KO_id            varchar(63)    not null,
    KO_description   varchar(255),
    EC_id            varchar(63)
);
