import DashboardComponent from "../components/DashboardComponent";
import { NoSsr } from "@material-ui/core";

// import { withRouter } from 'next/router'

// const Index = ({router}) => {
  const Index = () => {
  // router.push('/performance');
  return (
  <NoSsr>
    <DashboardComponent />
  </NoSsr>
  );
}
// export default withRouter(Index);

export default Index;