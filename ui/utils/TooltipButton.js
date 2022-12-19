import { Button,Tooltip } from "@material-ui/core";

export default function TooltipButton({ children, onClick, title, variant,...props }) {
  return (
    <Tooltip title={title} placement="top" arrow interactive >
      <Button variant={variant} onClick={onClick} {...props}>
        {children}
      </Button>
    </Tooltip>
  );
}