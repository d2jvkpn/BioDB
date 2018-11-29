----
load data local infile 'data_Taxonomy/Taxonomy.tsv' into table BioDB.Taxonomy 
ignore 1 lines (taxon_id, scientific_name, taxon_rank, parent_id, escape_name);
-- Query OK, 2003451 rows affected (15.47 sec)          
-- Records: 2003451  Deleted: 0  Skipped: 0  Warnings: 0

alter table BioDB.Taxonomy add index escape_name (escape_name);

load data local infile 'data_Taxonomy/Homotypic_synonym.tsv' into 
table BioDB.Taxonomy_homotypic ignore 1 lines (taxon_id, name);

alter table BioDB.Taxonomy_homotypic add index taxon_id (taxon_id);
alter table BioDB.Taxonomy_homotypic add index name (name);

----
load data local infile 'data_Pathway/Pathway.tsv' into table BioDB.Pathway 
ignore 1 lines;
-- Query OK, 0 rows affected, 65535 warnings (7 min 38.73 sec)
-- Records: 10815282  Deleted: 0  Skipped: 10815282  Warnings: 28708929

alter table BioDB.Pathway add index taxon_id (taxon_id);
alter table BioDB.Pathway add index pathway_id (pathway_id);

----
load data local infile 'data_Pathway/Pathway_Definition.tsv' into table 
BioDB.Pathway_definition ignore 1 lines;

alter table BioDB.Pathway_definition add index id (id);
alter table BioDB.Pathway_definition add index name (name);

----
load data local infile 'data_Genome/Genome.tsv' into table BioDB.Genome 
ignore 1 lines;
-- Query OK, 231504 rows affected, 221 warnings (32.45 sec)
-- Records: 231725  Deleted: 0  Skipped: 221  Warnings: 221

alter table BioDB.Genome add index taxon_id (taxon_id);
