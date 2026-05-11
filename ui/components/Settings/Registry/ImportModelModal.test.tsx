import React from 'react';
import { render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import ImportModelModal from './ImportModelModal';

const importModelReq = vi.fn();
const updateProgress = vi.fn();
let capturedRjsfProps: any;

vi.mock('@sistent/sistent', () => ({
  FormControlLabel: ({ children }) => <div>{children}</div>,
  Button: ({ children }) => <button type="button">{children}</button>,
  FormControl: ({ children }) => <div>{children}</div>,
  RadioGroup: ({ children }) => <div>{children}</div>,
  Radio: () => <div />,
  Typography: ({ children }) => <div>{children}</div>,
  ModalFooter: ({ children }) => <div>{children}</div>,
  Box: ({ children }) => <div>{children}</div>,
  Modal: ({ children, open }) => (open ? <div>{children}</div> : null),
  styled: (Component) => () => {
    const StyledComponent = ({ children, ...props }) => (
      <Component {...props}>{children}</Component>
    );
    StyledComponent.displayName = 'StyledSistentMock';
    return StyledComponent;
  },
  useTheme: () => ({
    palette: {
      background: {
        surfaces: 'white',
      },
    },
  }),
  importModelSchema: {
    title: 'Import model',
    properties: {
      uploadType: {
        enum: ['file', 'urlImport', 'csv'],
        enumNames: ['File Import', 'URL Import', 'CSV Import'],
      },
    },
    allOf: [
      {
        if: { properties: { uploadType: { const: 'file' } } },
        then: {
          properties: {
            modelFile: { type: 'string' },
            fileName: { type: 'string' },
          },
          required: ['modelFile', 'fileName'],
        },
      },
      {
        if: { properties: { uploadType: { const: 'urlImport' } } },
        then: {
          properties: {
            url: { type: 'string' },
          },
          required: ['url'],
        },
      },
    ],
  },
  importModelUiSchema: {
    uploadType: { 'ui:widget': 'radio' },
    modelFile: { 'ui:widget': 'file', 'ui:options': { accept: '.json,.yaml' } },
    modelCsv: { 'ui:widget': 'file' },
    componentCsv: { 'ui:widget': 'file' },
    relationshipCsv: { 'ui:widget': 'file' },
    'ui:order': ['uploadType', 'fileName', 'modelFile', 'url', 'modelCsv'],
  },
}));

vi.mock('@/components/General/Modals/Modal', () => ({
  RJSFModalWrapper: (props) => {
    capturedRjsfProps = props;
    return <div data-testid="rjsf-wrapper" />;
  },
}));

vi.mock('@/components/DesignLifeCycle/common', () => ({
  Loading: ({ message }) => <div>{message}</div>,
}));

vi.mock('@/components/NotificationCenter', () => ({
  NotificationCenterContext: React.createContext({
    operationsCenterActorRef: {
      on: () => ({
        unsubscribe: vi.fn(),
      }),
    },
  }),
}));

vi.mock('machines/operationsCenter', () => ({
  OPERATION_CENTER_EVENTS: {
    EVENT_RECEIVED_FROM_SERVER: 'EVENT_RECEIVED_FROM_SERVER',
  },
}));

vi.mock('@/components/NotificationCenter/formatters/model_registration', () => ({
  ModelImportedSection: () => <div />,
  ModelImportMessages: () => <div />,
}));

vi.mock('./Stepper/CSVStepper', () => ({
  default: () => <div />,
}));

vi.mock('./Stepper/style', () => ({
  StyledDocsRedirectLink: ({ children, ...props }) => <a {...props}>{children}</a>,
}));

vi.mock('@/utils/utils', () => ({
  getUnit8ArrayDecodedFile: vi.fn(() => [65]),
}));

vi.mock('@/rtk-query/meshModel', () => ({
  useImportMeshModelMutation: () => [importModelReq],
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: (...args) => updateProgress(...args),
}));

describe('ImportModelModal', () => {
  beforeEach(() => {
    capturedRjsfProps = undefined;
    importModelReq.mockReset();
    updateProgress.mockReset();
  });

  it('flattens canonical allOf properties into the rendered file/url import subset', () => {
    render(<ImportModelModal isImportModalOpen={true} setIsImportModalOpen={vi.fn()} />);

    expect(capturedRjsfProps.schema.properties).toEqual(
      expect.objectContaining({
        uploadType: expect.objectContaining({
          enum: ['file', 'urlImport', 'csv'],
        }),
        modelFile: expect.objectContaining({ type: 'string' }),
        fileName: expect.objectContaining({ type: 'string' }),
        url: expect.objectContaining({ type: 'string' }),
      }),
    );
    expect(capturedRjsfProps.schema.required).toEqual(['uploadType']);
    expect(capturedRjsfProps.schema.allOf).toEqual([
      expect.objectContaining({
        then: { required: ['url'] },
      }),
    ]);
  });

  it('preserves canonical file widget settings while applying modal-specific uiSchema overrides', () => {
    render(<ImportModelModal isImportModalOpen={true} setIsImportModalOpen={vi.fn()} />);

    expect(capturedRjsfProps.uiSchema.modelFile).toEqual(
      expect.objectContaining({
        'ui:widget': 'file',
        'ui:options': { accept: '.json,.yaml' },
      }),
    );
    expect(capturedRjsfProps.uiSchema.fileName).toEqual({ 'ui:widget': 'hidden' });
    expect(capturedRjsfProps.uiSchema.modelCsv).toEqual({ 'ui:widget': 'hidden' });
    expect(capturedRjsfProps.uiSchema.componentCsv).toEqual({ 'ui:widget': 'hidden' });
    expect(capturedRjsfProps.uiSchema.relationshipCsv).toEqual({ 'ui:widget': 'hidden' });
    expect(capturedRjsfProps.uiSchema['ui:order']).toEqual(['uploadType', 'modelFile', 'url']);
  });
});
