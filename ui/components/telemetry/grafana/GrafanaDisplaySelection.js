import React from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@sistent/sistent';
import { Chip, Box, styled } from '@sistent/sistent';
import MUIDataTable from 'mui-datatables';

const Root = styled(Box)(({ theme }) => ({
  padding: theme.spacing(5),
}));

const PanelChips = styled(Box)(() => ({
  display: 'flex',
  flexWrap: 'wrap',
}));

const StyledChip = styled(Chip)(({ theme }) => ({
  margin: theme.spacing(0.25),
  maxWidth: '200px',
  overflow: 'visible',
  position: 'relative',
  '& .MuiChip-label': {
    overflow: 'hidden',
    whiteSpace: 'nowrap',
    textOverflow: 'ellipsis',
    transition: 'all 0.3s ease',
  },
  '&:hover': {
    zIndex: 1000,
    '& .MuiChip-label': {
      overflow: 'visible',
      whiteSpace: 'normal',
      wordWrap: 'break-word',
      position: 'absolute',
      left: 0,
      top: 0,
      width: 'max-content',
      maxWidth: 'min(300px, 80vw)',
      backgroundColor: theme.palette.background.paper || '#fff',
      border: `1px solid ${theme.palette.divider || '#e0e0e0'}`,
      borderRadius: '4px',
      padding: '8px 12px',
      boxShadow: theme.shadows[4] || '0px 2px 8px rgba(0,0,0,0.15)',
      zIndex: 1001,
      // Ensure the tooltip stays within viewport bounds
      transform: 'translateX(max(-50%, min(0px, calc(100vw - 100% - 20px))))',
      // Add word breaking for very long words
      wordBreak: 'break-word',
      overflowWrap: 'break-word',
      hyphens: 'auto',
    },
  },
}));

const GrafanaDisplaySelection = ({ boardPanelConfigs, deleteSelectedBoardPanelConfig }) => {
  const selectedValsForDisplay = boardPanelConfigs.map((cf) => ({
    board: cf.board?.title || '',
    panels: (
      <PanelChips>
        {cf.panels.map((panel, ind) => (
          <StyledChip key={`${panel.id}_-_${ind}`} label={panel.title} />
        ))}
      </PanelChips>
    ),
    template_variables: (
      <PanelChips>
        {cf.templateVars
          ? cf.templateVars.map((tv, ind) =>
              tv && tv !== '' ? <StyledChip key={`${tv}-_-${ind}`} label={tv} /> : null,
            )
          : []}
      </PanelChips>
    ),
  }));

  const columns = [
    { name: 'board', label: 'Board' },
    { name: 'panels', label: 'Panels' },
    { name: 'template_variables', label: 'Template Variables' },
  ];
  const options = {
    filter: false,
    sort: false,
    search: false,
    filterType: 'textField',
    responsive: 'vertical',
    count: selectedValsForDisplay.length,
    print: false,
    download: false,
    pagination: false,
    viewColumns: false,
    onRowsDelete: (rowsDeleted) => {
      const delRows = rowsDeleted.data.map(({ dataIndex }) => dataIndex);
      deleteSelectedBoardPanelConfig(delRows);
      return false;
    },
  };
  return (
    <NoSsr>
      <Root>
        <MUIDataTable
          key={`gds_${new Date().getTime()}`}
          title="Meshery Results"
          data={selectedValsForDisplay}
          columns={columns}
          options={options}
        />
      </Root>
    </NoSsr>
  );
};

GrafanaDisplaySelection.propTypes = {
  classes: PropTypes.object.isRequired,
  boardPanelConfigs: PropTypes.array.isRequired,
  deleteSelectedBoardPanelConfig: PropTypes.func.isRequired,
};

export default GrafanaDisplaySelection;
