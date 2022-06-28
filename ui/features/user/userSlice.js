import { createSlice } from "@reduxjs/toolkit";

const userSlice = createSlice({
  name: "user",
  initialState: {},
  reducers: {
    updateUser: function (state, action) {
      return action.user;
    },
  },
});

export const { updateUser } = userSlice.actions;

export default userSlice.reducer;
