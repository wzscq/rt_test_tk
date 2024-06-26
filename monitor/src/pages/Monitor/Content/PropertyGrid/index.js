import { useRef, useState } from "react";

import './index.css';

var isMouseDown=false;
var mouseDownLeft=150;

export default function PropertyGrid({obj,title}){
  const [splitLeft,setSplitLeft]=useState(150);
  const refSplitBar=useRef();

  const onSplitMouseDown=(e)=>{
      //开始拉动，记录鼠标的起始位置
      isMouseDown=true;
      mouseDownLeft=e.clientX
      console.log(mouseDownLeft)
  }

  const onSplitonMouseMove=(e)=>{
      //鼠标移动过程中，改变分割条的位置
      if(refSplitBar.current&&isMouseDown===true){
          console.log(e)
          const diff=e.clientX-mouseDownLeft;
          refSplitBar.current.style.left=splitLeft+diff+'px';
      }
  }

  const onSplitMouseUp=(e)=>{
      //释放鼠标后移动完成，更新控件的位置
      if(refSplitBar.current&&isMouseDown===true){
          console.log(e)
          isMouseDown=false;
          const diff=e.clientX-mouseDownLeft;
          refSplitBar.current.style.left=splitLeft+diff+'px';
          setSplitLeft(splitLeft+diff);
      }
  }

  const rows=[];

  const isObject = obj => {
    return typeof obj === 'object' && obj !== null && !Array.isArray(obj)
  }

  Object.keys(obj).forEach(key=>{
    console.log('obj keys:',key);
    let value=obj[key];
    if(isObject(value)){
      value=JSON.stringify(value); 
    }

    if(value===null){
      value='';
    }

    if(value===undefined){
      value='';
    } 

    value=value.toString();

    rows.push((<div className='row' key={key}>
      <div className='col name' style={{width:(splitLeft-5)}}>{key}</div>
      <div className='col value' style={{width:'calc(100% - '+(splitLeft+6)+'px)'}}>{value}</div>
    </div>));
  });

  return (<div className="property-grid"
          onMouseMove={onSplitonMouseMove}
          onMouseUp={onSplitMouseUp}>
    <div className='property-grid-title'>{title}</div>
    <div className='col-title'>
        <div className='name' style={{width:splitLeft}}>name</div>
        <div className='value' style={{width:'calc(100% - '+(splitLeft+1)+'px)'}}>value</div>
    </div>
    <div className='property-grid-content'>
      {rows.length>0?rows:null}
      <div 
        ref={refSplitBar}
        className="property-grid-split"  
        style={{left:splitLeft}}
        onMouseDown={onSplitMouseDown}
      />
    </div>
  </div>);
}