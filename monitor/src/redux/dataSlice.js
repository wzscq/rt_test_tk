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

const iperf_result = {
  data:{
    exampleCode: "iperf_executor", 
    msg: null, 
    msg_type: "report", 
    testData: {
      msgTpye: "unregisteredReport", 
      layer: "radio", 
      event: "hangon", 
      measures: {}, 
      throughput: {}, 
      state: "NO-SERVICE", 
      gps: {
        gps_status: "Valid", 
        gps_ant_status: "OK", 
        datetime: "2024-04-05 15:59:24.473338", 
        longitude_direction: "E", 
        longitude: 118.8350601196289, 
        latitude_direction: "N", 
        latitude: 31.9070987701416
      }, 
      taskId: "iperf_executor", 
      imei: "867826050500420", 
      imsi: "460001933146754", 
      operator: "CHINA MOBILE"}
    },
    commandResult:{},
    pingRec:[]
}

const iperf_result_1={
  data:{
    exampleCode: "iperf_executor", 
    msg: null, 
    msg_type: "finally", 
    testData: {
      id: "460011895631209", 
      func_name: "iperf_executor", 
      type: "command", 
      result: "argument 2: <class 'TypeError'>: wrong type", 
      code: 400, 
      taskId: "iperf_executor", 
      imei: "867826050500420", 
      imsi: "460001933146754", 
      operator: "CHINA MOBILE"}
  },
  commandResult:{},
  pingRec:[]
}

const ftp = {
  data:{
    exampleCode:"",
    msg: null,
    msg_type:"",
    testData:{
      msgTpye: "",
      layer: "", 
      event: "",
      throughput: {
        downlink: "", 
        uplink: ""
      }, 
      state: "", 
      taskId: "", 
      imei: "", 
      imsi: "", 
      operator: "",
      measures: {
        duplex_mode: "",
        MCC: "",
        MNC: "", 
        cellID: "", 
        PCID: "", 
        TAC: "", 
        ARFCN: "", 
        band: "", 
        NR_DL_bandwidth: "", 
        RSRP: "", 
        RSRQ: "", 
        SINR: "", 
        srxlev: ""
      }
    }
  },
  commandResult:{},
  pingRec:[]
}

const initialState = iperf_result_1;

export const dataSlice = createSlice({
    name: 'data',
    initialState,
    reducers: {
      addDataItem:(state,action)=>{
        //if(action.payload.msg_type==='report'){
          state.data=action.payload;
        /*} else {
          state.data.msg_type=action.payload.msg_type;
          state.data.result_res=action.payload.testData?.result?.res;
          state.data.result_msg=JSON.stringify(action.payload.testData?.result?.msg);
          state.data.result_data_res=action.payload.testData?.result?.data?.res;
          state.data.result_data_msg=JSON.stringify(action.payload.testData?.result?.data?.msg);
        }*/
      },
      setCommandResult:(state,action)=>{
        state.commandResult=action.payload;
        state.data={};
        state.pingRec=[];
      },
      setPingRec:(state,action)=>{
        state.pingRec=[...state.pingRec,action.payload];
      }
    }
});

export const { 
  addDataItem,
  setCommandResult,
  setPingRec
} = dataSlice.actions

export default dataSlice.reducer