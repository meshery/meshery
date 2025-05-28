import { ErrorBoundary } from '@layer5/sistent';
import WorkspaceContent from '@/components/SpacesSwitcher/WorkspaceContent';

const WorkSpaceContentDataTable = ({ workspace }) => {
  return (
    <ErrorBoundary>
      <WorkspaceContent workspace={workspace} />
    </ErrorBoundary>
  );
};

export default WorkSpaceContentDataTable;
