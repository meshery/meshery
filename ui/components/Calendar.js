import React, { useState } from "react";
import Calendar from "react-calendar";

export function PerformanceCalendar() {
  const [value, onChange] = useState(new Date());

  return (
    <div>
      <header>
        <h1>Performance Calendar</h1>
      </header>
      <div>
        <Calendar onChange={onChange} showWeekNumbers value={value} />
      </div>
    </div>
  );
}
