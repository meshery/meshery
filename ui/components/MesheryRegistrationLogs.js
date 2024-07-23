import React, { useState, useEffect } from 'react';
import { CircularProgress, Box } from '@material-ui/core';
import axios from 'axios';

const Loading = () => {
  return (
    <Box sx={{ display: 'flex', justifyContent: 'center' }}>
      <CircularProgress />
    </Box>
  );
};

const EventsView = () => {
  const [data, setData] = useState('No Logs Available');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response1 = await axios.get(`/api/v2/events?category=["registration"]`);
        const path = response1.data.events[0].metadata.ViewLink;
        console.log('Path:', path);

        const response2 = await axios.get(`/api/system/fileView?file=${encodeURIComponent(path)}`);
        console.log('Response Data:', response2.data);

        setData(response2.data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return <Loading />;
  }

  if (error) {
    return <p>Error Loading Logs: {error}</p>;
  }

  return (
    <>
      <p>{JSON.stringify(data, null, 2)}</p>
    </>
  );
};

const MesheryRegistrationLogs = () => {
  return (
    <>
      <div>Meshery Registration Logs</div>
      <EventsView />
    </>
  );
};

export default MesheryRegistrationLogs;
