import { useEffect } from "react";

import { createRedirectMessage } from "../../../../utils/normalOperations";
import { DATA_TYPE } from "../../../../utils/constant";
import DialInfo from "./DialInfo";

export default function Dial({sendMessageToParent,frame}){
    const dialFunc=(op)=>{
        const frameParams={
            frameType:frame.item.frameType,
            frameID:frame.item.params.key,
            origin:frame.origin
        };
        sendMessageToParent(createRedirectMessage(frameParams,"rttkservice/device/"+op,DATA_TYPE.DIAL));
    }

    useEffect(()=>{  
        dialFunc('dialQuery');
    });

    console.log("Dial ...");

    return (<DialInfo dialFunc={dialFunc} />);
}