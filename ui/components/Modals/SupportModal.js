import React from 'react';
import Modal from '../Modal';
import PublicIcon from '@material-ui/icons/Public';
import { useNotification } from '@/utils/hooks/useNotification';
import { connect } from 'react-redux';
import { EVENT_TYPES } from '../../lib/event-types';
import axios from 'axios';

const helpModalSchema = {
  title: 'Support Form',
  properties: {
    subject: {
      type: 'string',
      title: 'Subject',
      description:
        'Enter a concise and descriptive title for your support request. This will help us quickly understand the nature of your inquiry.',
      minLength: 1,
      'x-rjsf-grid-area': '12',
    },
    description: {
      type: 'string',
      title: 'Description',
      description:
        'Please provide a detailed description of your issue or question. Include any relevant information that you think will help us assist you more effectively. The more details you provide, the better we can understand and address your concerns.',
      minLength: 10,
      format: 'textarea',
      'x-rjsf-grid-area': '12',
    },
    scope: {
      type: 'string',
      enum: ['Support', 'Community', 'Account', 'Commercial'],
      title: 'Scope of Questions',
      description: 'Select the category that best represents the nature of your inquiry.',
      default: 'Technical',
    },
  },
  required: ['subject', 'description'],
};

const helpModalUiSchema = {
  subject: {
    'ui:placeholder': 'Summary or title for your support request',
  },
  description: {
    'ui:placeholder': 'Detailed description of your support request',
  },
  scope: {
    'ui:widget': 'radio',
  },
};

// This modal is used in MeshMap also
function SupportModal(props) {
  const { open, handleClose, user } = props;
  const { notify } = useNotification();

  const handleSupportFormSubmission = async (data) => {
    try {
      await axios.post('https://hook.us1.make.com/r5qgpjel5tlhtyndcgjvkrdkoc65417y', {
        memberFormOne: {
          ...data,
          name: user?.first_name + ' ' + user?.last_name,
          email: user?.email,
        },
      });
      notify({
        message:
          'Your response has been recorded. We will endeavor to promptly contact you with a suitable solution',
        event_type: EVENT_TYPES.SUCCESS,
      });
    } catch (error) {
      notify({
        message: 'Sorry we are unable to submit your request',
        event_type: EVENT_TYPES.ERROR,
      });
    }
  };

  return (
    <Modal
      open={open}
      schema={helpModalSchema}
      uiSchema={helpModalUiSchema}
      title={`Help & Support`}
      handleClose={handleClose}
      handleSubmit={handleSupportFormSubmission}
      submitBtnText="Submit for Approval"
      submitBtnIcon={<PublicIcon data-cy="import-button" />}
      showInfoIcon={{
        text: 'Upon submitting your catalog item, an approval flow will be initiated.',
        link: 'https://docs.meshery.io/concepts/catalog',
      }}
    />
  );
}

const mapStateToProps = (state) => ({
  user: state.get('user')?.toObject(),
});

export default connect(mapStateToProps)(SupportModal);
