import {useSelector,useDispatch} from 'react-redux';
import { useEffect,useMemo,useRef,useState } from 'react';
import {Button} from "antd";
import { PlusSquareOutlined,MinusSquareOutlined,FullscreenOutlined  } from '@ant-design/icons';

import { setCurrentPos } from '../../../../../redux/dataSlice';
import Indicator from './Indicator';
import './index.css';



export default function Content({map}){
    const dispatch=useDispatch();
    const indicator=useSelector(state=>state.data.indicator);
    const points=useSelector(state=>state.data.points);
    const currentPos=useSelector(state=>state.data.currentPos);
    const [imageWidth,setImageWidth]=useState(0);

    const refWapper=useRef(null);
    const refImg=useRef(null);
    
    const pointsControl=useMemo(()=>{
        const setCurrentPosTo=(pos)=>{
            dispatch(setCurrentPos(pos));
        }

        const convertXY=(x,y)=>{
            let xRatio=1;
            let yRatio=1;
            if(refImg.current){
                const orgWidth=refImg.current.naturalWidth;
                xRatio=imageWidth/orgWidth;
    
                //const imageHeight=refImg.current.height;
                //const orgHeight=refImg.current.naturalHeight;
                yRatio=xRatio;//imageHeight/orgHeight;
            }
            return {x:x*xRatio,y:y*yRatio};
        }

        return points.map((dataItem,index)=>{
            if(dataItem){
                const {x:xO,y:yO,rgb,value}=dataItem;
                const {x,y}=convertXY(xO,yO);
                const isCurPoint=(index+1===currentPos)?true:false;
                if(isCurPoint===true){
                    return (
                        <>    
                            <div key={'label_'+index} className='map-point-label' style={{left:x,top:y-25}}>{value}</div>
                            <div key={'point_'+index} onClick={()=>setCurrentPosTo(index+1)} className='map-point' style={{left:x,top:y,backgroundColor:rgb}}></div>
                        </>
                    );
                }
    
                return (
                    <div key={'point_'+index} onClick={()=>setCurrentPosTo(index+1)} className='map-point' style={{left:x,top:y,backgroundColor:rgb}}></div>
                );
    
            } else {
                return null;
            }
        });
    },[imageWidth,currentPos,points,dispatch]); 

    useEffect(()=>{
        if(refWapper.current){
            console.log("setImageWidth",refWapper.current.clientWidth);
            setImageWidth(refWapper.current.clientWidth);
        }
    },[refWapper]);

    const resizeImage=(value)=>{
        if(value===0){
            if(refWapper.current){
                setImageWidth(refWapper.current.clientWidth);
            }
        } else {
            setImageWidth(imageWidth+value);
        }
    }

    return (
        <div ref={refWapper} className="monitor-map-content">
            <div className='image-wrapper'>
                <img ref={refImg}   style={{height:'auto',width:imageWidth}} src={map?.url} alt='' />
                {pointsControl}
            </div>
            <Indicator indicator={indicator}/>
            <Button onClick={()=>resizeImage(5)} size='small'  style={{float:'right',position:'absolute',top:5,right:5}} icon={<PlusSquareOutlined/>} />
            <Button onClick={()=>resizeImage(0)} size='small' style={{position:'absolute',top:30,right:5}} icon={<FullscreenOutlined />} />
            <Button onClick={()=>resizeImage(-5)} size='small' style={{position:'absolute',top:55,right:5}} icon={<MinusSquareOutlined />} />
        </div>
    );
}