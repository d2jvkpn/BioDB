#! /bin/bash

## Taxonomy, 2018-11-07, 2003451 records
mkdir -p data_Taxonomy; cd data_Taxonomy

wget -c https://ftp.ncbi.nlm.nih.gov/pub/taxonomy/new_taxdump/new_taxdump.tar.gz
wget -O new_taxdump_readme.txt \
https://ftp.ncbi.nlm.nih.gov/pub/taxonomy/new_taxdump/taxdump_readme.txt

tar -xf new_taxdump.tar.gz

sed 's/\t|\t/\t/g; s/|$//' names.dmp |
awk 'BEGIN{FS=OFS="\t"; print "taxon_id", "name"} 
$4=="synonym"{print $1, $2}' > Homotypic_synonym.0.tsv

python3 homotypic_synonym_urlquote.py

sed 's/\t|\t/\t/g; s/|$//' names.dmp |
awk 'BEGIN{FS=OFS="\t"} $4=="scientific name"{print $1, $2}' > id2name.txt

sed 's/\t|\t/\t/g; s/|$//' nodes.dmp |
awk 'BEGIN{FS=OFS="\t"} {print $1,$3,$2}' > id_rank_parent.txt

paste id2name.txt id_rank_parent.txt |
awk -F "\t" '$1!=$3'

paste id2name.txt id_rank_parent.txt | awk 'BEGIN{FS=OFS="\t"; 
print "taxon_id", "scientific_name", "taxon_rank", "parent_id"}
{print $1,$2,$4,$5}' > Taxonomy.0.tsv

# python3 scientific_name_url_quote.py
## urllib.parser.quote converts " " with "%20", "+" with "%2B"
## urllib.quote_plus converts " " with "+", "+" with "%2B"
## pandas with add double quote to filed when the field contains a quote inside
##    Nostoc sp. 'Peltigera sp. "hawaiensis" P1236 cyanobiont'
##    "Nostoc sp. 'Peltigera sp. ""hawaiensis"" P1236 cyanobiont'"

## Golang net/url.QueryEscape converts like

go run scientific_name_QueryEscape.go

rm -r id2name.txt id_rank_parent.txt Taxonomy.0.tsv

cd -
