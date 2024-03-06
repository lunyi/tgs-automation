package models

import "github.com/alibabacloud-go/tea/tea"

type QueryApiDomainListServiceResponseDomainList struct {
	// {"en":"Name of accelerated domain name", "zh_CN":"加速域名的名称"}
	DomainName *string `json:"domain-name,omitempty" xml:"domain-name,omitempty" require:"true"`
	// {"en":"The corresponding domain name ID: the domain name ID, used to perform the query and modification operations of the related domain name.", "zh_CN":"对应的域名ID：域名ID，用于执行相关域名的查询、修改操作等。"}
	DomainId *int `json:"domain-id,omitempty" xml:"domain-id,omitempty" require:"true"`
	// {"en":"Accelerated domain CNAME corresponding to CNAME, for example: 7nt6mrh7sdkslj.cdn30.com", "zh_CN":"加速域名对应的CNAME域名，例如：7nt6mrh7sdkslj.cdn30.com"}
	Cname *string `json:"cname,omitempty" xml:"cname,omitempty" require:"true"`
	// {"en":"Speed up the service type of the domain name, the value is:
	// Web/web-https: web acceleration / web acceleration - https
	// Wsa/wsa-https: Total Station Acceleration / Total Station Acceleration - https
	// Vodstream/vod-https: on-demand acceleration/on-demand acceleration - https
	// Download/dl-https: Download acceleration/download acceleration - https
	// livestream/live-https/cloudv-live: Live acceleration/Live acceleration - https/Cloud vedio for live
	// 1028: Content Acceleration;
	// 1115: Dynamic Web Acceleration;
	// 1369: Media Acceleration - RTMP
	// 1391: Download Acceleration
	// 1348: Media Acceleration Live Broadcast
	// 1551: Floodshield", "zh_CN":"加速域名的服务类型，取值：
	// web/web-https：网页加速/网页加速-https
	// wsa/wsa-https：全站加速/全站加速-https
	// vodstream/vod-https：点播加速/点播加速-https
	// download/dl-https：下载加速/下载加速-https
	// livestream/live-https/cloudv-live：直播加速/直播加速-https/云直播
	// 1028 : Content Acceleration;
	// 1115 : Dynamic Web Acceleration;
	// 1369 : Media Acceleration - RTMP
	// 1391 : Download Acceleration
	// 1348 : Media Acceleration Live Broadcast
	// 1551 : Floodshield"}
	ServiceType *string `json:"service-type,omitempty" xml:"service-type,omitempty" require:"true"`
	// {"en":"The deployment status of the accelerated domain name: Deployed indicates that the accelerated domain name configuration is complete; InProgress indicates that the deployment task for this accelerated domain name configuration is still in InProgress and may be in a queue, deploy, or fail in any one of the states", "zh_CN":"加速域名的部署状态：Deployed表示该加速域名配置完成部署；InProgress表示该加速域名配置的部署任务还在进行中，可能处于排队、部署中或失败任意一种状态"}
	Status *string `json:"status,omitempty" xml:"status,omitempty" require:"true"`
	// {"en":"Accelerate the CDN service status of the domain name: This is false when the accelerated domain name CDN service is canceled; this is true when the accelerated domain name CDN service is restored.", "zh_CN":"加速域名的CDN服务状态：当取消加速域名CDN服务后，此项为false；当恢复加速域名CDN服务后，此项为true"}
	CdnServiceStatus *string `json:"cdn-service-status,omitempty" xml:"cdn-service-status,omitempty" require:"true"`
	// {"en":"Accelerated domain activation: This is false when the accelerated domain name service is disabled; true when the accelerated domain name service is enabled", "zh_CN":"加速域名的启用状态：当禁用加速域名服务后，此项为false；当启用加速域名服务后，此项为true"}
	Enabled *string `json:"enabled,omitempty" xml:"enabled,omitempty" require:"true"`
}

func (s QueryApiDomainListServiceResponseDomainList) String() string {
	return tea.Prettify(s)
}

func (s QueryApiDomainListServiceResponseDomainList) GoString() string {
	return s.String()
}

func (s *QueryApiDomainListServiceResponseDomainList) SetDomainName(v string) *QueryApiDomainListServiceResponseDomainList {
	s.DomainName = &v
	return s
}

func (s *QueryApiDomainListServiceResponseDomainList) SetDomainId(v int) *QueryApiDomainListServiceResponseDomainList {
	s.DomainId = &v
	return s
}

func (s *QueryApiDomainListServiceResponseDomainList) SetCname(v string) *QueryApiDomainListServiceResponseDomainList {
	s.Cname = &v
	return s
}

func (s *QueryApiDomainListServiceResponseDomainList) SetServiceType(v string) *QueryApiDomainListServiceResponseDomainList {
	s.ServiceType = &v
	return s
}

func (s *QueryApiDomainListServiceResponseDomainList) SetStatus(v string) *QueryApiDomainListServiceResponseDomainList {
	s.Status = &v
	return s
}

func (s *QueryApiDomainListServiceResponseDomainList) SetCdnServiceStatus(v string) *QueryApiDomainListServiceResponseDomainList {
	s.CdnServiceStatus = &v
	return s
}

func (s *QueryApiDomainListServiceResponseDomainList) SetEnabled(v string) *QueryApiDomainListServiceResponseDomainList {
	s.Enabled = &v
	return s
}
