#!/bin/bash

find ./out -mindepth 1 -maxdepth 1 -type d -exec rm -r {} \;
./bin/wbtpl -db local/org.csv -tpl local/tpl.html -out ./out -days 7
