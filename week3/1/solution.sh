#!/bin/bash

for i in {1..20}
do
    input_file="./in/input$i.txt"
    output_file="./out/output$i.txt"
    ans_file="./ans/ans$i.txt"
    diff_file="./diff-res/res$i.txt"
    
    bash ./code.sh < "$input_file" > "$output_file"
    diff -b "$ans_file" "$output_file" > "$diff_file"
      
done
