import AceEditor from "react-ace";
import "ace-builds/src-noconflict/mode-json"
import "ace-builds/src-noconflict/snippets/json"
import "ace-builds/src-min-noconflict/ext-searchbox";
import "ace-builds/src-min-noconflict/ext-language_tools";
import "ace-builds/src-noconflict/theme-github"

import './index.css';

export default function JSONViewer({json,title}){
  console.log('json:',json);

  //const jsonStr='[{"msg":"Packet 1/5: 72 bytes from 192.168.100.109 (192.168.100.109): icmp_seq=1 time=1.50 ms"},{"msg":"Packet 2/5: 72 bytes from 192.168.100.109 (192.168.100.109): icmp_seq=2 time=2.00 ms"},{"msg":"Packet 3/5: 72 bytes from 192.168.100.109 (192.168.100.109): icmp_seq=3 time=2.33 ms"},{"msg":"Packet 4/5: 72 bytes from 192.168.100.109 (192.168.100.109): icmp_seq=4 time=2.15 ms"},{"msg":"Packet 5/5: 72 bytes from 192.168.100.109 (192.168.100.109): icmp_seq=5 time=2.07 ms"}]'

  return (
    <div className="json-viewer">
        <div className="json-viewer-title">{title}</div>
        <div className="json-viewer-content">
        <AceEditor
            style={{height:"100%",width:"100%",overflow:"auto"}}
            placeholder="Placeholder Text"
            mode="json"
            theme="github"
            name="funcScript"
            fontSize={12}
            showPrintMargin={true}
            showGutter={true}
            highlightActiveLine={true}
            value={json}
            setOptions={{
            enableBasicAutocompletion: true,
            enableLiveAutocompletion: true,
            enableSnippets: false,
            showLineNumbers: true,
            tabSize: 2,
        }}/>
        </div>
    </div>
  );
}