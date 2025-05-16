import { Typography } from '@layer5/sistent';

export function TabPanel(props) {
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
        <Typography component="div" style={{ paddingTop: 2 }}>
          {children}
        </Typography>
      )}
    </div>
  );
}
