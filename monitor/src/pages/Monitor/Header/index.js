import {useSelector} from 'react-redux';

import './index.css';

export default function Header(){
  const {mqttStatus}=useSelector(state=>state.mqtt);
  
  return (
    <div className='monitor-header'>
      {' MQTT: '+mqttStatus}
    </div>
  )
}