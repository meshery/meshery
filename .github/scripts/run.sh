#! /bin/bash

while IFS="," read line
do
  echo "$line"
 
  
done < <(cut -d "," -f26,3 ~/test.csv | tail -n +4)

