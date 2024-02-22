import { configureStore } from '@reduxjs/toolkit'

import mqttReducer from './mqttSlice';
import dataReducer from './dataSlice';
import frameReducer from './frameSlice';
import i18nReducer from './i18nSlice';

export default configureStore({
  reducer: {
    mqtt:mqttReducer,
    data:dataReducer,
    frame:frameReducer,
    i18n:i18nReducer
  }
});