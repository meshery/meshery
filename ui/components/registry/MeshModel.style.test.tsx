import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

// Make sure the styled() factory returns a passthrough renderer so we can
// verify the underlying tags render.
vi.mock('@sistent/sistent', async () => {
  const styled = (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) => {
      if (typeof Component === 'string') {
        return React.createElement(Component, props, children);
      }
      return React.createElement(Component, props, children);
    };
    StyledComponent.displayName = 'StyledSistentMock';
    return StyledComponent;
  };

  return {
    styled,
    Button: (props: any) => <button {...props}>{props.children}</button>,
    alpha: (color: string) => color,
  };
});

vi.mock('../shared/TreeView', () => ({
  TreeItem: ({ children, ...rest }: any) => (
    <div data-testid="tree-item" {...rest}>
      {children}
    </div>
  ),
  treeItemClasses: {
    content: 'content-class',
    expanded: 'expanded-class',
    selected: 'selected-class',
    focused: 'focused-class',
    label: 'label-class',
    group: 'group-class',
  },
}));

import {
  DisableButton,
  JustifyAndAlignCenter,
  StyledTreeItemRoot,
  StyledKeyValuePropertyDiv,
  StyledKeyValueProperty,
  StyledKeyValueFormattedValue,
  StyledTreeItemNameDiv,
  StyledTreeItemDiv,
  MesheryTreeViewWrapper,
} from './MeshModel.style';

describe('MeshModel.style', () => {
  it('exports all styled components', () => {
    expect(DisableButton).toBeDefined();
    expect(JustifyAndAlignCenter).toBeDefined();
    expect(StyledTreeItemRoot).toBeDefined();
    expect(StyledKeyValuePropertyDiv).toBeDefined();
    expect(StyledKeyValueProperty).toBeDefined();
    expect(StyledKeyValueFormattedValue).toBeDefined();
    expect(StyledTreeItemNameDiv).toBeDefined();
    expect(StyledTreeItemDiv).toBeDefined();
    expect(MesheryTreeViewWrapper).toBeDefined();
  });

  it('renders simple wrapper components without crashing', () => {
    const { container } = render(
      <MesheryTreeViewWrapper>
        <StyledKeyValuePropertyDiv>
          <StyledKeyValueProperty>label</StyledKeyValueProperty>
          <StyledKeyValueFormattedValue>value</StyledKeyValueFormattedValue>
        </StyledKeyValuePropertyDiv>
        <StyledTreeItemDiv>
          <StyledTreeItemNameDiv>name</StyledTreeItemNameDiv>
        </StyledTreeItemDiv>
        <JustifyAndAlignCenter>centered</JustifyAndAlignCenter>
      </MesheryTreeViewWrapper>,
    );

    expect(container.textContent).toContain('label');
    expect(container.textContent).toContain('value');
    expect(container.textContent).toContain('name');
    expect(container.textContent).toContain('centered');
  });

  it('renders DisableButton', () => {
    const { container } = render(<DisableButton>hello</DisableButton>);
    expect(container.textContent).toContain('hello');
  });
});
