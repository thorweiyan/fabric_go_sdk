# fabric_go_sdk
先确定fabric的bin工具位置，并使用其在fixtures下生成对应的artifacts和crypto-config文件夹。

dep ensure可以自动下载相对应的库文件到vendor下，但国内使用比较困难，可以使用go get或者git clone下载缺少的包放到vendor下。

test中描述了sdk的调用方法：
  Initialize初始化生成fabric的环境
  InstallAndInstantiateCC安装并实例化CC
  Invoke负责调用CC并生成Event，返回交易id
  Query负责查询CC并返回对应的payload
  UpdateCC尚未完成

chaincode中保存cc代码，所有的路径和参数都在test中设置
