import React from 'react';
import {
  Box,
  useTheme,
  Typography,
  Stack,
  AddIcon,
  IconButton,
  Card,
  CardContent,
  CustomTooltip,
  Button,
  DeleteIcon,
  DragIcon,
  Theme,
} from '@sistent/sistent';

import { iconMedium } from 'css/icons.styles';

type Widget = {
  key: string;
  title: string;
  thumbnail?: string;
  component?: React.ReactNode;
};

const layoutIconProps = (theme: Theme) => ({
  fill: theme.palette.background.neutral.default,
  primaryFill: theme.palette.background.neutral.default,
  width: '30',
  height: '30',
});

type AddWidgetsToLayoutPanelProps = {
  widgetsToAdd: Widget[];
  editMode: boolean;
  onAddWidget: (widget: Widget, key: string) => void;
};

export const AddWidgetsToLayoutPanel = ({
  widgetsToAdd,
  editMode,
  onAddWidget,
}: AddWidgetsToLayoutPanelProps) => {
  const theme = useTheme();
  const iconsProps = layoutIconProps(theme);

  if (!editMode) {
    return null;
  }
  return (
    <Box
      display="flex"
      flexDirection="column"
      gap="0.5rem"
      p={2}
      width="100%"
      bgcolor={theme.palette.background.card}
      boxShadow="0px 2px 10px rgba(0, 0, 0, 0.2)"
      marginBlock={'1rem'}
      minHeight={'17rem'}
    >
      <Typography variant="h5"> Widgets</Typography>
      {widgetsToAdd.length == 0 && (
        <Box marginInline={'auto'} marginTop={'5%'}>
          <Typography variant="h6"> All widgets added to the layout </Typography>
        </Box>
      )}

      <Box display="flex" flexWrap="wrap" gap="1rem">
        {widgetsToAdd.map(({ key, ...widget }) => (
          <StyledCard
            key={key}
            sx={{
              height: '18rem',
              width: '100%',
              maxWidth: '16rem',
            }}
          >
            <Stack direction="row" alignItems="center" justifyContent="space-between">
              <Typography variant="button">{widget.title}</Typography>
              <IconButton onClick={() => onAddWidget(widget, key)}>
                <AddIcon {...iconsProps} />
              </IconButton>
            </Stack>
            <img
              src={widget.thumbnail}
              alt={widget.title}
              style={{
                width: '100%',
                height: 'auto',
                objectFit: 'contain',
              }}
            />
          </StyledCard>
        ))}
      </Box>
    </Box>
  );
};

type StyledCardProps = {
  title?: React.ReactNode;
  icon?: React.ReactNode;
  children?: React.ReactNode;
  sx?: object;
  button?: React.ReactNode;
};

export const StyledCard = ({ title, icon, children, sx = {}, button }: StyledCardProps) => {
  const theme = useTheme();
  return (
    <>
      <Card
        sx={{
          minWidth: 275,
          height: '100%',
          ...sx,
          backgroundColor: theme.palette.background.elevatedComponents,
        }}
      >
        <CardContent>
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
            }}
          >
            <Box
              sx={{
                display: 'flex',
                mb: 1.5,
              }}
            >
              {icon} {/* Card Icon */}
              <Typography variant="h6" fontWeight="700" component="div" sx={{ mx: 1 }}>
                {title}
              </Typography>
            </Box>
            {button}
          </div>

          {children}
        </CardContent>
      </Card>
    </>
  );
};

type LayoutActionButtonProps = {
  Icon: React.ComponentType<{
    fill?: string;
    primaryFill?: string;
    width?: string;
    height?: string;
  }>;
  label: React.ReactNode;
  action: () => void;
  description: React.ReactNode;
  isShown: boolean;
};

export const LayoutActionButton = ({
  Icon,
  label,
  action,
  description,
  isShown,
}: LayoutActionButtonProps) => {
  const theme = useTheme();
  const iconsProps = layoutIconProps(theme);

  if (!isShown) {
    return null;
  }

  return (
    <>
      <Button
        variant="text"
        style={{ color: iconsProps.fill }}
        onClick={action}
        endIcon={<Icon {...iconsProps} />}
      >
        <CustomTooltip title={description} fontSize="1rem" variant="standard">
          <div>{label}</div>
        </CustomTooltip>
      </Button>
    </>
  );
};

type LayoutWidgetProps = {
  widget: Widget;
  removeWidget: (key: string) => void;
  isEditMode: boolean;
};

export const LayoutWidget = ({ widget, removeWidget, isEditMode }: LayoutWidgetProps) => {
  const theme = useTheme();
  const iconsProps = layoutIconProps(theme);

  return (
    <>
      {isEditMode && (
        <Box
          justifyContent="end"
          alignItems="center"
          gap="1"
          display="flex"
          backgroundColor={theme.palette.background.default}
        >
          <IconButton onClick={() => removeWidget(widget.key)}>
            <DeleteIcon {...iconsProps} {...iconMedium} />
          </IconButton>
          <IconButton className="react-grid-dragHandleExample">
            <DragIcon fill={iconsProps.fill} {...iconMedium} />
          </IconButton>
        </Box>
      )}
      {widget.component}
    </>
  );
};
