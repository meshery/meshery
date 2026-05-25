import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/assets/icons', () => {
  const makeIcon = (name: string) =>
    function Mock({ style }: any) {
      return <span data-testid={name} data-color={style?.color || ''} />;
    };
  return {
    AddCircle: makeIcon('AddCircle'),
    BuildRounded: makeIcon('BuildRounded'),
    DirectionsCar: makeIcon('DirectionsCar'),
    Filter: makeIcon('Filter'),
    SimCard: makeIcon('SimCard'),
    SupervisedUserCircle: makeIcon('SupervisedUserCircle'),
    TouchApp: makeIcon('TouchApp'),
    Explore: makeIcon('Explore'),
    FileCopy: makeIcon('FileCopy'),
    ListAlt: makeIcon('ListAlt'),
    Lock: makeIcon('Lock'),
  };
});

import NameToIcon from './NameToIcon';

const cases: Array<[string, string]> = [
  ['Application', 'TouchApp'],
  ['Kubernetes Service', 'BuildRounded'],
  ['AuthorizationPolicy', 'Lock'],
  ['DestinationRule', 'Explore'],
  ['EnvoyFilter', 'Filter'],
  ['Gateway', 'ListAlt'],
  ['PeerAuthentication', 'FileCopy'],
  ['Sidecar', 'DirectionsCar'],
  ['VirtualService', 'SupervisedUserCircle'],
  ['WorkloadEntry', 'SimCard'],
];

describe('NameToIcon', () => {
  it.each(cases)('returns the right icon for %s', (name, expectedTestId) => {
    render(<NameToIcon name={name} />);
    expect(screen.getByTestId(expectedTestId)).toBeInTheDocument();
  });

  it('falls back to AddCircle for unknown names', () => {
    render(<NameToIcon name="Unknown" />);
    expect(screen.getByTestId('AddCircle')).toBeInTheDocument();
  });

  it('applies the default color to the icon', () => {
    render(<NameToIcon name="Application" />);
    expect(screen.getByTestId('TouchApp')).toHaveAttribute('data-color', '#607D8B');
  });

  it('passes the custom color through to the icon style', () => {
    render(<NameToIcon name="Application" color="#ff00ff" />);
    expect(screen.getByTestId('TouchApp')).toHaveAttribute('data-color', '#ff00ff');
  });
});
