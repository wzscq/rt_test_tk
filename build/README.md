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

20240505
1、下发测试用例，增加了检查，如果已经存在测试用例，则不允许下发。
2、下发解码请求，增加了检查，如果已经存在未完成的解码任务，则不允许下发解码任务。
3、增加配置项目。
4、补充定时扫描程序，第一次启动时删除所有数据库数据，然后扫描磁盘文件添加到数据库，后续定时扫描，基于差异文件更新数据库。

20240512
1、补充测试停止功能，修改后台服务逻辑。
2、补充attach、detatch测试

20240519
1、增加一个弹出框选择测试参数，目前仅添加单次测试何循环测试参数
    a、修改模型配置：rt_project_test_case
    b、修改后台逻辑

20240520
1、修改logfile列表获取方式，改为通过解码器接口获取文件列表
   a、修改后台逻辑，增加新的接口从解码器获取文件列表
   b、修拍配置文件，增加到解码器文件列表获取url配置

   