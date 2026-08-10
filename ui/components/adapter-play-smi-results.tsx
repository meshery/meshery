import React from 'react';
import {
  Dialog,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Tooltip,
} from '@sistent/sistent';
import MUIDataTable from '@sistent/mui-datatables';
import Moment from 'react-moment';
import { AdapterTableHeader, SecondaryTable } from './adapter-play-styled';
import { useGetProviderCapabilitiesQuery } from '@/rtk-query/user';
import { isLocalProvider } from '@/utils/provider';

interface SmiResultRow {
  id?: string;
  date?: string;
  mesh_name?: string;
  mesh_version?: string;
  passing_percentage?: string;
  status?: string;
  more_details?: Array<{
    smi_specification?: string;
    assertions?: string;
    time?: string;
    smi_version?: string;
    capability?: string;
    status?: string;
    reason?: string;
  }>;
}

interface SmiResult {
  results?: SmiResultRow[];
  totalCount?: number;
}

interface AdapterSmiResultsDialogProps {
  open: boolean;
  onClose: () => void;
  adapterName: string;
  smiResult: SmiResult | unknown[];
  page: number;
  pageSize: number;
  search: string;
  sortOrder: string;
  fetchSMIResults: (
    adapterName: string,
    page: number,
    pageSize: number,
    search: string,
    sortOrder: string,
  ) => void;
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
  onSearchChange: (search: string) => void;
  onSortOrderChange: (sortOrder: string) => void;
}

/**
 * Renders the Service Mesh Interface (SMI) conformance results
 * dialog. Extracted verbatim from MesheryAdapterPlayComponent's
 * generateSMIResult() — markup and option callbacks are identical.
 */
const AdapterSmiResultsDialog: React.FC<AdapterSmiResultsDialogProps> = ({
  open,
  onClose,
  adapterName,
  smiResult,
  page,
  pageSize,
  search,
  sortOrder,
  fetchSMIResults,
  onPageChange,
  onPageSizeChange,
  onSearchChange,
  onSortOrderChange,
}) => {
  const { data: capabilitiesData } = useGetProviderCapabilitiesQuery();
  const smi_columns = [
    {
      name: 'ID',
      label: 'ID',
      options: {
        filter: true,
        sort: true,
        searchable: true,
        customHeadRender: ({ index, ...column }) => {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender: (value) => (
          <Tooltip title={value} placement="top">
            <div>{value.slice(0, 5) + '...'}</div>
          </Tooltip>
        ),
      },
    },
    {
      name: 'Date',
      label: 'Date',
      options: {
        filter: true,
        sort: true,
        searchable: true,
        customHeadRender: ({ index, ...column }) => {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender: (value) => <Moment format="LLLL">{value}</Moment>,
      },
    },
    {
      name: 'Service Mesh',
      label: 'Service Mesh',
      options: {
        filter: true,
        sort: true,
        searchable: true,
        customHeadRender: ({ index, ...column }) => {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      },
    },
    {
      name: 'Service Mesh Version',
      label: 'Service Mesh Version',
      options: {
        filter: true,
        sort: true,
        searchable: true,
        customHeadRender: ({ index, ...column }) => {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      },
    },
    {
      name: '% Passed',
      label: '% Passed',
      options: {
        filter: true,
        sort: true,
        searchable: true,
        customHeadRender: ({ index, ...column }) => {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      },
    },
    {
      name: 'status',
      label: 'Status',
      options: {
        filter: true,
        sort: true,
        searchable: true,
        customHeadRender: ({ index, ...column }) => {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
      },
    },
  ];

  const smi_options = {
    sort: !isLocalProvider(capabilitiesData),
    search: !isLocalProvider(capabilitiesData),
    filterType: 'textField',
    expandableRows: true,
    selectableRows: 'none',
    rowsPerPage: pageSize,
    rowsPerPageOptions: [10, 20, 25],
    fixedHeader: true,
    print: false,
    download: false,
    renderExpandableRow: (rowData, rowMeta) => {
      const column = [
        'Specification',
        'Assertions',
        'Time',
        'Version',
        'Capability',
        'Result',
        'Reason',
      ];
      const results = (smiResult as SmiResult)?.results;
      const data =
        results?.[rowMeta.dataIndex]?.more_details?.map((val) => {
          return [
            val.smi_specification,
            val.assertions,
            val.time,
            val.smi_version,
            val.capability,
            val.status,
            val.reason,
          ];
        }) || [];
      const colSpan = rowData.length + 1;
      return (
        <TableRow>
          <TableCell colSpan={colSpan}>
            <SecondaryTable>
              <Table aria-label="a dense table">
                <TableHead>
                  <TableRow>
                    {column.map((val) => (
                      <TableCell colSpan={colSpan} key={val}>
                        {val}
                      </TableCell>
                    ))}
                  </TableRow>
                </TableHead>
                <TableBody>
                  {data.map((row, idx) => (
                    <TableRow key={idx}>
                      {row.map((val, vIdx) => {
                        if (val && val.match(/[0-9]+m[0-9]+.+[0-9]+s/i) != null) {
                          const time = val.split(/m|s/);
                          return (
                            <TableCell colSpan={colSpan} key={`${idx}-${vIdx}`}>
                              {time[0] + 'm ' + parseFloat(time[1]).toFixed(1) + 's'}
                            </TableCell>
                          );
                        } else {
                          return (
                            <TableCell colSpan={colSpan} key={`${idx}-${vIdx}`}>
                              {val}
                            </TableCell>
                          );
                        }
                      })}
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </SecondaryTable>
          </TableCell>
        </TableRow>
      );
    },
    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
      let order = '';
      if (tableState.activeColumn) {
        order = `${smi_columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case 'changePage':
          fetchSMIResults(adapterName, tableState.page, pageSize, search, sortOrder);
          onPageChange(tableState.page);
          break;
        case 'changeRowsPerPage':
          fetchSMIResults(adapterName, page, tableState.rowsPerPage, search, sortOrder);
          onPageSizeChange(tableState.rowsPerPage);
          break;
        case 'search':
          if (search !== tableState.searchText) {
            fetchSMIResults(
              adapterName,
              page,
              pageSize,
              tableState.searchText !== null ? tableState.searchText : '',
              sortOrder,
            );
            onSearchChange(tableState.searchText);
          }
          break;
        case 'sort':
          if (sortInfo.length === 2) {
            if (sortInfo[1] === 'ascending') {
              order = `${smi_columns[tableState.activeColumn].name} asc`;
            } else {
              order = `${smi_columns[tableState.activeColumn].name} desc`;
            }
          }
          if (order !== sortOrder) {
            fetchSMIResults(adapterName, page, pageSize, search, order);
            onSortOrderChange(order);
          }
          break;
      }
    },
  };

  let data: unknown[] = [];
  const sr = smiResult as SmiResult;
  if (sr && sr.results) {
    data = sr.results.map((val) => {
      return [
        val.id,
        val.date,
        val.mesh_name,
        val.mesh_version,
        val.passing_percentage,
        val.status,
      ];
    });
  }

  return (
    <Dialog
      onClose={onClose}
      aria-labelledby="adapter-dialog-title"
      open={open}
      fullWidth
      maxWidth="md"
    >
      <MUIDataTable
        title={<AdapterTableHeader>Service Mesh Interface Conformance Results</AdapterTableHeader>}
        data={data}
        columns={smi_columns}
        options={smi_options}
      />
    </Dialog>
  );
};

export default AdapterSmiResultsDialog;
