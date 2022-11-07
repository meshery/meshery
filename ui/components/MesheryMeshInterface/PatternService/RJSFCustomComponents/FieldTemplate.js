import React from "react";

export const CustomFieldTemplate=(props) => {
  const {  classNames, children } = props;
  return (
    <div className={classNames}>
      {children}
    </div>
  )
}