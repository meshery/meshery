import React from "react";
import { Calendar, momentLocalizer } from "react-big-calendar";
import moment from "moment";
import PropTypes from "prop-types";
import "react-big-calendar/lib/css/react-big-calendar.css";

const localizer = momentLocalizer(moment);

MesheryCalendar.propTypes = {
  title: PropTypes.string.isRequired,
  events: PropTypes.array.isRequired,
  defaultDate: PropTypes.instanceOf(Date),
  defaultView: PropTypes.string,
  style: PropTypes.object.isRequired,
};

MesheryCalendar.defaultProps = {
  defaultDate: moment().toDate(),
  defaultView: "month",
};

export const MesheryCalendar = ({ title, events, defaultDate, defaultView, style }) => {
  return (
    <div>
      <h1>{title}</h1>

      <Calendar
        localizer={localizer}
        defaultDate={defaultDate}
        defaultView={defaultView}
        events={events}
        style={style}
      />
    </div>
  );
};
