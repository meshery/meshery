import React, { useEffect, useState } from 'react';
import { List, ListItem, Grid } from '@material-ui/core';

import KubernetesIcon from '../../../../assets/icons/KubernetesIcon';
import SMPIcon from '../../../../assets/icons/SMPIcon';
import LeftArrowIcon from '../../../../assets/icons/LeftArrowIcon';
import RightArrowIcon from '../../../../assets/icons/RightArrowIcon';
import {
  ButtonGrid,
  ListGrid,
  ListHeading,
  StyledCheckbox,
  StyledChip,
  StyledPaper,
  TransferButton,
} from './style';
import { Typography } from '@mui/material';
import { Tooltip } from '@mui/material';
import { TRANSFER_COMPONENT } from '../../../../utils/Enum';
import { Colors } from '@/themes/app';

function not(a, b) {
  return a.filter((value) => b.indexOf(value) === -1);
}

function intersection(a, b) {
  return a.filter((value) => b.indexOf(value) !== -1);
}

/**
 * Renders transfer component.
 *
 * @param {Object} props - The component props.
 * @param {String} props.name - This is the name of the data list.
 * @param {Array} props.assignableData - The assignable data list.
 * @param {Function} props.assignedData - The callback function to transfer assigned data list.
 * @param {Array} props.originalAssignedData - The already assigend data list.
 * @param {Element} props.emptyStateIconLeft - Icon for empty state of list left.
 * @param {String} props.emtyStateMessageLeft - Message for the empty state of the list left.
 * @param {Element} props.emptyStateIconRight - Icon for empty state of list right.
 * @param {String} props.emtyStateMessageRight - Message for the empty state of the list right.
 * @param {String} props.transferComponentType - Type of the component transfer (There is two types: chip and other).
 */

export default function TransferList({
  name,
  assignableData,
  assignedData,
  originalAssignedData,
  emptyStateIconLeft,
  emtyStateMessageLeft,
  emptyStateIconRight,
  emtyStateMessageRight,
  transferComponentType = TRANSFER_COMPONENT.OTHER,
  assignablePage,
  assignedPage,
  originalLeftCount,
  originalRightCount,
}) {
  const [checked, setChecked] = React.useState([]);
  const [left, setLeft] = useState([]);
  const [right, setRight] = useState([]);
  const [leftCount, setLeftCount] = useState(0);
  const [rightCount, setRightCount] = useState(0);

  useEffect(() => {
    setRight(originalAssignedData);
  }, [originalAssignedData]);

  useEffect(() => {
    setLeft(assignableData);
  }, [assignableData]);

  useEffect(() => {
    setLeftCount(originalLeftCount);
    setRightCount(originalRightCount);
  }, [originalLeftCount, originalRightCount]);

  const leftChecked = intersection(checked, left);
  const rightChecked = intersection(checked, right);

  useEffect(() => {
    assignedData(right);
  }, [right]);

  useEffect(() => {
    const handleScroll = (entries) => {
      const target = entries[0];
      if (target.isIntersecting) {
        assignablePage();
      }
    };

    const observer = new IntersectionObserver(handleScroll, { threshold: 1 });
    const sentinel = document.getElementById('leftList');
    if (sentinel) {
      observer.observe(sentinel);
    }

    return () => {
      if (sentinel) {
        observer.unobserve(sentinel);
      }
    };
  }, [assignablePage]);

  useEffect(() => {
    const handleScroll = (entries) => {
      const target = entries[0];
      if (target.isIntersecting) {
        assignedPage();
      }
    };

    const observer = new IntersectionObserver(handleScroll, { threshold: 1 });
    const sentinel = document.getElementById('rightList');
    if (sentinel) {
      observer.observe(sentinel);
    }

    return () => {
      if (sentinel) {
        observer.unobserve(sentinel);
      }
    };
  }, [assignedPage]);

  const handleToggle = (value) => () => {
    const currentIndex = checked.indexOf(value);
    const newChecked = [...checked];

    if (currentIndex === -1) {
      newChecked.push(value);
    } else {
      newChecked.splice(currentIndex, 1);
    }

    setChecked(newChecked);
  };

  const handleAllRight = () => {
    setRight(right.concat(left));
    setLeft([]);
    setLeftCount(0);
    setRightCount(originalLeftCount);
  };

  const handleCheckedRight = () => {
    setRight(right.concat(leftChecked));
    setLeft(not(left, leftChecked));
    setChecked(not(checked, leftChecked));
    setLeftCount((prevLeftCount) => prevLeftCount - leftChecked.length);
    setRightCount((prevRightCount) => prevRightCount + leftChecked.length);
  };

  const handleCheckedLeft = () => {
    setLeft(left.concat(rightChecked));
    setRight(not(right, rightChecked));
    setChecked(not(checked, rightChecked));
    setRightCount((prevRightCount) => prevRightCount - rightChecked.length);
    setLeftCount((prevLeftCount) => prevLeftCount + rightChecked.length);
  };

  const handleAllLeft = () => {
    setLeft(left.concat(right));
    setRight([]);
    setRightCount(0);
    setLeftCount(originalLeftCount);
  };

  const customList = (items, emptyStateIcon, emtyStateMessage, listId) => (
    <StyledPaper>
      <List dense component="div" role="list">
        {items?.length > 0 ? (
          items.map((item) => {
            const labelId = `transfer-list-item-${item.name}-label`;
            return (
              <ListItem
                key={item.id}
                role="listitem"
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  '&:hover': {
                    backgroundColor: '#00000010',
                  },
                }}
                onClick={handleToggle(item)}
              >
                {transferComponentType === TRANSFER_COMPONENT.CHIP ? (
                  <Tooltip title={item.name} placement="top">
                    <StyledChip
                      style={{ padding: '10px 0' }}
                      variant="outlined"
                      label={item.name}
                      onDelete={() => {}}
                      deleteIcon={<SMPIcon />}
                      icon={<KubernetesIcon />}
                    />
                  </Tooltip>
                ) : (
                  <Tooltip title={item.name} placement="top">
                    <Typography style={{ maxWidth: '230px', height: '1.5rem', overflow: 'hidden' }}>
                      {item.name}
                    </Typography>
                  </Tooltip>
                )}
                <StyledCheckbox
                  checked={checked.indexOf(item) !== -1}
                  tabIndex={-1}
                  disableRipple
                  inputProps={{
                    'aria-labelledby': labelId,
                  }}
                />
              </ListItem>
            );
          })
        ) : (
          <div
            style={{
              textAlign: 'center',
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'center',
              alignItems: 'center',
              height: '264px',
            }}
          >
            {emptyStateIcon}
            <Typography style={{ color: '#979797', padding: '24px 5px', lineHeight: 1 }}>
              {emtyStateMessage}
            </Typography>
          </div>
        )}
      </List>
      <div id={listId}></div>
    </StyledPaper>
  );

  return (
    <Grid container justifyContent="center" alignItems="center">
      <ListGrid>
        <ListHeading>
          Available {name} ({leftCount ? leftCount : 0})
        </ListHeading>
        {customList(left, emptyStateIconLeft, emtyStateMessageLeft, 'leftList')}
      </ListGrid>
      <ButtonGrid>
        <Grid container direction="column" alignItems="center">
          <TransferButton
            variant="outlined"
            size="small"
            color="primary"
            onClick={handleAllRight}
            disabled={left?.length === 0 || left.length < leftCount}
            aria-label="move all right"
          >
            <RightArrowIcon primaryFill={Colors.keppelGreen} width={18} height={18} />
            <RightArrowIcon
              primaryFill={Colors.keppelGreen}
              style={{ position: 'absolute', left: '27px' }}
              width={18}
              height={18}
            />
          </TransferButton>
          <TransferButton
            variant="outlined"
            size="small"
            color="primary"
            onClick={handleCheckedRight}
            disabled={leftChecked.length === 0}
            aria-label="move selected right"
          >
            <RightArrowIcon primaryFill={Colors.keppelGreen} width={18} height={18} />
          </TransferButton>
          <TransferButton
            variant="outlined"
            size="small"
            color="primary"
            onClick={handleCheckedLeft}
            disabled={rightChecked.length === 0}
            aria-label="move selected left"
          >
            <LeftArrowIcon primaryFill={Colors.keppelGreen} width={18} height={18} />
          </TransferButton>
          <TransferButton
            variant="outlined"
            size="small"
            color="primary"
            onClick={handleAllLeft}
            disabled={right.length === 0 || right.length < rightCount}
            aria-label="move all left"
          >
            <LeftArrowIcon primaryFill={Colors.keppelGreen} width={18} height={18} />
            <LeftArrowIcon
              primaryFill={Colors.keppelGreen}
              style={{ position: 'absolute', left: '27px' }}
              width={18}
              height={18}
            />
          </TransferButton>
        </Grid>
      </ButtonGrid>
      <ListGrid>
        <ListHeading>
          Assigned {name} ({rightCount ? rightCount : 0})
        </ListHeading>
        {customList(right, emptyStateIconRight, emtyStateMessageRight, 'rightList')}
      </ListGrid>
    </Grid>
  );
}
