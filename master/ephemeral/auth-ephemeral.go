/*Package ephemeral -
Future:
 - Stored in cache/database for a short amount of time.
*/
package ephemeral

// AuthEphemeral - An authentication token
type AuthEphemeral struct {
	Token string
}

// ToMap - Returns a map representation of an authentication ephemeral
func (aT AuthEphemeral) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"authToken": aT.Token,
	}
}
