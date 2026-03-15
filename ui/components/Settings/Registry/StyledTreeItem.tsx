import React, { useState } from 'react';
import { Box, Typography, useTheme, Checkbox } from '@sistent/sistent';
import SearchBar from '@/utils/custom-search';
import debounce from '@/utils/debounce';
import { StyledTreeItemRoot as StyledTreeItemRootBase } from './MeshModel.style';
import { useWindowDimensions } from '@/utils/dimension';

const StyledTreeItemRoot = StyledTreeItemRootBase as any;

/**
 * Customized item component in mui-x-tree
 */
type StyledTreeItemProps = {
  labelText: React.ReactNode;
  nodeId: string;
  root?: boolean;
  search?: boolean;
  setSearchText?: (_value: string) => void;
  check?: boolean;
  [key: string]: any;
};

const StyledTreeItem = React.forwardRef<HTMLLIElement, StyledTreeItemProps>(
  function StyledTreeItem(props, ref) {
    const [checked, setChecked] = useState(false);
    const [hover, setHover] = useState(false);
    const { labelText, nodeId, root, search, setSearchText, check, ...other } = props;
    const theme = useTheme();
    const { width } = useWindowDimensions();
    const [isSearchExpanded, setIsSearchExpanded] = useState(false);

    return (
      <StyledTreeItemRoot
        nodeId={nodeId}
        theme={theme}
        onMouseEnter={() => setHover(true)}
        onMouseLeave={() => setHover(false)}
        root={root}
        lineColor={theme.palette.text.default || theme.palette.text.primary}
        label={
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              py: 1.5,
              px: 0,
            }}
          >
            {width < 1370 && isSearchExpanded ? null : (
              <div
                style={{
                  display: 'flex',
                  flexDirection: 'row',
                  alignItems: 'center',
                }}
              >
                <Typography variant="body1" style={{ color: `${root}` }}>
                  {labelText}
                </Typography>
              </div>
            )}

            {check && (
              <Checkbox
                onClick={(e) => {
                  e.stopPropagation();
                  setChecked((prevcheck) => !prevcheck);
                }}
                size="small"
                checked={checked}
                style={{
                  visibility: hover || checked ? 'hidden' : 'hidden', //TODO: make it visible when bulk status change is supported
                }}
              />
            )}
            {search && (
              <SearchBar
                onSearch={debounce((value) => setSearchText(value), 200)}
                expanded={isSearchExpanded}
                setExpanded={setIsSearchExpanded}
                setModelsFilters={() => {}}
                placeholder="Search"
              />
            )}
          </Box>
        }
        {...other}
        ref={ref}
      />
    );
  },
);

export default StyledTreeItem;
