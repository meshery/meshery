import React, { useRef } from "react";
import { Avatar, IconButton, TableCell, TableSortLabel } from "@mui/material";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import DoneAllIcon from "@mui/icons-material/DoneAll";
import EditIcon from "@mui/icons-material/Edit";
import UndeployIcon from "../../public/static/img/UndeployIcon";

function MesheryFilterTable({
  filters = [],
  setSelectedRowData,
  handleModalOpen,
  user,
  page = 0,
  count = 0,
  pageSize = 10,
  sortOrder,
  fetchFilters,
  showModal,
  handleDelete,
  search = "",
  setSearch,
}) {
  const searchTimeout = useRef(null);
  const columns = [
    {
      name: "name",
      label: "Name",
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: "created_at",
      label: "Upload Timestamp",
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: "updated_at",
      label: "Update Timestamp",
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: "Actions",
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(_, tableMeta) {
          const rowData = filters[tableMeta.rowIndex];
          return (
            <>
              <IconButton>
                <EditIcon
                  title="Config"
                  aria-label="config"
                  color="inherit"
                  onClick={() => setSelectedRowData(filters[tableMeta.rowIndex])}
                />
              </IconButton>
              <IconButton title="Deploy" onClick={() => handleModalOpen(rowData.filter_file, true)}>
                <DoneAllIcon data-cy="deploy-button" />
              </IconButton>
              <IconButton title="Undeploy" onClick={() => handleModalOpen(rowData.filter_file, false)}>
                <UndeployIcon fill="rgba(0, 0, 0, 0.54)" data-cy="undeploy-button" />
              </IconButton>
            </>
          );
        },
      },
    },
  ];

  const options = {
    filter: false,
    sort: !(user && user.user_id === "meshery"),
    search: !(user && user.user_id === "meshery"),
    filterType: "textField",
    // responsive : "scrollFullHeight",
    resizableColumns: true,
    serverSide: true,
    count,
    rowsPerPageOptions: [10, 20, 25],
    fixedHeader: true,
    page,
    print: false,
    download: false,
    textLabels: {
      selectedRows: {
        text: "filter(s) selected",
      },
    },
    onCellClick: (_, meta) => meta.colIndex !== 3 && setSelectedRowData(filters[meta.rowIndex]),

    onRowsDelete: async function handleDeleteRow(row) {
      let response = await showModal(Object.keys(row.lookup).length);
      console.log(response);
      if (response === "Yes") {
        const fid = Object.keys(row.lookup).map((idx) => filters[idx]?.id);
        fid.forEach((fid) => handleDelete(fid));
      }
      if (response === "No") fetchFilters(page, pageSize, search, sortOrder);
    },

    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(" : ") : [];
      let order = "";
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case "changePage":
          setPage(tableState.page);
          break;
        case "changeRowsPerPage":
          setPageSize(tableState.rowsPerPage);
          break;
        case "search":
          if (searchTimeout.current) clearTimeout(searchTimeout.current);
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) setSearch(tableState.searchText);
          }, 500);
          break;
        case "sort":
          if (sortInfo.length == 2) {
            if (sortInfo[1] === "ascending") {
              order = `${columns[tableState.activeColumn].name} asc`;
            } else {
              order = `${columns[tableState.activeColumn].name} desc`;
            }
          }
          if (order !== sortOrder) setSortOrder(order);
          break;
      }
    },
  };

  return <MUIDataTable title={<div>Filters</div>} data={filters} columns={columns} options={options} />;
}

export default MesheryFilterTable;

