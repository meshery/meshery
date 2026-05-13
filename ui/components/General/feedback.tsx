import { useHandleFeedbackFormSubmissionMutation } from '@/rtk-query/user';
import { useNotification } from '@/utils/hooks/useNotification';
import { FeedbackButton, styled } from '@sistent/sistent';
import { EVENT_TYPES } from 'lib/event-types';
import { useRouter } from 'next/router';
import { type CSSProperties, useCallback } from 'react';

const FeedbackWrapper = styled('div')(({ theme }) => ({
  [theme.breakpoints.down('sm')]: {
    display: 'none',
  },
}));

type FeedbackPayload = {
  label?: string;
  message?: string;
};

type StickyFeedbackButtonProps = {
  defaultMessage?: string;
  defaultOpen?: boolean;
  containerStyles?: CSSProperties;
};

const validateFeedback = (feedback?: FeedbackPayload) =>
  Boolean(feedback?.label?.trim() && feedback.message?.trim());

export const StickyFeedbackButton = ({
  defaultMessage,
  defaultOpen,
  containerStyles,
}: StickyFeedbackButtonProps) => {
  const router = useRouter();
  const [submitFeedback] = useHandleFeedbackFormSubmissionMutation();

  const { notify } = useNotification();

  const onSubmit = useCallback(
    async (feedback: FeedbackPayload) => {
      if (!validateFeedback(feedback)) {
        notify({
          message:
            'We are unable to process your feedback. Did you include a message in your submission?',
          event_type: EVENT_TYPES.WARNING,
        });
        return;
      }

      const resp = await submitFeedback({
        userFeedbackRequestBody: {
          scope: feedback.label,
          message: feedback.message,
          page_location: router.pathname,
          metadata: {},
        },
      });

      if (resp.error) {
        notify({
          message: 'Error submitting feedback. Check your Internet connection and try again.',
          event_type: EVENT_TYPES.ERROR,
        });
        return;
      }

      notify({
        message: 'Thank you! We have received your feedback.',
        event_type: EVENT_TYPES.SUCCESS,
      });
    },
    [notify, router.pathname, submitFeedback],
  );

  return (
    <FeedbackWrapper>
      <FeedbackButton
        defaultMessage={defaultMessage}
        containerStyles={containerStyles}
        renderPosition="right-middle"
        onSubmit={onSubmit}
        defaultOpen={defaultOpen}
      />
    </FeedbackWrapper>
  );
};
