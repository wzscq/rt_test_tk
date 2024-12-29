import { useEffect } from "react";

import { createRedirectMessage } from "../../../../utils/normalOperations";
import { DATA_TYPE } from "../../../../utils/constant";
import RATInfo from "./RATInfo";

export default function RAT({sendMessageToParent,frame}){
    const getRAT=()=>{
        const frameParams={
            frameType:frame.item.frameType,
            frameID:frame.item.params.key,
            origin:frame.origin
        };
        sendMessageToParent(createRedirectMessage(frameParams,"rttkservice/device/getRAT",DATA_TYPE.RAT));
    }

    const setRAT=()=>{
        const frameParams={
            frameType:frame.item.frameType,
            frameID:frame.item.params.key,
            origin:frame.origin
        };
        sendMessageToParent(createRedirectMessage(frameParams,"rttkservice/device/setRAT",DATA_TYPE.RAT));
    }

    useEffect(()=>{  
        getRAT();
    });

    console.log("RAT ...");

    return (<RATInfo getRAT={getRAT} setRAT={setRAT}/>);
}