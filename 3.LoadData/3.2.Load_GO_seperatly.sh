#! /bin/bash

set -eu -o pipefail

date
pigz -dc data_GO/GO.tsv.gz | sed '1d' |
split -l 50000000 --additional-suffix=.GO.tsv
# xaa.GO.tsv xab.GO.tsv...

date

{
  for i in x*.GO.tsv; do
    echo "select concat ('Loading $i, ', now());"
    echo "load data local infile '$i' into table BioDB.GO (taxon_id, genes, GO_id);"
    echo "show warnings;"
    echo
  done

  echo "alter table BioDB.GO add index taxon_id (taxon_id);"
} > load_GO.sql

mysql -u hello < load_GO.sql

date

rm load_GO.sql x*.GO.tsv

exit

locad xaa.GO.tsv
-- Query OK, 49999635 rows affected, 365 warnings (5 min 42.88 sec)
-- Records: 50000000  Deleted: 0  Skipped: 365  Warnings: 365
