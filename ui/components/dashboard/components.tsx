import React from 'react';
import {
  Box,
  useTheme,
  Typography,
  IconButton,
  Card,
  CardContent,
  CustomTooltip,
  DeleteIcon,
  DragIcon,
  type Theme,
} from '@sistent/sistent';

import { iconMedium } from 'css/icons.styles';
import { ActionButton } from './style';

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

const actionIconProps = (theme: Theme) => ({
  fill: theme.palette.background.neutral.default,
  primaryFill: theme.palette.background.neutral.default,
  ...iconMedium,
});

// Widget is the internal shape of a dashboard widget as used by LayoutWidget.
// Consumers should rely on WidgetItem from @sistent/sistent for the picker API.

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
            {icon}
            <Typography variant="h6" fontWeight="700" component="div" sx={{ mx: 1 }}>
              {title}
            </Typography>
          </Box>
          {button}
        </div>

        {children}
      </CardContent>
    </Card>
  );
};

type LayoutActionButtonProps = {
  Icon: React.ComponentType<{
    fill?: string;
    primaryFill?: string;
    width?: string | number;
    height?: string | number;
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
  const iconsProps = actionIconProps(theme);

  if (!isShown) {
    return null;
  }

  return (
    <ActionButton variant="text" onClick={action} endIcon={<Icon {...iconsProps} />}>
      <CustomTooltip title={description} variant="standard">
        <div>{label}</div>
      </CustomTooltip>
    </ActionButton>
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
