export default function Indicator({indicator}){
  const legend=indicator?.legend?.list;

  const legendItems=legend?.map(item=>{
    let startVal=item.start;
    if(startVal===null){
      startVal='   ';
    }
    let endVal=item.end;
    if(endVal===null){
      endVal='   ';
    }
    
    const legendTitle='('+startVal+' , '+endVal+']';
    return (<div className="legend-item" key={item.id}>
      <div className="legend-rgb" style={{backgroundColor:item.rgb}}></div>
      <div className="legend-title">{legendTitle}</div>
    </div>);
  });

  return (
    <div className="indicator">
      {legendItems}
    </div>
  );
}