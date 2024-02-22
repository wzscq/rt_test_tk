import { createSlice } from '@reduxjs/toolkit';

import { getMqttServer } from '../api';

const initialState = {
    mqttConfLoaded:false,
    mqttConf:{
        broker:"121.37.185.248",
        wsPort:9101,
        user:"mosquitto",
        password:"123456",
        uploadMeasurementMetrics:"realtime_measurement_reporting/"
    }
}

export const mqttSlice = createSlice({
    name: 'mqtt',
    initialState,
    reducers: {
        
    },
    extraReducers: (builder) => {
        //获取MQTT配置信息
        builder.addCase(getMqttServer.pending, (state, action) => {
            state.mqttConfLoaded=true;
        });
        builder.addCase(getMqttServer.fulfilled, (state, action) => {
            console.log("getMqttServer fulfilled:",action);
            console.log(action);
            state.mqttConf=action.payload.result;
            console.log(action.payload.result);
        });
        builder.addCase(getMqttServer.rejected , (state, action) => {
            console.log("getMqttServer return error:",action);
            if(action.error&&action.error.message){
              console.log(action.error.message);
            } else {
              console.log("获取MQTT服务配置失败");
            }
        });
    }
});

/*export const { 
    
} = mqttSlice.actions
*/

export default mqttSlice.reducer