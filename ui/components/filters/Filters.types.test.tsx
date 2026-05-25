import { describe, expect, it } from 'vitest';
import type {
  TypeView,
  TooltipIconProps,
  FilterSubmitPayload,
  YAMLEditorProps,
  ImportModalProps,
  PublishModalProps,
} from './Filters.types';

// Types module - the only meaningful test is "this compiles" but we can still
// assert the shape with values to ensure no accidental rename.
describe('Filters.types', () => {
  it('TypeView accepts "grid" and "table"', () => {
    const a: TypeView = 'grid';
    const b: TypeView = 'table';
    expect(a).toBe('grid');
    expect(b).toBe('table');
  });

  it('TooltipIconProps and other types are usable', () => {
    const tooltipProps: TooltipIconProps = {
      children: null,
      onClick: () => {},
      title: 'hello',
    };
    expect(tooltipProps.title).toBe('hello');

    const filterPayload: FilterSubmitPayload = {
      data: 'x',
      id: 'i',
      name: 'n',
      type: 't',
    };
    expect(filterPayload.name).toBe('n');

    const yamlEditorProps: YAMLEditorProps = {
      filter: {},
      onClose: () => {},
      onSubmit: () => {},
    };
    expect(yamlEditorProps).toBeTruthy();

    const importProps: ImportModalProps = {
      handleClose: () => {},
      handleImportFilter: () => {},
    };
    expect(importProps).toBeTruthy();

    const publishProps: PublishModalProps = {
      handleClose: () => {},
      handleSubmit: () => {},
      title: 'Publish',
    };
    expect(publishProps.title).toBe('Publish');
  });
});
