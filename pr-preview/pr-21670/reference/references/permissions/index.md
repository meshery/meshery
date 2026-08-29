# Permission Keys

> List of default permissions.

Source: /pr-preview/pr-21670/reference/references/permissions/

Permissions are represented as <b>keys</b>, each serving as a unique identifier for a specific permission. One or more keys can be grouped together and assigned to a <b>keychain</b>.<br>

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Customizable Permissions</div>

Default permissions can be easily customized by simply creating your own keychains and roles.</div>


<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Contributing to Permission Keys</div>

To contribute permission keys, see the <a href='/pr-preview/pr-21670/reference/extensibility/authorization/#adding-a-new-permission-key'>Adding a New Permission Key Guide</a>. This guide explains Meshery's extensible authorization system.</div>


<div style="overflow-x:auto;">
<style>
  tbody {
    display: table-row-group;
    vertical-align: middle;
    unicode-bidi: isolate;
    border-color: inherit;
  }

  td {
    padding: 1rem;
    border-top: 2px solid var(--color-primary-dark);
  }

  tr {
    cursor: default;
  }

  th {
    cursor: default;
  }

  @media screen and (min-width: 1366px) {
    table {
      width: 100%;
      table-layout: fixed;
    }

    td {
      word-wrap: break-word;
    }

    #tablesContainer {
      max-width: 1200px;
      margin: auto;
    }
  }


  tr:nth-child(even) {
    background-color: var(--color-primary-medium) !important;
  }
</style>
<div id="tablesContainer"></div>

<script type="module">
  window.onload = function () {
    const xhr = new XMLHttpRequest();
    xhr.open('GET', '/excel/keys.csv', true);
    xhr.responseType = 'arraybuffer';

    xhr.onload = function () {
      if (xhr.status === 200) {
        const data = new Uint8Array(xhr.response);
        const workbook = XLSX.read(data, { type: 'array' });
        const sheet = workbook.Sheets[workbook.SheetNames[0]];
        const jsonData = XLSX.utils.sheet_to_json(sheet, { header: 1 });

        const colHeaders = jsonData[1];

        const categoryIdx = colHeaders.indexOf("Category");
        const functionIdx = colHeaders.indexOf("Function");
        const featureIdx = colHeaders.indexOf("Feature");
        const keyIdIdx = colHeaders.indexOf("Key ID");

        if (categoryIdx === -1 || functionIdx === -1 || featureIdx === -1 || keyIdIdx === -1) {
          console.error("One or more required columns ('Category', 'Function', 'Feature', 'Key ID') not found in the permissions data. Please check the 'keys.csv' file format.");
          return;
        };

        
        jsonData.shift();
        jsonData.shift();

        

        const processedData = jsonData.map(row => [row[categoryIdx], row[functionIdx], row[featureIdx], row[keyIdIdx]]);

        
        const categories = {};
        processedData.forEach(row => {
          if (!row[0]) return; 
          const category = row[0].replace(/\s+/g, '');
          if (!categories[category]) {
            categories[category] = [];
          }
          categories[category].push(row);
        });

        const headers = ["Category", "Key Name", "Description", "Key ID"];

        
        const container = document.getElementById('tablesContainer');
        Object.keys(categories).forEach(category => {
          const table = document.createElement('table');
          table.style.marginBottom = '2rem';

          const thead = document.createElement('thead');
          const headerRow = document.createElement('tr');
          headers.forEach(header => {
            const th = document.createElement('th');
            th.style.textAlign = 'left';
            th.style.padding = '1rem';
            th.textContent = header;
            headerRow.appendChild(th);
          });
          thead.appendChild(headerRow);
          table.appendChild(thead);

          const tbody = document.createElement('tbody');
          categories[category].forEach((row) => {
            const tr = document.createElement('tr');

            
            row.forEach((cell, index) => {
              if (index < 4) {
                const td = document.createElement('td');
                td.style.padding = '1rem';
                td.textContent = cell ?? '';
                tr.appendChild(td);
              }
            });

            tbody.appendChild(tr);
          });
          table.appendChild(tbody);

          const categoryHeader = document.createElement('h2');
          categoryHeader.textContent = `${category} Permissions`;
          container.appendChild(categoryHeader);
          container.appendChild(table);
        });
      } else {
        console.error('Failed to load Excel file! Status code: ' + xhr.status);
      }
    };

    xhr.onerror = function () {
      console.error('Failed to load Excel file!');
    };

    xhr.send();
  };

</script>

<script src="https://cdnjs.cloudflare.com/ajax/libs/xlsx/0.18.4/xlsx.full.min.js"></script>
</div>
