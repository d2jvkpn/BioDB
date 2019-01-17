#! /usr/bin/python3

__author__ = 'd2jvkpn'
__version__ = '1.4'
__release__ = '2019-01-17'
__project__ = 'https://github.com/d2jvkpn/BioDB'
__license__ = 'GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)'

import os
import requests
from urllib.parse import urlparse
import string
import time
from ftplib import FTP
from bs4 import BeautifulSoup
import pandas as pd
from biomart import BiomartServer

'''
http://www.ensembl.org/index.html
http://pre.ensembl.org/index.html
http://asia.ensembl.org/index.html

http://metazoa.ensembl.org/index.html
http://plants.ensembl.org/index.html
http://fungi.ensembl.org/index.html
http://bacteria.ensembl.org/index.html
http://protists.ensembl.org/index.html

http://asia.ensembl.org/info/about/species.html
http://metazoa.ensembl.org/species.html
http://plants.ensembl.org/species.html
http://fungi.ensembl.org/species.html
http://bacteria.ensembl.org/species.html
http://protists.ensembl.org/species.html
'''

HELP = '''
Search species genome in Ensembl by providing scientific name, and get genome
  ftp  download links and achive gene annotation (GO, kegg_enzyme, entrez) 
  from biomart by provide Ensembl genome address.

Usage:
    python3 BioDB_Ensemblgenome.py search "species scientific name"
    e.g. "Mus musculus", Mus_musculus

    python3 BioDB_Ensemblgenome.py getftp/biomart "ensembl genome address"
    e.g. http://asia.ensembl.org/Mus_musculus/Info/Index

    python3 BioDB_Ensemblgenome.py search getftp biomart "species scientific name"

Note:
    Please use Python3.6 or higher.
'''

if len(os.sys.argv) == 1 or os.sys.argv[1] in ['-h', '--help']:
    print(HELP)

    _ = '\nauthor: {}\nversion: {}\nrelease: {}\nproject: {}\nlicense: {}\n'
    __ = [__author__,  __version__, __release__, __project__, __license__]
    print (_.format (*__))

    os.sys.exit(0)


####
def formatSpeciesName(s):
    wds = s.replace ("+", " ").split ()

    for i in range(len(wds)):
        a = False not in [i in string.ascii_letters for i in wds[i]]
        b = False not in [i in string.ascii_uppercase for i in wds[i]]
        if a and not b: wds[i] = wds[i].lower()

    wds[0] = wds[0].capitalize()
    return(' '.join(wds))


def query(url):
    query = requests.get(url)
    if not query.ok: os.sys.exit('Failed to request "%s"' % url)
    bs = BeautifulSoup(query.text, 'html.parser')

    _  = bs.find('span', class_ = 'header').select('a')[0]
    version = _.text.strip(')').split(' (', 1)[1].replace(") (", "_")

    dnaFTP = ''

    for i in bs.find_all('a', class_='nodeco'):
        if i.get('href').startswith('') \
        and i.text == 'Download DNA sequence': dnaFTP = i.get('href')

    netloc = urlparse(dnaFTP).netloc
    path = urlparse(dnaFTP).path

    _ = path.strip("/").split("/")
    ensembl = _[_.index('fasta')-1].replace('release', 'Ensembl')
    ScentificName = _[-2]
    ScentificName = ScentificName[0].upper() + ScentificName[1:]

    loca = '__'.join([ensembl, ScentificName, version])
    return(netloc, path, loca)


def getftp(netloc, path, loca, url):
    at = time.strftime('%Y-%m-%d %H:%M:%S %z')

    ensembl, ScentificName, version = loca.split('__')

    ftp = FTP(netloc)
    ftp.login()
    
    for i in ftp.nlst(path):
        if i.endswith('.dna_sm.toplevel.fa.gz'): dna = 'ftp://' + netloc + i
    
    for i in ftp.nlst(path.replace('/dna/', '/cdna/')):
        if i.endswith('.all.fa.gz'): cdna = 'ftp://' + netloc + i
    
    for i in ftp.nlst(path.replace('/dna/', '/ncrna/')):
        if i.endswith('.ncrna.fa.gz'): ncrna = 'ftp://' + netloc + i
    
    for i in ftp.nlst(path.replace('/dna/', '/pep/')):
        if i.endswith('.pep.all.fa.gz'): pep = 'ftp://' + netloc + i
    
    for i in ftp.nlst(path.replace('/dna/', '').replace('/fasta/', '/gtf/')):
        if i.endswith('.gtf.gz') and not i.endswith('.abinitio.gtf.gz'):
            gtf = 'ftp://' + netloc + i
    
    ftp.close()
    
    os.system('mkdir -p %s'  % loca)
    
    with open (loca + '/genome.infor.txt', 'w') as f:
        f.write('URL: %s\n' % url)
        f.write('Acess time: %s\n' % at)
        f.write('Scentific name: %s\n' % ScentificName.replace('_', ' '))
        f.write('Assembly version: %s\n' % version)
        f.write('Ensembl version: %s\n\n' % ensembl)
        if 'dna' in locals(): f.write('DNA fasta:\n    %s\n\n' % dna)
        if 'cdna' in locals(): f.write('cdna fasta:\n    %s\n\n' % cdna)
        if 'ncrna' in locals(): f.write('ncrna fasta:\n    %s\n\n' % ncrna)
        if 'pep' in locals(): f.write('pep fasta:\n    %s\n\n' % pep)
        if 'gtf' in locals(): f.write('annotation gtf:\n    %s\n' % gtf)
    
    wget = 'wget -c -O {0} {1} -o {0}.download.logging &&\nrm {0}.download.logging'

    with open (loca + '/download.sh', 'w') as f:
        f.write('#! /bin/bash\n\n## URL: %s\n' % url)
        f.write('## Species: %s\n' % ScentificName.replace('_', ' '))
        f.write('## Acess time: %s\n' % at)
        if 'dna' in locals(): f.write("\n{\n" + wget.format('genomic.fa.gz', dna) + "\n} &\n")
        if 'cdna' in locals(): f.write("\n{\n" + wget.format('cdna.fa.gz', cdna) + "\n} &\n")
        if 'ncrna' in locals(): f.write("\n{\n" + wget.format('ncrna.fa.gz', ncrna) + "\n} &\n")
        if 'pep' in locals(): f.write("\n{\n" + wget.format('pep.fa.gz', pep) + "\n} &\n")
        if 'gtf' in locals(): f.write("\n{\n" + wget.format('genomic.gtf.gz', gtf) + "\n} &\n")
        f.write('\nwait\n')

    print(loca)


def biomart_anno(url, loca):
    urlp = urlparse(url)

    species = urlp.path.split('/')[1]
    code = species.split('_')[0][0].lower() + species.split('_')[1]

    server = BiomartServer('%s://%s/biomart' % (urlp.scheme, urlp.netloc))
    datasets = server.datasets

    print("Connecting to Ensembl biomart...")
    
    _ = ['metazoa', 'plants', 'fungi', 'bacteria', 'protists']
    dn = code + ('_eg_gene'  if urlp.netloc.split('.')[0] in _ else '_gene_ensembl')
    ds = datasets[dn]

    os.system('mkdir -p %s' % loca)

    S = ds.search({'attributes': ['ensembl_gene_id', 'gene_biotype', \
    'external_gene_name', 'description', 'chromosome_name', 'start_position', \
    'end_position', 'strand']})

    gene_infor = pd.DataFrame.from_records(
    [str(i, encoding = 'utf-8').split('\t') for i in S.iter_lines()], \
    columns = ['gene', 'gene_biotype', 'gene_name', 'gene_description', \
    'chromosome_name', 'start_position', 'end_position', 'strand'])

    m = {"-1":"-", "1":"+"}
    gene_infor['strand'] = [ m[i] if i in m else str(i) for i in gene_infor['strand']]

    gene_infor['gene_position'] = gene_infor.loc[:, ['chromosome_name', \
    'start_position', 'end_position', 'strand']].apply(\
    lambda x: ':'.join(x), axis=1)

    gene_infor.drop(['chromosome_name', 'start_position', 'end_position', \
    'strand'], axis = 1, inplace=True)

    ####
    gene2GO = getGO(ds, loca)
    if gene2GO.shape[0] != 0:
        g = gene2GO.groupby('gene')['GO_id'].apply(lambda x: ', '.join(x))
        gene_infor['GO_id'] = [ g[i] if i in g else '' for i in gene_infor['gene']]

    ####
    gene2kegg = getKEGG(ds, loca)
    if gene2kegg.shape[0] != 0:
        k = gene2kegg.groupby('gene')['KEGG_enzyme'].apply(lambda x: ', '.join(x))
        gene_infor['KEGG_enzyme'] = [ k[i] if i in k else '' for i in gene_infor['gene']]


    ####
    gene2entrez = getEntrezgene(ds, loca)
    if gene2entrez.shape[0] != 0:
        e = gene2entrez.groupby('gene')['entrez'].apply(lambda x: ', '.join(x))
        gene_infor['entrez'] = [ e[i] if i in e else '' for i in gene_infor['gene']]


    ####
    gene2swissprot = getSwissProt(ds, loca)
    if gene2swissprot.shape[0] != 0:
        s = gene2swissprot.groupby('gene')['SwissProt'].apply(lambda x: ', '.join(x))
        gene_infor['SwissProt'] = [ s[i] if i in s else '' for i in gene_infor['gene']]

    gene_infor.to_csv(loca + '/gene.infor.tsv', sep='\t', index=False)

    print('Saved %d records to %s/gene.infor.tsv' % (gene_infor.shape[0], loca))


def search (name):
    msg = 'NotFound'
    name = name.replace(' ', '_').replace('-', '_')
 
    UA = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 ' + \
    '(KHTML, like Geko) chrom/61.0.3163.100'

    for i in ['www', 'pre', 'metazoa', 'plants', 'fungi', 'bacteria', 'protists']:
        url = 'http://%s.ensembl.org/%s' % (i, name)
        
        query = requests.get(url,
        headers = {'User-Agent': UA, 'Referer': 'http://asia.ensembl.org'})

        # print('Search "%s" in http://%s.ensembl.org' % (name, i))
        if query.status_code != 200: continue
        msg = query.url; break

    return (msg)


def getGO(ds, loca):
    gene2GO = pd.DataFrame()

    try:
        s = ds.search({'attributes': ['ensembl_gene_id', 'go_id', 'name_1006', \
        'namespace_1003']})
    except:
        return gene2GO


    gene2GO = pd.DataFrame.from_records(
    [str(i, encoding = 'utf-8').split('\t') for i in s.iter_lines()], 
    columns = ['gene', 'GO_id', "name", "namespace"])
    
    gene2GO = gene2GO.loc[gene2GO['GO_id'] != '', :]
    gene2GO.drop_duplicates(inplace=True)
    gene2GO.to_csv(loca + '/gene2GO.tsv', sep='\t', index=False)

    print('Saved %d records to %s/gene2GO.tsv' % (gene2GO.shape[0], loca))

    return gene2GO

def getEntrezgene(ds, loca):
    gene2entrez = pd.DataFrame()
    try:
        s = ds.search({'attributes': ['ensembl_gene_id', 'entrezgene']})
    except:
        return gene2entrez

    gene2entrez = pd.DataFrame.from_records(
    [str(i, encoding = 'utf-8').split('\t') for i in s.iter_lines()], 
    columns = ['gene', 'entrez'])
    
    gene2entrez = gene2entrez.loc[gene2entrez['entrez'] != '', :]
    gene2entrez.drop_duplicates(inplace=True)
    gene2entrez.to_csv(loca + '/gene2entrez.tsv', sep='\t', index=False)

    print('Saved %d records to %s/gene2entrez.tsv' % (gene2entrez.shape[0], loca))
    return gene2entrez


def getKEGG(ds, loca):
    gene2kegg = pd.DataFrame()

    try:
        s = ds.search({'attributes': ['ensembl_gene_id', 'kegg_enzyme']})
    except:
         return gene2kegg

    gene2kegg = pd.DataFrame.from_records(
    [str(i, encoding = 'utf-8').split('\t') for i in s.iter_lines()], 
    columns = ['gene', 'KEGG_enzyme'])

    gene2kegg = gene2kegg.loc[gene2kegg['KEGG_enzyme'] != '', :]
    gene2kegg.to_csv(loca + '/gene2KEGG_enzyme.tsv', sep='\t', index=False)

    print('Saved %d records to %s/gene2KEGG_enzyme.tsv' % (gene2kegg.shape[0], loca))
    return gene2kegg


def getSwissProt(ds, loca):
    gene2swissprot = pd.DataFrame()

    try:
        s = ds.search({'attributes': ['ensembl_gene_id', 'uniprotswissprot']})
        # uniprotsptrembl
    except:
        return gene2swissprot

    gene2swissprot = pd.DataFrame.from_records(
    [str(i, encoding = 'utf-8').split('\t') for i in s.iter_lines()], 
    columns = ['gene', 'SwissProt'])

    gene2swissprot = gene2swissprot.loc[gene2swissprot['SwissProt'] != '', :]
    gene2swissprot.drop_duplicates(inplace=True)
    gene2swissprot.to_csv(loca + '/gene2swissprot.tsv', sep='\t', index=False)

    print('Saved %d records to %s/gene2SwissProt.tsv' % (gene2swissprot.shape[0], loca))
    return gene2swissprot


args = os.sys.argv
val = args[-1]

for cmd in args[1:-1]:
    if cmd == 'search':
        out =  search (formatSpeciesName (val))
        if out != "NotFound":
            print(out); val = out
        else:
            os.sys.exit(out)

    elif cmd == 'getftp':
        netloc, path, loca = query(val)
        getftp (netloc, path, loca, val)

    elif cmd == 'biomart':
        netloc, path, loca = query(val)
        biomart_anno(val, loca)
        
    else:
        os.sys.exit("invalid subcommand {}".format(cmd))
