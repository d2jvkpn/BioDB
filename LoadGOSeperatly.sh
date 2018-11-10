#! /bin/bash

set -eu -o pipefail

date
pigz -dc GO.tsv.gz | sed '1d' | split -l 50000000 --additional-suffix=.GO.tsv
# xaa.GO.tsv ... xak.GO.tsv

date
{
  for i in x*.GO.tsv; do
    echo "load data local infile '$i' into table BioDB.GO (taxon_id, genes, GO_id);"
  done
  echo "alter table BioDB.Taxonomy_GO add index taxon_id (taxon_id);"
} > LoadGO.sql

mysql -u hello < LoadGO.sql
date

exit

-- Query OK, 49999635 rows affected, 365 warnings (5 min 42.88 sec)
-- Records: 50000000  Deleted: 0  Skipped: 365  Warnings: 365

show warnings;
