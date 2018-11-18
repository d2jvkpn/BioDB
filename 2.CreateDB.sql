-- author: d2jvkpn
-- version: 0.2
-- release: 2018-11-12
-- project: https://github.com/d2jvkpn/BioDB
-- license: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)

-- sudo mysql -u hello -p;
create database BioDB;

insert into mysql.user (User, Password, Host)
    values ("hello", password(""), "localhost");
-- set password for "hello"@"localhost" = PASSWORD("");

grant all privileges on BioDB.* to 'hello'@'localhost';
flush privileges;

---- mysql -u hello BioDB;

use BioDB;

create table Taxonomy_homotypic (
    taxon_id   int           not null,
    name       varchar(256)  not null
);

create table Taxonomy (
    taxon_id         int           not null,
    scientific_name  varchar(255)  not null,
    taxon_rank       varchar(255)  not null,
    parent_id        int           not null,
    escape_name      varchar(255)  not null,
    primary key      (taxon_id)
) ENGINE = InnoDB;


----
create table GO (
    taxon_id     int            not null,
    genes        varchar(1024)  not null,
    GO_id        nchar(10)      not null,
	constraint GO_taxon_id
        foreign key (taxon_id) references Taxonomy (taxon_id)
        on delete cascade
        on update restrict
) ENGINE = InnoDB;


-- show create table GO;

create table Pathway (
    taxon_id           int           not null,
    pathway_id         varchar(32)   not null,
    gene_id            varchar(64)   not null,
    KO_id              char(6)       not null,
    KO_description     varchar(256),
    EC_ids             varchar(128),
	constraint Pathway_taxon_id
        foreign key (taxon_id) references Taxonomy (taxon_id)
        on delete cascade
        on update restrict
) ENGINE = InnoDB;


create table Pathway_definition (
    id             char(6)       not null,
    name           varchar(255)  not null,
    parent_id      char(6)       not null,
    primary key    (id)
) ENGINE = InnoDB;

-- select * from BioDB.Pathway_Def where id not regexp '^[ABC][0-9]{5}$'
