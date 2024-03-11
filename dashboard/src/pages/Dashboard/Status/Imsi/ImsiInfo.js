import { useSelector } from "react-redux";
import { RedoOutlined } from '@ant-design/icons';
import { Button,Tooltip } from 'antd';

export default function Imsi({getImsi}){
    const imsi=useSelector(state=>state.data.imsi);

    return (
        <div style={{display:'grid',border:'1px solid gray',width:'300px',lineHeight:'30px'}}>
            <div style={{gridColumnStart:1,gridColumnEnd:3,gridRowStart:1,gridRowEnd:2,paddingLeft:'5px',borderBottom:'1px solid gray',color:'white',backgroundColor:'#1677ff'}} >
                <span>IMSI Info</span>
                <Tooltip title="Refresh">
                    <Button type='primary' icon={<RedoOutlined />} style={{float:'right'}} onClick={()=>{getImsi();}} />
                </Tooltip>
            </div>
            <div style={{gridColumnStart:1,gridColumnEnd:2,gridRowStart:2,gridRowEnd:3,paddingLeft:'5px',borderRight:'1px solid gray'}}>imsi:</div>
            <div style={{gridColumnStart:2,gridColumnEnd:3,gridRowStart:2,gridRowEnd:3,paddingLeft:'5px'}}>{imsi?.res?.imsi}</div>
        </div>
    );
}