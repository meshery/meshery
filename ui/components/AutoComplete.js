// import React from 'react'
// //import ReactDOM from 'react-dom'
// import Downshift from 'downshift'
// import {all as starWarsNames} from 'starwars-names'
// import matchSorter from 'match-sorter'

// const items = starWarsNames.map(name => ({
//   value: name,
//   id: name.toLowerCase(),
// }))

// const getItems = value =>
//   value ? matchSorter(items, value, {keys: ['value']}) : items

// const itemToString = item => (item ? item.value : '')

// const stateReducer = (state, changes) => {
//   if (changes.type === Downshift.stateChangeTypes.blurButton) {
//     return {...changes, isOpen: true}
//   }
//   return changes
// }

// export class AutoComplete extends React.Component {
//   render() {
//     return (
//       <div>
//         <h1>Autocomplete rocks!</h1>
//         <div>
//           <Downshift stateReducer={stateReducer} itemToString={itemToString}>
//             {({
//               getLabelProps,
//               getInputProps,
//               getMenuProps,
//               getItemProps,
//               getToggleButtonProps,

//               clearSelection,

//               highlightedIndex,
//               selectedItem,
//               isOpen,
//               inputValue,
//             }) => (
//               <div>
//                 <label {...getLabelProps()}>Select a Star Wars Character</label>
//                 <input {...getInputProps()} />
//                 <button {...getToggleButtonProps()}>
//                   {isOpen ? 'close' : 'open'}
//                 </button>
//                 {selectedItem ? (
//                   <button onClick={clearSelection}>x</button>
//                 ) : null}
//                 <ul
//                   {...getMenuProps({
//                     style: {height: 200, overflowY: 'scroll'},
//                   })}
//                 >
//                   {isOpen
//                     ? getItems(inputValue).map((item, index) => (
//                         // eslint-disable-next-line react/jsx-key
//                         <li
//                           {...getItemProps({
//                             item,
//                             key: item.id,
//                             style: {
//                               backgroundColor:
//                                 index === highlightedIndex ? 'gray' : null,
//                             },
//                           })}
//                         >
//                           {item.value}
//                         </li>
//                       ))
//                     : null}
//                 </ul>
//               </div>
//             )}
//           </Downshift>
//         </div>
//       </div>
//     )
//   }
// }

import * as React from "react";
//import {render} from 'react-dom'
import Downshift from "downshift";

const items = [
  { value: "apple" },
  { value: "pear" },
  { value: "orange" },
  { value: "grape" },
  { value: "banana" },
  { value: "kiwi" },
];

export class AutoComplete extends React.Component {
  render() {
    return (
      <div>
        <Downshift
          onChange={(selection) => alert(selection ? `You selected ${selection.value}` : "Selection Cleared")}
          itemToString={(item) => (item ? item.value : "")}
        >
          {({
            getInputProps,
            getItemProps,
            getLabelProps,
            getMenuProps,
            isOpen,
            inputValue,
            highlightedIndex,
            selectedItem,
            getRootProps,
          }) => (
            <div>
              <h2> Auto Complete Component</h2>
              <div>
                <label {...getLabelProps()}>Enter a fruit</label>
              </div>
              <div style={{ display: "inline-block" }} {...getRootProps({}, { suppressRefError: true })}>
                <input {...getInputProps()} />
              </div>
              <ul {...getMenuProps()}>
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
  }
}
