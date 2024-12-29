import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    data:{},
    imsi:{},
    dial:{},
    attachStatus:{},
    rat:{}
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
      },
      setAttachStatus: (state, action) => {
        state.attachStatus=action.payload;
      },
      setRATStatus: (state, action) => {
        state.rat=action.payload;
      }
    }
});

export const { 
  setData,
  setImsi,
  setDial,
  setAttachStatus,
  setRATStatus
} = dataSlice.actions

export default dataSlice.reducer