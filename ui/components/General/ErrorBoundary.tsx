import {
  Box,
  Modal,
  helpAndSupportModalSchema,
  helpAndSupportModalUiSchema,
  useTheme,
} from '@sistent/sistent';
import { useState } from 'react';
import { RJSFModalWrapper } from './Modals/Modal';
import SupportIcon from '@/assets/icons/support';
import { useNotification } from '@/utils/hooks/useNotification';
import { useSupportWebHookMutation } from '@/rtk-query/webhook';
import { EVENT_TYPES } from 'lib/event-types';
import { useGetLoggedInUserQuery, useGetProviderCapabilitiesQuery } from '@/rtk-query/user';

import {
  EditButton,
  FallbackWrapper,
  TextButton,
  ToolBarButtonContainer,
  TryAgainButton,
} from './style';
import { StickyFeedbackButton } from './feedback';

/**
 * CustomErrorFallback component can be use to show error message to users
 * This components can be passed to error boundary to have custom fallback component
 */
const CustomErrorFallback = (props) => {
  const [openSupportModal, setOpenSupportModal] = useState(false);

  const { error } = props;
  const theme = useTheme();
  const { notify } = useNotification();
  const [triggerWebhook] = useSupportWebHookMutation();
  const { data: userData } = useGetLoggedInUserQuery(undefined);
  const { data: providerData } = useGetProviderCapabilitiesQuery(undefined);
  const showSupportBasedOnProvider = providerData?.provider_type === 'remote';

  const handleOpenSupportModal = () => {
    setOpenSupportModal(true);
  };

  const handleCloseSupportModal = () => {
    setOpenSupportModal(false);
  };

  const pageUrl = window.location.href;
  const timestamp = new Date().toLocaleString();

  const errorMessage = `An error occurred on the page at ${pageUrl} on ${timestamp}.
  Error Details: ${error?.message}.`;

  const handleSupportFormSubmission = async (data) => {
    triggerWebhook({
      body: {
        memberFormOne: {
          ...data,
          firstname: userData?.first_name,
          lastname: userData?.last_name,
          email: userData?.email,
        },
      },
      type: 'support',
    })
      .unwrap()
      .then(() => {
        notify({
          message:
            'Your response has been recorded. We will endeavor to promptly contact you with a suitable solution',
          event_type: EVENT_TYPES.SUCCESS,
        });
      })
      .catch(() => {
        notify({
          message: 'Sorry we are unable to submit your request',
          event_type: EVENT_TYPES.ERROR,
        });
      });
  };

  return (
    <FallbackWrapper>
      <Box {...props}>
        <>
          {showSupportBasedOnProvider ? (
            <ToolBarButtonContainer style={{ marginTop: '0.7rem' }}>
              <EditButton
                variant="contained"
                style={{ marginRight: '0.7rem' }}
                onClick={handleOpenSupportModal}
              >
                <TextButton>Get Help</TextButton>
              </EditButton>
            </ToolBarButtonContainer>
          ) : null}
        </>

        <TryAgainButton color="primary" onClick={props.resetErrorBoundary}>
          <TextButton
            style={{
              color: theme.palette.text.default,
            }}
          >
            Try Again
          </TextButton>
        </TryAgainButton>

        <Modal
          open={openSupportModal}
          closeModal={handleCloseSupportModal}
          title="Help & Support"
          headerIcon={<SupportIcon width={24} height={24} fill="#51636B" />}
        >
          <RJSFModalWrapper
            schema={helpAndSupportModalSchema}
            uiSchema={helpAndSupportModalUiSchema}
            handleClose={handleCloseSupportModal}
            handleSubmit={handleSupportFormSubmission}
            handleNext={() => {}}
            title="Help & Support"
            helpText=""
            submitBtnText="Submit"
          />
        </Modal>
        {showSupportBasedOnProvider ? (
          <StickyFeedbackButton
            containerStyles={{ zIndex: 11 }}
            defaultMessage={errorMessage}
            defaultOpen={true}
          />
        ) : null}
      </Box>
    </FallbackWrapper>
  );
};

export default CustomErrorFallback;
