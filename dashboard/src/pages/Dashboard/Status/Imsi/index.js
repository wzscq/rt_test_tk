import { useEffect } from "react";

import { createRedirectMessage } from "../../../../utils/normalOperations";
import { DATA_TYPE } from "../../../../utils/constant";
import ImsiInfo from "./ImsiInfo";

export default function Imsi({sendMessageToParent,frame}){
    const getImsi=()=>{
        const frameParams={
            frameType:frame.item.frameType,
            frameID:frame.item.params.key,
            origin:frame.origin
        };
        sendMessageToParent(createRedirectMessage(frameParams,"rttkservice/device/getImsi",DATA_TYPE.IMSI));
    }

    useEffect(()=>{  
        getImsi();
    });

    console.log("Imsi ...");

    return (<ImsiInfo getImsi={getImsi}/>);
}