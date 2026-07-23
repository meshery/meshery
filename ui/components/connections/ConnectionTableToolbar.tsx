import React from 'react';
import {
  CustomColumnVisibilityControl,
  SearchBar,
  UniversalFilter,
  DataTableToolbar,
} from '@sistent/sistent';
import { CreateButton } from './styles';
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
  tabs?: React.ReactNode;
};

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
  tabs,
}: ConnectionTableToolbarProps) => {
  return (
    <DataTableToolbar
      primaryActions={
        <CreateButton>
          <ConnectionWizardLauncher />
        </CreateButton>
      }
      search={
        <div data-testid="ConnectionTable-search">
          <SearchBar
            onSearch={onSearch}
            placeholder="Search Connections..."
            expanded={isSearchExpanded}
            setExpanded={setIsSearchExpanded}
          />
        </div>
      }
      filter={
        <UniversalFilter
          id="connection-table-filter"
          filters={filters}
          selectedFilters={selectedFilters}
          setSelectedFilters={setSelectedFilters}
          handleApplyFilter={handleApplyFilter}
        />
      }
      columnVisibility={
        <CustomColumnVisibilityControl
          style={{ zIndex: 1300 }}
          id="connection-table-column-visibility"
          columns={getVisibilityColums(columns)}
          customToolsProps={{
            columnVisibility,
            setColumnVisibility,
          }}
        />
      }
      tabs={tabs}
    />
  );
};
