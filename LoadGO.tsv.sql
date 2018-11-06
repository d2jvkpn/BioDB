-- pigz -dc goa_uniprot_all.b2018-07-16.gaf.gz |
-- awk 'BEGIN{FS=OFS="\t"} $11!=""{sub("taxon:", "", $13);
-- print $5,$8,$9, $11,$13}' > GO.tsv

load data local infile 'GO.tsv' into table BioDB.GO 
(GO_id, prot_id, class, genes, taxon_id);


-- ## or split into parts and load one by one
-- split -l 20000000 GO.tsv --additional-suffix=.GO.tsv
