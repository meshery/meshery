import csv
import os

def create_files_from_csv(csv_file_path, output_folder):
    with open(csv_file_path, 'r') as csv_file:
        csv_reader = csv.reader(csv_file)
        for i, row in enumerate(csv_reader):
            # Assuming the first column of each row contains the filename
            filename = row[0]
            content = row[1] # Joining all columns as content

            # Creating the output folder if it doesn't exist
            if not os.path.exists(output_folder):
                os.makedirs(output_folder)

            # Creating a text file for each line
            output_file_path = os.path.join(output_folder, f'{filename}.json')
            with open(output_file_path, 'w') as output_file:
                output_file.write(content)

            print(f'File {i + 1} created: {output_file_path}')

# Example usage:
csv_file_path = '/Users/lee/Downloads/input - Sheet1.csv'  # Replace with the path to your CSV file
output_folder = 'output_files'  # Replace with the desired output folder

create_files_from_csv(csv_file_path, output_folder)
