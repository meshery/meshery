import { NoSsr, Paper, withStyles } from "@material-ui/core";
import MesheryFilters from "../../components/Filters";
import { updatepagepath } from "../../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux'
import Head from 'next/head';
import { getPath } from "../../lib/path";
// import { resolveHref } from "next/dist/next-server/lib/router/router";
// import dataFetch from "../../lib/data-fetch";

const styles = {
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  }
}

class NewFilters extends React.Component {
  componentDidMount () {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  // async componentsDidMount() {
  //   await new Promise(resolve => {
  //     dataFetch(
  //       `/api/experimental/filter${query}`,
  //       {
  //         credentials: "include",
  //       },
  //       (result) => {
  //         console.log("FilterFile API", `/api/experimental/filter${query}`);
  //         updateProgress({ showProgress: false });
  //         if (result) {
  //           setFilters(result.filters || []);
  //           setPage(result.page || 0);
  //           setPageSize(result.page_size || 0);
  //           setCount(result.total_count || 0);
  //         }
  //       },
  //       handleError
  //     );
  //   })
  // }

  render () {
    return (
      <NoSsr>
        <Head>
          <title>Filters | Meshery</title>
        </Head>
        <Paper className={this.props.classes.paper}>
          <MesheryFilters />
        </Paper>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch)
})

export default withStyles(styles)(connect(
  null,
  mapDispatchToProps
)(NewFilters));