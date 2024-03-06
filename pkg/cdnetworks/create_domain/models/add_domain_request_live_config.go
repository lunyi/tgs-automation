package models

import "github.com/alibabacloud-go/tea/tea"

type AddCdnDomainRequestLiveConfig struct {
	// {"en":"The live push-pull stream type, the optional values are pull and push, pull means pull flow; push means push flow.", "zh_CN":"直播推拉流类型，可选值为pull和push，pull表示拉流；   push表示推流。"}
	StreamType *string `json:"stream-type,omitempty" xml:"stream-type,omitempty"`
	// {"en":"The push-pull domain name is used to set the push-flow domain name corresponding to the live streaming domain name. When the stream-type is pull, at least one of the source IP address and the corresponding push-stream domain name is not empty. When the stream-type is push, Incoming.", "zh_CN":"配套推流域名，用于设置直播拉流域名对应的推流域名，当stream-type为pull时，源站IP和配套推流域名至少一个不为空；当stream-type为push时，无需传入。"}
	OriginPushHost *string `json:"origin-push-host,omitempty" xml:"origin-push-host,omitempty"`
	// {"en":"Source station IP. When the stream-type is pull, at least one of the source station IP and the companion push stream domain name is not empty.
	// 1. If it is a push-pull flow package, fill in 127.0.0.1, and the system will also default to 127.0.0.1.
	// 2. If it is directly returning to the source, fill in the source IP of the source pull stream.
	// ", "zh_CN":"源站IP，当stream-type为pull时，源站IP和配套推流域名至少一个不为空。
	// 1、如果是推拉流配套，则填写127.0.0.1，不传系统也默认为127.0.0.1
	// 2、如果是直接回源拉流，则填写回源拉流的源站IP"}
	LiveConfigOriginIps *string `json:"origin-ips,omitempty" xml:"origin-ips,omitempty"`
}

func (s AddCdnDomainRequestLiveConfig) String() string {
	return tea.Prettify(s)
}

func (s AddCdnDomainRequestLiveConfig) GoString() string {
	return s.String()
}

func (s *AddCdnDomainRequestLiveConfig) SetStreamType(v string) *AddCdnDomainRequestLiveConfig {
	s.StreamType = &v
	return s
}

func (s *AddCdnDomainRequestLiveConfig) SetOriginPushHost(v string) *AddCdnDomainRequestLiveConfig {
	s.OriginPushHost = &v
	return s
}

func (s *AddCdnDomainRequestLiveConfig) SetLiveConfigOriginIps(v string) *AddCdnDomainRequestLiveConfig {
	s.LiveConfigOriginIps = &v
	return s
}
