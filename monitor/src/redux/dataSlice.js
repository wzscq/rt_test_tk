import { createSlice } from '@reduxjs/toolkit';

const finalData = {
  exampleCode: "ftp_executor", 
  msg: null, 
  msg_type: "finally", 
  testData: {
    id: "460011895631209", 
    func_name: "ftp_executor", 
    type: "command", 
    result: {
      res: true, 
      data: {
        res: false, 
        msg: ["NotReachable"]
      }, 
      msg: null
    }, 
    code: 100, 
    taskId: "tc001", 
    imei: "867826050500420", 
    imsi: "460001933146754", 
    operator: "CHINA MOBILE"
  }
}

const initialState = {
  data:{
    exampleCode:"ftp_executor",
    msg: null,
    msg_type:"report",
    testData:{
      msgTpye: "nr5g-saReport",
      layer: "radio", 
      event: "hangon",
      throughput: {
        downlink: "0.0 bytes/s", 
        uplink: "0.0 bytes/s"
      }, 
      state: "idle", 
      taskId: "tc001", 
      imei: "867826050500420", 
      imsi: "460001933146754", 
      operator: "CHINA MOBILE",
      measures: {
        duplex_mode: "TDD",
        MCC: "460",
        MNC: "00", 
        cellID: "87D008", 
        PCID: "18", 
        TAC: "100000", 
        ARFCN: "504990", 
        band: "41", 
        NR_DL_bandwidth: "12", 
        RSRP: "-62", 
        RSRQ: "-11", 
        SINR: "26", 
        srxlev: "7488"
      }
    }
  },
  commandResult:{}
}

export const dataSlice = createSlice({
    name: 'data',
    initialState,
    reducers: {
      addDataItem:(state,action)=>{
        if(action.payload.msg_type==='report'){
          state.data=action.payload;
        } else {
          state.data.msg_type=action.payload.msg_type;
          state.data.result_res=action.payload.testData?.result?.res;
          state.data.result_msg=JSON.stringify(action.payload.testData?.result?.msg);
          state.data.result_data_res=action.payload.testData?.result?.data?.res;
          state.data.result_data_msg=JSON.stringify(action.payload.testData?.result?.data?.msg);
        }
      },
      setCommandResult:(state,action)=>{
        state.commandResult=action.payload;
      }
    }
});

export const { 
  addDataItem,
  setCommandResult
} = dataSlice.actions

export default dataSlice.reducer