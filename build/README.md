# rtservice

#更新镜像流程
1、停止并删除原有镜像实例
    docker stop rtservice
    docker rm rtservice
2、删除原有镜像
    docker rmi rtservice:0.1.0
3、加载新的镜像
    docker load -i rtserviceXXXXX.tar
4、启动实例
    docker run 。。。

#导出镜像包命令
docker save -o rtservice.tar wangzhsh/rtservice:0.1.0
#导入镜像包命令
docker load -i rtservice.tar

#run rtservice in docker
docker run -d --name rtservice -p8300:80 -v /root/rtservice/localcache:/services/rtservice/localcache -v /root/rtservice/conf:/services/rtservice/conf wangzhsh/rtservice:0.1.0

修改记录
20230919
1、monitor：修正BUG，调整页面布局。
2、play:修正BUG，调整页面布局。
3、补充对指标图例的配置页面，增加以下表和对应配置
    rt_indicator
    rt_indicator_legend
4、增加对应菜单项。

20230920
2、修改play页面，补充了对指标的选择和图例的显示功能。