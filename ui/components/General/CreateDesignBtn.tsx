import React from 'react';
import { AddCircleIcon as AddIcon, Button } from '@sistent/sistent';
import Link from 'next/link';
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
