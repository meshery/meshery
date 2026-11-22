// Description: Contains functions related to DOM elements and react components.

import React from 'react';

// recursively check if element or any of its parent has the class
export const hasClass = (element, className) => {
  if (!element || typeof className !== 'string' || !className.trim()) {
    return false;
  }
  try {
    if (typeof element.className === 'string' && element.className.includes(className)) {
      return true;
    }
    if (element.parentElement) {
      return hasClass(element.parentElement, className);
    }
  } catch (e) {
    console.error(`Error in hasClass while checking for ${className} in `, element, e);
  }
  return false;
};

// recursively go through component and its children and add the class to each of them
// This is required to prevent the clickaway listener from blocking the click event
// on the notification center IconButton create it as a HOC and use react.cloneElement to add the class
export const AddClassRecursively = ({ children, className }) => {
  if (!className || typeof className !== 'string') {
    return children;
  }
  return React.Children.map(children, (child) => {
    if (React.isValidElement(child)) {
      const existingClassName =
        typeof (child.props as { className?: unknown })?.className === 'string'
          ? (child.props as { className?: string }).className.trim()
          : '';
      const combinedClassName = existingClassName ? `${existingClassName} ${className}` : className;
      return React.cloneElement(child, {
        className: combinedClassName,
        children: AddClassRecursively({
          children: (child.props as { children?: React.ReactNode })?.children,
          className,
        }),
      });
    }

    // if child is a svg or animated svg string
    return child;
  });
};
