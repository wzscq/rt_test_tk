import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    data:[
      /*{
        data: [
          {
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
          {
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
          }
        ],
        robot_info: {
          id: 1,
          robot_id: "2bee174b7d7c36e9b98bd8772e66af5e",
          map_id: "9ec5a9c61a374be4b0570077d7d6dd05",
          pixel_x: 616.645741420779,
          pixel_y: 285.72210123115497,
          pixel_theta: -110.04722431915482,
          record: "2023-08-18 21:27:05"
        },
        pcTime: 1692365376
      } */     
    ],
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
    currentPos:1,
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
    },
    robot_map:{},
    robot_map_record:undefined,
    indicator:null,
    points:[
      {pixel_x:10,pixel_y:10,rgb:'#FF0000',value:123},
      {pixel_x:20,pixel_y:10,rgb:'#00FF00',value:456},
      {pixel_x:30,pixel_y:10,rgb:'#0000FF',value:789}
    ],
    pointsLoaded:false,
}

export const dataSlice = createSlice({
    name: 'data',
    initialState,
    reducers: {
      addDataItem:(state,action)=>{
        state.data=[...state.data,action.payload];
        //更新robotinfo
        state.currentRobotInfo={...action.payload.robot_info,pcTime:action.payload.pcTime};
    
        //收敛event和message，设置当前UE信息
        action.payload.data.forEach(dataItem => {
          if(dataItem.event?.length>0){
            const imsi=dataItem.case_progress?.imsi;
            let eventPool=state.event[imsi];
            if(eventPool&&eventPool.length>0){
              const lastEventItem=eventPool[eventPool.length-1];
              dataItem.event.forEach(eventItem=>{
                if(eventItem.EventTime>lastEventItem.EventTime){
                  eventPool.push(eventItem);
                }
              });
            } else {
              eventPool=[...dataItem.event];
            }
            state.event={...state.event,[imsi]:[...eventPool]};
          }

          if(dataItem.msg?.length>0){
            const imsi=dataItem.case_progress?.imsi;
            let msgPool=state.message[imsi];
            if(msgPool&&msgPool.length>0){
              const lastMsgItem=msgPool[msgPool.length-1];
              dataItem.msg.forEach(msgItem=>{
                if(msgItem.MsgTime>lastMsgItem.MsgTime){
                  msgPool.push(msgItem);
                }
              });
            } else {
              msgPool=[...dataItem.msg];
            }
            state.message={...state.message,[imsi]:[...msgPool]};
          }

          //设置当前UE信息
          state.currentUes[dataItem.case_progress.imsi]={...dataItem};
        });
      },
      setDevice:(state,action)=>{
        console.log("setDevice:",action.payload);
        state.device=action.payload;
        state.deviceLoaded=true;
      },
      setCurrentPos:(state,action)=>{
        state.currentPos=action.payload;
      },
      setTestFileContent:(state,action)=>{
        if(action.payload[0]){
          const content=JSON.parse(action.payload[0]);
          //这里暂时不保留历史数据
          state.data=[content];
          //更新robotinfo
          state.currentRobotInfo={...content.robot_info,pcTime:content.pcTime};
          //更新robot_map
          if(state.robot_map.robot_id!==content.robot_info.robot_id||
            state.robot_map.map_id!==content.robot_info.map_id){
            state.robot_map={robot_id:content.robot_info.robot_id,map_id:content.robot_info.map_id};
          }
          //收敛event和message，设置当前UE信息
          content.data.forEach(dataItem => {
            if(dataItem.event?.length>0){
              const imsi=dataItem.case_progress?.imsi;
              let eventPool=state.event[imsi];
              if(eventPool&&eventPool.length>0){
                const lastEventItem=eventPool[eventPool.length-1];
                dataItem.event.forEach(eventItem=>{
                  if(eventItem.EventTime>lastEventItem.EventTime){
                    eventPool.push(eventItem);
                  }
                });
              } else {
                eventPool=[...dataItem.event];
              }
              state.event={...state.event,[imsi]:[...eventPool]};
            }

            if(dataItem.msg?.length>0){
              const imsi=dataItem.case_progress?.imsi;
              let msgPool=state.message[imsi];
              if(msgPool&&msgPool.length>0){
                const lastMsgItem=msgPool[msgPool.length-1];
                dataItem.msg.forEach(msgItem=>{
                  if(msgItem.MsgTime>lastMsgItem.MsgTime){
                    msgPool.push(msgItem);
                  }
                });
              } else {
                msgPool=[...dataItem.msg];
              }
              state.message={...state.message,[imsi]:[...msgPool]};
            }

            //设置当前UE信息
            state.currentUes[dataItem.case_progress.imsi]={...dataItem};
          });
        }
      },
      setRobotMapRecord:(state,action)=>{
        state.robot_map_record=action.payload;
      },
      setIndicator:(state,action)=>{
        state.indicator=action.payload;
      },
      setPoints:(state,action)=>{
        state.points=action.payload;
        state.pointsLoaded=true;
      },
    }
});

export const { 
  addDataItem,
  setDevice,
  setCurrentPos,
  setTestFileContent,
  setRobotMapRecord,
  setIndicator,
  setPoints
} = dataSlice.actions

export default dataSlice.reducer