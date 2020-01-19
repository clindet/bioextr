#!/bin/bash
# extract from pubmed abstract
bget api ncbi -q "Galectins control MTOR and AMPK in response to lysosomal damage to induce autophagy OR MTOR-independent autophagy induced by interrupted endoplasmic reticulum-mitochondrial Ca2+ communication: a dead end in cancer cells. OR The PARK10 gene USP24 is a negative regulator of autophagy and ULK1 protein stability OR Coordinate regulation of autophagy and the ubiquitin proteasome system by MTOR." | bioctl cvrt --xml2json pubmed - | bioextr --mode pubmed -w 'MTOR,AMPK,autophagy' --call-cor - > final.json

bget api ncbi -q "Galectins control MTOR and AMPK in response to lysosomal damage to induce autophagy OR MTOR-independent autophagy induced by interrupted endoplasmic reticulum-mitochondrial Ca2+ communication: a dead end in cancer cells. OR The PARK10 gene USP24 is a negative regulator of autophagy and ULK1 protein stability OR Coordinate regulation of autophagy and the ubiquitin proteasome system by MTOR." | bioctl cvrt --xml2json pubmed - > test0.json

for i in {1..100}
do
cp test0.json test${i}.json
done

bioextr test*json --mode pubmed -w 'MTOR,AMPK,autophagy' --call-cor -t 30 > final2.json

rm final.json final2.json test*json

# extract from sra json
bget api ncbi -d 'sra' -q PRJNA527714 | bioctl cvrt --xml2json sra - | bioextr --mode sra --call-cor -w "Chromatin,mouse" -

