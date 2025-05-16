import { ClickAwayListener, styled } from '@layer5/sistent';
import React, { useEffect, useState } from 'react';

export const DesignerWrapper = styled('div')(({ theme }) => ({
  width: '100%',
  color: '#fff',
  position: 'fixed',
  top: 80,
  backgroundColor: theme.palette.mode === 'dark' ? '#222222' : '#477E96', // `theme.palette.mode` is used for dark/light mode
  zIndex: 1,
  marginLeft: '4px',
  padding: '4px 50px',
  transform: 'translateX(-40px)',
}));

export const StyledSpan = styled('span')(() => ({
  color: '#fff',
  fontStyle: 'italic',
  '&:hover': {
    cursor: 'pointer',
    textDecoration: 'underline',
  },
}));

export const StyledInput = styled('input')(() => ({
  background: 'transparent',
  border: 'none',
  color: '#fff',
  textDecoration: 'underline',
  '&:focus': {
    outline: 'none',
    border: 'none',
  },
}));

function CustomBreadCrumb({ title, onBack, titleChangeHandler }) {
  const [editing, setEditing] = useState(false);
  const [name, setName] = useState(title);

  useEffect(() => {
    const timer = setTimeout(() => {
      titleChangeHandler(name.trim());
    }, 400);

    return () => clearTimeout(timer);
  }, [name]);

  useEffect(() => {
    setName(title);
  }, [title]);

  const handleInputChange = (event) => {
    setName(event.target.value);
  };

  return (
    <DesignerWrapper>
      {'> '}
      <StyledSpan onClick={onBack}>Designs</StyledSpan>
      {' > '}

      {editing ? (
        <ClickAwayListener onClickAway={() => setEditing(false)}>
          <StyledInput value={name} onChange={handleInputChange} autoFocus />
        </ClickAwayListener>
      ) : (
        <StyledSpan onClick={() => setEditing(true)}>{title}</StyledSpan>
      )}
    </DesignerWrapper>
  );
}
export default CustomBreadCrumb;
