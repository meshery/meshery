import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';

const useListConnectionDefinitionsQuerySpy = vi.fn();
const useConnectionWizardSpy = vi.fn();

vi.mock('@meshery/schemas/mesheryApi', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@meshery/schemas/mesheryApi')>()),
  useListConnectionDefinitionsQuery: (...args: unknown[]) =>
    useListConnectionDefinitionsQuerySpy(...args),
}));

vi.mock('react-redux', () => ({
  useSelector: (selector: (state: unknown) => unknown) =>
    selector({ ui: { connectionMetadataState: {} } }),
}));

vi.mock('@/components/shared/Modal', () => ({
  Modal: ({ children, open }: { children: React.ReactNode; open: boolean }) =>
    open ? <div data-testid="modal">{children}</div> : null,
}));

vi.mock('./wizard/useConnectionWizard', () => ({
  useConnectionWizard: (...args: unknown[]) => useConnectionWizardSpy(...args),
}));

vi.mock('./wizard/WizardStepper', () => ({
  default: () => <div data-testid="wizard-stepper" />,
}));

vi.mock('./wizard/ConnectionWizard', () => ({
  default: () => <div data-testid="connection-wizard" />,
}));

vi.mock('./ConnectionWizard.helpers', () => ({
  buildConnectionWizardKindConfigs: () => [],
  DEFAULT_CONNECTION_DOCS_URL: 'https://docs.meshery.io',
}));

import ConnectionWizardModal from './ConnectionWizardModal';

describe('ConnectionWizardModal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useListConnectionDefinitionsQuerySpy.mockReturnValue({
      data: { connectionDefinitions: [] },
      isFetching: false,
    });
    useConnectionWizardSpy.mockReturnValue({
      ctx: { data: {} },
      steps: [],
      activeStep: 0,
      renderStepContent: () => null,
      handleNext: () => {},
      handleBack: () => {},
      canProceed: true,
      isSubmitting: false,
      isTransitioning: false,
      completedSteps: new Set(),
      reset: () => {},
      resetWizard: () => {},
    });
  });

  it('queries connection definitions with pageSize all when open', () => {
    render(<ConnectionWizardModal isOpen={true} onClose={() => {}} />);

    expect(useListConnectionDefinitionsQuerySpy).toHaveBeenCalledWith(
      expect.objectContaining({ pageSize: 'all' }),
      { skip: false },
    );
  });

  it('skips querying connection definitions when modal is closed', () => {
    render(<ConnectionWizardModal isOpen={false} onClose={() => {}} />);

    expect(useListConnectionDefinitionsQuerySpy).toHaveBeenCalledWith(
      expect.objectContaining({ pageSize: 'all' }),
      { skip: true },
    );
  });
});
