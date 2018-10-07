-- author: d2jvkpn
-- version: 0.1
-- release: 2018-10-01
-- project: https://github.com/d2jvkpn/BioDB
-- license: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)


create database BioDB;
use BioDB;

create table GO (
    id           int              auto_increment,
    GO_id        nchar(10)        not null,
    prot_id      varchar(255),
    class        nchar(1),
    -- product   varchar(255),
    genes        varchar(255)     not null,
    taxon_id     int              not null,
    check        (taxon_id > 0),
    primary key  (id)
);


create table Taxon (
    taxon_id         int           not null,
    scientific_name  varchar(255)  not null,
    -- common_name      varchar(255),
    lineage          varchar(255),
    primary key      (taxon_id)
);

create table Pathway_code (
    code         varchar(63)   not null,
    oragnism     varchar(255)  not null,
    primary key  (code)
);


create table Pathway_infor (
    code             varchar(63)     not null,
    C_id             varchar(63)     not null,
    C_name           varchar(255),       
    A_id             varchar(63),
    A_name           varchar(255),
    B_id             varchar(63),
    B_name           varchar(255),
    gene_id          varchar(63)     not null,
    gene_name        varchar(63),
    gene_desciption  varchar(255),
    KO_id            varchar(63),
    KO_description   varchar(255),
    EC_id            varchar(63)
);
