import { withRouter } from 'next/router'

const Index = ({router}) => {
  router.push('/performance');
  return('');
}
export default withRouter(Index);