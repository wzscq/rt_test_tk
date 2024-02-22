import './index.css';

export default function Header({mapInfo}){
    return (
        <div className="monitor-map-header">
            {'building:'+mapInfo.building_code+' floor:'+mapInfo.floor+' pic:'+mapInfo.picture_name}
        </div>
    );
}