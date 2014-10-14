#!/bin/bash
while read -r dep
do 
  echo getting ${dep}
  go get -u ${dep}
done < "deps"
