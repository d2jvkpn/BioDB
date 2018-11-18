
use BioDB;

select database();
show variables like 'max_connections';
select @@datadir;

---
show columns from GO;
show index from GO;

select * from GO limit 10;
select * from GO where id > (select max(id) from GO) - 10;
select * from GO order by taxon_id asc;

select count(taxon_id) from GO;
select max(taxon_id) from GO;
select count(distinct taxon_id) from GO;
select count(*) from GO where taxon_id = 9606;

select genes, GO_id from GO where taxon_id = 9606 into 
outfile "/tmp/GO.9606.tsv";

insert into GO (taxon_id, genes, GO_id) values (
-- insert into GO values (
	"0",
	"J517_1173",
    "GO:0008152"
);

delete from GO where taxon_id = "0";

drop table GO;
truncate table GO;


----
select count(distinct(escape_name)) from Taxonomy;

delete from Pathway where taxon_id not in 
(select taxon_id from Taxonomy);

alter table Pathway add foreign key (taxon_id) references
Taxonomy(taxon_id);

select * from Pathway order by taxon_id limit 100;
create table P2 like Pathway;
insert into P2 select * from Pathway order by taxon_id;
drop table Pathway;
rename table P2 to Pathway;
alter table Pathway change KO_description KO_information varchar(256);
