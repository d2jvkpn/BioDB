#! /bin/bash

# wget -c ftp://ftp.ebi.ac.uk/pub/databases/GO/goa/UNIPROT/goa_uniprot_all.gaf.gz

# wget ftp://ftp.ncbi.nlm.nih.gov/pub/taxonomy/taxdump.tar.gz
# mkdir taxdump
# tar -xf taxdump.tar.gz -C taxdump

# 130542438 records
pigz -dc goa_uniprot_all.gaf.gz | go run Table_GO_v0.4.go


# 1858425 records
sed 's/\t|\t/\t/g; s/|$//' taxdump/names.dmp |
awk 'BEGIN{FS=OFS="\t"} $4=="scientific name"{print $1, $2}' |
go run Table_Taxon_v0.2.go
