import AttachInfo from "./AttachInfo";
import {useEffect} from 'react';
import { createRedirectMessage } from "../../../../utils/normalOperations";
import { DATA_TYPE } from "../../../../utils/constant";

export default function Attach({sendMessageToParent,frame}){
    const attach=()=>{
        const frameParams={
            frameType:frame.item.frameType,
            frameID:frame.item.params.key,
            origin:frame.origin
        };
        sendMessageToParent(createRedirectMessage(frameParams,"rttkservice/device/attach",DATA_TYPE.ATTACH));
    }

    const detach=()=>{
        const frameParams={
            frameType:frame.item.frameType,
            frameID:frame.item.params.key,
            origin:frame.origin
        };
        sendMessageToParent(createRedirectMessage(frameParams,"rttkservice/device/detach",DATA_TYPE.DEATCH));
    }

    const attach_query=()=>{
        const frameParams={
            frameType:frame.item.frameType,
            frameID:frame.item.params.key,
            origin:frame.origin
        };
        sendMessageToParent(createRedirectMessage(frameParams,"rttkservice/device/attach_query",DATA_TYPE.ATTACH_QUERY));
    }

    useEffect(()=>{  
        attach_query();
    });


    return (
        <AttachInfo attach={attach} detach={detach} attach_query={attach_query}/>
    );
}