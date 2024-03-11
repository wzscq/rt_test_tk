import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    data:{},
    imsi:{},
    dial:{},
}

export const dataSlice = createSlice({
    name: 'data',
    initialState,
    reducers: {
      setData: (state, action) => {
        state.data=action.payload;
      },
      setImsi: (state, action) => {
        state.imsi=action.payload;
      },
      setDial: (state, action) => {
        state.dial=action.payload;
      }
    }
});

export const { 
  setData,
  setImsi,
  setDial
} = dataSlice.actions

export default dataSlice.reducer