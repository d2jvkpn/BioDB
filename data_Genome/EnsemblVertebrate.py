import pandas as pd
import csv

d = pd.read_csv("Species.csv", header=0, na_values="-", quotechar="\"")

d1 = d.iloc[:, [1,2,3,4,8]]
d1.columns = ["organism_name", "taxon_id", "assembly", "accession", "pre_assembly"]

d1 = d1.loc[pd.notna(d1.taxon_id), :]
d1.taxon_id = d1.taxon_id.astype(int)

d1["source"] = "EnsemblVertebrate"

suffix = [i.replace(" ", "_").replace("-", "_") for i in d1.loc[:, "organism_name"]]
prefix = ["http://ensembl.org/"  if pd.isna(i) else "http://pre.ensembl.org/" \
for i in d1.pre_assembly]

d1["path"] = [ prefix[i]+suffix[i] for i in range(0, len(suffix))]

d1.loc[:, ["taxon_id", "organism_name", "path", "assembly", "accession", "source"]].\
to_csv("Ensembl_genome_Vertebrate.tsv", sep="\t", index=False,
quoting = csv.QUOTE_NONE, na_rep="")
