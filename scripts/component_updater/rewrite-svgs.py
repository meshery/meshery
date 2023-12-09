import csv
import argparse
import json


argParser = argparse.ArgumentParser()
argParser.add_argument("-f", "--file", help="CSV file with components")

args = argParser.parse_args()
print("args=%s" % args)

print("args.name=%s" % args.file)

# Replace 'your_input_file.csv' with the actual path to your CSV file
input_file_path = args.file

# Replace 'your_output_file.csv' with the desired output file path
output_file_path = 'your_output_file.csv'

# Read from the CSV file, modify the "svgColor" and "svgWhite" properties, and write to the output file
with open(input_file_path, 'r', encoding='utf-8') as input_file, open(output_file_path, 'w', encoding='utf-8', newline='') as output_file:
    csv_reader = csv.reader(input_file)
    csv_writer = csv.writer(output_file)

    # Write the header to the output file
    header = next(csv_reader)
    csv_writer.writerow(header)

    for row in csv_reader:
        # Assuming the JSON string is in the second column
        json_data = json.loads(row[1])

        print("before"+json_data["metadata"]["svgColor"])
        # Encode the "svgColor" property in UTF-8
        json_data["metadata"]["svgColor"] = json_data["metadata"]["svgColor"].encode('utf-8').decode("unicode-escape")
        
        print("after"+json_data["metadata"]["svgColor"])
        
        # print("encode:")
        # print(x)

        # x = json_data["metadata"]["svgColor"].encode('utf-8').decode('utf-8')

        # print(x)
        # # Encode the "svgWhite" property in UTF-8
        # json_data["metadata"]["svgWhite"] = json_data["metadata"]["svgWhite"].encode('utf-8').decode('utf-8')

        # Update the second column with the modified JSON string
        row[1] = json.dumps(json_data, ensure_ascii=False)

        # Write the modified row to the output file
        csv_writer.writerow(row)

print("Modification completed. Check the output file at:", output_file_path)
