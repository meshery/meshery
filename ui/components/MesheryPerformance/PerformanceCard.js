import React, { useEffect, useState } from 'react';
import moment from 'moment';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import {
  CustomTooltip,
  Typography,
  Avatar,
  Button,
  useTheme,
  Grid,
  IconButton,
  Link,
  Table,
  TableCell,
  TableRow,
} from '@layer5/sistent';
import FlipCard from '../FlipCard';
import PerformanceResults from './PerformanceResults';
import { MESHERY_CLOUD_PROD } from '../../constants/endpoints';
import { iconMedium } from '../../css/icons.styles';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useGetUserByIdQuery } from '@/rtk-query/user';
import useTestIDsGenerator from '@/components/hooks/useTestIDs';
import { BottomPart, CardButton, ResultContainer } from './style';

function PerformanceCard({
  profile,
  handleDelete,
  handleEdit,
  handleRunTest,
  handleProfile,
  requestFullSize,
  requestSizeRestore,
}) {
  const theme = useTheme();
  const [userAvatar, setUserAvatar] = useState(null);
  const {
    data: userData,
    isSuccess: isUserDataFetched,
    isError: userError,
  } = useGetUserByIdQuery(profile.user_id);
  const dataTestIDs = useTestIDsGenerator('performanceProfileCard');

  useEffect(() => {
    if (isUserDataFetched && userData && userData.avatar_url) {
      setUserAvatar(userData.avatar_url);
    } else if (userError) {
      console.error('Failed to fetch user profile with ID');
    }
  }, [isUserDataFetched, userError, userData]);

  const {
    id,
    name,
    endpoints,
    load_generators: loadGenerators,
    total_results: results,
    duration: testRunDuration,
    concurrent_request: concurrentRequest,
    qps,
    service_mesh: serviceMesh,
    content_type: contentType,
    request_body: requestBody,
    request_cookies: requestCookies,
    request_headers: requestHeaders,
    last_run: lastRun,
    metadata,
  } = profile;

  const [renderTable, setRenderTable] = useState(false);
  const tableData = [
    {
      name: 'Endpoints',
      value: endpoints?.join(', '),
    },
    {
      name: 'Load Generators',
      value: loadGenerators?.join(', '),
    },
    {
      name: 'Additional Option',
      value:
        metadata?.additional_options && metadata?.additional_options[0] !== ''
          ? JSON.parse(metadata?.additional_options[0])
          : '',
    },
    {
      name: 'Certifcate Name',
      value: metadata?.ca_certificate ? metadata.ca_certificate.name : '',
    },
    {
      name: 'Running Duration',
      value: testRunDuration,
    },
    {
      name: 'Concurrent Requests',
      value: concurrentRequest,
    },
    {
      name: 'queries/second',
      value: qps,
    },
    {
      name: 'Service Mesh',
      value: serviceMesh,
    },
    {
      name: 'Content Type',
      value: contentType,
      omitEmpty: true,
    },
    {
      name: 'Request Body',
      value: requestBody,
      omitEmpty: true,
    },
    {
      name: 'Cookies',
      value: requestCookies,
      omitEmpty: true,
    },
    {
      name: 'Request Headers',
      value: requestHeaders,
      omitEmpty: true,
    },
    {
      name: 'Created At',
      value: profile.created_at ? moment(profile.created_at).format('LLL') : 'unknown',
    },
    {
      name: 'Last Updated',
      value: profile.updated_at ? moment(profile.updated_at).format('LLL') : 'unknown',
    },
    {
      name: 'Last Run',
      value: profile.last_run ? moment(profile.last_run).format('LLL') : 'unknown',
    },
  ];

  function genericClickHandler(ev, fn) {
    ev.stopPropagation();
    fn();
  }

  return (
    <FlipCard
      onClick={() => {
        setRenderTable(false);
        requestSizeRestore();
      }}
      duration={600}
    >
      {/* FRONT PART */}
      <>
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <Typography variant="h6" component="div">
            {name}
          </Typography>
          <img
            src={`/static/img/load-test/${loadGenerators[0]}.svg`}
            alt="load-generator"
            height="24px"
          />
        </div>
        <ResultContainer>
          <div>
            <Typography
              variant="h2"
              component="div"
              style={{
                marginRight: '0.75rem',
                color: `${theme.palette.mode === 'dark' ? '#fff' : '#647881'}`,
              }}
            >
              {(results || '0').toLocaleString('en')}
            </Typography>
            <Typography
              variant="body1"
              sx={{
                color: theme.palette.text.disabled,
              }}
              component="div"
            >
              Results
            </Typography>
          </div>
        </ResultContainer>
        <div style={{}}>
          <BottomPart>
            <Link href={`${MESHERY_CLOUD_PROD}/user/${profile.user_id}`} target="_blank">
              <Avatar alt="profile-avatar" src={userAvatar} />
            </Link>
            <div
              style={{
                marginRight: '0.5rem',
                marginLeft: '0.5rem',
              }}
            >
              {lastRun && (
                <Typography
                  variant="caption"
                  style={{
                    fontStyle: 'italic',
                    color: `${
                      theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.7)' : '#647881'
                    }`,
                  }}
                >
                  Last Run: {moment(lastRun).format('LLL')}
                </Typography>
              )}
            </div>
          </BottomPart>
          <CardButton>
            <Button
              variant="outlined"
              onClick={(ev) =>
                genericClickHandler(ev, () => {
                  setRenderTable((renderTable) => {
                    if (renderTable) {
                      requestSizeRestore();
                      return false;
                    }

                    requestFullSize();
                    return true;
                  });
                })
              }
              disabled={!CAN(keys.VIEW_RESULTS.action, keys.VIEW_RESULTS.subject)}
              sx={{ marginRight: '0.5rem' }}
            >
              {renderTable ? 'Hide' : 'View'} Results
            </Button>
            <Button
              color="primary"
              variant="contained"
              onClick={(ev) => genericClickHandler(ev, handleProfile)}
              disabled={!CAN(keys.RUN_TEST.action, keys.RUN_TEST.subject)}
              sx={{ marginRight: '0.5rem' }}
            >
              Edit Profile
            </Button>
            <Button
              color="primary"
              variant="contained"
              onClick={(ev) => genericClickHandler(ev, handleRunTest)}
              disabled={!CAN(keys.RUN_TEST.action, keys.RUN_TEST.subject)}
            >
              Run Test
            </Button>
          </CardButton>
        </div>
        {renderTable ? (
          <div onClick={(ev) => ev.stopPropagation()} style={{ marginTop: '0.5rem' }}>
            <PerformanceResults
              // @ts-ignore
              CustomHeader={<Typography variant="h6">Test Results</Typography>}
              // @ts-ignore
              endpoint={`/api/user/performance/profiles/${id}/results`}
              // @ts-ignore
              elevation={0}
            />
          </div>
        ) : null}
      </>

      {/* BACK PART */}
      <>
        <Grid
          sx={{
            marginBottom: '0.25rem',
            minHeight: '6rem',
          }}
          container
          spacing={1}
          alignContent="space-between"
          alignItems="center"
        >
          <Grid item xs={8}>
            <Typography variant="h6" gutterBottom>
              {name} Details
            </Typography>
          </Grid>
          <Grid item xs={4}>
            <div
              style={{
                width: 'fit-content',
                margin: '0 0 0 auto',
              }}
            >
              <CustomTooltip title="Edit">
                <IconButton
                  onClick={(ev) => genericClickHandler(ev, handleEdit)}
                  data-testid={dataTestIDs('edit')}
                  disabled={
                    !CAN(keys.EDIT_PERFORMANCE_TEST.action, keys.EDIT_PERFORMANCE_TEST.subject)
                  }
                >
                  <EditIcon style={iconMedium} />
                </IconButton>
              </CustomTooltip>
              <CustomTooltip title="Delete">
                <IconButton
                  onClick={(ev) => genericClickHandler(ev, handleDelete)}
                  data-testid={dataTestIDs('delete')}
                  disabled={
                    !CAN(keys.DELETE_PERFORMANCE_TEST.action, keys.DELETE_PERFORMANCE_TEST.subject)
                  }
                >
                  <DeleteIcon style={iconMedium} />
                </IconButton>
              </CustomTooltip>
            </div>
          </Grid>
        </Grid>
        <Table size="small" dense>
          {tableData.map(function renderDesignTableRow(data) {
            const { name, value, omitEmpty } = data;
            return <DetailsTable key={name} rowKey={name} value={value} omitEmpty={omitEmpty} />;
          })}
        </Table>
      </>
    </FlipCard>
  );
}

// @ts-ignore
export default PerformanceCard;

function DetailsTable({ rowKey, value, omitEmpty }) {
  const [isExpanded, setIsExpanded] = useState(false);
  if (omitEmpty && (value === undefined || value === null)) {
    return null;
  }

  const MAX_TEXT_LENGTH = 150;
  const shouldShowButton =
    rowKey === 'Additional Option' && value && JSON.stringify(value).length > MAX_TEXT_LENGTH;
  const displayText =
    isExpanded && rowKey === 'Additional Option'
      ? JSON.stringify(value)
      : JSON.stringify(value)?.slice(0, MAX_TEXT_LENGTH);

  const handleExpandClick = (e) => {
    setIsExpanded((prevExpanded) => !prevExpanded);
    e.stopPropagation();
  };

  return (
    <TableRow>
      <TableCell>
        <b>{rowKey}</b>
      </TableCell>
      <TableCell style={{ maxWidth: '300px', overflow: 'hidden' }}>
        <p>{rowKey === 'Additional Option' ? displayText : value || 'none'}</p>
        {shouldShowButton && (
          <Link onClick={handleExpandClick}>{isExpanded ? 'Show Less' : 'Show More'}</Link>
        )}
      </TableCell>
    </TableRow>
  );
}
