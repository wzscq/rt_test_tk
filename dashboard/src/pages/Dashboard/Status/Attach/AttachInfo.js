import React from 'react';
import {useSelector} from 'react-redux';
import { RedoOutlined } from '@ant-design/icons';
import { Button,Tooltip } from 'antd';


export default function AttachInfo({attach, detach, attach_query}){
    const attachStatus=useSelector(state=>state.data.attachStatus);

    return (
        <div style={{display:'grid',border:'1px solid gray',width:'300px',lineHeight:'30px'}}>
            <div style={{gridColumnStart:1,gridColumnEnd:3,gridRowStart:1,gridRowEnd:2,paddingLeft:'5px',borderBottom:'1px solid gray',color:'white',backgroundColor:'#1677ff'}} >
                <span>Attach Status</span>
                <Tooltip title="Refresh">
                    <Button type='primary' icon={<RedoOutlined />} style={{float:'right'}} onClick={()=>{attach_query();}} />
                </Tooltip>
            </div>
            <div style={{gridColumnStart:1,gridColumnEnd:2,gridRowStart:2,gridRowEnd:3,paddingLeft:'5px',borderRight:'1px solid gray',borderBottom:'1px solid gray'}}>status:</div>
            <div style={{gridColumnStart:2,gridColumnEnd:3,gridRowStart:2,gridRowEnd:3,paddingLeft:'5px',borderBottom:'1px solid gray'}}>{attachStatus?.res}</div>
            <div style={{gridColumnStart:1,gridColumnEnd:2,gridRowStart:3,gridRowEnd:4,paddingLeft:'5px',borderRight:'1px solid gray'}}>
                
            </div>
            <div style={{gridColumnStart:2,gridColumnEnd:3,gridRowStart:3,gridRowEnd:4,paddingLeft:'5px'}}>
                <Button type='primary' style={{float:'left'}} onClick={()=>{attach();}}>Attach</Button>
                <Button type='primary' style={{float:'right'}} onClick={()=>{detach();}}>Detach</Button>
            </div>
        </div>
    );
}