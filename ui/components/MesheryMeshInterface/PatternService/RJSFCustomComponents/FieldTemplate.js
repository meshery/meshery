import React from "react";

export const CustomFieldTemplate=(props) => {
  const {  classNames, children } = props;
  return (
    <div className={classNames}>
      <div style={{ marginTop : "5px" }}> {" "}  </div>
      {children}
    </div>
  )
}