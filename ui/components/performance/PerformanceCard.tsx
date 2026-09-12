import React, { useState } from 'react';
import moment from 'moment';
import { Delete as DeleteIcon } from '@/assets/icons';
import {
  CustomTooltip,
  Typography,
  Avatar,
  Button,
  useTheme,
  Grid2,
  IconButton,
  Link,
  Table,
  TableCell,
  TableRow,
  EditIcon,
} from '@sistent/sistent';
import FlipCard from '../general/FlipCard';
import PerformanceResults from './PerformanceResults';
import { MESHERY_CLOUD_PROD } from '../../constants/endpoints';
import { iconMedium } from '../../css/icons.styles';

import { Keys } from '@meshery/schemas/permissions';
import { useResourceOwner } from '@/utils/hooks/useResourceOwner';
import useTestIDsGenerator from '@/utils/hooks/useTestIDs';
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
  const { owner, hasCloudProfile } = useResourceOwner(profile.owner);
  const dataTestIDs = useTestIDsGenerator('performanceProfileCard');

  const {
    id,
    name,
    endpoints,
    loadGenerators,
    totalResults: results,
    duration: testRunDuration,
    concurrentRequest,
    qps,
    serviceMesh,
    contentType,
    requestBody,
    requestCookies,
    requestHeaders,
    lastRun,
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
      value: profile.createdAt ? moment(profile.createdAt).format('LLL') : 'unknown',
    },
    {
      name: 'Last Updated',
      value: profile.updatedAt ? moment(profile.updatedAt).format('LLL') : 'unknown',
    },
    {
      name: 'Last Run',
      value: profile.lastRun ? moment(profile.lastRun).format('LLL') : 'unknown',
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
          {loadGenerators?.[0] && (
            <img
              src={`/static/img/load-test/${loadGenerators[0]}.svg`}
              alt="load-generator"
              height="24px"
            />
          )}
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
            {hasCloudProfile ? (
              <Link
                href={`${MESHERY_CLOUD_PROD}/user/${profile.owner}`}
                target="_blank"
                rel="noopener noreferrer"
              >
                <Avatar alt="profile-avatar" src={owner?.avatarUrl} />
              </Link>
            ) : (
              <Avatar alt="profile-avatar" src={owner?.avatarUrl} />
            )}
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
              permissionKey={Keys.PerformanceManagementViewResults}
              sx={{ marginRight: '0.5rem' }}
            >
              {renderTable ? 'Hide' : 'View'} Results
            </Button>
            <Button
              color="primary"
              variant="contained"
              onClick={(ev) => genericClickHandler(ev, handleProfile)}
              permissionKey={Keys.PerformanceManagementEditPerformanceTest}
              sx={{ marginRight: '0.5rem' }}
            >
              Edit Profile
            </Button>
            <Button
              color="primary"
              variant="contained"
              onClick={(ev) => genericClickHandler(ev, handleRunTest)}
              permissionKey={Keys.PerformanceManagementRunTest}
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
        <Grid2
          sx={{
            marginBottom: '0.25rem',
            minHeight: '6rem',
          }}
          container
          spacing={1}
          alignContent="space-between"
          alignItems="center"
          size="grow"
        >
          <Grid2 size={{ xs: 8 }}>
            <Typography variant="h6" gutterBottom>
              {name} Details
            </Typography>
          </Grid2>
          <Grid2 size={{ xs: 4 }}>
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
                  permissionKey={Keys.PerformanceManagementEditPerformanceTest}
                >
                  <EditIcon style={iconMedium} fill={theme?.palette?.icon?.default} />
                </IconButton>
              </CustomTooltip>
              <CustomTooltip title="Delete">
                <IconButton
                  onClick={(ev) => genericClickHandler(ev, handleDelete)}
                  data-testid={dataTestIDs('delete')}
                  permissionKey={Keys.PerformanceManagementDeletePerformanceTest}
                >
                  <DeleteIcon style={iconMedium} fill={theme?.palette?.icon?.default} />
                </IconButton>
              </CustomTooltip>
            </div>
          </Grid2>
        </Grid2>
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
