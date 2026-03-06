touch __intermediate_file.yml
echo $FILE_PATH

append_with_newline() {
  file=$1
  cat "$file" >> __intermediate_file.yml
  if [ -n "$(tail -c 1 "$file")" ]; then
    echo "" >> __intermediate_file.yml
  fi
  echo "---" >> __intermediate_file.yml
}

if [[ -d $FILE_PATH ]]; then
  for file in "$FILE_PATH"/*; do
    if [[ $file == *.yaml || $file == *.yml ]]; then
      append_with_newline "$file"
    fi
  done
else
  append_with_newline "$FILE_PATH"
fi