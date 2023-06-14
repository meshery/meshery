import React from "react";
import PropTypes from  "prop-types";
import { Tab as MuiTab, Box } from "@mui/material";
import Tabs from '@mui/material/Tabs';
import Typography from '@mui/material/Typography';

function TabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  );
}

function a11yProps(index) {
  return {
    id: `simple-tab-${index}`,
    'aria-controls': `simple-tabpanel-${index}`,
  };
}

export const Tab = ({ textColor, disabled, ...rest }) => {
  const [value, setValue] = React.useState(0);

  const handleChange = (event, newValue) => {
    setValue(newValue);
  };
  return (
    <Box sx={{ width: '100%' }}>
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs 
          value={value} 
          onChange={handleChange} 
          textColor={textColor}
          {...rest}
          aria-label={`tabs example`}>
          <MuiTab label="Item One" {...a11yProps(0)} />
          <MuiTab label="Item Two" {...a11yProps(1)} disabled={disabled} />
          <MuiTab label="Item Three" {...a11yProps(2)} />
        </Tabs>
      </Box>
      <TabPanel value={value} index={0}>
        Item One
      </TabPanel>
      <TabPanel value={value} index={1}>
        Item Two
      </TabPanel>
      <TabPanel value={value} index={2}>
        Item Three
      </TabPanel>
    </Box>
  )
}

Tab.propTypes = {
  textColor : PropTypes.string,
  indicatorColor : PropTypes.string,
  disabled: PropTypes.bool,
  orientation: PropTypes.oneOf(['vertical'])
}

Tab.defaultProps = {
  textColor : "primary",
  indicatorColor : "primary",
  disabled: false,
}