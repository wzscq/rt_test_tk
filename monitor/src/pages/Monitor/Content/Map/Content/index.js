import {useSelector} from 'react-redux';

import './index.css';

export default function Content({map}){
    const data=useSelector(state=>state.data.data);

    const points=data.map(dataItem=>{
        const {pixel_x,pixel_y}=dataItem.robot_info;
        return (
            <div className='map-point' style={{left:pixel_x,top:pixel_y}}></div>
        );
    });

    return (
        <div className="monitor-map-content">
            <img src={map?.url} alt='' />
            {points}
        </div>
    );
}