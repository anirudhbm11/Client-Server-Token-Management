#! /bin/sh

client -create -id 1234

client -create -id 41

client -read -id 41

client -read -id 1234

client -create -id 100 

client -write -id 100 -name abc -low 0 -mid 10 -high 100

client -read -id 100

client -write -id 100 -name abc -low 20 -mid 50 -high 200

client -read -id 1234

