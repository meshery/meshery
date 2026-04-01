const grafanaDateRangeToDate = (dt, startDate) => {
  const dto = new Date();
  switch (dt) {
    case 'now-2d':
      dto.setDate(dto.getDate() - 2);
      break;
    case 'now-7d':
      dto.setDate(dto.getDate() - 7);
      break;
    case 'now-30d':
      dto.setDate(dto.getDate() - 30);
      break;
    case 'now-90d':
      dto.setDate(dto.getDate() - 90);
      break;
    case 'now-6M':
      dto.setMonth(dto.getMonth() - 6);
      break;
    case 'now-1y':
      dto.setFullYear(dto.getFullYear() - 1);
      break;
    case 'now-2y':
      dto.setFullYear(dto.getFullYear() - 2);
      break;
    case 'now-5y':
      dto.setFullYear(dto.getFullYear() - 5);
      break;
    case 'now-1d/d':
      dto.setDate(dto.getDate() - 1);
      if (startDate) {
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case 'now-2d/d':
      dto.setDate(dto.getDate() - 2);
      if (startDate) {
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case 'now-7d/d':
      dto.setDate(dto.getDate() - 7);
      if (startDate) {
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case 'now-1w/w':
      dto.setDate(dto.getDate() - 6 - ((dto.getDay() + 8) % 7));
      if (startDate) {
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setDate(dto.getDate() + 6);
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case 'now-1M/M':
      dto.setMonth(dto.getMonth() - 1);
      if (startDate) {
        dto.setDate(1);
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setMonth(dto.getMonth());
        dto.setDate(0);
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case 'now-1y/y':
      dto.setFullYear(dto.getFullYear() - 1);
      if (startDate) {
        dto.setMonth(0);
        dto.setDate(1);
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setMonth(12);
        dto.setDate(0);
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case 'now/d':
      dto.setDate(dto.getDate() - 6 - ((dto.getDay() + 8) % 7));
      if (startDate) {
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case 'now':
      break;
    case 'now/w':
      dto.setDate(dto.getDate() - ((dto.getDay() + 7) % 7));
      if (startDate) {
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setDate(dto.getDate() + 6);
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case 'now/M':
      if (startDate) {
        dto.setDate(1);
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setMonth(dto.getMonth() + 1);
        dto.setDate(0);
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case 'now/y':
      if (startDate) {
        dto.setMonth(0);
        dto.setDate(1);
        dto.setHours(0);
        dto.setMinutes(0);
        dto.setSeconds(0);
        dto.setMilliseconds(0);
      } else {
        dto.setMonth(12);
        dto.setDate(0);
        dto.setHours(23);
        dto.setMinutes(59);
        dto.setSeconds(59);
        dto.setMilliseconds(999);
      }
      break;
    case 'now-5m':
      dto.setMinutes(dto.getMinutes() - 5);
      break;
    case 'now-15m':
      dto.setMinutes(dto.getMinutes() - 15);
      break;
    case 'now-30m':
      dto.setMinutes(dto.getMinutes() - 30);
      break;
    case 'now-1h':
      dto.setHours(dto.getHours() - 1);
      break;
    case 'now-3h':
      dto.setHours(dto.getHours() - 3);
      break;
    case 'now-6h':
      dto.setHours(dto.getHours() - 6);
      break;
    case 'now-12h':
      dto.setHours(dto.getHours() - 12);
      break;
    case 'now-24h':
      dto.setHours(dto.getHours() - 24);
      break;
    default:
      return new Date(parseFloat(dt));
  }
  return dto;
};

export default grafanaDateRangeToDate;
