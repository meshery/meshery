import React from 'react';
import { Button } from '@layer5/sistent';
import Link from 'next/link';
import AddIcon from '@mui/icons-material/AddCircleOutline';
import { iconMedium } from 'css/icons.styles';

function CreateDesignBtn() {
  return (
    <Link href="/configuration/design">
      <Button
        type="submit"
        variant="contained"
        color="primary"
        size="large"
        style={{ margin: '0.5rem 0.5rem', whiteSpace: 'nowrap' }}
      >
        <AddIcon style={{ marginRight: '0.5rem', ...iconMedium }} />
        Create Design
      </Button>
    </Link>
  );
}

export default CreateDesignBtn;
