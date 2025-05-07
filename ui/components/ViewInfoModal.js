import {
  Avatar,
  Chip,
  CircularProgress,
  FormLabel,
  Typography,
  ViewIcon,
  getFullFormattedTime,
  styled,
  useTheme,
} from '@layer5/sistent';
import React, { useState } from 'react';
import _ from 'lodash';
import { Box, Modal, ModalBody, ModalFooter } from '@layer5/sistent';
import { useGetViewQuery, useUpdateViewVisibilityMutation } from '@/rtk-query/view';
import { useGetLoggedInUserQuery, useGetUserProfileSummaryByIdQuery } from '@/rtk-query/user';
import { iconLarge } from 'css/icons.styles';
import { VisibilityChipMenu } from '@layer5/sistent';
import RJSFWrapper from './MesheryMeshInterface/PatternService/RJSF_wrapper';
import { MDEditor } from './Markdown';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { ModalButtonSecondary } from '@layer5/sistent';
import { handleUpdateViewVisibility, viewPath } from './SpacesSwitcher/hooks';
import { ModalButtonPrimary } from '@layer5/sistent';
import rehypeSanitize from 'rehype-sanitize';
import { Lock, Public } from '@mui/icons-material';
import { VIEW_VISIBILITY } from '@/utils/Enum';

const Row = styled('div')(({ justifyContent = 'space-between' }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: '0.5rem',
  flexWrap: 'wrap',
  justifyContent,
  paddingInline: '1rem',
  paddingBlock: '0.5rem',
}));

const Title = styled(Typography)(({ theme }) => ({
  color: theme.palette.text.primary,
  fontSize: '1.1rem',
  fontWeight: '500',
}));

const InfoFormSchema = {
  title: 'View Info',
  description: 'View Info',
  properties: {
    notes: {
      type: 'string',
      title: 'Notes',
      'x-rjsf-grid-area': 12,
      default: '',
    },
  },
};

const UIFormSchema = {
  notes: {
    'ui:autofocus': true,
    'ui:widget': 'markdown',
  },
};

const StyledRJSFWrapper = styled('div')({
  width: '100%',
  '& .rjsf': {
    width: '100%',
    margin: '0',
    padding: '0',
  },
});

export const ActionBox = styled(Box)(() => ({
  display: 'flex',
  justifyContent: 'end',
  width: '100%',
  gap: '1rem',
}));

export const ViewsInfoModal = ({ open, closeModal, view_id, view_name, metadata }) => {
  const [formState, setFormState] = useState(metadata);

  const viewRes = useGetViewQuery(
    { viewId: view_id },
    {
      skip: !view_id,
    },
  );
  const userRes = useGetLoggedInUserQuery();
  const view = viewRes.data;
  const user = userRes.data;

  const isLoading =
    viewRes.isLoading || userRes.isLoading || viewRes.isFetching || userRes.isFetching;
  const canRenderInfo = viewRes.isSuccess && userRes.isSuccess;

  const formRef = React.useRef(null);
  const [saving, setSaving] = useState(false);
  const [updateView] = useUpdateViewVisibilityMutation();

  const handleSave = () => {
    setSaving(true);
    updateView({
      id: view.id,
      body: {
        metadata: formState,
      },
    })
      .unwrap()
      .then(() => {
        setSaving(false);
        closeModal();
      });
  };
  const canEdit = view?.user_id === user?.id;
  const uiSchema = _.merge({}, UIFormSchema, {
    'ui:readonly': !canEdit,
  });
  const viewExists = (v) => Boolean(v && v?.id);
  const theme = useTheme();
  return (
    <Modal
      open={open}
      closeModal={closeModal}
      title={view_name}
      headerIcon={<ViewIcon {...iconLarge} fill={theme.palette.common.white} />}
      maxWidth="sm"
    >
      <ModalBody>
        {isLoading && (
          <div>
            <CircularProgress />
          </div>
        )}

        {!isLoading && canRenderInfo && (
          <div>
            <Row style={{ paddingInline: '0rem' }}>
              <Row justifyContent="start">
                <Title>Owner: </Title>
                <UserChip user_id={view?.user_id} />
              </Row>
              <Row justifyContent="start">
                <Title>Visibility: </Title>
                <ViewVisibilityMenu view={view} />
              </Row>
            </Row>
            <RJSFWrapper
              hideTitle
              jsonSchema={InfoFormSchema}
              formData={formState}
              formRef={formRef}
              uiSchema={uiSchema}
              widgets={{
                markdown: MarkdownInput,
              }}
              onChange={(formData) => {
                setFormState(formData);
              }}
              RJSFWrapperComponent={StyledRJSFWrapper}
            />

            <Row justifyContent="start">
              <Title>Created: </Title>
              <Typography>{getFullFormattedTime(view?.created_at)}</Typography>
            </Row>
            <Row justifyContent="start">
              <Title>Updated: </Title>
              <Typography>{getFullFormattedTime(view?.updated_at)}</Typography>
            </Row>
          </div>
        )}
      </ModalBody>
      <ModalFooter variant="filled">
        <ActionBox>
          <CopyLinkButton link={view?.id ? viewPath(view) : ''} disabled={!viewExists(view)} />
          {canEdit && <SaveButton isSaving={saving} onClick={handleSave} />}
        </ActionBox>
      </ModalFooter>
    </Modal>
  );
};

export const useViewsInfoModal = () => {
  const [data, setData] = React.useState({
    open: false,
    view_id: null,
  });

  const [formState, setFormState] = useState({});

  const closeModal = () => {
    setData({ open: false });
    setFormState({});
  };

  const openViewInfo = (view) => {
    setData({ open: true, view_id: view.id, view_name: view.name });
    setFormState(view?.metadata);
  };

  return {
    open: data.open,
    view_id: data.view_id,
    view_name: data.view_name,
    closeModal,
    openViewInfo,
    formState,
    setFormState,
  };
};

const StyledChip = styled(Chip)(({ theme }) => ({
  backgroundColor: theme.palette.background.paper,
}));

export const UserChip = ({ user_id }) => {
  const userProfileRes = useGetUserProfileSummaryByIdQuery({ id: user_id });

  if (userProfileRes.isError || userProfileRes.isLoading) {
    return null;
  }
  const { avatar_url } = userProfileRes.data || {};
  const userName = formatUsername(userProfileRes.data);

  return <StyledChip avatar={<Avatar src={avatar_url} />} label={userName} variant="outlined" />;
};

const formatUsername = ({ first_name, last_name }) => {
  return `${first_name || ''} ${last_name || ''}`.trim();
};

const ViewVisibilityMenu = ({ view }) => {
  const { data: userData } = useGetLoggedInUserQuery();
  const [updateView] = useUpdateViewVisibilityMutation();
  return (
    <VisibilityChipMenu
      value={view.visibility}
      onChange={(value) =>
        handleUpdateViewVisibility({ value: value, updateView: updateView, selectedResource: view })
      }
      enabled={view?.user_id === userData?.id}
      options={[
        [VIEW_VISIBILITY.PUBLIC, Public],
        [VIEW_VISIBILITY.PRIVATE, Lock],
      ]}
    />
  );
};

const MarkdownInput = (props) => {
  const preview = props.readonly ? 'preview' : 'edit';
  const hideToolbar = props.readonly ? true : false;
  const theme = useTheme();
  return (
    <div
      data-color-mode={theme.palette.mode}
      style={{
        width: '100%',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        gap: '0.5rem',
      }}
    >
      <FormLabel>{props.label}</FormLabel>
      <MDEditor
        value={props.value}
        onChange={props.onChange}
        preview={preview}
        hideToolbar={hideToolbar}
        previewOptions={{
          rehypePlugins: [[rehypeSanitize]],
        }}
      />
    </div>
  );
};

const CopyLinkButton = ({ onClick, link, ...props }) => {
  const { notify } = useNotification();

  const handleClick = () => {
    navigator.clipboard.writeText(link);
    notify({
      message: 'Link copied to clipboard',
      event_type: EVENT_TYPES.INFO,
    });
    if (onClick) onClick();
  };
  return (
    <ModalButtonSecondary onClick={handleClick} {...props}>
      Copy Link
    </ModalButtonSecondary>
  );
};

const SaveButton = ({ onClick, isSaving, ...props }) => (
  <ModalButtonPrimary onClick={onClick} {...props} disabled={props?.disabled || isSaving}>
    {isSaving ? (
      <div style={{ display: 'flex', alignItems: 'center' }}>
        <CircularProgress size={20} style={{ marginRight: '0.5rem', color: '#fff' }} />
        <Typography variant="body1"> Saving...</Typography>
      </div>
    ) : (
      'Save'
    )}
  </ModalButtonPrimary>
);
