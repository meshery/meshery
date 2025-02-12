import * as React from 'react';
import TroubleshootingModal from './TroubleshootingModalComponent';
import { RJSFModalWrapper } from './Modal';
import {
  Modal,
  helpAndSupportModalSchema,
  helpAndSupportModalUiSchema,
  Button,
} from '@layer5/sistent';
import { useNotification } from '@/utils/hooks/useNotification';
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import SupportIcon from '@/assets/icons/support';
import { EVENT_TYPES } from 'lib/event-types';

import { useSupportWebHookMutation } from '@/rtk-query/webhook';
import { StickyFeedbackButton } from './General/feedback';

const Troubleshoot = ({ errorMessage, ...props }) => {
  const [open, setOpen] = React.useState(true);
  const [openForm, setOpenForm] = React.useState(false);
  const { notify } = useNotification();
  const [triggerWebhook] = useSupportWebHookMutation();
  const handleOpen = () => {
    setOpen(true);
  };

  const { data: userData } = useGetLoggedInUserQuery();

  const handleOpenForm = () => {
    setOpenForm(true);
  };

  const handleSupportFormClose = () => {
    setOpenForm(false);
  };

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
    <div>
      <Button variant="contained" color="primary" size="large" onClick={handleOpen}>
        Troubleshooting Guide
      </Button>
      <Button
        variant="contained"
        color="primary"
        size="large"
        onClick={handleOpenForm}
        style={{
          marginLeft: '0.31rem',
        }}
      >
        Get Help
      </Button>
      <TroubleshootingModal
        viewDataErrorMessage={props?.viewDataErrorMessage}
        viewHeaderErrorMessage={props?.viewHeaderErrorMessage}
        open={open}
        setOpen={setOpen}
      />

      <Modal
        open={openForm}
        closeModal={handleSupportFormClose}
        title="Help & Support"
        headerIcon={<SupportIcon style={{ height: '24px', width: '24px' }} />}
      >
        <RJSFModalWrapper
          schema={helpAndSupportModalSchema}
          uiSchema={helpAndSupportModalUiSchema}
          handleClose={handleSupportFormClose}
          handleSubmit={handleSupportFormSubmission}
          submitBtnText="Submit"
        />
      </Modal>

      <StickyFeedbackButton
        containerStyles={{ zIndex: 11 }}
        defaultMessage={errorMessage}
        defaultOpen={true}
      />
    </div>
  );
};

export default Troubleshoot;
