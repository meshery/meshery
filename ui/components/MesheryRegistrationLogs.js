import React, { useState, useEffect } from 'react';
import { useSelector, Provider } from 'react-redux';
import { selectEvents } from '../store/slices/events';
import { CircularProgress, Box } from '@material-ui/core';
import axios from 'axios';
import { store } from '../store';

const Loading = () => {
  return (
    <Box sx={{ display: 'flex', justifyContent: 'center' }}>
      <CircularProgress />
    </Box>
  );
};

const EventsView = () => {
  const events = useSelector(selectEvents) || [];
  const logsUrl = events.find((event) => event.category === 'registration')?.metadata?.ViewLink;
  const [data, setData] = useState('No Logs Available');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(
          `/api/system/fileView?file=${encodeURIComponent(logsUrl)}`,
        );
        setData(response.data);
      } catch (err) {
        setError(error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [events]);

  if (loading) {
    return <Loading />;
  }

  if (error) {
    return <p>Error Loading Logs: {error}</p>;
  }

  return (
    <>
      <pre style={{ whiteSpace: 'pre-line' }}>{JSON.stringify(data)}</pre>
    </>
  );
};

const MesheryRegistrationLogs = () => {
  return (
    <>
      <div>MesheryRegistrationLogs</div>
      <Provider store={store}>
        <EventsView />
      </Provider>
    </>
  );
};

export default MesheryRegistrationLogs;
