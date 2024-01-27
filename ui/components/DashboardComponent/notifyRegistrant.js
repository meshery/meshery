import React, { useEffect } from 'react';

const NotifyRegistrant = () => {
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
    fetchData();
  }, []);

  return <div></div>;
};

export default NotifyRegistrant;
