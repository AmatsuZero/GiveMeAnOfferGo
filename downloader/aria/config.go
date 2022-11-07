package aria

import (
	"GiveMeAnOffer/downloader"
	"GiveMeAnOffer/logger"
	"GiveMeAnOffer/parse"
	"GiveMeAnOffer/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

type FileAllocType string

const (
	FileAllocTypeNone     FileAllocType = "none"
	FileAllocTypePreAlloc FileAllocType = "prealloc"
	FileAllocTypeTrunc    FileAllocType = "trunc"
	FileAllocTypeFalloc   FileAllocType = "falloc"
)

type TLSVersion string

const (
	TLSV1_1  TLSVersion = "TLSv1.1"
	TLVSV1_2 TLSVersion = "TLSv1.2"
	TLVSV1_3 TLSVersion = "TLSv1.3"
)

type FollowTorrentStrategy string

const (
	SaveOnDisk FollowTorrentStrategy = "true"
	SaveOnMem  FollowTorrentStrategy = "mem"
	DontSave   FollowTorrentStrategy = "false"
)

type CryptoLevel string

const (
	CryptoLevelNone CryptoLevel = "plain"
	CryptoLevelArc4 CryptoLevel = "arc4"
)

type EventPollType string

const (
	EventPollTypeEPoll  EventPollType = "epoll"
	EventPollTypeKqueue EventPollType = "kqueue"
	EventPollTypePort   EventPollType = "port"
	EventPollTypePoll   EventPollType = "poll"
	EventPollTypeSelect EventPollType = "select"
)

type LogLevel string

const (
	LogLevelDebug  LogLevel = "debug"
	LogLevelInfo   LogLevel = "info"
	LogLevelNotice LogLevel = "notice"
	LogLevelWarn   LogLevel = "warn"
	LogLevelError  LogLevel = "error"
)

type Config struct {
	Trackers []string
	// 下载目录。可使用绝对路径或相对路径, 默认: 当前启动位置
	SaveDir string `json:"dir"`
	/*
		磁盘缓存, 0 为禁用缓存，默认:16M
		磁盘缓存的作用是把下载的数据块临时存储在内存中，然后集中写入硬盘，以减少磁盘 I/O ，提升读写性能，延长硬盘寿命。
		建议在有足够的内存空闲情况下适当增加，但不要超过剩余可用内存空间大小。
		此项值仅决定上限，实际对内存的占用取决于网速(带宽)和设备性能等其它因素。
	*/
	DiskCache int64
	/*
		文件预分配方式, 可选：none, prealloc, trunc, falloc, 默认:prealloc
		预分配对于机械硬盘可有效降低磁盘碎片、提升磁盘读写性能、延长磁盘寿命。
		机械硬盘使用 ext4（具有扩展支持），btrfs，xfs 或 NTFS（仅 MinGW 编译版本）等文件系统建议设置为 falloc
		若无法下载，提示 fallocate failed.cause：Operation not supported 则说明不支持，请设置为 none
		prealloc 分配速度慢, trunc 无实际作用，不推荐使用。
		固态硬盘不需要预分配，只建议设置为 none ，否则可能会导致双倍文件大小的数据写入，从而影响寿命。
	*/
	FileAllocation FileAllocType
	// 文件预分配大小限制。小于此选项值大小的文件不预分配空间，单位 K 或 M，默认：5M
	NoFileAllocationLimit int64
	// 断点续传
	Continue bool
	// 始终尝试断点续传，无法断点续传则终止下载，默认：true
	AlwaysResume bool
	/*
		不支持断点续传的 URI 数值，当 always-resume=false 时生效。
		达到这个数值从将头开始下载，值为 0 时所有 URI 不支持断点续传时才从头开始下载。
	*/
	MaxResumeFailureTries int
	// 获取服务器文件时间，默认:false
	RemoteTime bool
	// 从会话文件中读取下载任务
	InputFile string
	/*
		会话文件保存路径
		Aria2 退出时或指定的时间间隔会保存`错误/未完成`的下载任务到会话文件
	*/
	SaveSession string
	/*
		# 任务状态改变后保存会话的间隔时间（秒）, 0 为仅在进程正常退出时保存, 默认:0
		# 为了及时保存任务状态、防止任务丢失，此项值只建议设置为 1
	*/
	SaveSessionInterval int
	/*
		自动保存任务进度到控制文件(*.aria2)的间隔时间（秒），0 为仅在进程正常退出时保存，默认：60
		此项值也会间接影响从内存中把缓存的数据写入磁盘的频率
		想降低磁盘 IOPS (每秒读写次数)则提高间隔时间
		想在意外非正常退出时尽量保存更多的下载进度则降低间隔时间
		非正常退出：进程崩溃、系统崩溃、SIGKILL 信号、设备断电等
	*/
	AutoSaveInterval int
	/*
		强制保存，即使任务已完成也保存信息到会话文件, 默认:false
		开启后会在任务完成后保留 .aria2 文件，文件被移除且任务存在的情况下重启后会重新下载。
		关闭后已完成的任务列表会在重启后清空。
	*/
	ForceSave bool
	/*
		文件未找到重试次数，默认:0 (禁用)
		重试时同时会记录重试次数，所以也需要设置 max-tries 这个选项
	*/
	MaxFileNotFound int
	// 最大尝试次数，0 表示无限，默认:5
	MaxTries int
	// 重试等待时间（秒）, 默认:0 (禁用)
	RetryWait int
	// 连接超时时间（秒）。默认：60
	ConnectionTimeout int
	// 超时时间（秒）。默认：60
	Timeout int
	// 最大同时下载任务数, 运行时可修改, 默认:5
	MaxConCurrentDownload int
	/*
		单服务器最大连接线程数, 任务添加时可指定, 默认:1
		最大值为 16 (增强版无限制), 且受限于单任务最大连接线程数(split)所设定的值。
	*/
	MaxConnectionPerServe int
	// 单任务最大连接线程数, 任务添加时可指定, 默认:5
	Split int
	/*
		文件最小分段大小, 添加时可指定, 取值范围 1M-1024M (增强版最小值为 1K), 默认:20M
		比如此项值为 10M, 当文件为 20MB 会分成两段并使用两个来源下载, 文件为 15MB 则只使用一个来源下载。
		理论上值越小使用下载分段就越多，所能获得的实际线程数就越大，下载速度就越快，但受限于所下载文件服务器的策略。
	*/
	MinSplitSize int64
	// HTTP/FTP 下载分片大小，所有分割都必须是此项值的倍数，最小值为 1M (增强版为 1K)，默认：1M
	PieceLength int64
	/*
		允许分片大小变化。默认：false
		false：当分片大小与控制文件中的不同时将会中止下载
		true：丢失部分下载进度继续下载
	*/
	AllowPieceLengthChange bool
	// 最低下载速度限制。当下载速度低于或等于此选项的值时关闭连接（增强版本为重连），此选项与 BT 下载无关。单位 K 或 M ，默认：0 (无限制)
	LowestSpeedLimit int
	// 全局最大下载速度限制, 运行时可修改, 默认：0 (无限制)
	MaxOverallDownloadLimit int
	// 单任务下载速度限制, 默认：0 (无限制)
	MaxDownloadLimit int
	// 禁用 IPv6, 默认:false
	DisableIpv6 bool
	// GZip 支持，默认:false
	HttpAcceptGzip bool
	// URI 复用，默认: true
	ReuseUri bool
	// 禁用 netrc 支持，默认:false
	NoNetrc bool
	// 允许覆盖，当相关控制文件(.aria2)不存在时从头开始重新下载。默认:false
	AllowOverwrite bool
	// 文件自动重命名，此选项仅在 HTTP(S)/FTP 下载中有效。新文件名在名称之后扩展名之前加上一个点和一个数字（1..9999）。默认:true
	AutoFileRenaming bool
	// 使用 UTF-8 处理 Content-Disposition ，默认:false
	ContentDispositionDefaultUtf8 bool `json:"content-disposition-default-utf8"`
	// 最低 TLS 版本，可选：TLSv1.1、TLSv1.2、TLSv1.3 默认:TLSv1.2
	MinTLSVersion TLSVersion
	/*
		BT 监听端口(TCP), 默认:6881-6999
		直通外网的设备，比如 VPS ，务必配置防火墙和安全组策略允许此端口入站
		内网环境的设备，比如 NAS ，除了防火墙设置，还需在路由器设置外网端口转发到此端口
	*/
	ListenPort int
	/*
		DHT 网络与 UDP tracker 监听端口(UDP), 默认:6881-6999
		因协议不同，可以与 BT 监听端口使用相同的端口，方便配置防火墙和端口转发策略。
	*/
	DHTListenPort int `conf:"dht-listen-port"`
	// 启用 IPv4 DHT 功能, PT 下载(私有种子)会自动禁用, 默认:true
	EnableDHT bool
	/*
		启用 IPv6 DHT 功能, PT 下载(私有种子)会自动禁用，默认:false
		在没有 IPv6 支持的环境开启可能会导致 DHT 功能异常
	*/
	EnableDHT6 bool `conf:"enable-dht-6"`
	/*
		指定 BT 和 DHT 网络中的 IP 地址
		使用场景：在家庭宽带没有公网 IP 的情况下可以把 BT 和 DHT 监听端口转发至具有公网 IP 的服务器，在此填写服务器的 IP ，可以提升 BT 下载速率。
	*/
	BTExternalIP string `conf:"bt-external-ip"`
	// IPv4 DHT 文件路径，默认：$HOME/.aria2/dht.dat
	DHTFilePath string
	// IPv6 DHT 文件路径，默认：$HOME/.aria2/dht6.dat
	DHT6FilePath string `conf:"dht-6-file-path"`
	// IPv4 DHT 网络引导节点
	DHTEntryPoint string `conf:"dht-entry-point"`
	// IPv6 DHT 网络引导节点
	DHTEntryPoint6 string `json:"dht-entry-point6"`
	// 本地节点发现, PT 下载(私有种子)会自动禁用 默认:false
	BTEnableLPD bool `json:"bt-enable-lpd"`
	/*
		指定用于本地节点发现的接口，可能的值：接口，IP地址
		如果未指定此选项，则选择默认接口。
	*/
	BTLPDInterface string `conf:"bt-lpd-interface"`
	// 启用节点交换, PT 下载(私有种子)会自动禁用, 默认:true
	EnablePeerExchange bool `conf:"enable-peer-exchange"`
	/*
		BT 下载最大连接数（单任务），运行时可修改。0 为不限制，默认:55
		理想情况下连接数越多下载越快，但在实际情况是只有少部分连接到的做种者上传速度快，其余的上传慢或者不上传。
		如果不限制，当下载非常热门的种子或任务数非常多时可能会因连接数过多导致进程崩溃或网络阻塞。
		进程崩溃：如果设备 CPU 性能一般，连接数过多导致 CPU 占用过高，因资源不足 Aria2 进程会强制被终结。
		网络阻塞：在内网环境下，即使下载没有占满带宽也会导致其它设备无法正常上网。因远古低性能路由器的转发性能瓶颈导致
	*/
	BTMaxPeers int `conf:"bt-max-peers"`
	/*
		BT 下载期望速度值（单任务），运行时可修改。单位 K 或 M 。默认:50K
		BT 下载速度低于此选项值时会临时提高连接数来获得更快的下载速度，不过前提是有更多的做种者可供连接。
		实测临时提高连接数没有上限，但不会像不做限制一样无限增加，会根据算法进行合理的动态调节。
	*/
	BTRequestPeerSpeedLimit int64 `conf:"bt-request-peer-speed-limit"`
	/*
		全局最大上传速度限制, 运行时可修改, 默认:0 (无限制)
		设置过低可能影响 BT 下载速度
	*/
	MaxOverallUploadLimit int64 `conf:"max-overall-upload-limit"`
	// 单任务上传速度限制, 默认:0 (无限制)
	MaxUploadLimit int64 `conf:"max-upload-limit"`
	/*
		最小分享率。当种子的分享率达到此选项设置的值时停止做种, 0 为一直做种, 默认:1.0
		强烈建议您将此选项设置为大于等于 1.0
	*/
	SeedRatio float64 `conf:"seed-ratio"`
	// 最小做种时间（分钟）。设置为 0 时将在 BT 任务下载完成后停止做种。
	SeedTime int `conf:"seed-time"`
	// 做种前检查文件哈希, 默认:true
	BTHashCheckSeed bool `conf:"bt-hash-check-seed"`
	// 继续之前的BT任务时, 无需再次校验, 默认:false
	BTSeedUnverified bool `conf:"bt-seed-unverified"`
	/*
		BT tracker 服务器连接超时时间（秒）。默认：60
		建立连接后，此选项无效，将使用 bt-tracker-timeout 选项的值
	*/
	BTTrackerConnectTimeout int `json:"bt-tracker-connect-timeout"`
	// BT tracker 服务器超时时间（秒）。默认：60
	BTTrackerTimeout int `conf:"bt-tracker-timeout"`
	// BT 服务器连接间隔时间（秒）。默认：0 (自动)
	BTTrackerInterval int `conf:"bt-tracker-interval"`
	// BT 下载优先下载文件开头或结尾
	BTPrioritizePiece struct {
		Head, Tail int64
	} `conf:"bt-prioritize-piece"`
	/*
		保存通过 WebUI(RPC) 上传的种子文件(.torrent)，默认:true
		所有涉及种子文件保存的选项都建议开启，不保存种子文件有任务丢失的风险。
		通过 RPC 自定义临时下载目录可能不会保存种子文件。
	*/
	RPCSaveUploadMetadata bool `conf:"rpc-save-upload-metadata"`
	/*
		下载种子文件(.torrent)自动开始下载, 默认:true，可选：false|mem
		true：保存种子文件
		false：仅下载种子文件
		mem：将种子保存在内存中
	*/
	FollowTorrent FollowTorrentStrategy `conf:"follow-torrent"`
	/*
	   种子文件下载完后暂停任务，默认：false
	   在开启 follow-torrent 选项后下载种子文件或磁力会自动开始下载任务进行下载，而同时开启当此选项后会建立相关任务并暂停。
	*/
	PauseMetadata bool
	// 加载已保存的元数据文件(.torrent)，默认:false
	BTSaveMetadata bool
	// 删除 BT 下载任务中未选择文件，默认:false
	BTRemoveUnselectedFile bool `json:"bt-remove-unselected-file"`
	/*
		BT强制加密, 默认: false
		启用后将拒绝旧的 BT 握手协议并仅使用混淆握手及加密。可以解决部分运营商对 BT 下载的封锁，且有一定的防版权投诉与迅雷吸血效果。
		此选项相当于后面两个选项(bt-require-crypto=true, bt-min-crypto-level=arc4)的快捷开启方式，但不会修改这两个选项的值。
	*/
	BTForceEncryption bool
	/*
	   BT加密需求，默认：false
	   启用后拒绝与旧的 BitTorrent 握手协议(\19BitTorrent protocol)建立连接，始终使用混淆处理握手。
	*/
	BTRequireCrypto bool
	// BT最低加密等级，可选：plain（明文），arc4（加密），默认：plain
	BTMinCryptoLevel CryptoLevel
	/*
		分离仅做种任务，默认：false
		从正在下载的任务中排除已经下载完成且正在做种的任务，并开始等待列表中的下一个任务。
	*/
	BTDetachSeedOnly bool
	// 自定义 User Agent
	UserAgent               string
	PeerAgent, PeerIdPrefix string
	/*
		下载停止后执行的命令
		从 正在下载 到 删除、错误、完成 时触发。暂停被标记为未开始下载，故与此项无关。
	*/
	OnDownloadStop string
	/*
		下载完成后执行的命令
		此项未定义则执行 下载停止后执行的命令 (on-download-stop)
	*/
	OnDownloadComplete string
	/*
		下载错误后执行的命令
		此项未定义则执行 下载停止后执行的命令 (on-download-stop)
	*/
	OnDownloadError string
	// 下载暂停后执行的命令
	OnDownloadPause string
	// 下载开始后执行的命令
	OnDownloadStart string
	// BT 下载完成后执行的命令
	OnBTDownloadComplete string `json:"on-bt-download-complete"`
	// 启用 JSON-RPC/XML-RPC 服务器, 默认:false
	EnableRPC bool
	// 接受所有远程请求, 默认:false
	RPCAllowOriginAll bool
	// 允许外部访问, 默认:false
	RPCListenAll bool
	// RPC 监听端口, 默认:6800
	RPCListenPort int
	// RPC 密钥
	RPCSecret string
	// RPC 最大请求大小
	RPCMaxRequestSize int64
	/*
		RPC 服务 SSL/TLS 加密, 默认：false
		启用加密后必须使用 https 或者 wss 协议连接
		不推荐开启，建议使用 web server 反向代理，比如 Nginx、Caddy ，灵活性更强。
	*/
	RPCSecure bool
	// 在 RPC 服务中启用 SSL/TLS 加密时的证书文件(.pem/.crt)
	RPCCertificate string
	// 在 RPC 服务中启用 SSL/TLS 加密时的私钥文件(.key)
	RPCPrivateKey string
	// 事件轮询方式, 可选：epoll, kqueue, port, poll, select, 不同系统默认值不同
	EventPoll EventPollType
	// 启用异步 DNS 功能。默认：true
	AsyncDNS bool
	// 指定异步 DNS 服务器列表，未指定则从 /etc/resolv.conf 中读取
	AsyncDNSServer []string
	/*
		指定单个网络接口，可能的值：接口，IP地址，主机名
		如果接口具有多个 IP 地址，则建议指定 IP 地址。
		已知指定网络接口会影响依赖本地 RPC 的连接的功能场景，即通过 localhost 和 127.0.0.1 无法与 Aria2 服务端进行讯通。
	*/
	Interface string
	/*
		指定多个网络接口，多个值之间使用逗号(,)分隔。
		使用 interface 选项时会忽略此项。
	*/
	MultipleInterface []string
	// 日志文件保存路径，忽略或设置为空为不保存，默认：不保存
	Log string
	// 日志级别，可选 debug, info, notice, warn, error 。默认：debug
	LogLevel LogLevel
	// 控制台日志级别，可选 debug, info, notice, warn, error ，默认：notice
	ConsoleLogLevel LogLevel
	// 安静模式，禁止在控制台输出日志，默认：false
	Quiet bool
	// 下载进度摘要输出间隔时间（秒），0 为禁止输出。默认：60
	SummaryInterval int
	Refer           string
	BTMaxOpenFiles  int
}

func DefaultConfig(dir string) *Config {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return nil
		}
	}

	return &Config{
		Trackers: []string{
			"https://trackerslist.com/all_aria2.txt",
			"https://cdn.staticaly.com/gh/XIU2/TrackersListCollection@master/all_aria2.txt",
			"https://trackers.p3terx.com/all_aria2.txt",
		},
		FileAllocation:          FileAllocTypeNone,
		DiskCache:               64 * utils.MB,
		NoFileAllocationLimit:   64 * utils.MB,
		Continue:                true,
		AlwaysResume:            true,
		RemoteTime:              false,
		InputFile:               filepath.Join(dir, "aria2.session"),
		SaveSession:             filepath.Join(dir, "aria2.session"),
		SaveSessionInterval:     1,
		AutoSaveInterval:        20,
		MaxFileNotFound:         10,
		RetryWait:               10,
		ConnectionTimeout:       10,
		MaxConCurrentDownload:   5,
		Split:                   64,
		MinSplitSize:            4 * utils.MB,
		PieceLength:             utils.MB,
		AllowPieceLengthChange:  true,
		DisableIpv6:             true,
		HttpAcceptGzip:          true,
		NoNetrc:                 true,
		AutoFileRenaming:        true,
		MinTLSVersion:           TLVSV1_2,
		ListenPort:              51413,
		DHTListenPort:           51413,
		EnableDHT:               true,
		DHTFilePath:             filepath.Join(dir, "DHT", "dht.dat"),
		DHT6FilePath:            filepath.Join(dir, "DHT", "dht6.dat"),
		DHTEntryPoint:           "dht.transmissionbt.com:6881",
		DHTEntryPoint6:          "dht.transmissionbt.com:6881",
		BTEnableLPD:             true,
		EnablePeerExchange:      true,
		BTMaxPeers:              128,
		BTRequestPeerSpeedLimit: 10 * utils.MB,
		MaxOverallUploadLimit:   2 * utils.MB,
		SeedRatio:               1.0,
		BTHashCheckSeed:         true,
		BTTrackerConnectTimeout: 10,
		BTTrackerTimeout:        10,
		BTPrioritizePiece: struct{ Head, Tail int64 }{
			Head: 32 * utils.MB,
			Tail: 32 * utils.MB,
		},
		RPCSaveUploadMetadata:  true,
		FollowTorrent:          SaveOnDisk,
		BTSaveMetadata:         true,
		BTRemoveUnselectedFile: true,
		BTForceEncryption:      true,
		BTRequireCrypto:        true,
		BTMinCryptoLevel:       CryptoLevelArc4,
		BTDetachSeedOnly:       true,
		UserAgent:              "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36 Edg/93.0.961.47",
		PeerAgent:              "Deluge 1.3.15",
		PeerIdPrefix:           "-DE13F0-",
		EnableRPC:              true,
		RPCAllowOriginAll:      true,
		OnDownloadStop:         filepath.Join(dir, "scripts", "delete.sh"),
		OnDownloadComplete:     filepath.Join(dir, "scripts", "clean.sh"),
		RPCListenAll:           true,
		RPCSecret:              "P3TERX",
		RPCMaxRequestSize:      10 * utils.MB,
		EventPoll:              EventPollTypeSelect,
		AsyncDNS:               true,
		AsyncDNSServer:         []string{"119.29.29.29", "223.5.5.5", "8.8.8.8", "1.1.1.1"},
		LogLevel:               LogLevelNotice,
		BTMaxOpenFiles:         16,
	}
}

func (c *Config) StartUp(ctx context.Context, client *http.Client, logger logger.AppLogger) {
	wg := &sync.WaitGroup{}
	// 更新 Trackers
	wg.Add(1)
	go func() {
		t := &parse.ParserTask{
			TaskName: "trackers",
			DstPath:  c.SaveDir,
			Ctx:      ctx,
			Client:   client,
			Logger:   logger,
		}
		// 下载 trackers
		dsts := make([]string, 0, len(c.Trackers))
		for i := 0; i < len(c.Trackers); i++ {
			dsts = append(dsts, fmt.Sprintf("%v.txt", i))
		}

		d := &downloader.CommonDownloader{}
		err := d.StartDownload(t, c.Trackers, dsts...)
		if err != nil && logger != nil {
			logger.LogErrorf("更新 trackers 失败: %v", err)
		} else {
			logger.LogInfo("更新 trackers 成功")
		}
		wg.Done()
	}()

	// 下载脚本
	wg.Add(1)
	go func() {
		t := &parse.ParserTask{
			TaskName: "scripts",
			DstPath:  c.SaveDir,
			Ctx:      ctx,
			Client:   client,
			Logger:   logger,
		}

		dsts := []string{"delete.sh", "clean.sh"}
		d := &downloader.CommonDownloader{}

		scripts := []string{
			"https://raw.githubusercontent.com/P3TERX/aria2.conf/master/delete.sh",
			"https://raw.githubusercontent.com/P3TERX/aria2.conf/master/clean.sh",
		}
		err := d.StartDownload(t, scripts, dsts...)
		if err != nil && logger != nil {
			logger.LogErrorf("更新脚本失败: %v", err)
		} else {
			logger.LogInfo("更新脚本成功")
		}
		wg.Done()
	}()

	// 下载 DHT
	wg.Add(1)
	go func() {
		t := &parse.ParserTask{
			TaskName: "DHT",
			DstPath:  c.SaveDir,
			Ctx:      ctx,
			Client:   client,
			Logger:   logger,
		}

		dsts := []string{"dht.dat", "dht6.dat"}
		d := &downloader.CommonDownloader{}

		urls := []string{
			"https://github.com/P3TERX/aria2.conf/blob/master/dht.dat?raw=true",
			"https://github.com/P3TERX/aria2.conf/blob/master/dht6.dat?raw=true",
		}

		err := d.StartDownload(t, urls, dsts...)
		if err != nil && logger != nil {
			logger.LogErrorf("更新脚本失败: %v", err)
		} else {
			logger.LogInfo("更新脚本成功")
		}
		wg.Done()
	}()

	wg.Wait()
}

func (c *Config) GenerateConfigFile() (string, error) {
	f, err := os.Create(filepath.Join(c.SaveDir, "ari2.conf"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	t := reflect.TypeOf(c)
	v := reflect.ValueOf(c)

	for k := 0; k < t.NumField(); k++ {
		lex := t.Field(k).Tag.Get("conf")
		val := v.Field(k).Interface()
		_, err = f.WriteString(fmt.Sprintf("%v=%v\n", lex, val))
		if err != nil {
			return "", err
		}
	}

	return f.Name(), nil
}
