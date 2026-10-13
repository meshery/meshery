import React from 'react';
import { render, screen } from '@testing-library/react';
import MesheryTreeView from '../ui/components/registry/MesheryTreeView';

// Mock components
jest.mock('../ui/components/Spinner', () => () => <div>Spinner</div>);
jest.mock('../ui/components/EmptyStateMessage', () => () => <div>EmptyStateMessage</div>);
jest.mock('../ui/components/DataDisplay', () => () => <div>DataDisplay</div>);

describe('MesheryTreeView Component', () => {
  test('renders Spinner when isLoading is true', () => {
    const { getByText } = render(<MesheryTreeView isLoading={true} data={[]} searchText="" />);
    expect(getByText('Spinner')).toBeInTheDocument();
  });

  test('renders EmptyStateMessage when data is empty and searchText is empty', () => {
    const { getByText } = render(<MesheryTreeView isLoading={false} data={[]} searchText="" />);
    expect(getByText('EmptyStateMessage')).toBeInTheDocument();
  });

  test('renders DataDisplay when data is not empty', () => {
    const { getByText } = render(<MesheryTreeView isLoading={false} data={[{ id: 1 }]} searchText="" />);
    expect(getByText('DataDisplay')).toBeInTheDocument();
  });

  test('renders DataDisplay when searchText is not empty', () => {
    const { getByText } = render(<MesheryTreeView isLoading={false} data={[]} searchText="search" />);
    expect(getByText('DataDisplay')).toBeInTheDocument();
  });

  test('renders Spinner when isLoading is true regardless of data and searchText', () => {
    const { getByText } = render(<MesheryTreeView isLoading={true} data={[{ id: 1 }]} searchText="search" />);
    expect(getByText('Spinner')).toBeInTheDocument();
  });
});