import { Button } from '@material-ui/core';
import { CustomTextTooltip } from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';

export default function TooltipButton({ children, onClick, title, variant, ...props }) {
  return (
    <CustomTextTooltip title={title} placement="top" interactive>
      <Button variant={variant} onClick={onClick} {...props}>
        {children}
      </Button>
    </CustomTextTooltip>
  );
}
