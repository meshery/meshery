import type { ReactNode } from 'react';
import Header from './Header';

const layoutStyle = { margin: 20, padding: 20, border: '1px solid #DDD' };

interface LayoutProps {
  children?: ReactNode;
}

const Layout = (props: LayoutProps) => (
  <div style={layoutStyle}>
    <Header
      onDrawerToggle={() => {}}
      onDrawerCollapse={false}
      contexts={{}}
      activeContexts={[]}
      setActiveContexts={(_arg?: any) => {}}
      searchContexts={(_arg?: any) => {}}
      abilityUpdated={false}
    />
    {props.children}
  </div>
);
export default Layout;
