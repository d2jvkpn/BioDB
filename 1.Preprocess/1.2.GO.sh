#! /bin/bash

set -eu -o pipefail

mkdir -p data_GO; cd data_GO

#### Gene Ontology
wget -c ftp://ftp.ebi.ac.uk/pub/databases/GO/goa/UNIPROT/goa_uniprot_all.gaf.gz
wget -c http://purl.obolibrary.org/obo/go.obo

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

sed 's/: /\t/' go.obo  |
awk 'BEGIN{FS=OFS="\t"; print "id", "name", "namespace", "def" > "go.def.tsv";
print "id", "relation", "target" > "go.rel.tsv"}
$1=="id"{id=$2} $1=="name"{name=$2} $1=="namespace"{ns=$2}
$1=="def"{print id, name, ns, $2 > "go.def.tsv"; next}
NF==2 && id!="" && $1!="id" && $1!="name" && $1!="namespace" && $1!="def" &&
id~"^GO:" {if($2~"GO:" && $2~" ! ") sub(" ! .*$", "", $2);
print id, $1, $2 > "go.rel.tsv"}'

# TSV_fileds_maxlen go.def.tsv
# GO_id, name, namespace, definition
# 10, 287, 18, 1364
# 10, 256, 18, 2048

# TSV_fileds_maxlen go.rel.tsv
# GO_id, relation, target
# 10, 15, 1943
# 10, 32, 2048

cd -
