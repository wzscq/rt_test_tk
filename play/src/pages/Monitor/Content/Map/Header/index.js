import {useState,useEffect} from 'react';
import {Select} from 'antd';
import {useSelector,useDispatch} from 'react-redux';

import {FRAME_MESSAGE_TYPE} from '../../../../../utils/constant';
import {createGetTestFilePointsMessage} from '../../../../../utils/normalOperations';
import { setIndicator } from '../../../../../redux/dataSlice';

import './index.css';

const { Option } = Select;
const queryFields=[
    {field:'id'},
    {field:'name'},
    {field:'extract_path'},
    {
        field:'legend',
        fieldType:'one2many',
        relatedModelID:'rt_indicator_legend',
        relatedField:'indicator_id',
        fields:[
            {field:'id'},
            {field:'rgb'},
            {field:'indicator_id'},
            {field:'start'},
            {field:'end'},
            {field:'sn'}
        ],
        sorter:[{field:'sn',order:'asc'}]
    },
];

export default function Header({mapInfo,sendMessageToParent}){
    const {origin,item:frameItem}=useSelector(state=>state.frame);
    const device=useSelector(state=>state.data.device.list?state.data.device.list[0]:undefined);

    const [options,setOptions]=useState([]);
    const [selectValue,setSelectValue]=useState('');
    const dispatch=useDispatch();

    const getQueryParams=(value)=>{
        const filter={name:'%'+value.replace("'","")+'%'};
        return {
            modelID:"rt_indicator",
            fields:queryFields,
            filter:filter,
            pagination:{current:1,pageSize:500}
        }
    }
    
    const onSearch=(value)=>{
        const queryParams=getQueryParams(value);
        if(queryParams){
            const frameParams={
                frameType:frameItem.frameType,
                frameID:frameItem.params.key,
                dataKey:"queryIndicatorResponse",
                origin:origin
            };

            const message={
                type:FRAME_MESSAGE_TYPE.QUERY_REQUEST,
                data:{
                    frameParams:frameParams,
                    queryParams:queryParams
                }
            }
            
            sendMessageToParent(message);
        }
    }

    const onChange=(value)=>{
        setSelectValue(value);
        const indicator=options.find(item=>item.id===value);
        dispatch(setIndicator(indicator));

        if(device){
            const frameParams={
                frameType:frameItem.frameType,
                frameID:frameItem.params.key,
                origin:origin
            };
            sendMessageToParent(createGetTestFilePointsMessage(frameParams,{deviceID:device.device_id,timestamp:device.timestamp,indicator:indicator}));
        }
    }

    const onFocus=()=>{
        onSearch('');
    }

    useEffect(()=>{
        const queryResponse=(event)=>{
            const {type,dataKey,data}=event.data;
            if(type===FRAME_MESSAGE_TYPE.QUERY_RESPONSE&&
                dataKey==="queryIndicatorResponse"){
                setOptions(data.list);
            }
        }
        window.addEventListener("message",queryResponse);
        return ()=>{
            window.removeEventListener("message",queryResponse);
        }
    },[setOptions]);

    const optionControls=options?options.map((item,index)=>{
        return (<Option key={item.id} value={item.id}>{item.name}</Option>);
    }):[];

    const selectIndicatorControl=(<Select
            style={{width:'100%'}}  
            size='small'
            value={selectValue} 
            showSearch
            onSearch={onSearch}
            onChange={onChange}
            onFocus={onFocus}
            filterOption={(input, option) =>
                option.children?.toLowerCase().indexOf(input.toLowerCase()) >= 0||
                option.value?.toLowerCase().indexOf(input.toLowerCase()) >= 0
            }>
                {optionControls}
            </Select>);

    return (
        <div className="monitor-map-header">
            {selectIndicatorControl}
            <div className='title'>
                {'building:'+mapInfo?.building_code+' floor:'+mapInfo?.floor+' pic:'+mapInfo?.picture_name}
            </div>
        </div>
    );
}