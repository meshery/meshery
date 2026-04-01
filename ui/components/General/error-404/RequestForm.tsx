import React, { useEffect } from 'react';
import { EVENT_TYPES } from 'lib/event-types';
import { useNotification } from 'utils/hooks/useNotification';
import { useGetOrgsQuery } from 'rtk-query/organization';
import OrgIcon from 'assets/icons/OrgIcon';
import { ErrorBoundary, FormControl, FormGroup, MenuItem, useTheme, NoSsr } from '@sistent/sistent';
import {
  OrgName,
  StyledSelect,
  StyledFormControlLabel,
  StyledTextField,
  ErrorSectionContent,
  StyledTypography,
  StyledFormButton,
} from './styles';
import { useGetCurrentAbilities } from 'rtk-query/ability';
import CustomErrorFallback from '../ErrorBoundary';
import { useDispatch, useSelector } from 'react-redux';
import { setKeys, setOrganization } from '@/store/slices/mesheryUi';

const RequestForm = () => {
  const {
    data: orgsResponse,
    isSuccess: isOrgsSuccess,
    isError: isOrgsError,
    error: orgsError,
  } = useGetOrgsQuery({});

  const theme = useTheme();
  let orgs = orgsResponse?.organizations || [];
  const { organization } = useSelector((state) => state.ui);
  const dispatch = useDispatch();
  const abilitiesResult = useGetCurrentAbilities(organization);

  useEffect(() => {
    if (abilitiesResult?.currentData?.keys) {
      dispatch(setKeys({ keys: abilitiesResult.currentData.keys }));
    }
  }, [abilitiesResult?.currentData?.keys]);
  const { notify } = useNotification();

  useEffect(() => {
    if (isOrgsError) {
      notify({
        message: `There was an error fetching available data ${orgsError?.data}`,
        event_type: EVENT_TYPES.ERROR,
      });
    }
  }, [orgsError]);

  const handleOrgSelect = (e) => {
    const id = e.target.value;
    const selected = orgs.find((org) => org.id === id);
    dispatch(setOrganization({ organization: selected }));
  };

  return (
    <NoSsr>
      <form>
        <ErrorSectionContent>
          <div>
            <StyledTypography variant="h6" component="h6">
              Request More Role(s)
            </StyledTypography>
            <StyledTextField id="request-message" minRows={3} maxRows={4} fullWidth multiline />
          </div>
          <FormControl fullWidth component="fieldset">
            <StyledTypography variant="h6" component="h6">
              Select Recipient
            </StyledTypography>
            <FormGroup>
              <StyledFormControlLabel
                key="SelectRecipient"
                control={
                  <StyledSelect
                    fullWidth
                    value={organization?.id ? organization.id : ''}
                    onChange={handleOrgSelect}
                    SelectDisplayProps={{
                      style: { display: 'flex' },
                    }}
                  >
                    {isOrgsSuccess &&
                      orgs &&
                      orgs?.map((org) => {
                        return (
                          <MenuItem key={org.id} value={org.id}>
                            <OrgIcon
                              width="24"
                              height="24"
                              secondaryFill={theme.palette.icon.secondary}
                            />
                            <OrgName>{org.name}</OrgName>
                          </MenuItem>
                        );
                      })}
                  </StyledSelect>
                }
              />
            </FormGroup>
          </FormControl>
          <StyledFormButton variant="outlined">Request Role(s)</StyledFormButton>
        </ErrorSectionContent>
      </form>
    </NoSsr>
  );
};

const RequestFormWithErrorBoundary = (props) => {
  return (
    <NoSsr>
      <ErrorBoundary customFallback={CustomErrorFallback}>
        <RequestForm {...props} />
      </ErrorBoundary>
    </NoSsr>
  );
};

export default RequestFormWithErrorBoundary;
