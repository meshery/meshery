import { useHandleFeedbackFormSubmissionMutation } from '@/rtk-query/user';
import { useNotification } from '@/utils/hooks/useNotification';
import { FeedbackButton, styled } from '@layer5/sistent';
import { EVENT_TYPES } from 'lib/event-types';
import _ from 'lodash';
import { useRouter } from 'next/router';

const FeedbackWrapper = styled('div')(({ theme }) => ({
  [theme.breakpoints.down('sm')]: {
    display: 'none',
  },
}));

const validateFeedback = (feedback) => {
  if (_.isEmpty(feedback) || _.isEmpty(feedback.label) || _.isEmpty(feedback.message)) {
    return false;
  }
  return true;
};

export const StickyFeedbackButton = ({ defaultMessage, defaultOpen, containerStyles }) => {
  const router = useRouter();
  const [submitFeedback] = useHandleFeedbackFormSubmissionMutation();

  const { notify } = useNotification();

  const onSubmit = async (feedback) => {
    if (!validateFeedback(feedback)) {
      notify({
        message: `validate error`,
        event_type: EVENT_TYPES.ERROR,
        // details: error_msg.toString(),
      });
      return;
    }
    const path = router.pathname;
    const userFeedbackRequestBody = {
      scope: feedback?.label,
      message: feedback?.message,
      page_location: path,
      metadata: {},
    };
    const resp = await submitFeedback({
      userFeedbackRequestBody,
    });

    if (resp.error) {
      notify({
        message: `error`,
        event_type: EVENT_TYPES.ERROR,
        // details: error_msg.toString(),
      });
      return;
    }
    notify({
      message: `SUCCESS`,
      event_type: EVENT_TYPES.SUCCESS,
      // details: error_msg.toString(),
    });
  };

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
