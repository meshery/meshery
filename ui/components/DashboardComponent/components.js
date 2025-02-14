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
} from '@layer5/sistent';

import { iconMedium } from 'css/icons.styles';

const layoutIconProps = (theme) => ({
  fill: theme.palette.background.neutral.default,
  primaryFill: theme.palette.background.neutral.default,
  width: '30',
  height: '30',
});

export const AddWidgetsToLayoutPanel = ({ widgetsToAdd, editMode, onAddWidget }) => {
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
      width="22rem"
      bgcolor={theme.palette.background.card}
      boxShadow="0px 2px 10px rgba(0, 0, 0, 0.2)"
      marginBlock={'1rem'}
    >
      <Typography variant="h5"> Widgets</Typography>
      {widgetsToAdd.length == 0 && (
        <Typography variant="h6"> All widgets added to the layout </Typography>
      )}

      {widgetsToAdd.map(({ key, ...widget }) => (
        <StyledCard
          key={key}
          sx={{
            height: '18rem',
            width: '100%',
            minWidth: '16rem',
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
  );
};

export const StyledCard = ({ title, icon, children, sx = {}, button }) => {
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

export const LayoutActionButton = ({ Icon, label, action, description, isShown }) => {
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

// render the widget inside the layout
export const LayoutWidget = ({ widget, removeWidget, isEditMode }) => {
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
