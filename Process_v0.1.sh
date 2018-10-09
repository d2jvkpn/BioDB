#! /bin/bash

# wget -c ftp://ftp.ebi.ac.uk/pub/databases/GO/goa/UNIPROT/goa_uniprot_all.gaf.gz

# wget ftp://ftp.ncbi.nlm.nih.gov/pub/taxonomy/taxdump.tar.gz
# mkdir taxdump
# tar -xf taxdump.tar.gz names.dmp

# 130,542,438 records
# 513,365,944
pigz -dc goa_uniprot_all.gaf.gz | go run BuildGO_v0.5.go


# 1,858,425 records
sed 's/\t|\t/\t/g; s/|$//' names.dmp |
awk 'BEGIN{FS=OFS="\t"} $4=="scientific name"{print $1, $2}' |
go run BuildTaxon_v0.2.go
