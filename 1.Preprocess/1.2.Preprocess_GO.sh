#! /bin/bash

set -eu -o pipefail

mkdir -p data_GO; cd data_GO

#### Gene Ontology
wget -c ftp://ftp.ebi.ac.uk/pub/databases/GO/goa/UNIPROT/goa_uniprot_all.gaf.gz

pigz -dc goa_uniprot_all.gaf.gz | awk 'BEGIN{FS=OFS="\t"}
$11!="" && $11!~" "{sub("taxon:", "", $13); 
if ($13!~"|") { print $13, $11, $5; next };
split($13,x,"|"); for(i in x) print x[i], $11, $5}' | uniq > GO0.tsv
# 379,646,174 records

sort -k1,1n -t $'\t' --parallel=16 GO0.tsv | uniq |
awk 'BEGIN{print "taxon_id", "genes", "GO_id"} {print}' | pigz -c > GO.tsv.gz
rm GO0.tsv

pigz -dc GO.tsv.gz | awk 'BEGIN{a=0} 
NR>1{if(length($2)>a) a=length($2)} END{print a, NR}'
# 538 350040327

cd -
