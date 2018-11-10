#! /bin/bash
wget -c ftp://ftp.ebi.ac.uk/pub/databases/GO/goa/UNIPROT/goa_uniprot_all.gaf.gz

pigz -dc goa_uniprot_all.gaf.gz | awk 'BEGIN{FS=OFS="\t"}
$11!="" && $11!~" "{sub("taxon:", "", $13); 
if ($13!~"|") { print $13, $11, $5; next };
split($13,x,"|"); for(i in x) print x[i], $11, $5}' | uniq > GO0.tsv
# 350039671

sort -k1,1n -t $'\t' --parallel=16 GO0.tsv | uniq |
awk 'BEGIN{print "taxon_id", "genes", "GO_id"} {print}' | pigz -c > GO.tsv.gz
rm GO0.tsv

pigz -dc GO.tsv.gz | awk 'BEGIN{a=0} 
NR>1{if(length($2)>a) a=length($2)} END{print a, NR}'
# 538 350040327


## 2018-11-07, 2003451 records
wget -c https://ftp.ncbi.nlm.nih.gov/pub/taxonomy/new_taxdump/new_taxdump.tar.gz
mkdir new_taxdump
tar -xf new_taxdump.tar.gz -C new_taxdump
wget -O new_taxdump/taxdump_readme.txt \
https://ftp.ncbi.nlm.nih.gov/pub/taxonomy/new_taxdump/taxdump_readme.txt

sed 's/\t|\t/\t/g; s/|$//' new_taxdump/names.dmp |
awk 'BEGIN{FS=OFS="\t"} $4=="scientific name"{print $1, $2}' > id2name.txt

sed 's/\t|\t/\t/g; s/|$//' new_taxdump/nodes.dmp |
awk 'BEGIN{FS=OFS="\t"} {print $1,$3,$2}' > id_rank_parent.txt

paste id2name.txt id_rank_parent.txt |
awk -F "\t" '$1!=$3'

paste id2name.txt id_rank_parent.txt | awk 'BEGIN{FS=OFS="\t"; 
print "taxon_id", "scientific_name", "taxon_rank", "parent_id"}
{print $1,$2,$4,$5}' > Taxonomy0.tsv

python3 scientific_name_url_quote.py
rm id2name.txt id_rank_parent.txt Taxonomy0.tsv
