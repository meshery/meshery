import React, { useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { loadEvents } from '../../store/slices/events';

const NotifyRegistrant = () => {
  const dispatch = useDispatch();

  const fetchData = async () => {
    try {
      const response = await fetch('api/meshmodels/registrants');
      const data = await response.json();

      if (data.total_count > 0) {
        for (const registrant of data.registrants) {
          console.log(registrant);
          const { hostname } = registrant;
          await fetch('api/meshmodels/nonRegisterEntity', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ Hostname: hostname }),
          });
        }
      }
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  };

  useEffect(() => {
    const isInitialVisit = localStorage.getItem('isInitialVisit') === null;

    if (isInitialVisit) {
      dispatch(loadEvents(fetchData, 1, {}));
      localStorage.setItem('isInitialVisit', 'true');
    } else {
      fetchData();
    }
  }, [dispatch]);

  return <div></div>;
};

export default NotifyRegistrant;
