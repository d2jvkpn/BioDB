-- author: d2jvkpn
-- version: 0.4
-- release: 2018-12-06
-- project: https://github.com/d2jvkpn/BioDB
-- license: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)

-- sudo mysql -u hello -p;
create database BioDB;

insert into mysql.user (User, Password, Host)
    values ("hello", password(""), "localhost");
-- set password for "hello"@"localhost" = PASSWORD("");

grant all privileges on BioDB.* to 'hello'@'localhost';
flush privileges;

insert into mysql.user (User, Password, Host)
    values ("world", password(""), "localhost");

grant select on BioDB.* to 'world'@'localhost';
flush privileges;

---- mysql -u hello;

create table BioDB.Taxonomy (
    taxon_id         int           not null,
    scientific_name  varchar(255)  not null,
    taxon_rank       varchar(255)  not null,
    parent_id        int           not null,
    escape_name      varchar(255)  not null,
    primary key      (taxon_id)
) ENGINE = InnoDB;


create table BioDB.Taxonomy_homotypic (
    taxon_id   int           not null,
    name       varchar(256)  not null
);

----
create table BioDB.GO (
    taxon_id     int            not null,
    genes        varchar(1024)  not null,
    GO_id        nchar(10)      not null,
	constraint GO_taxon_id
        foreign key (taxon_id) references Taxonomy (taxon_id)
        on delete cascade
        on update restrict
) ENGINE = InnoDB;


create table BioDB.GO_definition (
    GO_id       char(10)       not null,
    name        varchar(256)   not null,
    namespace   char(18)       not null,
    definition  varchar(2048)  not null
) ENGINE = InnoDB;


create table BioDB.GO_relation (
    GO_id     char(10)       not null,
    relation  varchar(32)    not null,
    target    varchar(2048)  not null
) ENGINE = InnoDB;


create table BioDB.Pathway (
    taxon_id           int           not null,
    pathway_id         varchar(32)   not null,
    gene_id            varchar(64)   not null,
	gene_information   varchar(1024),
    KO_id              char(6)       not null,
    KO_information     varchar(256),
    EC_ids             varchar(128),
	constraint Pathway_taxon_id
        foreign key (taxon_id) references Taxonomy (taxon_id)
        on delete cascade
        on update restrict
) ENGINE = InnoDB;


-- create table BioDB.Pathway_definition (
--    id             char(6)       not null,
--    name           varchar(255)  not null,
--    parent_id      char(6)       not null,
--    primary key    (id)
-- ) ENGINE = InnoDB;
--

create table BioDB.Pathway_definition (
    C_id             char(6)       not null,
    C_name           varchar(128)  not null,
    B_id             char(6)       not null,
    B_name           varchar(128)  not null,
    A_id             char(6)       not null,
    A_name           varchar(128)  not null,
    primary key    (C_id)
) ENGINE = InnoDB;


create table BioDB.Genome (
    taxon_id           int           not null,
    organism_name      varchar(256)  not null,
    URL                varchar(256)  not null,
    information        varchar(512),
	constraint Genome_taxon_id
        foreign key (taxon_id) references Taxonomy (taxon_id)
        on delete cascade
        on update restrict
) ENGINE = InnoDB;
