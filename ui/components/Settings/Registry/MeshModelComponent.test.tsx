import React from 'react';
import { render, screen } from '@testing-library/react';
import { TabBar } from './MeshModelComponent';
import CAN from '@/utils/can';
import { MODELS } from '../../../constants/navigator';
import { describe, it, expect, vi } from 'vitest';

// Mock the CAN utility
vi.mock('@/utils/can', () => ({
  default: vi.fn(),
}));

// Mock components that might import SVGs or heavy dependencies
vi.mock('./MesheryTreeView', () => ({ default: () => <div /> }));
vi.mock('./MeshModelDetails', () => ({ default: () => <div /> }));
vi.mock('./ImportModelModal', () => ({ default: () => <div /> }));
vi.mock('./CreateModelModal', () => ({ default: () => <div /> }));
vi.mock('@/components/RelationshipBuilder/CreateRelationshipModal', () => ({
  default: () => <div />,
}));

describe('TabBar Component RBAC', () => {
  it('disables Create Model and Import Model buttons when user lacks permission', () => {
    // Mock CAN to return false (simulating missing RBAC capability)
    (CAN as any).mockReturnValue(false);

    render(
      <TabBar
        openImportModal={() => {}}
        openCreateModal={() => {}}
        view={MODELS}
        openRelationshipModal={() => {}}
      />,
    );

    const createModelBtn = screen.getByTestId('TabBar-Button-CreateModel');
    const importModelBtn = screen.getByTestId('TabBar-Button-ImportModel');

    expect(createModelBtn).toBeDisabled();
    expect(importModelBtn).toBeDisabled();
  });

  it('enables Create Model and Import Model buttons when user has permission', () => {
    // Mock CAN to return true (simulating present RBAC capability)
    (CAN as any).mockReturnValue(true);

    render(
      <TabBar
        openImportModal={() => {}}
        openCreateModal={() => {}}
        view={MODELS}
        openRelationshipModal={() => {}}
      />,
    );

    const createModelBtn = screen.getByTestId('TabBar-Button-CreateModel');
    const importModelBtn = screen.getByTestId('TabBar-Button-ImportModel');

    expect(createModelBtn).not.toBeDisabled();
    expect(importModelBtn).not.toBeDisabled();
  });
});
