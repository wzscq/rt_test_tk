import { createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios';

export const getHost=()=>{
  const rootElement=document.getElementById('root');
  const host=rootElement?.getAttribute("host");
  console.log("host:"+host);
  return host;
}

const host=getHost()+process.env.REACT_APP_SERVICE_API_PREFIX; //'/frameservice';

export const getMqttServer = createAsyncThunk(
  'getMqttServer',
  async () => {
    const config={
      url:host+'/flow/getMqttServer',
      method:'post'
    }
    const response =await axios(config);
    return response.data;
  }
);