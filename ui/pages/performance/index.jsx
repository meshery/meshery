import { NoSsr } from "@mui/material";
import MesheryPerformanceComponent from "../../components/MesheryPerformance/Dashboard";
import Head from "next/head";

const Performance = () => {
  return (
    <NoSsr>
      <Head>
        <title>Performance Dashboard | Meshery</title>
      </Head>
      <MesheryPerformanceComponent />
    </NoSsr>
  );
};

export default Performance;
