import React, { useState } from 'react';
import { CircularProgress, Box, Button, Typography } from '@material-ui/core';
import useStyles from '../assets/styles/general/tool.styles';
import axios from 'axios';

const Loading = () => {
  return (
    <Box sx={{ display: 'flex', justifyContent: 'center', marginTop: '2rem' }}>
      <CircularProgress />
    </Box>
  );
};

const EventsView = ({ fetchData, buttonText }) => {
  const StyleClass = useStyles();
  const customInlineStyle = {
    marginBottom: '0.5rem',
    marginTop: '1rem',
    alignItems: 'center',
  };

  return (
    <div className={StyleClass.toolWrapper} style={customInlineStyle}>
      <div>
        <Typography align="center" variant="h6">
          Registry Logs
        </Typography>
      </div>
      <div>
        <Button variant="contained" color="primary" size="large" onClick={fetchData}>
          <Typography align="center" variant="subtitle2">
            {buttonText}
          </Typography>
        </Button>
      </div>
    </div>
  );
};

const TextFormatter = ({ rawData }) => {
  const addNewLinesBefore = (rawData) => {
    if (!rawData) return ''; // Check if rawData is undefined or null
    const pattern = /(\[|Components:)/g;
    let result = rawData.replace(pattern, '<br>$1');
    return result;
  };

  const formattedText = addNewLinesBefore(rawData);
  const textArray = formattedText.split('<br>');

  return (
    <div>
      {textArray.map((line, index) => (
        <Typography key={index}>{line}</Typography>
      ))}
    </div>
  );
};

const MesheryRegistrationLogs = () => {
  const StyleClass = useStyles();
  const [data, setData] = useState();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [buttonText, setButtonText] = useState('Display Logs');

  const fetchData = async () => {
    setLoading(true);
    try {
      const response1 = await axios.get(`/api/v2/events?category=["registration"]`);
      const path = response1.data.events[0].metadata.ViewLink;
      const response2 = await axios.get(`/api/system/fileView?file=${encodeURIComponent(path)}`);
      const rawData = response2.data;
      setData(rawData);
      setButtonText('Refresh Logs');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };
  return (
    <>
      <EventsView fetchData={fetchData} buttonText={buttonText} />
      {loading ? (
        <Loading />
      ) : error ? (
        <p>Error Loading Logs: {error}</p>
      ) : data ? ( // Conditionally render the Box only when data is available
        <Box sx={{ padding: '1rem' }} className={StyleClass.mainContainer}>
          <Box
            sx={{
              padding: '0.5rem',
              width: '100%',
              maxHeight: '400px',
              overflowY: 'auto',
            }}
          >
            <TextFormatter rawData={data} />
          </Box>
        </Box>
      ) : null}
    </>
  );
};

export default MesheryRegistrationLogs;
