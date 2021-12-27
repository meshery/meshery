/* eslint-disable react/prop-types */
import * as React from "react";
import Downshift from "downshift";

export const AutoComplete = ({ items, handleSelectionChange }) => {
  return (
    <div>
      <Downshift onChange={handleSelectionChange} itemToString={(item) => (item ? item.value : "")}>
        {({ getInputProps, getItemProps, isOpen, inputValue, highlightedIndex, selectedItem, getRootProps }) => (
          <div>
            <div style={{ display: "inline-block" }} {...getRootProps({}, { suppressRefError: true })}>
              <input {...getInputProps()} />
            </div>
            <ul>
              {isOpen
                ? items
                    .filter((item) => !inputValue || item.value.includes(inputValue))
                    .map((item, index) => (
                      // eslint-disable-next-line react/jsx-key
                      <li
                        {...getItemProps({
                          key: item.value,
                          index,
                          item,
                          style: {
                            backgroundColor: highlightedIndex === index ? "lightgray" : "white",
                            fontWeight: selectedItem === item ? "bold" : "normal",
                          },
                        })}
                      >
                        {item.value}
                      </li>
                    ))
                : null}
            </ul>
          </div>
        )}
      </Downshift>
    </div>
  );
};
