/**
 * View info modal for the workspaces flow.
 *
 * Shows a view's owner, visibility, notes (markdown), and timestamps with a
 * "Copy Link" / "Save" action row in the footer. Composed on the shared
 * `Modal` primitive so the action bar plugs into the standard modal footer
 * rather than the body.
 */
import React, { useState, FC } from 'react';
import {
  Avatar,
  Chip,
  CircularProgress,
  FormLabel,
  Typography,
  ViewIcon,
  getFullFormattedTime,
  ModalButtonPrimary,
  ModalButtonSecondary,
  VisibilityChipMenu,
} from '@sistent/sistent';
import _ from 'lodash';
import { styled, useTheme } from '@/theme';
import { Modal } from '@/components/shared/Modal';
import { useGetViewQuery, useUpdateViewVisibilityMutation } from '@/rtk-query/view';
import { useGetLoggedInUserQuery, useGetUserProfileSummaryByIdQuery } from '@/rtk-query/user';
import { iconLarge } from 'css/icons.styles';
import RJSFWrapper from '../meshery-mesh-interface/PatternService/RJSF_wrapper';
import { MDEditor } from '../general/Markdown';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { handleUpdateViewVisibility, viewPath } from './SpacesSwitcher/hooks';
import rehypeSanitize from 'rehype-sanitize';
import { Lock, Public } from '@/assets/icons';
import { VIEW_VISIBILITY } from '@/utils/Enum';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';

type RowProps = { justifyContent?: string };

const Row = styled('div', {
  shouldForwardProp: (prop) => prop !== 'justifyContent',
})<RowProps>(({ justifyContent = 'space-between' }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: '0.5rem',
  flexWrap: 'wrap',
  justifyContent,
  paddingInline: '1rem',
  paddingBlock: '0.5rem',
}));

const OuterRow = styled(Row)({
  paddingInline: '0rem',
});

const Title = styled(Typography)(({ theme }) => ({
  color: theme.palette.text.primary,
  fontSize: '1.1rem',
  fontWeight: 500,
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

const StyledChip = styled(Chip)(({ theme }) => ({
  backgroundColor: theme.palette.background.paper,
}));

const LoaderShell = styled('div')({
  display: 'flex',
  justifyContent: 'center',
  padding: '1rem',
});

const MarkdownShell = styled('div')({
  width: '100%',
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  gap: '0.5rem',
});

const SavingShell = styled('div')({
  display: 'flex',
  alignItems: 'center',
});

const SavingSpinner = styled(CircularProgress)(({ theme }) => ({
  marginRight: '0.5rem',
  color: theme.palette.common.white,
}));

export interface ViewInfoModalProps {
  open: boolean;
  closeModal: () => void;
  viewId: string;
  viewName: string;
  metadata: unknown;
  refetch?: () => void;
}

export const ViewInfoModal_: FC<ViewInfoModalProps> = ({
  open,
  closeModal,
  viewId,
  viewName,
  metadata,
  refetch,
}) => {
  const [formState, setFormState] = useState(metadata);
  const viewRes = useGetViewQuery({ viewId }, { skip: !viewId });
  const userRes = useGetLoggedInUserQuery();
  const view = viewRes.data;
  const user = userRes.data;

  const isLoading =
    viewRes.isLoading || userRes.isLoading || viewRes.isFetching || userRes.isFetching;
  const formRef = React.useRef(null);
  const [saving, setSaving] = useState(false);
  const [updateView] = useUpdateViewVisibilityMutation();
  const theme = useTheme();

  const handleSave = async () => {
    setSaving(true);
    try {
      await updateView({
        id: view?.id,
        body: {
          metadata: formState,
        },
      }).unwrap();
      refetch && refetch();
      closeModal();
    } catch (error) {
      console.error('Failed to update view metadata:', error);
    } finally {
      setSaving(false);
    }
  };
  const canEdit = view?.userId === user?.id;
  const uiSchema = _.merge({}, UIFormSchema, {
    'ui:readonly': !canEdit,
  });
  const hasView = Boolean(view?.id);

  return (
    <Modal
      isOpen={open}
      onClose={closeModal}
      title={viewName}
      headerIcon={<ViewIcon {...iconLarge} fill={theme.palette.common.white} />}
      size="sm"
      actions={
        <>
          <CopyLinkButton link={hasView ? viewPath(view) : ''} disabled={!hasView} />
          {canEdit && <SaveButton isSaving={saving} onClick={handleSave} />}
        </>
      }
    >
      {isLoading ? (
        <LoaderShell>
          <CircularProgress />
        </LoaderShell>
      ) : (
        <div>
          <OuterRow>
            <Row justifyContent="start">
              <Title>Owner: </Title>
              <UserChip userId={view?.userId} />
            </Row>
            <Row justifyContent="start">
              <Title>Visibility: </Title>
              <ViewVisibilityMenu view={view} />
            </Row>
          </OuterRow>
          <RJSFWrapper
            hideTitle
            jsonSchema={InfoFormSchema}
            formData={formState}
            formRef={formRef}
            uiSchema={uiSchema}
            widgets={{
              markdown: MarkdownInput,
            }}
            onChange={(formData: unknown) => {
              setFormState(formData);
            }}
            RJSFWrapperComponent={StyledRJSFWrapper}
          />
          <Row justifyContent="start">
            <Title>Created: </Title>
            <Typography>{getFullFormattedTime(view?.createdAt)}</Typography>
          </Row>
          <Row justifyContent="start">
            <Title>Updated: </Title>
            <Typography>{getFullFormattedTime(view?.updatedAt)}</Typography>
          </Row>
        </div>
      )}
    </Modal>
  );
};

export const ViewInfoModal: FC<ViewInfoModalProps> = (props) => {
  return (
    <ProviderStoreWrapper>
      <ViewInfoModal_ {...props} />
    </ProviderStoreWrapper>
  );
};

export const UserChip: FC<{ userId?: string }> = ({ userId }) => {
  const userProfileRes = useGetUserProfileSummaryByIdQuery({ id: userId }, { skip: !userId });

  if (userProfileRes.isError || userProfileRes.isLoading) {
    return null;
  }
  const { avatarUrl } = userProfileRes?.data || {};
  const userName = formatUsername(userProfileRes?.data || {});

  return (
    <StyledChip avatar={<Avatar src={avatarUrl} />} label={userName || ''} variant="outlined" />
  );
};

const formatUsername = ({ firstName, lastName }: { firstName?: string; lastName?: string }) => {
  return `${firstName || ''} ${lastName || ''}`.trim();
};

type ViewLike = { id?: string; userId?: string; visibility?: string } | undefined;

const ViewVisibilityMenu: FC<{ view: ViewLike }> = ({ view }) => {
  const { data: userData } = useGetLoggedInUserQuery();
  const [updateView] = useUpdateViewVisibilityMutation();
  return (
    <VisibilityChipMenu
      value={view?.visibility || VIEW_VISIBILITY.PUBLIC}
      onChange={(value: string) =>
        handleUpdateViewVisibility({
          value: value,
          updateView: updateView,
          selectedResource: view,
        })
      }
      enabled={view?.userId === userData?.id}
      options={[
        [VIEW_VISIBILITY.PUBLIC, Public],
        [VIEW_VISIBILITY.PRIVATE, Lock],
      ]}
    />
  );
};

type MarkdownInputProps = {
  readonly?: boolean;
  value?: string;
  label?: string;
  onChange?: (value: string) => void;
};

const MarkdownInput: FC<MarkdownInputProps> = ({ readonly, value, label, onChange }) => {
  const theme = useTheme();

  const handleChange = (next: string | undefined) => {
    if (onChange && next !== value) {
      onChange(next || '');
    }
  };

  return (
    <MarkdownShell data-color-mode={theme.palette.mode}>
      <FormLabel>{label}</FormLabel>
      <MDEditor
        value={value || ''}
        onChange={handleChange}
        preview={readonly ? 'preview' : 'edit'}
        hideToolbar={Boolean(readonly)}
        previewOptions={{
          rehypePlugins: [[rehypeSanitize]],
        }}
      />
    </MarkdownShell>
  );
};

const CopyLinkButton: FC<{ link: string; disabled?: boolean; onClick?: () => void }> = ({
  onClick,
  link,
  ...rest
}) => {
  const { notify } = useNotification();

  const handleClick = async () => {
    try {
      await navigator.clipboard.writeText(link);
      notify({
        message: 'Link copied to clipboard',
        event_type: EVENT_TYPES.INFO,
      });
      if (onClick) onClick();
    } catch (err) {
      notify({
        message: 'Failed to copy link to clipboard',
        event_type: EVENT_TYPES.ERROR,
      });
    }
  };
  return (
    <ModalButtonSecondary onClick={handleClick} {...rest}>
      Copy Link
    </ModalButtonSecondary>
  );
};

const SaveButton: FC<{ onClick: () => void; isSaving: boolean; disabled?: boolean }> = ({
  onClick,
  isSaving,
  disabled,
}) => (
  <ModalButtonPrimary onClick={onClick} disabled={disabled || isSaving}>
    {isSaving ? (
      <SavingShell>
        <SavingSpinner size={20} />
        <Typography variant="body1"> Saving...</Typography>
      </SavingShell>
    ) : (
      'Save'
    )}
  </ModalButtonPrimary>
);
