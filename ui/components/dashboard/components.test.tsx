import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { LayoutActionButton, LayoutWidget, StyledCard } from './components';

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
  IconButton: ({ children, onClick, 'aria-label': ariaLabel }: any) => (
    <button type="button" aria-label={ariaLabel} onClick={onClick}>
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
  WidgetPicker: ({ widgetsToAdd, onAddWidget, onClose }: any) => (
    <div data-testid="widget-picker">
      <span>Widgets</span>
      {onClose && (
        <button type="button" aria-label="Close widget picker" onClick={onClose}>
          Close
        </button>
      )}
      {widgetsToAdd.length === 0 && <span>All widgets added to the layout.</span>}
      {widgetsToAdd.map((widget: any) => (
        <div key={widget.key}>
          <span>{widget.title}</span>
          {widget.thumbnail && <img src={widget.thumbnail} alt={widget.title} />}
          <button
            type="button"
            aria-label={`Add ${widget.title} widget`}
            onClick={() => {
              const { key, ...rest } = widget;
              onAddWidget(rest, key);
            }}
          >
            Add
          </button>
        </div>
      ))}
    </div>
  ),
  DashboardLayout: ({ children, isSidebarOpen, sidebarContent }: any) => (
    <div data-testid="dashboard-layout" data-sidebar-open={String(isSidebarOpen)}>
      <div data-testid="dashboard-children">{children}</div>
      {isSidebarOpen && <div data-testid="dashboard-sidebar">{sidebarContent}</div>}
    </div>
  ),
  WidgetEmptyState: ({ message, icon, action }: any) => (
    <div data-testid="widget-empty-state" role="status">
      {icon}
      <span>{message}</span>
      {action && <button onClick={action.onClick}>{action.label}</button>}
    </div>
  ),
  styled: (Component: any) => () => (props: any) => {
    if (typeof Component === 'string') {
      const Tag = Component as any;
      return <Tag {...props} />;
    }
    return <Component {...props} />;
  },
  Paper: ({ children }: any) => <div>{children}</div>,
  Tab: ({ children }: any) => <div>{children}</div>,
  Tabs: ({ children }: any) => <div>{children}</div>,
  gray: {},
  charcoal: {},
}));

vi.mock('css/icons.styles', () => ({
  iconMedium: {},
}));

// NOTE: WidgetPicker is exported from @sistent/sistent. Because this file
// mocks the entire @sistent/sistent module, any test that imports and renders
// WidgetPicker from here only exercises the local vi.mock stub — not the real
// component — and cannot catch regressions in Sistent. Real WidgetPicker
// coverage belongs in layer5io/sistent's own test suite. If Meshery-specific
// integration behavior needs asserting (e.g. the onClose callback restores
// orgDashboardLayout), write an integration test against the Dashboard component
// directly rather than the WidgetPicker sub-component in isolation.

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
