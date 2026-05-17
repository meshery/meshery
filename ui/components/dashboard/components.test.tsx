import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import {
  AddWidgetsToLayoutPanel,
  LayoutActionButton,
  LayoutWidget,
  StyledCard,
} from './components';

vi.mock('@sistent/sistent', () => ({
  Box: ({ children }: any) => <div data-testid="box">{children}</div>,
  useTheme: () => ({
    palette: {
      background: {
        card: '#fff',
        neutral: { default: '#000' },
        default: '#eee',
        elevatedComponents: '#ddd',
      },
    },
  }),
  Typography: ({ children, variant }: any) => <span data-variant={variant}>{children}</span>,
  Stack: ({ children }: any) => <div>{children}</div>,
  AddIcon: (props: any) => <svg data-testid="add-icon" data-fill={props.fill} />,
  IconButton: ({ children, onClick }: any) => (
    <button type="button" onClick={onClick}>
      {children}
    </button>
  ),
  Card: ({ children }: any) => <div data-testid="card">{children}</div>,
  CardContent: ({ children }: any) => <div data-testid="card-content">{children}</div>,
  CustomTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
  Button: ({ children, onClick, endIcon }: any) => (
    <button type="button" onClick={onClick}>
      {children}
      {endIcon}
    </button>
  ),
  DeleteIcon: (props: any) => <svg data-testid="delete-icon" data-fill={props.fill} />,
  DragIcon: (props: any) => <svg data-testid="drag-icon" data-fill={props.fill} />,
}));

vi.mock('css/icons.styles', () => ({
  iconMedium: {},
}));

describe('AddWidgetsToLayoutPanel', () => {
  it('returns null when not in edit mode', () => {
    const { container } = render(
      <AddWidgetsToLayoutPanel widgetsToAdd={[]} editMode={false} onAddWidget={vi.fn()} />,
    );
    expect(container.firstChild).toBeNull();
  });

  it('shows the empty state when in edit mode but with no widgets to add', () => {
    render(<AddWidgetsToLayoutPanel widgetsToAdd={[]} editMode={true} onAddWidget={vi.fn()} />);
    expect(screen.getByText(/All widgets added to the layout/i)).toBeInTheDocument();
  });

  it('renders widget cards and invokes onAddWidget on click', async () => {
    const user = userEvent.setup();
    const onAddWidget = vi.fn();
    render(
      <AddWidgetsToLayoutPanel
        widgetsToAdd={[{ key: 'OVERVIEW', title: 'Overview', thumbnail: '/a.png' }]}
        editMode={true}
        onAddWidget={onAddWidget}
      />,
    );

    expect(screen.getByText('Overview')).toBeInTheDocument();
    const img = screen.getByRole('img', { name: /overview/i });
    expect(img).toHaveAttribute('src', '/a.png');

    await user.click(screen.getByRole('button'));
    expect(onAddWidget).toHaveBeenCalledTimes(1);
    expect(onAddWidget).toHaveBeenCalledWith(
      { title: 'Overview', thumbnail: '/a.png' },
      'OVERVIEW',
    );
  });
});

describe('LayoutActionButton', () => {
  const FakeIcon = (props: any) => <svg data-testid="fake-icon" {...props} />;

  it('renders nothing when isShown is false', () => {
    const { container } = render(
      <LayoutActionButton
        Icon={FakeIcon}
        label="Edit"
        action={vi.fn()}
        description="Edit layout"
        isShown={false}
      />,
    );
    expect(container.firstChild).toBeNull();
  });

  it('renders label, tooltip, and triggers action on click', async () => {
    const user = userEvent.setup();
    const action = vi.fn();
    render(
      <LayoutActionButton
        Icon={FakeIcon}
        label="Edit"
        action={action}
        description="Edit layout"
        isShown={true}
      />,
    );
    expect(screen.getByText('Edit')).toBeInTheDocument();
    expect(screen.getByTestId('tooltip')).toHaveAttribute('data-title', 'Edit layout');
    await user.click(screen.getByRole('button'));
    expect(action).toHaveBeenCalledTimes(1);
  });
});

describe('LayoutWidget', () => {
  const widget = { key: 'KEY', title: 'My Widget', component: <div>widget-body</div> };

  it('renders only the widget component when not in edit mode', () => {
    render(<LayoutWidget widget={widget} removeWidget={vi.fn()} isEditMode={false} />);
    expect(screen.getByText('widget-body')).toBeInTheDocument();
    expect(screen.queryByTestId('delete-icon')).not.toBeInTheDocument();
  });

  it('renders delete + drag controls in edit mode and removes by key', async () => {
    const user = userEvent.setup();
    const removeWidget = vi.fn();
    render(<LayoutWidget widget={widget} removeWidget={removeWidget} isEditMode={true} />);

    expect(screen.getByTestId('delete-icon')).toBeInTheDocument();
    expect(screen.getByTestId('drag-icon')).toBeInTheDocument();

    const buttons = screen.getAllByRole('button');
    // First button is delete
    await user.click(buttons[0]);
    expect(removeWidget).toHaveBeenCalledWith('KEY');
  });
});

describe('StyledCard', () => {
  it('renders title, icon, button, and children', () => {
    render(
      <StyledCard
        title="Widget Title"
        icon={<span>icon-x</span>}
        button={<button type="button">Action</button>}
      >
        <div>card-body</div>
      </StyledCard>,
    );

    expect(screen.getByText('Widget Title')).toBeInTheDocument();
    expect(screen.getByText('icon-x')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Action' })).toBeInTheDocument();
    expect(screen.getByText('card-body')).toBeInTheDocument();
  });
});
