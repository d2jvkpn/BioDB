import pandas as pd
from urllib.parse import quote_plus
import csv

td = pd.read_csv("Taxonomy.0.tsv", sep="\t", header=0, index_col=0, dtype=str)
td["escape_name"] = [quote_plus(i.lower()) for i in td["scientific_name"]]
td.to_csv("Taxonomy.tsv", sep="\t",  quoting = csv.QUOTE_NONE)
