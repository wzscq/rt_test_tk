import { useSelector } from "react-redux";
import { SplitPane } from "react-collapse-pane";
import PropertyGrid from '../PropertyGrid';
import ListGrid from '../ListGrid';

export default function UEContent({imsi}){
    const dataItem=useSelector(state=>state.data.currentUes[imsi]);
    const event=useSelector(state=>state.data.event[imsi]);
    const message=useSelector(state=>state.data.message[imsi]);

    console.log('UEContent',imsi,dataItem,event,message);

    return (
        <SplitPane key={imsi} dir='ltr' initialSizes={[50,50]} split="vertical" collapse={false}>
            <SplitPane split="horizontal" collapse={false}>
                <PropertyGrid obj={dataItem?.radio?.measures_lte} title="lte measures"/>
                <PropertyGrid obj={dataItem?.radio?.measures_nr} title="nr measures"/>
            </SplitPane>
            <SplitPane split="horizontal" collapse={false}>
                <ListGrid list={event} title={'events'} name={'event'} nameField={'EventTime'} timeField={'name'} />
                <ListGrid list={message} title={'messages'} name={'message'} nameField={'MsgTime'} timeField={'name'}/>
            </SplitPane>
        </SplitPane>
    );

    /*return (
        <SplitPane key={imsi} dir='ltr' initialSizes={[30,30,40]} split="vertical" collapse={false}>
            <SplitPane split="horizontal" collapse={false}>
                <PropertyGrid obj={dataItem?.radio?.measures_lte} title="lte measures"/>
                <PropertyGrid obj={dataItem?.radio?.measures_nr} title="nr measures"/>
            </SplitPane>
            <SplitPane split="horizontal" collapse={false}>
                <ListGrid list={event} title={'events'} name={'event'} nameField={'name'} timeField={'EventTime'} />
                <ListGrid list={message} title={'messages'} name={'message'} nameField={'name'} timeField={'MsgTime'}/>
            </SplitPane>
            <SplitPane split="horizontal" collapse={false}>
                <PropertyGrid obj={dataItem?.radio?.measures_common} title="common measures"/>
                <PropertyGrid obj={dataItem?.case_progress} title="case progress"/>
            </SplitPane>
        </SplitPane>
    );*/
}