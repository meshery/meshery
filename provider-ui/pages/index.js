import ProviderComponent from "../components/ProviderComponent";
import { NoSsr } from "@material-ui/core";

class Index extends React.Component {
  componentDidMount () {
  }

  render () {
    return (
      <NoSsr>
        <ProviderComponent />
      </NoSsr>
    );
  }
}


export default Index;