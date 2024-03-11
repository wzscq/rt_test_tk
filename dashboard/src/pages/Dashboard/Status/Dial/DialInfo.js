import { useSelector } from "react-redux";
import { RedoOutlined,CaretRightOutlined } from '@ant-design/icons';
import { Button,Tooltip } from 'antd';

export default function DialInfo({dialFunc}){
    const dial=useSelector(state=>state.data.dial);

    return (
        <div style={{display:'grid',border:'1px solid gray',width:'300px',lineHeight:'30px'}}>
            <div style={{gridColumnStart:1,gridColumnEnd:3,gridRowStart:1,gridRowEnd:2,paddingLeft:'5px',borderBottom:'1px solid gray',color:'white',backgroundColor:'#1677ff'}} >
                <span>Dial Info</span>
                <Tooltip title="Refresh">
                    <Button type='primary' icon={<RedoOutlined />} style={{float:'right'}} onClick={()=>{dialFunc('dialQuery');}} />
                </Tooltip>
                <Tooltip title="Start">
                    <Button type='primary' icon={<CaretRightOutlined />} style={{float:'right'}} onClick={()=>{dialFunc('dialTrigger');}} />
                </Tooltip>
            </div>
            <div style={{gridColumnStart:1,gridColumnEnd:2,gridRowStart:2,gridRowEnd:3,paddingLeft:'5px',borderRight:'1px solid gray',borderBottom:'1px solid gray'}}>status:</div>
            <div style={{gridColumnStart:2,gridColumnEnd:3,gridRowStart:2,gridRowEnd:3,paddingLeft:'5px',borderBottom:'1px solid gray'}}>{dial?.dial_status?'true':'false'}</div>
            <div style={{gridColumnStart:1,gridColumnEnd:2,gridRowStart:3,gridRowEnd:4,paddingLeft:'5px',borderRight:'1px solid gray'}}>IP:</div>
            <div style={{gridColumnStart:2,gridColumnEnd:3,gridRowStart:3,gridRowEnd:4,paddingLeft:'5px'}}>{dial?.IP}</div>
        </div>
    );
}