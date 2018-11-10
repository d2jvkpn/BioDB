#! /bin/bash

set -eu -o pipefail

pigz -dc GO.tsv.gz | sed '1d' | split -l 50000000 --additional-suffix=.GO.tsv
# xaa.GO.tsv ... xak.GO.tsv

{
  for i in x*.GO.tsv; do
    echo "load data local infile '$i' into table BioDB.GO (taxon_id, genes, GO_id);"
  done
  echo "alter table BioDB.Taxonomy_GO add index taxon_id (taxon_id);"
} > LoadGO.sql

mysql -u hello < LoadGO.sql
