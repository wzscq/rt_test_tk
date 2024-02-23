import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    data:{},
    event:{
      /*["460011895631208"]:[
        {
          Logfile: "testRecoder_0818-212920",
          pointIndex: 204,
          EventIndex: 11,
          EventTime: "2023-08-18 21:29:29.121",
          EventCode: "0x8D",
          name: "FTP Download TCPSlow"
        }
      ],
      ["460011895631209"]:[
        {
          Logfile: "testRecoder_0818-212920",
          pointIndex: 204,
          EventIndex: 11,
          EventTime: "2023-08-18 21:29:29.121",
          EventCode: "0x8D",
          name: "FTP Download TCPSlow"
        }
      ]*/
    },
    message:{
      /*["460011895631208"]:[
        {
          Logfile: "testRecoder_0818-212920",
          pointIndex: 430,
          MsgTime: "2023-08-18 21:29:33.894",
          MsgCode: "0x50000200",
          name: "NR->MeasurementReport"
        }
      ],
      ["460011895631209"]:[
        {
          Logfile: "testRecoder_0818-212920",
          pointIndex: 430,
          MsgTime: "2023-08-18 21:29:33.894",
          MsgCode: "0x50000200",
          name: "NR->MeasurementReport"
        }
      ]*/
    },
    deviceLoaded:false,
    device:null,
    currentPoint:0,
    currentRobotInfo: {
      /*id: 1,
      robot_id: "2bee174b7d7c36e9b98bd8772e66af5e",
      map_id: "9ec5a9c61a374be4b0570077d7d6dd05",
      pixel_x: 616.645741420779,
      pixel_y: 285.72210123115497,
      pixel_theta: -110.04722431915482,
      record: "2023-08-18 21:27:05"*/
    },
    currentUes:{
      /*["460011895631208"]:{
        radio: {
          measures_common: {
            Current_Network_Type: "",
            IMSI: "460011895631209",
            IMEI: "864451045593183",
            imsi: "460011895631209",
            imei: "864451045593183",
            duplex_mode: "",
            net_Tye: "NR_SA",
            plmn: "46001",
            operator: "\u4e2d\u56fd\u8054\u901a",
            operator_country: "cn",
            operator_short: "CHN-UNICOM",
            phone_name: "Redmi K30 5G Speed"
          },
          measures_nr: {
            SS_SINR:"11.37",
            SS_RSRP:"-84.59"
          }
        },
        event: [
          {
            Logfile: "testRecoder_0818-212920",
            pointIndex: 204,
            EventIndex: 11,
            EventTime: "2023-08-18 21:29:29.121",
            EventCode: "0x8D",
            name: "FTP Download TCPSlow"
          }
        ],
        msg: [
          {
            Logfile: "testRecoder_0818-212920",
            pointIndex: 430,
            MsgTime: "2023-08-18 21:29:33.894",
            MsgCode: "0x50000200",
            name: "NR->MeasurementReport"
          }
        ],
        time: 1692365376,
        case_progress: {
          session_id: 161920005569,
          full_path: "d:\\test_recoder\\testReport_001_161920005569\\progress_460011895631209.csv",
          imsi: "460011895631209",
          DevId: "28e79a383c89153af87780ec745e2481",
          Status: "FTPDownload",
          Times: "1/1",
          Progress: "8/30s",
          FailTimes: 0,
          OtherInfo: "Inst:2,729.30 kbps  Avg:2,799.36 kbps SuccRatio:100%(1/1)"
        }
      },
      ["460011895631209"]:{
        radio: {
          measures_common: {
            Current_Network_Type: "",
            IMSI: "460011895631209",
            IMEI: "864451045593183",
            imsi: "460011895631208",
            imei: "864451045593183",
            duplex_mode: "",
            net_Tye: "NR_SA",
            plmn: "46001",
            operator: "\u4e2d\u56fd\u8054\u901a",
            operator_country: "cn",
            operator_short: "CHN-UNICOM",
            phone_name: "Redmi K30 5G Speed"
          },
          measures_nr: {
            SS_SINR:"11.37",
            SS_RSRP:"-84.59"
          }
        },
        event: [
          {
            Logfile: "testRecoder_0818-212920",
            pointIndex: 204,
            EventIndex: 11,
            EventTime: "2023-08-18 21:29:29.121",
            EventCode: "0x8D",
            name: "FTP Download TCPSlow"
          }
        ],
        msg: [
          {
            Logfile: "testRecoder_0818-212920",
            pointIndex: 430,
            MsgTime: "2023-08-18 21:29:33.894",
            MsgCode: "0x50000200",
            name: "NR->MeasurementReport"
          }
        ],
        time: 1692365376,
        case_progress: {
          session_id: 161920005569,
          full_path: "d:\\test_recoder\\testReport_001_161920005569\\progress_460011895631209.csv",
          imsi: "460011895631208",
          DevId: "28e79a383c89153af87780ec745e2481",
          Status: "FTPDownload",
          Times: "1/1",
          Progress: "8/30s",
          FailTimes: 0,
          OtherInfo: "Inst:2,729.30 kbps  Avg:2,799.36 kbps SuccRatio:100%(1/1)"
        }
      }*/
    }
}

export const dataSlice = createSlice({
    name: 'data',
    initialState,
    reducers: {
      setData: (state, action) => {
        state.data=action.payload;
      },
    }
});

export const { 
  setData
} = dataSlice.actions

export default dataSlice.reducer