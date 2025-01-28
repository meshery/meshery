import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import {
  Timeline,
  TimelineItem,
  TimelineSeparator,
  TimelineConnector,
  TimelineContent,
  TimelineDot,
} from '@mui/lab';

const useStyles = makeStyles((theme) => ({
  paper: { padding: '6px 16px' },
  secondaryTail: {
    // backgroundColor: "black",
  },
  chipIcon: { width: theme.spacing(2.5) },
}));

export const ScrollIndicator = ({ items, handleClick, activeIndex }) => {
  const classes = useStyles();
  // let image = "/static/img/meshery-logo.png";

  return (
    <Timeline>
      {items &&
        items.map((item, index) => (
          <TimelineItem key={item + index}>
            <TimelineSeparator>
              <TimelineDot
                variant={activeIndex === index ? 'outlined' : 'default'}
                style={{ cursor: 'pointer' }}
                onClick={handleClick(index)}
              >
                <img
                  src={activeIndex === index ? item.activeIcon : item.inactiveIcon}
                  className={classes.chipIcon}
                />
              </TimelineDot>
              {index === items.length - 1 ? null : (
                <TimelineConnector className={classes.secondaryTail} />
              )}
            </TimelineSeparator>
            <TimelineContent></TimelineContent>
          </TimelineItem>
        ))}
    </Timeline>
  );
};
