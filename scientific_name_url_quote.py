import pandas as pd
from urllib.parse import quote

td = pd.read_csv("Taxonomy0.tsv", sep="\t", header=0, index_col=0)

td["escape_name"] = [quote(i.lower()) for i in td["scientific_name"]]

td.to_csv("Taxonomy.tsv", sep="\t")
