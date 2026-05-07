package ldap_check

import (
	"axia4/ldap/ldap_conn"
)

func Run(ldapId int32) error {

	ldapConn, _, err := ldap_conn.ConnectAndBind(ldapId)
	if err != nil {
		return err
	}
	defer ldapConn.Close()

	return nil
}
