import { Button,Tooltip } from "@mui/material";

export default function TooltipButton({ children, onClick, title, variant,...props }) {
  return (
    <Tooltip title={title} placement="top" arrow interactive >
      <Button variant={variant} onClick={onClick} {...props}>
        {children}
      </Button>
    </Tooltip>
  );
}