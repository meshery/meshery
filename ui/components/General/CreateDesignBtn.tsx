import { AddCircleIcon as AddIcon, Button } from '@sistent/sistent';
import Link from 'next/link';
import { iconMedium } from 'css/icons.styles';

const buttonSx = { m: 1, whiteSpace: 'nowrap' } as const;

function CreateDesignBtn() {
  return (
    <Link href="/configuration/design">
      <Button type="submit" variant="contained" color="primary" size="large" sx={buttonSx}>
        <AddIcon style={{ marginRight: '0.5rem', ...iconMedium }} />
        Create Design
      </Button>
    </Link>
  );
}

export default CreateDesignBtn;
