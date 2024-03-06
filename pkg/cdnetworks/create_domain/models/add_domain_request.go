package models

import "github.com/alibabacloud-go/tea/tea"

type AddCdnDomainRequest struct {
	// {"en":"The current version is 1.0.0", "zh_CN":"版本号，当前版本号1.0.0"}
	Version *string `json:"version,omitempty" xml:"version,omitempty" require:"true"`
	// {"en":"CDN accelerated domain name. Support generic domain name, witch starts with a dot, such as: .example.com, If the domain example.com have China ICP, then the domain name xx.example.com have ICP too.", "zh_CN":"需要接入CDN的域名。支持泛域名，以符号“.”开头，如：.example.com，
	//   泛域名也包含多级 a.b.example.com。 如果example.com已备案，那么域名xx.example.com则不需要备案。"}
	DomainName *string `json:"domain-name,omitempty" xml:"domain-name,omitempty" require:"true"`
	// {"en":"The value should be an effective domain. The added domain will copy the config of this domain.",
	//   "zh_CN":"指定参考域名。当使用参考域名来新增域名是，新域名将复用参考域名的配置"}
	ReferencedDomainName *string `json:"referenced-domain-name,omitempty" xml:"referenced-domain-name,omitempty"`
	// {"en":"Configure single template. If you want to add a new accelerated domain name with a specified configuration, you can use a specified configured template. For details, please consult the corresponding customer support staff.", "zh_CN":"配置单模板，特定的使用场景下，如果希望新增的加速域名参照某些指定配置时，可以指定配置单模板，具体使用请咨询对应的客户负责人。"}
	ConfigFormId *string `json:"config-form-id,omitempty" xml:"config-form-id,omitempty"`
	// {"en":"The optional value is: true or false.", "zh_CN":"是否纯海外加速，入参范围：true、false"}
	AccelerateNoChina *string `json:"accelerate-no-china,omitempty" xml:"accelerate-no-china,omitempty"`
	// {"en":"The id of contract", "zh_CN":"合同号"}
	ContractId *string `json:"contract-id,omitempty" xml:"contract-id,omitempty" require:"true"`
	// {"en":"The id of item", "zh_CN":"产品号"}
	ItemId *string `json:"item-id,omitempty" xml:"item-id,omitempty" require:"true"`
	// {"en":"An alias for public service cname. When you have multiple domains that need to share a cname, then you may specify a cname-label here. If domains use the same cname-label, they will share the same cname and has the same dns coverage. Note: 1. Restrictions on domains with the same cname-label: they should have the same acceleration type,  certificate (if use ssl), service areas. 2. Multiple http domains can share the same cname; multiple sni https domains can share the same cname as wll. 3. When a single domain uses a cname-label, it can be cancelled acceleration; while multiple domains are not allowed to cancel acceleration for part of them.", "zh_CN":"共用一级别名，客户存在较多一级域名共用的需求，因此在接口中引入cname-label标识，即拥有相同cname-label的一组域名，共用一级cname。当加速域名的加速区域、加速类型、资源链是一致的时候，建议使用共用一级别名，相同的业务使用相同一级别名cname。 注意： 1、拥有相同cname-label的域名共用一级cname，且有完全一致的dns覆盖 2、共用一级的约束：配置一致、加速类型一致、证书id一致（如果有证书）、加速区域一致、平台套餐一致 3、多个http域名可共用一级，多个sni https域名可共用一级，多个共享证书域名可共用一级 4、单个域名使用cname-label时，域名可取消加速；多个域名共用一级时，不允许取消加速这些域名 5、支持通过修改cname-label达到修改cname的目的"}
	CnameLabel *string `json:"cname-label,omitempty" xml:"cname-label,omitempty"`
	// {"en":"Remarks, up to 1000 characters", "zh_CN":"备注信息，最大限制1000个字符"}
	Comment *string `json:"comment,omitempty" xml:"comment,omitempty"`
	// {"en":"Pass the response header of client IP. The optional values are Cdn-Src-Ip and X-Forwarded-For. The default value is Cdn-Src-Ip.", "zh_CN":"传递客户端ip的响应头部，可选值为Cdn-Src-Ip和X-Forwarded-For，默认值为Cdn-Src-Ip"}
	HeaderOfClientIp *string `json:"header-of-clientip,omitempty" xml:"header-of-clientip,omitempty"`
	// {"en":"Back-to-origin policy setting, which is used to set the origin site information and the back-to-origin policy of the accelerated domain name", "zh_CN":"回源策略设置，用于设置加速域名的源站信息和回源策略"}
	OriginConfig *AddCdnDomainRequestOriginConfig `json:"origin-config,omitempty" xml:"origin-config,omitempty" require:"true" type:"Struct"`
	// {"en":"Live domain configuration, used to set the livestream acceleration domain origin config.
	// Note: In addition to the API call permission, you need to contact the dedicated customer service to apply for the corresponding API client template.", "zh_CN":"直播域名配置，用于设置直播加速域名的推拉流（使用需申请）
	// 注意：该节点下的相关参数配置，除开通API调用权限外，还需要联系专属客服申请开通对应的API客户模板"}
	LiveConfig *AddCdnDomainRequestLiveConfig `json:"live-config,omitempty" xml:"live-config,omitempty" type:"Struct"`
	// {"en":"Livestream domain settings. Set the publishing point of the live push-pull domain.
	// note:
	// 1. The pull stream and the corresponding push stream domain must be configured with the same publishing point.
	// 2. If you are not going to modify the publishing point, please do not pass this param.
	// 3. The publishing point adopts the overlay update. Each time you modify, you need to submit all the publishing points. You cannot submit only the parts that need to be modified.", "zh_CN":"设置直播推拉流域名的发布点
	// 注意：
	// 1、拉流和对应的推流域名，必须配置相同的发布点；
	// 2、不想修改发布点时，不要传入该节点及以下入参；
	// 3、发布点采用覆盖式更新，每次修改时，需要提交全部发布点，不能仅提交需要修改的部分。"}
	PublishPoints []*AddCdnDomainRequestPublishPoints `json:"publish-points,omitempty" xml:"publish-points,omitempty" type:"Repeated"`
}

func (s AddCdnDomainRequest) String() string {
	return tea.Prettify(s)
}

func (s AddCdnDomainRequest) GoString() string {
	return s.String()
}

func (s *AddCdnDomainRequest) SetVersion(v string) *AddCdnDomainRequest {
	s.Version = &v
	return s
}

func (s *AddCdnDomainRequest) SetDomainName(v string) *AddCdnDomainRequest {
	s.DomainName = &v
	return s
}

func (s *AddCdnDomainRequest) SetReferencedDomainName(v string) *AddCdnDomainRequest {
	s.ReferencedDomainName = &v
	return s
}

func (s *AddCdnDomainRequest) SetConfigFormId(v string) *AddCdnDomainRequest {
	s.ConfigFormId = &v
	return s
}

func (s *AddCdnDomainRequest) SetAccelerateNoChina(v string) *AddCdnDomainRequest {
	s.AccelerateNoChina = &v
	return s
}

func (s *AddCdnDomainRequest) SetContractId(v string) *AddCdnDomainRequest {
	s.ContractId = &v
	return s
}

func (s *AddCdnDomainRequest) SetItemId(v string) *AddCdnDomainRequest {
	s.ItemId = &v
	return s
}

func (s *AddCdnDomainRequest) SetCnameLabel(v string) *AddCdnDomainRequest {
	s.CnameLabel = &v
	return s
}

func (s *AddCdnDomainRequest) SetComment(v string) *AddCdnDomainRequest {
	s.Comment = &v
	return s
}

func (s *AddCdnDomainRequest) SetHeaderOfClientIp(v string) *AddCdnDomainRequest {
	s.HeaderOfClientIp = &v
	return s
}

func (s *AddCdnDomainRequest) SetOriginConfig(v *AddCdnDomainRequestOriginConfig) *AddCdnDomainRequest {
	s.OriginConfig = v
	return s
}

func (s *AddCdnDomainRequest) SetLiveConfig(v *AddCdnDomainRequestLiveConfig) *AddCdnDomainRequest {
	s.LiveConfig = v
	return s
}

func (s *AddCdnDomainRequest) SetPublishPoints(v []*AddCdnDomainRequestPublishPoints) *AddCdnDomainRequest {
	s.PublishPoints = v
	return s
}
