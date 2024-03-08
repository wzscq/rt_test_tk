import React, { useEffect, useRef } from 'react';
import {useSelector,useDispatch} from 'react-redux';
import Map from 'ol/Map'
import View from 'ol/View'
import TileLayer from 'ol/layer/Tile'
import VectorLayer from 'ol/layer/Vector';
import VectorSource from 'ol/source/Vector';
import XYZ from 'ol/source/XYZ'
import {Point} from 'ol/geom.js';
import Feature from 'ol/Feature.js';
import {Circle, Fill, Stroke, Style} from 'ol/style.js';
import { transform } from "ol/proj";

import './index.css';

const g_MQTTClient=null;

export default function MapWrapper(){
  const dispatch=useDispatch();
  const {mapConf,mqttConf}=useSelector(state=>state.mqtt);
  const mapElement = useRef(null);
  const tipEle = useRef(null);
  const map = useRef(null); //地图全局变量
  

  useEffect(()=>{
    const connectMqtt=()=>{
      console.log("connectMqtt ... ");
      if(g_MQTTClient!==null){
          g_MQTTClient.end();
          g_MQTTClient=null;
      }
  
      const server='ws://'+mqttConf.broker+':'+mqttConf.wsPort;
      const options={
          username:mqttConf.user,
          password:mqttConf.password,
          keepalive:3600,
          reconnectPeriod:60
      }
      console.log("connect to mqtt server ... "+server+" with options:",options);
      g_MQTTClient  = mqtt.connect(server,options);
      g_MQTTClient.on('connect', () => {
          setMqttStatus("connected to mqtt server "+server+".");
          const topic=mqttConf.uploadMeasurementMetrics;
          g_MQTTClient.subscribe(topic, (err) => {
              if(!err){
                  setMqttStatus("subscribe topics success.");
                  console.log("topic:",topic);
              } else {
                  setMqttStatus("subscribe topics error :"+err.toString());
              }
          });
      });
      g_MQTTClient.on('message', (topic, payload, packet) => {
          console.log("receiconsolleve message topic :"+topic+" content :"+payload.toString());
          dispatch(addDataItem(JSON.parse(payload.toString())));
      });
      g_MQTTClient.on('close', () => {
        setMqttStatus("mqtt client is closed.");
      });
    }

    connectMqtt();
  },[dispatch,mqttConf]);

  useEffect(()=>{
    if(map.current!==null) return; //防止重复渲染地图
    const center=transform(mapConf.center,'EPSG:4326','EPSG:3857')
    map.current = new Map({
        view: new View({
            center: center,//地图中心位置
            zoom: mapConf.zoom,//地图初始层级
            maxZoom: mapConf.maxZoom,
            minZoom: mapConf.minZoom
        }),
        target: mapElement.current
    });
    let tileLayer = new TileLayer({
        source: new XYZ({
          tileUrlFunction:(coordinate)=>{
            const z = coordinate[0];
            const x = coordinate[1];
            const y = coordinate[2];
            const file= mapConf.url.replace('{z}',z).replace('{x}',x).replace('{y}',y);
            console.log("file",file);
            return file;
          }
        })
    });
    map.current.addLayer(tileLayer)

    const circleFeature = new Feature({
      geometry: new Point(transform(mapConf.center,'EPSG:4326','EPSG:3857')),
    });

    const fill = new Fill({
      color: 'rgba(255,0,0,1)',
    });
    const stroke = new Stroke({
      color: '#3399CC',
      width: 1.25,
    });

    const iconStyle=new Style({
      image: new Circle({
        fill: fill,
        stroke: stroke,
        radius: 5,
      }),
      fill: fill,
      stroke: stroke,
    })
    
    circleFeature.setStyle(iconStyle);
    circleFeature.set('TIP_TEXT', 'This is the center of the map');

    const ectorLayer=new VectorLayer({
      source: new VectorSource({
        features: [circleFeature],
      }),
    });

    map.current.addLayer(ectorLayer);

    let currentFeature;
    const displayFeatureInfo = function (pixel, target) {
      const feature = target.closest('.ol-control')
        ? undefined
        : map.current.forEachFeatureAtPixel(pixel, function (feature) {
            return feature;
          });
      if (feature) {
        tipEle.current.style.left = pixel[0] + 'px';
        tipEle.current.style.top = pixel[1] + 'px';
        if (feature !== currentFeature) {
          tipEle.current.style.visibility = 'visible';
          tipEle.current.innerText = feature.get('TIP_TEXT');
        }
      } else {
        tipEle.current.style.visibility = 'hidden';
      }
      currentFeature = feature;
    };

    map.current.on('pointermove', function (evt) {
      if (evt.dragging) {
        tipEle.current.style.visibility = 'hidden';
        currentFeature = undefined;
        return;
      }
      const pixel = map.current.getEventPixel(evt.originalEvent);
      displayFeatureInfo(pixel, evt.originalEvent.target);
    });

    map.current.on('click', function (evt) {
      displayFeatureInfo(evt.pixel, evt.originalEvent.target);
    });

    map.current.getTargetElement().addEventListener('pointerleave', function () {
      currentFeature = undefined;
      tipEle.current.style.visibility = 'hidden';
    });

  });

  return (
    <div ref={mapElement} className="map-container">
      <div ref={tipEle} id="info"></div>
    </div>
  );
}