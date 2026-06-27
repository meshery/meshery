import React from 'react';
import { CustomColumnVisibilityControl, SearchBar, UniversalFilter } from '@sistent/sistent';
import { CreateButton } from './styles';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import { styled } from '@/theme';
import ConnectionWizardLauncher from './ConnectionWizardLauncher';
import { getVisibilityColums } from '../../utils/utils';
import type { SelectedFilters } from './ConnectionTable.types';

type ConnectionTableToolbarProps = {
  isSearchExpanded: boolean;
  setIsSearchExpanded: (expanded: boolean) => void;
  onSearch: (value: string) => void;
  filters: Record<string, { name: string; options: Array<{ label: string; value: string }> }>;
  selectedFilters: SelectedFilters;
  setSelectedFilters: (filters: SelectedFilters) => void;
  handleApplyFilter: () => void;
  columns: Array<{ name: string; label?: string; options?: { display?: boolean } }>;
  columnVisibility: Record<string, boolean | undefined>;
  setColumnVisibility: (visibility: Record<string, boolean | undefined>) => void;
};

const ToolbarActions = styled('div')(() => ({
  display: 'flex',
  borderRadius: '0.5rem 0.5rem 0 0',
  width: '100%',
  justifyContent: 'flex-end',
}));

export const ConnectionTableToolbar = ({
  isSearchExpanded,
  setIsSearchExpanded,
  onSearch,
  filters,
  selectedFilters,
  setSelectedFilters,
  handleApplyFilter,
  columns,
  columnVisibility,
  setColumnVisibility,
}: ConnectionTableToolbarProps) => {
  return (
    <ToolWrapper style={{ marginBottom: '5px', marginTop: '-30px' }}>
      <CreateButton>
        <ConnectionWizardLauncher />
      </CreateButton>
      <ToolbarActions>
        <div data-testid="ConnectionTable-search">
          <SearchBar
            onSearch={onSearch}
            placeholder="Search Connections..."
            expanded={isSearchExpanded}
            setExpanded={setIsSearchExpanded}
          />
        </div>

        <UniversalFilter
          id="ref"
          filters={filters}
          selectedFilters={selectedFilters}
          setSelectedFilters={setSelectedFilters}
          handleApplyFilter={handleApplyFilter}
        />

        <CustomColumnVisibilityControl
          style={{ zIndex: 1300 }}
          id="ref"
          columns={getVisibilityColums(columns)}
          customToolsProps={{ columnVisibility, setColumnVisibility }}
        />
      </ToolbarActions>
    </ToolWrapper>
  );
};
