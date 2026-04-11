// Description: Contains functions related to DOM elements and react components.

import React from 'react';

// recursively check if element or any of its parent has the class
export const hasClass = (element, className) => {
  try {
    if (typeof element?.className == 'string' && element?.className?.includes(className)) {
      return true;
    }
    if (element?.parentElement) {
      return hasClass(element.parentElement, className);
    }
  } catch (e) {
    console.error(`Error in hasClass while checking for ${className} in `, element, e);
  }
  return false;
};

// recursively got throught component and its children and add the class to each of them
// This is required to prevent the clickaway listner from blocking the click event
// on the notification center IconButton create it as a HOC and use react.cloneElement to add the class
export const AddClassRecursively = ({ children, className }) => {
  return React.Children.map(children, (child) => {
    if (React.isValidElement(child)) {
      return React.cloneElement(child, {
        className: `${child.props.className} ${className}`,
        children: AddClassRecursively({ children: child.props.children, className }),
      });
    }

    // if child is a svg or animated svg string
    return child;
  });
};
