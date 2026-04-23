import { ErrorBoundary } from '@sistent/sistent';
import UnifiedSwitcherContent from './switcher/UnifiedSwitcherContent';

const WorkSpaceContentDataTable = ({ workspace }) => {
  return (
    <ErrorBoundary>
      <UnifiedSwitcherContent variant="workspace" workspace={workspace} />
    </ErrorBoundary>
  );
};

export default WorkSpaceContentDataTable;
