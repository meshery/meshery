import { ErrorBoundary } from '@sistent/sistent';
import WorkspaceContent from '@/components/workspaces/SpacesSwitcher/WorkspaceContent';

const WorkSpaceContentDataTable = ({ workspace }: { workspace: unknown }) => {
  return (
    <ErrorBoundary>
      <WorkspaceContent workspace={workspace} />
    </ErrorBoundary>
  );
};

export default WorkSpaceContentDataTable;
