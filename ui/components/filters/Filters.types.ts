import type React from 'react';
import type { Key } from '@meshery/schemas/permissions';

export type TypeView = 'grid' | 'table';

export type TooltipIconProps = {
  children: React.ReactNode;
  onClick: (_event: React.MouseEvent) => void;
  title: string;
  disabled?: boolean;
  placement?: string;
  permissionKey?: Key;
  permissionAction?: 'showShield' | 'hide';
};

export type FilterSubmitPayload = {
  data: string;
  id: string;
  name: string;
  type: string;
  catalogData?: any;
};

export type YAMLEditorProps = {
  filter: any;
  onClose: () => void;
  onSubmit: (_payload: FilterSubmitPayload) => void;
};

export type ImportModalProps = {
  handleClose: () => void;
  handleImportFilter: (_data: any) => void;
};

export type PublishModalProps = {
  handleClose: () => void;
  handleSubmit: (_data: any) => void;
  title: string;
};
