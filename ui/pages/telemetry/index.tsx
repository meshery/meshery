import { useEffect } from 'react';
import { useRouter } from 'next/router';

// /telemetry has no content of its own; it lands on the Charts sub-page.
const TelemetryIndex = () => {
  const router = useRouter();
  useEffect(() => {
    router.replace('/telemetry/charts');
  }, [router]);
  return null;
};

export default TelemetryIndex;
