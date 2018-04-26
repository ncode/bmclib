package m1000e

// /cgi-bin/webcgi/loginSecurity
type LoginSecurityParams struct {
	SessionToken               string `url:"ST2"`                          //7bdaaa73307ebb471d0e71a9cecc44fb most likely the auth token
	EnforcedIpBlockEnable      bool   `url:"ENFORCED_IPBLOCK_enable,int"`  //1
	EnforcedIpBlockFailcount   int    `url:"ENFORCED_IPBLOCK_failcount"`   //5
	EnforcedIpBlockFailwindow  int    `url:"ENFORCED_IPBLOCK_failwindow"`  //60
	EnforcedIpBlockPenaltyTime int    `url:"ENFORCED_IPBLOCK_penaltytime"` //300
}

// cgi-bin/webcgi/datetime
type Datetime struct {
	SessionToken          string `url:"ST2"`                      //ST2=ba9a6bbf88764c829ca4f49146fd4817
	NtpEnable             bool   `url:"NTP_enable,int"`           //NTP_enable=1
	NtpServer1            string `url:"NTP_server1"`              //NTP_server1=ntp0.dev.booking.com
	NtpServer2            string `url:"NTP_server2"`              //NTP_server2=ntp2.dev.booking.com
	NtpServer3            string `url:"NTP_server3"`              //NTP_server3=ntp3.dev.booking.com
	DateTimeChanged       bool   `url:"datetimeChanged,int"`      //datetimeChanged=0
	CmcTimeTimezoneString string `url:"CMC_TIME_timezone_string"` //CMC_TIME_timezone_string=CET
	TzChanged             bool   `url:"tzChanged,int"`            //tzChanged=1
}

// Manages user login account parameters,
// notes:
// 1. the url parameter ?id=<int> and UserID form value must match,
// 2. the CMC_GROUP, Privilege params must match,
//    4095 = Administrator - full access
//    3801 = Power user - no access to chassis, user config, debug commands
//    1    = Guest user - login access only
// Endpoint /webcgi/user?id=1
type UserParams struct {
	SessionToken    string `url:"ST2"`                //ST2=ba9a6bbf88764c829ca4f49146fd4817
	Privilege       int    `url:"Privilege"`          //Privilege=4095
	UserId          int    `url:"UserID"`             //UserID=1
	EnableUser      bool   `url:"EnableUser,int"`     //EnableUser=1
	UserName        string `url:"UserName"`           //UserName=Test
	ChangePassword  bool   `url:"ChangePassword,int"` //ChangePassword=1
	Password        string `url:"Password"`           //Password=foobar
	ConfirmPassword string `url:"ConfirmPassword"`    //ConfirmPassword=foobar
	CmcGroup        int    `url:"CMC_GROUP"`          //CMC_GROUP=4095
	Login           bool   `url:"login,int"`          //login=1
	Cfg             bool   `url:"cfg,int"`            //cfg=1
	CfgUser         bool   `url:"cfguser,int"`        //cfguser=1
	ClearLog        bool   `url:"clearlog,int"`       //clearlog=1
	ChassisControl  bool   `url:"chassiscontrol,int"` //chassiscontrol=1
	SuperUser       bool   `url:"superuser,int"`      //superuser=1
	ServerAdmin     bool   `url:"serveradmin,int"`    //serveradmin=1
	TestAlert       bool   `url:"testalert,int"`      //testalert=1
	Debug           bool   `url:"debug,int"`          //debug=1
	AFabricAdmin    bool   `url:"afabricadmin,int"`   //afabricadmin=1
	BFabricAdmin    bool   `url:"bfabricadmin,int"`   //bfabricadmin=1
	CFabricAcminc   bool   `url:"cfabricadmin,int"`   //cfabricadmin=1
}

// /cgi-bin/webcgi/interfaces
type InterfaceParams struct {
	SessionToken                     string `url:"ST2"`                                      //ST2=2754be61766abf5808085b3f2dd7bd94
	SerialEnable                     bool   `url:"SERIAL_enable,int"`                        //SERIAL_enable=1
	SerialRedirect                   bool   `url:"SERIAL_redirect_enable,int"`               //SERIAL_redirect_enable=1
	SerialTimeout                    int    `url:"SERIAL_timeout"`                           //SERIAL_timeout=1800
	SerialBaudrate                   int    `url:"SERIAL_baudrate"`                          //SERIAL_baudrate=115200
	SerialConsoleNoAuth              bool   `url:"SERIAL_console_no_auth,int"`               //SERIAL_console_no_auth=0
	SerialQuitKey                    string `url:"SERIAL_quit_key"`                          //SERIAL_quit_key=^\
	SerialHistoryBufSize             int    `url:"SERIAL_history_buf_size"`                  //SERIAL_history_buf_size=8192
	SerialLoginCommand               string `url:"SERIAL_login_command"`                     //SERIAL_login_command=
	WebserverEnable                  bool   `url:"WEBSERVER_enable,int"`                     //WEBSERVER_enable=1
	WebserverMaxSessions             int    `url:"WEBSERVER_maxsessions"`                    //WEBSERVER_maxsessions=4
	WebserverTimeout                 int    `url:"WEBSERVER_timeout"`                        //WEBSERVER_timeout=1800
	WebserverHttpPort                int    `url:"WEBSERVER_http_port"`                      //WEBSERVER_http_port=80
	WebserverHttpsPort               int    `url:"WEBSERVER_https_port"`                     //WEBSERVER_https_port=443
	SshEnable                        bool   `url:"SSH_enable,int"`                           //SSH_enable=1
	SshMaxSessions                   int    `url:"SSH_maxsessions"`                          //SSH_maxsessions=4
	SshTimeout                       int    `url:"SSH_timeout`                               //SSH_timeout=1800
	SshPort                          int    `url:"SSH_port"`                                 //SSH_port=22
	TelnetEnable                     bool   `url:"TELNET_enable,int"`                        //TELNET_enable=1
	TelnetMaxSessions                int    `url:"TELNET_maxsessions"`                       //TELNET_maxsessions=4
	TelnetTimeout                    int    `url:"TELNET_timeout"`                           //TELNET_timeout=1800
	TelnetPort                       int    `url:"TELNET_port"`                              //TELNET_port=23
	RacadmEnable                     bool   `url:"RACADM_enable,int"`                        //RACADM_enable=1
	RacadmMaxSessions                int    `url:"RACADM_maxsessions"`                       //RACADM_maxsessions=4
	RacadmTimeout                    int    `url:"RACADM_timeout"`                           //RACADM_timeout=60
	SnmpEnable                       bool   `url:"SNMP_enable,int"`                          //SNMP_enable=1
	SnmpCommunityNameGet             string `url:"SNMP_COMMUNITYNAME_get"`                   //SNMP_COMMUNITYNAME_get=public
	SnmpProtocol                     int    `url:"SNMP_Protocol"`                            //SNMP_Protocol=0
	SnmpDiscoveryPortSet             int    `url:"SNMP_DiscoveryPort_set"`                   //SNMP_DiscoveryPort_set=161
	ChassisLoggingRemoteSyslogEnable bool   `url:"CHASSIS_LOGGING_remote_syslog_enable,int"` //CHASSIS_LOGGING_remote_syslog_enable=1
	ChassisLoggingRemoteSyslogHost1  string `url:"CHASSIS_LOGGING_remote_syslog_host_1"`     //CHASSIS_LOGGING_remote_syslog_host_1=provision.anycast.prod.booking.com
	ChassisLoggingRemoteSyslogHost2  string `url:"CHASSIS_LOGGING_remote_syslog_host_2"`     //CHASSIS_LOGGING_remote_syslog_host_2=
	ChassisLoggingRemoteSyslogHost3  string `url:"CHASSIS_LOGGING_remote_syslog_host_3"`     //CHASSIS_LOGGING_remote_syslog_host_3=
	ChassisLoggingRemoteSyslogPort   int    `url:"CHASSIS_LOGGING_remote_syslog_port"`       //CHASSIS_LOGGING_remote_syslog_port=514

}

// /cgi-bin/webcgi/nic
//type NicParams struct {
//	SessionToken string `url:"ST2"` //ST2=2754be61766abf5808085b3f2dd7bd94
//NETWORK_NIC_enable=1
//DNS_register_cmc=1
//DNS_cmc_name=cmc-GF85C92
//DNS_use_dhcp_domain=1
//DNS_register_interval=0
//NETWORK_NIC_TUNE_auto_neg=1
//NETWORK_NIC_TUNE_mtu=1500
//NETWORK_NIC_TUNE_redundant=0
//NETWORK_NIC_ipv4_enable=1
//NETWORK_NIC_dhcp_enable=1
//DNS_dhcp_enable=1
//NETWORK_NIC_IPV6_enable=1
//NETWORK_NIC_IPV6_autoconfig_enable=1
//DNS_IPV6_dhcp_enable=1
//FIPS_Mode=0
//hidden_NETWORK_NIC_ipaddr=192.168.0.120
//hidden_NETWORK_NIC_gateway=192.168.0.1
//hidden_NETWORK_NIC_netmask=255.255.255.0
//hidden_NETWORK_NIC_TUNE_speed=100
//hidden_NETWORK_NIC_TUNE_fullduplex=1
//hidden_NETWORK_NIC_TUNE_auto_neg=1
//hidden_NETWORK_NIC_TUNE_mtu=1500
//hidden_DNS_server1=0.0.0.0
//hidden_DNS_server2=0.0.0.0
//hidden_DNS_cmc_name=cmc-GF85C92
//hidden_DNS_domain_name=
//hidden_DNS_register_interval=0
//hidden_NETWORK_NIC_enable=1
//hidden_NETWORK_NIC_dhcp_enable=1
//hidden_DNS_dhcp_enable=1
//hidden_DNS_register_cmc=1
//hidden_DNS_use_dhcp_domain=1
//hidden_NETWORK_NIC_TUNE_redundant=0
//hidden_NETWORK_NIC_ipv4_enable=1
//hidden_NETWORK_NIC_IPV6_enable=1
//hidden_NETWORK_NIC_IPV6_autoconfig_enable=1
//hidden_NETWORK_NIC_IPV6_ipaddr=%3A%3A
//hidden_NETWORK_NIC_IPV6_prefix_length=64
//hidden_NETWORK_NIC_IPV6_gateway=%3A%3A
//hidden_DNS_IPV6_dhcp_enable=1
//hidden_DNS_IPV6_server1=%3A%3A
//hidden_DNS_IPV6_server2=%3A%3A
//}
