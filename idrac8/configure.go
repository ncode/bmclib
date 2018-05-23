package idrac8

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.booking.com/go/bmc/cfgresources"
	"reflect"
	"runtime"
	"strconv"
)

// returns the calling function.
func funcName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func (i *IDrac8) ApplyCfg(config *cfgresources.ResourcesConfig) (err error) {
	cfg := reflect.ValueOf(config).Elem()

	//Each Field in ResourcesConfig struct is a ptr to a resource,
	//Here we figure the resources to be configured, i.e the ptr is not nil
	for r := 0; r < cfg.NumField(); r++ {
		resourceName := cfg.Type().Field(r).Name
		if cfg.Field(r).Pointer() != 0 {
			switch resourceName {
			case "User":
				//retrieve users resource values as an interface
				userAccounts := cfg.Field(r).Interface()

				//assert userAccounts interface to its actual type - A slice of ptrs to User
				for id, user := range userAccounts.([]*cfgresources.User) {

					//the dells have user id 1 set to a anon user, so we start with 2.
					userId := id + 2
					//setup params to post
					err := i.applyUserParams(user, userId)
					if err != nil {
						log.WithFields(log.Fields{
							"step":     "ApplyCfg",
							"Resource": cfg.Field(r).Kind(),
							"IP":       i.ip,
							"Error":    err,
						}).Warn("Unable to set user config.")
					}
				}

			case "Syslog":
				syslogCfg := cfg.Field(r).Interface().(*cfgresources.Syslog)
				err := i.applySyslogParams(syslogCfg)
				if err != nil {
					log.WithFields(log.Fields{
						"step":     "ApplyCfg",
						"resource": cfg.Field(r).Kind(),
						"IP":       i.ip,
						"Error":    err,
					}).Warn("Unable to set Syslog config.")
				}
			case "Network":
				fmt.Printf("%s: %v : %s\n", resourceName, cfg.Field(r), cfg.Field(r).Kind())
			case "Ntp":
				ntpCfg := cfg.Field(r).Interface().(*cfgresources.Ntp)
				err := i.applyNtpParams(ntpCfg)
				if err != nil {
					log.WithFields(log.Fields{
						"step":     "ApplyCfg",
						"resource": cfg.Field(r).Kind(),
						"IP":       i.ip,
					}).Warn("Unable to set NTP config.")
				}
			case "Ldap":
				ldapCfg := cfg.Field(r).Interface().(*cfgresources.Ldap)
				i.applyLdapParams(ldapCfg)
			case "Ssl":
				fmt.Printf("%s: %v : %s\n", resourceName, cfg.Field(r), cfg.Field(r).Kind())
			default:
				log.WithFields(log.Fields{
					"step": "ApplyCfg",
				}).Warn("Unknown resource.")
				//fmt.Printf("%v\n", cfg.Field(r))

			}
		}
	}

	return err
}

// Encodes the string the way idrac expects credentials to be sent
// foobar == @066@06f@06f@062@061@072
// convert ever character to its hex equiv, and prepend @0
func encodeCred(s string) string {
	r := ""
	for _, c := range s {
		r += fmt.Sprintf("@0%x", c)
	}

	return r
}

// escapeLdapString escapes ldap parameters strings
func escapeLdapString(s string) string {
	r := ""
	for _, c := range s {
		if c == '=' || c == ',' {
			r += fmt.Sprintf("\\%c", c)
		} else {
			r += string(c)
		}
	}

	return r
}

// Return bool value if the role is valid.
func (i *IDrac8) isRoleValid(role string) bool {

	validRoles := []string{"admin", "user"}
	for _, v := range validRoles {
		if role == v {
			return true
		}
	}

	return false
}

// attempts to add the user
// if the user exists, update the users password.
func (i *IDrac8) applyUserParams(cfg *cfgresources.User, Id int) (err error) {

	if cfg.Name == "" {
		log.WithFields(log.Fields{
			"step": "apply-user-cfg",
		}).Fatal("User resource expects parameter: Name.")
	}

	if cfg.Password == "" {
		log.WithFields(log.Fields{
			"step": "apply-user-cfg",
		}).Fatal("User resource expects parameter: Password.")
	}

	if !i.isRoleValid(cfg.Role) {
		log.WithFields(log.Fields{
			"step":     "apply-user-cfg",
			"Username": cfg.Name,
		}).Warn("User resource Role must be declared and a must be a valid role: 'admin' OR 'user'.")
		return
	}

	var enable string
	if cfg.Enable == false {
		enable = "Disabled"
	} else {
		enable = "Enabled"
	}

	user := User{UserName: encodeCred(cfg.Name), Password: encodeCred(cfg.Password), Enable: enable, SolEnable: "Enabled"}

	switch cfg.Role {
	case "admin":
		user.Privilege = "511"
		user.IpmiLanPrivilege = "Administrator"
	case "user":
		user.Privilege = "497"
		user.IpmiLanPrivilege = "Operator"

	}

	data := make(map[string]User)
	data["iDRAC.Users"] = user

	payload, err := json.Marshal(data)
	if err != nil {
		log.WithFields(log.Fields{
			"step": funcName(),
		}).Warn("Unable to marshal syslog payload.")
		return err
	}

	endpoint := fmt.Sprintf("sysmgmt/2012/server/configgroup/iDRAC.Users.%d", Id)
	response, err := i.put(endpoint, payload, false)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": endpoint,
			"step":     funcName(),
			"response": string(response),
		}).Warn("PUT request failed.")
		return err
	}

	log.WithFields(log.Fields{
		"IP":   i.ip,
		"User": user.UserName,
	}).Info("User parameters applied.")

	return err
}

func (i *IDrac8) applySyslogParams(cfg *cfgresources.Syslog) (err error) {

	var port int
	enable := "Enabled"

	if cfg.Server == "" {
		log.WithFields(log.Fields{
			"step": funcName(),
		}).Warn("Syslog resource expects parameter: Server.")
		return
	}

	if cfg.Port == 0 {
		log.WithFields(log.Fields{
			"step": funcName(),
		}).Debug("Syslog resource port set to default: 514.")
		port = 514
	} else {
		port = cfg.Port
	}

	if cfg.Enable != true {
		enable = "Disabled"
		log.WithFields(log.Fields{
			"step": funcName(),
		}).Debug("Syslog resource declared with enable: false.")
	}

	data := make(map[string]Syslog)

	data["iDRAC.SysLog"] = Syslog{
		Port:    strconv.Itoa(port),
		Server1: cfg.Server,
		Enable:  enable,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.WithFields(log.Fields{
			"step": funcName(),
		}).Warn("Unable to marshal syslog payload.")
		return err
	}

	endpoint := "sysmgmt/2012/server/configgroup/iDRAC.SysLog"
	response, err := i.put(endpoint, payload, false)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": endpoint,
			"step":     funcName(),
			"response": string(response),
		}).Warn("PUT request failed.")
		return err
	}

	log.WithFields(log.Fields{
		"IP": i.ip,
	}).Info("Syslog parameters applied.")

	return err
}

func (i *IDrac8) applyNtpParams(cfg *cfgresources.Ntp) (err error) {

	if cfg.Server1 == "" {
		log.WithFields(log.Fields{
			"step": "apply-ntp-cfg",
		}).Warn("NTP resource expects parameter: server1.")
		return
	}

	if cfg.Timezone == "" {
		log.WithFields(log.Fields{
			"step": "apply-ntp-cfg",
		}).Warn("NTP resource expects parameter: timezone.")
		return
	}

	i.applyTimezoneParam(cfg.Timezone)
	i.applyNtpServerParam(cfg)

	return err
}

func (i *IDrac8) applyNtpServerParam(cfg *cfgresources.Ntp) {

	var enable int
	if cfg.Enable != true {
		log.WithFields(log.Fields{
			"step": funcName(),
		}).Debug("Ntp resource declared with enable: false.")
		enable = 0
	} else {
		enable = 1
	}

	//https://10.193.251.10/data?set=tm_ntp_int_opmode:1, \\
	//                               tm_ntp_str_server1:ntp0.lhr4.prod.booking.com, \\
	//                               tm_ntp_str_server2:ntp0.ams4.prod.booking.com, \\
	//                               tm_ntp_str_server3:ntp0.fra4.prod.booking.com
	queryStr := fmt.Sprintf("set=tm_ntp_int_opmode:%d,", enable)
	queryStr += fmt.Sprintf("tm_ntp_str_server1:%s,", cfg.Server1)
	queryStr += fmt.Sprintf("tm_ntp_str_server2:%s,", cfg.Server2)
	queryStr += fmt.Sprintf("tm_ntp_str_server3:%s,", cfg.Server3)

	//GET - params as query string
	//ntp servers

	endpoint := fmt.Sprintf("data?%s", queryStr)
	response, err := i.get(endpoint, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"IP":       i.ip,
			"endpoint": endpoint,
			"step":     funcName(),
			"response": string(response),
		}).Warn("GET request failed.")
	}

	log.WithFields(log.Fields{
		"IP": i.ip,
	}).Info("NTP servers param applied.")

}

//applies ldap config parameters
func (i *IDrac8) applyLdapParams(cfg *cfgresources.Ldap) {
	// LDAP settings
	// Notes: - all non-numeric,alphabetic characters are escaped
	//        - the idrac posts each payload twice?
	//        - requests can be either POST or GET except for the final one - postset?ldapconf

	//Set ldap groups

	r := i.applyLdapServerParam(cfg)
	if r != 0 {
		return
	}

	r = i.applyLdapSearchFilterParam(cfg)
	if r != 0 {
		return
	}

	r = i.applyLdapGroupParam(cfg)
	if r != 0 {
		return
	}

	r = i.applyLdapRoleGroupPrivParam(cfg)
	if r != 0 {
		return
	}

}

// Applies ldap server param
// https://10.193.251.10/data?set=xGLServer:ldaps.prod.blah.com
func (i *IDrac8) applyLdapServerParam(cfg *cfgresources.Ldap) int {

	if cfg.Server == "" {
		log.WithFields(log.Fields{
			"step": "applyLdapServerParam",
		}).Warn("Ldap resource parameter Server required but not declared.")
		return 1
	}

	endpoint := fmt.Sprintf("data?set=xGLServer:%s", cfg.Server)
	response, err := i.get(endpoint, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"IP":       i.ip,
			"endpoint": endpoint,
			"step":     funcName(),
			"response": string(response),
		}).Warn("GET request failed.")
		return 1
	}

	log.WithFields(log.Fields{
		"IP": i.ip,
	}).Info("Ldap server param set.")

	return 0
}

// Applies ldap search filter param.
// set=xGLSearchFilter:objectClass\=posixAccount
func (i *IDrac8) applyLdapSearchFilterParam(cfg *cfgresources.Ldap) int {
	if cfg.SearchFilter == "" {
		log.WithFields(log.Fields{
			"step": "applyLdapSearchFilterParam",
		}).Warn("Ldap resource parameter SearchFilter required but not declared.")
		return 1
	}

	endpoint := fmt.Sprintf("data?set=xGLSearchFilter:%s", escapeLdapString(cfg.SearchFilter))
	response, err := i.get(endpoint, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"IP":       i.ip,
			"endpoint": endpoint,
			"step":     funcName(),
			"response": string(response),
		}).Warn("GET request failed.")
		return 1
	}

	log.WithFields(log.Fields{
		"IP": i.ip,
	}).Info("Ldap search filter param applied.")

	return 0
}

//applies ldap group params
//data?set=xGLGroup1Name:cn\=bmcAdmins\,ou\=Group\,dc\=activehotels\,dc\=com
func (i *IDrac8) applyLdapGroupParam(cfg *cfgresources.Ldap) int {

	if cfg.Group == "" {
		log.WithFields(log.Fields{
			"step": "applyLdapGroupParam",
		}).Warn("Ldap resource parameter Group required but not declared.")
		return 1
	}

	if cfg.GroupBaseDn == "" {
		log.WithFields(log.Fields{
			"step": "applyLdapGroupParam",
		}).Warn("Ldap resource parameter GroupBaseDn required but not declared.")
		return 1
	}

	groupDn := fmt.Sprintf("cn=%s,%s", cfg.Group, cfg.GroupBaseDn)

	groupDn = escapeLdapString(groupDn)

	endpoint := fmt.Sprintf("data?set=xGLGroup1Name:%s", groupDn)
	response, err := i.get(endpoint, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"IP":       i.ip,
			"endpoint": endpoint,
			"step":     funcName(),
			"response": string(response),
		}).Warn("GET request failed.")
		return 1
	}

	log.WithFields(log.Fields{
		"IP": i.ip,
	}).Info("Ldap GroupDN config applied.")

	return 0
}

//TODO - refactor to allow multiple role groups
// Apply ldap group privileges
//https://10.193.251.10/postset?ldapconf
// data=LDAPEnableMode:3,xGLNameSearchEnabled:0,xGLBaseDN:ou%5C%3DPeople%5C%2Cdc%5C%3Dactivehotels%5C%2Cdc%5C%3Dcom,xGLUserLogin:uid,xGLGroupMem:memberUid,xGLBindDN:,xGLCertValidationEnabled:1,xGLGroup1Priv:511,xGLGroup2Priv:97,xGLGroup3Priv:0,xGLGroup4Priv:0,xGLGroup5Priv:0,xGLServerPort:636
func (i *IDrac8) applyLdapRoleGroupPrivParam(cfg *cfgresources.Ldap) int {

	if cfg.Port == 0 {
		log.WithFields(log.Fields{
			"step": "applyLdapRoleGroupPrivParam",
		}).Warn("Ldap resource parameter Port required but not declared.")
		return 1
	}

	if cfg.Role == "" {
		log.WithFields(log.Fields{
			"step": "applyLdapRoleGroupPrivParam",
		}).Warn("Ldap resource parameter Role required but not declared.")
		return 1
	}

	if cfg.Group == "" {
		log.WithFields(log.Fields{
			"step": "applyLdapRoleGroupPrivParam",
		}).Warn("Ldap resource parameter Group required but not declared.")
		return 1
	}

	if cfg.BaseDn == "" {
		log.WithFields(log.Fields{
			"step": "applyLdapRoleGroupPrivParam",
		}).Warn("Ldap resource parameter BaseDn required but not declared.")
		return 1
	}

	if cfg.UserAttribute == "" {
		log.WithFields(log.Fields{
			"step": "applyLdapRoleGroupPrivParam",
		}).Warn("Ldap resource parameter userAttribute required but not declared.")
		return 1
	}

	if cfg.GroupAttribute == "" {
		log.WithFields(log.Fields{
			"step": "applyLdapRoleGroupPrivParam",
		}).Warn("Ldap resource parameter groupAttribute required but not declared.")
		return 1
	}

	if !i.isRoleValid(cfg.Role) {
		log.WithFields(log.Fields{
			"step": "applyLdapRoleGroupPrivParam",
			"role": cfg.Role,
		}).Warn("Ldap resource Role must be a valid role: admin OR user.")
		return 1
	}

	//511 == full privileges
	privId := "0"
	if cfg.Role == "admin" {
		privId = "511"
	}

	baseDn := escapeLdapString(cfg.BaseDn)
	payload := "data=LDAPEnableMode:3,"  //setup generic ldap
	payload += "xGLNameSearchEnabled:0," //lookup ldap server from dns
	payload += fmt.Sprintf("xGLBaseDN:%s,", baseDn)
	payload += fmt.Sprintf("xGLUserLogin:%s,", cfg.UserAttribute)
	payload += fmt.Sprintf("xGLGroupMem:%s,", cfg.GroupAttribute)
	payload += "xGLBindDN:,xGLCertValidationEnabled:1," //we may want to be able to set this from config
	payload += fmt.Sprintf("xGLGroup1Priv:%s,", privId)
	payload += "xGLGroup2Priv:0,xGLGroup3Priv:0,xGLGroup4Priv:0,xGLGroup5Priv:0," //for now we have just one group.
	payload += "xGLServerPort:636"

	endpoint := "postset?ldapconf"
	responseCode, responseBody, err := i.post(endpoint, []byte(payload))
	if err != nil || responseCode != 200 {
		log.WithFields(log.Fields{
			"IP":           i.ip,
			"endpoint":     endpoint,
			"step":         funcName(),
			"responseCode": responseCode,
			"response":     string(responseBody),
		}).Warn("POST request failed.")
		return 1
	}

	log.WithFields(log.Fields{
		"IP": i.ip,
	}).Info("Ldap Group role privileges applied.")

	return 0
}

func (i *IDrac8) applyTimezoneParam(timezone string) {
	//POST - params as query string
	//timezone
	//https://10.193.251.10/data?set=tm_tz_str_zone:CET

	endpoint := fmt.Sprintf("data?set=tm_tz_str_zone:%s", timezone)
	response, err := i.get(endpoint, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"IP":       i.ip,
			"endpoint": endpoint,
			"step":     funcName(),
			"response": string(response),
		}).Warn("GET request failed.")
	}

	log.WithFields(log.Fields{
		"IP": i.ip,
	}).Info("Timezone param applied.")

}
