#! /bin/bash

set -eu -o pipefail

mkdir data_Pathway; cd data_Pathway

awk 'BEGIN{FS=OFS="\t"; print "name", "taxon_id"}
FNR>1{print $2,$1}' Taxonomy.tsv \
../data_Taxonomy/Homotypic_synonym.0.tsv > name2taxon_id.tsv

sed '1d' KEGG_organism.tsv | cut -f2,3 | sed 's/ (.*$//' |
awk 'BEGIN{FS=OFS="\t"} NR==FNR{if(NR>1) a[$1]=$2; next}
!a[$2]{print}' name2taxon_id.tsv - > orgcode_notmatch.tsv

sed '1d' KEGG_organism.tsv | cut -f2,3 | sed 's/ (.*$//' |
awk 'BEGIN{FS=OFS="\t"; print "orgcode", "taxon_id"}
NR==FNR{if(NR>1) a[$1]=$2; next}
$1=="osa" && !a[$2]{print $1, 4530}
a[$2]{print $1,a[$2]}' name2taxon_id.tsv - > orgcode2taxon_id.tsv

rm name2taxon_id.tsv

mkdir Pathway_keg
## https://github.com/d2jvkpn/BioinformaticsAnalysis/blob/master/Pathway/Download_All_Pathway.sh
tar -xf Pathway_keg.tar -C Pathway_keg

sed '1d' orgcode2taxon_id.tsv | while read c t; do
    keg=Pathway_keg/${c}00001.keg.gz
    test -f $keg || { echo "$keg not available" 1>&2; continue; }

	Pathway tsv $keg | awk -v t=$t 'BEGIN{FS=OFS="\t"}
    $1~"^PATH:"{sub("PATH:", "", $1); print t,$1,$2,$3,$4,$5,$6}'
done | awk 'BEGIN{FS=OFS="\t"; print "taxon_id", "pathway_id", "gene_id",
"gene_information","KO_id", "KO_information", "EC_ids"} {print}' > Pathway.tsv

# awk -F "\t" '{if(a<length($4)) a=length($4)} END{print a}' Pathway.tsv
# 878

sed '1d' orgcode2taxon_id.tsv | while read c t; do
    keg=Pathway_keg/${c}00001.keg.gz
    test -f $keg || { echo "$keg not available" 1>&2; continue; }

	Pathway tsv $keg 
done | awk 'BEGIN{FS=OFS="\t"} /^#C/{sub("^#", "", $1);
if(++x[$1]==1) print}' > Pathway_Definition.0.tsv

awk 'BEGIN{FS=OFS="\t"} {print $1,$2,$5} ++x[$3]==1{print $3,$4,""}
++x[$5]==1{print $5,$6,$3}' Pathway_Definition.0.tsv |
sort | sed '1i id\tname\tparent_id' > Pathway_Definition1.tsv

awk 'BEGIN{FS=OFS="\t"} NR==1{print "C_id", "C_name", "B_id", "B_name",\
"A_id", "A_name"; next} $1~"^A"{a[$1]=$2} $1~"^B"{a[$1]=$2; b[$1]=$3}
$1~"^C"{print $1,$2,$3,a[$3],b[$3],a[b[$3]]}' Pathway_Definition1.tsv \
> Pathway_Definition.tsv

rm Pathway_Definition.0.tsv

cd -
