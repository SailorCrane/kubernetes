sequenceDiagram

# 从kubelet cmd内部直接开始cmd/kubelet/app/server.go/ NewKubeletCommand()

server ->> option : kubeletFlags := options.NewKubeletFlags()
    Note right of server : kubeletFlags, kubelet非dynamic config中的配置
    Note right of server : kubeletFlags, 不与其它kubelet共享的配置
    Note right of server : 都是一些本地的配置, 比如cert证书的地址, docker shim的socket所在等等
    Note right of server : 不同机器可能这些配置略有差异
option ->> server : return KubeletFlags

server ->> option : kubeletConfig := options.NewKubeletConfiguration()
    Note right of server : 不同机器间共享的配置(这些配置放在配置文件中)
option ->> server : KubeletConfiguration

server ->> pflag : cleanFlagSet := pflag.NewFlagSet(componentKubelet, pflag.ContinueOnError)
    Note right of pflag : 创建pflag.FlagSet, 用于Parse命令行args中的值
pflag ->> server : return FlagSet

server ->> option : kubeletFlags.AddFlags(cleanFlagSet)
    Note right of option : 使cleanFlagSet Parse()命令行后的一些值, 存入kubeletFlags中

server ->> option : options.AddKubeletConfigFlags(cleanFlagSet, kubeletConfig)
    Note right of option : 使cleanFlagSet Parse()命令行后的一些值, 存入kubeletConfig中



Note right of nothing1 : 只是为了mermaid显示完整所加, 不是源代码内容
Note right of nothing2 : 只是为了mermaid显示完整所加, 不是源代码内容
Note right of nothing3 : 只是为了mermaid显示完整所加, 不是源代码内容
