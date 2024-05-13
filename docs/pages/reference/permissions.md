---
layout: default
title: Permissions
permalink: reference/permissions
redirect_from: reference/permissions
language: en
list: exclude
abstract: List of default permissions.
---

<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Excel to HTML Table</title>
</head>
<body>

<table id="excelTable" border="1">
  
</table>

<script>

window.onload = function() {
  const xhr = new XMLHttpRequest();
  xhr.open('GET', '../../assets/excel/keys.csv', true);
  xhr.responseType = 'arraybuffer';

  xhr.onload = function() {
    if (xhr.status === 200) {
      const data = new Uint8Array(xhr.response);
      const workbook = XLSX.read(data, { type: 'array' });
      const sheet = workbook.Sheets[workbook.SheetNames[0]];
      const htmlTable = XLSX.utils.sheet_to_html(sheet);

      // Remove the first row from the HTML table
      const tableElement = document.createElement('div');
      tableElement.innerHTML = htmlTable;
      tableElement.querySelector('table').deleteRow(0);
      // Remove the second, fourth column from the HTML table
      const rows = tableElement.querySelectorAll('tr');
      rows.forEach(row => {
        row.deleteCell(1);
        row.deleteCell(2);
      });
      const modifiedHtmlTable = tableElement.innerHTML;

      document.getElementById('excelTable').innerHTML = modifiedHtmlTable;
    } else {
      console.error('Failed to load Excel file! Status code: ' + xhr.status);
    }
  };

  xhr.onerror = function() {
    console.error('Failed to load Excel file!');
  };

  xhr.send();
};
</script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/xlsx/0.18.4/xlsx.full.min.js"></script>
</body>
</html>