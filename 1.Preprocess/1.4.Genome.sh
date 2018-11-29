#! /bin/bash

set -eu -o pipefail

# https://ftp.ncbi.nih.gov/genomes/refseq/vertebrate_mammalian/Homo_sapiens/latest_assembly_versions/GCF_000001405.38_GRCh38.p12/

wget -c -O NCBI_assembly_summary_refseq.txt https://ftp.ncbi.nih.gov/genomes/refseq/assembly_summary_refseq.txt

awk 'BEGIN{FS=OFS="\t"} { if($5=="na") $5="";
infor=$8"\t"$20"\t"$5"\t"$9"\t"$12"\t"$15"\t"$16}
NR==2{print $6, infor} NR>2{print $6,infor; 
if($7!=$6) print $7,infor}' NCBI_assembly_summary_refseq.txt > NCBI_genome.tsv

wget -c -O Ensembl_species.txt ftp://ftp.ensemblgenomes.org/pub/current/species.txt

awk 'BEGIN{FS=OFS="\t"; print "taxon_id", "organism_name", "path", "assembly", 
"accession", "source"} NR>1{s=$3; sub("Ensembl", "", s);
s=tolower(s); $2="http://"s".ensembl.org/"$2;
print $4,$1,$2,$5,$6,$3}' Ensembl_species_list.txt > Ensembl_genome_notVertebrate.tsv

## ftp://ftp.ensembl.org/pub/release-94/
## http://asia.ensembl.org/info/about/species.html
python3 EnsemblVertebrate.py

awk 'BEGIN{FS=OFS="\t"} NR==1 || FNR>1{print}' Ensembl_genome_notVertebrate.tsv \
Ensembl_genome_Vertebrate.tsv > Ensembl_genome.tsv


{
  awk 'BEGIN{FS=OFS="\t"; 
  print "taxon_id", "organism_name", "URL", "information"}'

{
  awk 'BEGIN{FS=OFS="\t"}NR==1 {for(i=4; i<= NF; i++) a[i]=$i}
  NR>1{line=$1"\t"$2"\t"$3"\t"a[5]" \""$5"\"";
  for(i=5; i<=NF; i++) {if ($i!="") line=line"; "a[i]" \""$i"\""};
  print line}' NCBI_genome.tsv

  awk 'BEGIN{FS=OFS="\t"}NR==1 {for(i=4; i<= NF; i++) a[i]=$i}
  NR>1{line=$1"\t"$2"\t"$3"\t"a[5]" \""$5"\"";
  for(i=5; i<=NF; i++) {if ($i!="") line=line"; "a[i]" \""$i"\""};
  print line}' Ensembl_genome.tsv
} | sort -k1,1n
} > Genome.tsv

go run TSV_fileds_maxlen.go Genome.tsv

rm Ensembl_genome_notVertebrate.tsv Ensembl_genome_Vertebrate.tsv

cd -
