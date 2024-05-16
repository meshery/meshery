---
layout: default
title: Permissions
permalink: reference/permissions
language: en
type: Reference
abstract: List of default permissions.
---

<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Reference Permissions</title>
</head>
<body>
Permissions are represented as <b>keys</b>, each serving as a unique identifier for a specific permission. One or more keys can be grouped together and assigned to a <b>keychain</b>.<br>

<table style="margin:auto;padding-right:25%; padding-left:20%; margin-top:3rem;" id="excelTable" border="1">
  <thead>
  <tr>
    <th style="text-align:left">Category</th>
    <th style="text-align:left">Function</th>
    <th style="text-align:left">Feature</th>
    <th style="text-align:left">Keychain ID</th>
    <th style="text-align:left">Key ID</th>
  </tr>
</thead>
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

      // Remove the first, second row
      const tableElement = document.createElement('div');
      tableElement.innerHTML = htmlTable;
      tableElement.querySelector('table').deleteRow(0);
      tableElement.querySelector('table').deleteRow(0);

      // Remove the columns
      const rows = tableElement.querySelectorAll('tr');
      rows.forEach(row => {
        row.deleteCell(3);
        row.deleteCell(3);
        row.deleteCell(3);
        row.deleteCell(3);
        row.deleteCell(3);
        row.deleteCell(3);
        row.deleteCell(3);
        row.deleteCell(3);
        row.deleteCell(4);
        row.deleteCell(5);
        row.deleteCell(6);
        row.deleteCell(5);
      });
      
      document.getElementById('excelTable').innerHTML += tableElement.innerHTML;
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
