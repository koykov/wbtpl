#!/bin/bash

rm -f out/*.html
./bin/wptpl -db local/org.csv -tpl local/tpl.html -out ./out
