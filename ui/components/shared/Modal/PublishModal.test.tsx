import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import PublishModal from './PublishModal';

let modalProps: any = null;
let mockSchemaData: any = { rjsfSchema: { type: 'object' }, uiSchema: {} };
let mockIsSuccess = true;
const mockGetMeshModels = vi.fn();
const mockModifyRJSFSchema = vi.fn((schema: any, _path: string, enumValues: string[]) => ({
  ...schema,
  modified: enumValues,
}));

vi.mock('./Modal', () => ({
  default: (props: any) => {
    modalProps = props;
    return (
      <div
        data-testid="rjsf-modal"
        data-title={props.title}
        data-submit-btn={props.submitBtnText}
        data-open={String(props.open)}
      />
    );
  },
}));

vi.mock('@sistent/sistent', () => ({
  PublicIcon: (props: any) => <svg data-testid="public-icon" {...props} />,
}));

vi.mock('../../../api/meshmodel', () => ({
  getMeshModels: () => mockGetMeshModels(),
}));

vi.mock('../../../utils/utils', () => ({
  modifyRJSFSchema: (...args: any[]) => mockModifyRJSFSchema(...args),
}));

vi.mock('@/rtk-query/schema', () => ({
  useGetSchemaQuery: () => ({ data: mockSchemaData, isSuccess: mockIsSuccess }),
}));

describe('PublishModal', () => {
  beforeEach(() => {
    modalProps = null;
    mockGetMeshModels.mockReset();
    mockGetMeshModels.mockResolvedValue({
      models: [{ displayName: 'kubernetes' }, { displayName: 'istio' }],
    });
    mockSchemaData = { rjsfSchema: { type: 'object' }, uiSchema: {} };
    mockIsSuccess = true;
  });

  it('renders an underlying RJSFModal with the publish submit text', () => {
    render(
      <PublishModal
        open={true}
        title="Publish Design"
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    expect(screen.getByTestId('rjsf-modal')).toHaveAttribute('data-title', 'Publish Design');
    expect(screen.getByTestId('rjsf-modal')).toHaveAttribute(
      'data-submit-btn',
      'Submit for Approval',
    );
  });

  it('processes the schema when schema data succeeds', async () => {
    render(
      <PublishModal
        open={true}
        title="Publish Design"
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    // Allow microtask
    await new Promise((resolve) => setTimeout(resolve, 10));

    expect(mockGetMeshModels).toHaveBeenCalled();
    expect(mockModifyRJSFSchema).toHaveBeenCalledWith(
      mockSchemaData.rjsfSchema,
      'properties.compatibility.items.enum',
      ['istio', 'kubernetes'],
    );
  });

  it('falls back to the original schema when fetching meshmodels fails', async () => {
    mockGetMeshModels.mockRejectedValue(new Error('fail'));
    const errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    render(
      <PublishModal
        open={true}
        title="Publish Design"
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    await new Promise((resolve) => setTimeout(resolve, 10));
    expect(errorSpy).toHaveBeenCalled();
    errorSpy.mockRestore();
  });

  it('does not process the schema if isSuccess is false', async () => {
    mockIsSuccess = false;
    mockSchemaData = null;

    render(
      <PublishModal
        open={false}
        title="Publish Design"
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    await new Promise((resolve) => setTimeout(resolve, 10));
    expect(mockGetMeshModels).not.toHaveBeenCalled();
  });

  it('passes open prop through to the Modal', () => {
    render(
      <PublishModal
        open={false}
        title="Publish Design"
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    expect(screen.getByTestId('rjsf-modal')).toHaveAttribute('data-open', 'false');
  });
});
