/**
* @param {object} obj check obj is empty or not
* @returns {boolean} true or false based on object is empty
*/

export const isEmptyObj= (obj)=> {
    for (let i in obj) return true;
    return false;
  }