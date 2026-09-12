import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { StickyFeedbackButton } from './feedback';

const submitFeedback = vi.fn();
const notify = vi.fn();

vi.mock('@/rtk-query/user', () => ({
  useHandleFeedbackFormSubmissionMutation: () => [submitFeedback],
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('lib/event-types', () => ({
  EVENT_TYPES: {
    WARNING: 'warning',
    ERROR: 'error',
    SUCCESS: 'success',
  },
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ pathname: '/test-path' }),
}));

vi.mock('@sistent/sistent', () => ({
  FeedbackButton: (props: any) => (
    <div data-testid="feedback-btn">
      <button
        type="button"
        data-testid="submit-valid"
        onClick={() => props.onSubmit({ label: 'general', message: 'works great' })}
      >
        submit valid
      </button>
      <button
        type="button"
        data-testid="submit-invalid"
        onClick={() => props.onSubmit({ label: '', message: '' })}
      >
        submit invalid
      </button>
    </div>
  ),
  styled: (_tag: any) => () => {
    const Wrapper = ({ children }: any) => <div>{children}</div>;
    Wrapper.displayName = 'StyledWrapper';
    return Wrapper;
  },
}));

describe('StickyFeedbackButton', () => {
  beforeEach(() => {
    submitFeedback.mockReset();
    notify.mockReset();
  });

  it('warns when the feedback payload is empty', async () => {
    const user = userEvent.setup();
    render(<StickyFeedbackButton />);
    await user.click(screen.getByTestId('submit-invalid'));
    expect(notify).toHaveBeenCalledWith(
      expect.objectContaining({
        event_type: 'warning',
        message: expect.stringMatching(/unable to process your feedback/i),
      }),
    );
    expect(submitFeedback).not.toHaveBeenCalled();
  });

  it('notifies success on a successful submission', async () => {
    submitFeedback.mockResolvedValue({});
    const user = userEvent.setup();
    render(<StickyFeedbackButton />);
    await user.click(screen.getByTestId('submit-valid'));
    await waitFor(() => {
      expect(submitFeedback).toHaveBeenCalledWith(
        expect.objectContaining({
          userFeedbackRequestBody: expect.objectContaining({
            scope: 'general',
            message: 'works great',
            page_location: '/test-path',
          }),
        }),
      );
      expect(notify).toHaveBeenCalledWith(expect.objectContaining({ event_type: 'success' }));
    });
  });

  it('notifies an error when submission fails', async () => {
    submitFeedback.mockResolvedValue({ error: { message: 'boom' } });
    const user = userEvent.setup();
    render(<StickyFeedbackButton />);
    await user.click(screen.getByTestId('submit-valid'));
    await waitFor(() => {
      expect(notify).toHaveBeenCalledWith(expect.objectContaining({ event_type: 'error' }));
    });
  });
});
