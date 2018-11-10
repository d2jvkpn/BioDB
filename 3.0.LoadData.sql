load data local infile 'Taxonomy.tsv' into table BioDB.Taxonomy ignore 1 lines
(taxon_id, scientific_name, taxon_rank, parent_id, escape_name);
-- Query OK, 2003451 rows affected (15.47 sec)          
-- Records: 2003451  Deleted: 0  Skipped: 0  Warnings: 0

alter table BioDB.Taxonomy add index escape_name (escape_name);
-- Query OK, 0 rows affected, 1 warning (8.97 sec)     
-- Records: 0  Duplicates: 0  Warnings: 1

load data local infile 'GO.tsv' into table BioDB.GO 
(taxon_id, genes, GO_id);
-- Query OK, 49999635 rows affected, 365 warnings (5 min 42.88 sec)
-- Records: 50000000  Deleted: 0  Skipped: 365  Warnings: 365
