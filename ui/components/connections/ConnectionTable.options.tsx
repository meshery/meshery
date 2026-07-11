import React, { useMemo } from 'react';
import { Button, Grid2, Table, TableCell, TableRow, DeleteIcon, useTheme } from '@sistent/sistent';
import { ContentContainer, InnerTableContainer } from './styles';
import { iconMedium } from '../../css/icons.styles';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import FormatConnectionMetadata from './metadata';
import type { ConnectionRow, ExpansionFlags, SelectedRows } from './ConnectionTable.types';

type UseConnectionTableOptionsArgs = {
  totalCount?: number;
  page: number;
  pageSize: number;
  setPage: (page: number) => void;
  setPageSize: (size: number) => void;
  sortOrder: string;
  setSortOrder: (order: string) => void;
  rowsExpanded: number[];
  setRowsExpanded: (rows: number[]) => void;
  columns: Array<{ name: string }>;
  filteredConnections: ConnectionRow[];
  meshsyncControllerState: unknown;
  selectedConnectionId?: string;
  updateUrlWithConnectionId?: (connectionId: string) => void;
  expansionFlags: React.MutableRefObject<ExpansionFlags>;
  handleDeleteConnections: (selected: SelectedRows) => void | Promise<void>;
};

export const useConnectionTableOptions = ({
  totalCount,
  page,
  pageSize,
  setPage,
  setPageSize,
  sortOrder,
  setSortOrder,
  rowsExpanded,
  setRowsExpanded,
  columns,
  filteredConnections,
  meshsyncControllerState,
  selectedConnectionId,
  updateUrlWithConnectionId,
  expansionFlags,
  handleDeleteConnections,
}: UseConnectionTableOptionsArgs) => {
  const theme = useTheme();

  return useMemo(
    () => ({
      filter: false,
      viewColumns: false,
      search: false,
      responsive: 'standard',
      resizableColumns: true,
      serverSide: true,
      count: totalCount,
      rowsPerPage: pageSize,
      fixedHeader: true,
      page,
      print: false,
      download: false,
      textLabels: {
        selectedRows: {
          text: 'connection(s) selected',
        },
      },
      sortOrder: {
        name: sortOrder.split(' ')[0],
        direction: sortOrder.split(' ')[1],
      },
      customToolbarSelect: (selected) => (
        <Button
          color="error"
          variant="contained"
          size="large"
          onClick={() => handleDeleteConnections(selected)}
          sx={{ backgroundColor: `${theme.palette.error.dark} !important`, marginRight: '10px' }}
          disabled={
            !CAN(
              Keys.LifecycleManagementDeleteAConnection.id,
              Keys.LifecycleManagementDeleteAConnection.function,
            )
          }
          data-testid="Button-delete-connections"
        >
          <DeleteIcon style={iconMedium} fill={theme.palette.common.white} />
          Delete
        </Button>
      ),
      enableNestedDataAccess: '.',
      onTableChange: (action, tableState) => {
        const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
        let order = '';

        if (typeof tableState.activeColumn === 'number') {
          order = `${columns[tableState.activeColumn].name} desc`;
        }

        switch (action) {
          case 'changePage':
            setPage(tableState.page);
            break;
          case 'changeRowsPerPage':
            setPageSize(tableState.rowsPerPage);
            break;
          case 'sort':
            if (sortInfo.length === 2 && typeof tableState.activeColumn === 'number') {
              order =
                sortInfo[1] === 'ascending'
                  ? `${columns[tableState.activeColumn].name} asc`
                  : `${columns[tableState.activeColumn].name} desc`;
            }

            if (order && order !== sortOrder) {
              setSortOrder(order);
            }
            break;
        }
      },
      expandableRows: true,
      expandableRowsHeader: false,
      expandableRowsOnClick: true,
      rowsExpanded,
      isRowExpandable: () => true,
      onRowExpansionChange: (_, allRowsExpanded) => {
        if (expansionFlags.current.isUrlExpansion) {
          return;
        }

        expansionFlags.current.isHandlingExpansion = true;
        const expandedRows = allRowsExpanded.slice(-1);
        // mui-datatables fires `onRowExpansionChange` for every internal
        // re-render where the expanded set is unchanged. Bailing out on a
        // value-equal write avoids enqueuing redundant setState commits which
        // can cascade into the React #185 update-depth limit on slow renders.
        const nextExpandedRows = expandedRows.map((item) => item.index);
        const hasExpandedRowsChanged =
          nextExpandedRows.length !== rowsExpanded.length ||
          nextExpandedRows.some((rowIndex, index) => rowIndex !== rowsExpanded[index]);

        if (hasExpandedRowsChanged) {
          setRowsExpanded(nextExpandedRows);
        }

        if (expandedRows.length > 0) {
          const index = expandedRows[0].index;
          const connection = filteredConnections[index];

          if (
            connection &&
            updateUrlWithConnectionId &&
            (!expansionFlags.current.isInitialLoad || connection.id !== selectedConnectionId)
          ) {
            updateUrlWithConnectionId(connection.id);
          }
        } else if (
          updateUrlWithConnectionId &&
          !expansionFlags.current.isInitialLoad &&
          selectedConnectionId
        ) {
          updateUrlWithConnectionId('');
        }

        expansionFlags.current.isHandlingExpansion = false;
      },
      renderExpandableRow: (rowData, tableMeta) => {
        const connection = filteredConnections[tableMeta.rowIndex];
        return (
          <TableCell colSpan={rowData.length}>
            <InnerTableContainer>
              <Table>
                <TableRow style={{ padding: 0 }}>
                  <TableCell style={{ overflowX: 'hidden', padding: 0 }}>
                    <Grid2 container style={{ textTransform: 'lowercase' }} size="grow">
                      <ContentContainer size={{ xs: 12 }}>
                        <FormatConnectionMetadata
                          connection={connection}
                          meshsyncControllerState={meshsyncControllerState}
                        />
                      </ContentContainer>
                    </Grid2>
                  </TableCell>
                </TableRow>
              </Table>
            </InnerTableContainer>
          </TableCell>
        );
      },
    }),
    [
      columns,
      expansionFlags,
      filteredConnections,
      handleDeleteConnections,
      meshsyncControllerState,
      page,
      pageSize,
      rowsExpanded,
      selectedConnectionId,
      setPage,
      setPageSize,
      setRowsExpanded,
      setSortOrder,
      sortOrder,
      theme.palette.common.white,
      theme.palette.error.dark,
      totalCount,
      updateUrlWithConnectionId,
    ],
  );
};
