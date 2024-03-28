import * as React from 'react';
import Button from '@material-ui/core/Button';
import TroubleshootingModal from './TroubleshootingModalComponent';
import SupportModal from './Modals/SupportModal';

const Troubleshoot = (props) => {
  const [open, setOpen] = React.useState(true);
  const [openSupportModal, setOpenSupportModal] = React.useState(false);
  const handleOpen = () => {
    setOpen(true);
  };
  const handleOpenSupportModal = () => {
    setOpenSupportModal(true);
  };

  const handleCloseSupportModal = () => {
    setOpenSupportModal(false);
  };

  return (
    <div>
      <Button variant="contained" color="primary" size="large" onClick={handleOpen}>
        Troubleshooting Guide
      </Button>
      <Button
        variant="contained"
        color="primary"
        size="large"
        onClick={handleOpenSupportModal}
        style={{ marginLeft: '1rem' }}
      >
        Get Help
      </Button>
      <SupportModal open={openSupportModal} handleClose={handleCloseSupportModal} />
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
