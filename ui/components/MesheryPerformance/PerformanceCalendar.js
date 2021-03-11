import React from "react";
import { Calendar, momentLocalizer } from "react-big-calendar";
import moment from "moment";
import "react-big-calendar/lib/css/react-big-calendar.css";

const localizer = momentLocalizer(moment);

const ColoredDateCellWrapper = ({ children }) => {
  return React.cloneElement(React.Children.only(children), {
    style: {
      backgroundColor: "lightblue",
    },
  });
};

/**
 * PerformanceCalendar renders a calendar with details like tests run
 * @param {{ style?: React.CSSProperties }} props 
 * @returns 
 */
function PerformanceCalendar({ style }) {
  return (
    <div style={style}>
      <Calendar
        events={[]}
        views={['month', 'week', 'day']}
        step={60}
        showMultiDayTimes
        defaultDate={new Date()}
        components={{
          timeSlotWrapper: ColoredDateCellWrapper,
        }}
        localizer={localizer}
        style={{
          height: "100%",
        }}
      />
    </div>
  );
}

export default PerformanceCalendar;
