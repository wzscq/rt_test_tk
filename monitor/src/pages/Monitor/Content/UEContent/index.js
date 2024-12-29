import { useSelector } from "react-redux";
import { SplitPane } from "react-collapse-pane";
import PropertyGrid from '../PropertyGrid';
import JSONViewer from "../JSONViewer";

export default function UEContent(){
    const {data,commandResult,pingRec}=useSelector(state=>state.data);

    const useCase={...data};
    delete useCase.testData;

    const deviceInfo={...data.testData};
    delete deviceInfo.measures;
    delete deviceInfo.throughput;

    const gps=deviceInfo?.gps;
    deviceInfo.longitude=gps?.longitude??'';
    deviceInfo.latitude=gps?.latitude??'';
    delete deviceInfo.gps;

    const measures={...data?.testData?.measures};
    const throughput={...data?.testData?.throughput};
    let result=JSON.stringify(deviceInfo?.result??'');
    try {
        result=JSON.stringify(JSON.parse(result),undefined,4);
    } catch (error) {
        result=JSON.stringify(deviceInfo?.result??'');
    }

    delete deviceInfo.result;

    return (
        <SplitPane dir='ltr' initialSizes={[50,50]} split="vertical" collapse={false}>
            <SplitPane split="horizontal" initialSizes={[20,20,20,40]} collapse={false}>
                <PropertyGrid obj={commandResult??{}} title="测试执行结果"/>
                <PropertyGrid obj={useCase} title="测试用例"/>
                <PropertyGrid obj={deviceInfo} title="设备信息"/>
                <JSONViewer json={result} title="结果"/>
            </SplitPane>
            <SplitPane split="horizontal" initialSizes={[20,40,40]} collapse={false}>
                <PropertyGrid obj={throughput??{}} title="速率"/>
                <PropertyGrid obj={measures??{}} title="测量指标"/>
                <JSONViewer json={JSON.stringify(pingRec??'',undefined,4)+''} title="Ping/Attach记录"/>
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