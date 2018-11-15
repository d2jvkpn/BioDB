import pandas as pd
from urllib.parse import quote_plus
import csv

td = pd.read_csv("Homotypic_synonym.0.tsv", sep="\t", header=0, index_col=0, dtype=str)
td["name"] = [quote_plus(i.lower()) for i in td["name"]]
td.to_csv("Homotypic_synonym.tsv", sep="\t",  quoting = csv.QUOTE_NONE)
