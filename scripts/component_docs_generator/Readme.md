# Meshery Integration Docs Generator README

This Golang script is designed to fetch data from integration docs csv, filter records, generate integration README files, and download SVG images. The README files are created for each integration with overview and feature details, and the resulting files are stored in specific directories.

## Prerequisites

To run this script in local, ensure you have the following prerequisites:

- [Go (Golang)](https://golang.org/dl/)
- Integration docs spreadsheet URL

## Installation

1. Fork the Meshery repository and clone the repo to your local:

   ```shell
   git clone https://github.com/yourusername/meshery.git
   ```

2. Change to the project directory:

   ```shell
   cd meshery/scripts/component_docs_generator/
   ```

## Usage

Run the `main.go` file with the spreadsheet URL as an argument to execute the script.

```shell
go run main.go "https://example.com/spreadsheet-url"
```

## Workflow

1. The script fetches csv from the provided spreadsheet URL and stores it locally.
2. Records are filtered based on the "Publish?" flags. Records with true flags are included, while false flags are removed.
3. SVG images are downloaded based on the URLs provided in the CSV.
4. README files are generated for each integration, including overview and feature details.
5. Integration README files are stored in the `../../docs/pages/integrations` directory.
6. Downloaded SVG images are saved in the `../../docs/assets/img/integrations` directory.
7. Temporary CSV files are deleted after processing.

### Configuration

This file `index.json` have the parameter to column of the csv file:

```json
{
  "flag-index": 35,
  "svg-index": 16,
  "index-start": 24,
  "index-end": 33,
  "name-index": 0
}
```


### Flags

- `-flag-index`: The column index containing the "Publish?" flags in the CSV. Records with true flags will be included in README generation.
- `-svg-index`: The column index containing URLs to SVG images for integration.
- `-index-start` and `-index-end`: The column indices specifying the range for feature details in the CSV.
- `-name-index`: The column index for integration names.

## File Locations

- Integration READMEs: `../../docs/pages/integrations`
- Downloaded SVG images: `../../docs/assets/img/integrations`

Please replace placeholders (e.g., `https://example.com/spreadsheet-url`, `yourusername/your-repo`) with the actual details of your project and spreadsheet URL.