import { Button,Popconfirm } from "antd";

import { createRedirectMessage } from "../../../../utils/normalOperations";


export default function DeviceControl({sendMessageToParent,frame}){
  const rebootDevice=()=>{
    const frameParams={
      frameType:frame.item.frameType,
      frameID:frame.item.params.key,
      origin:frame.origin
    };
    sendMessageToParent(createRedirectMessage(frameParams,"rttkservice/device/reboot","no result"));
  }

  return (
    <div>
      <Popconfirm
        title="重启设备"
        description="请确认是否要重启设备?"
        onConfirm={rebootDevice}
        okText="确定"
        cancelText="取消"
      >
        <Button type='primary'>重启设备</Button>
      </Popconfirm>
    </div>
  )
}