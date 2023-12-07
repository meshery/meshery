import * as React from 'react';
import Button from '@material-ui/core/Button';
import TroubleshootingModal from './TroubleshootingModalComponent';

const Troubleshoot = (props) => {
  const [open, setOpen] = React.useState(true);
  const handleOpen = () => {
    setOpen(true);
  };

  return (
    <div>
      <Button variant="contained" color="primary" size="large" onClick={handleOpen}>
        Troubleshooting Guide
      </Button>
      <TroubleshootingModal
        viewDataErrorMessage={props?.viewDataErrorMessage}
        viewHeaderErrorMessage={props?.viewHeaderErrorMessage}
        open={open}
        setOpen={setOpen}
      />
    </div>
  );
};

export default Troubleshoot;
