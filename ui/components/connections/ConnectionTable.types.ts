import type { ReactNode } from 'react';

export type ConnectionTableProps = {
  selectedFilter?: string;
  selectedConnectionId?: string;
  updateUrlWithConnectionId?: (connectionId: string) => void;
  tabs?: ReactNode;
};

export type EnvironmentOption = {
  label: string;
  value: string;
  __isNew__?: boolean;
};

export type ConnectionRow = {
  id: string;
  kind: string;
  status: string;
  name?: string;
  kindLogo?: string;
  nextStatus?: string[];
  metadata?: Record<string, any>;
  environments?: Array<{ id: string; name: string }>;
  [key: string]: any;
};

export type SelectedFilters = {
  status: string;
  kind: string;
};

export type SelectedRows = {
  data?: Array<{ index: number }>;
};

export type RowData = {
  rowIndex: number;
};

export type ExpansionFlags = {
  isHandlingExpansion: boolean;
  isInitialLoad: boolean;
  isUrlExpansion: boolean;
  lastProcessedId: string | null;
};
