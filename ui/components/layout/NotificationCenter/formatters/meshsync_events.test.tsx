import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('../../../connections/styles', () => ({
  ChipWrapper: ({ label, href, target, clickable, component, ...rest }: any) => (
    <a
      data-testid="chip-wrapper"
      href={href}
      target={target}
      data-clickable={String(Boolean(clickable))}
      data-component={component}
      {...rest}
    >
      {label}
    </a>
  ),
}));

vi.mock('../../../data-formatter', () => ({
  KeyValue: ({ Key, Value }: any) => (
    <div data-testid="key-value">
      <span data-testid="kv-key">{Key}</span>
      <span data-testid="kv-value">{Value}</span>
    </div>
  ),
}));

import {
  humanizeFieldName,
  ConnectionFieldFormatter,
  MeshSyncPropertyFormatters,
} from './meshsync_events';

describe('humanizeFieldName', () => {
  it('inserts spaces before capital letters and capitalises the result', () => {
    expect(humanizeFieldName('connectionID')).toBe('Connection I D');
    expect(humanizeFieldName('k8sContextName')).toBe('K8s Context Name');
    expect(humanizeFieldName('foo')).toBe('Foo');
  });

  it('handles already capitalised words', () => {
    expect(humanizeFieldName('Foo')).toBe('Foo');
  });
});

describe('ConnectionFieldFormatter', () => {
  it('returns null when value is missing', () => {
    const { container } = render(
      <ConnectionFieldFormatter value={undefined as any} fieldName="connectionID" />,
    );
    expect(container).toBeEmptyDOMElement();
  });

  it('renders a clickable chip with the search-text URL', () => {
    render(<ConnectionFieldFormatter value="my-conn id" fieldName="connectionID" />);
    const chip = screen.getByTestId('chip-wrapper');
    expect(chip).toHaveTextContent('my-conn id');
    expect(chip).toHaveAttribute(
      'href',
      `/management/connections?tab=connections&searchText=${encodeURIComponent('my-conn id')}`,
    );
    expect(chip).toHaveAttribute('data-clickable', 'true');
    expect(chip).toHaveAttribute('data-component', 'a');
    expect(chip).toHaveAttribute('target', '_self');
  });

  it('humanises the supplied fieldName', () => {
    render(<ConnectionFieldFormatter value="ctx-1" fieldName="k8sContextName" />);
    // ConnectionFieldFormatter wraps the chip inside a KeyValue, but in this
    // module we only mocked KeyValue indirectly through MeshSyncPropertyFormatters;
    // here we verify the chip renders with the value.
    expect(screen.getByText('ctx-1')).toBeInTheDocument();
    expect(screen.getByTestId('kv-key')).toHaveTextContent('K8s Context Name');
  });
});

describe('MeshSyncPropertyFormatters', () => {
  it('formats connectionID/k8sContextID/k8sContextName via ConnectionFieldFormatter', () => {
    const cases: Array<[keyof typeof MeshSyncPropertyFormatters, string, string]> = [
      ['connectionID', 'a', 'Connection I D'],
      ['k8sContextID', 'b', 'K8s Context I D'],
      ['k8sContextName', 'c', 'K8s Context Name'],
    ];

    for (const [key, value, expectedKey] of cases) {
      const { container } = render(MeshSyncPropertyFormatters[key](value));
      expect(container.querySelector('[data-testid="kv-key"]')).toHaveTextContent(expectedKey);
      expect(container.querySelector('[data-testid="chip-wrapper"]')).toHaveTextContent(value);
    }
  });

  it('formats deployment, operator and broker fields with humanised keys', () => {
    const cases: Array<[keyof typeof MeshSyncPropertyFormatters, string, string]> = [
      ['meshsyncDeploymentMode', 'mode-x', 'Meshsync Deployment Mode'],
      ['operatorStatus', 'running', 'Operator Status'],
      ['brokerEndpoint', 'nats://1.2.3.4:4222', 'Broker Endpoint'],
    ];

    for (const [key, value, expectedKey] of cases) {
      const { container } = render(MeshSyncPropertyFormatters[key](value));
      expect(container.querySelector('[data-testid="kv-key"]')).toHaveTextContent(expectedKey);
      expect(container.querySelector('[data-testid="kv-value"]')).toHaveTextContent(value);
    }
  });
});
