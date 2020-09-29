import { NoSsr } from "@material-ui/core";
import MesherySMIResults from "../components/MesherySMIResults";
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux'
import Head from 'next/head';
import { getPath } from "../lib/path";

class SMIResults extends React.Component {
  componentDidMount () {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  render () {
    return (
      <NoSsr>
        <Head>
          <title>SMI Results | Meshery</title>
        </Head>
        <MesherySMIResults />
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch)
})

export default connect(
  null,
  mapDispatchToProps
)(SMIResults);