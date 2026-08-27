package auth

// LocalAuthUser is a User whose password this platform stores and verifies itself, as opposed to
// one authenticated by an external issuer.
//
// Reached by an unguarded assertion in the login and signup services
// (laatoomodules/security/dev/plugins/dblogin/src/server/go/dbloginservice.go:84,
// signup/src/server/go/accountregister.go:56), so the configured user object must implement it or
// those paths panic. user.DefaultUser does (defaultuser.go:110-116).
type LocalAuthUser interface {
	User

	// GetPassword returns whatever is in the password field right now, which is a different thing
	// at different points in a request — the interface carries no marker saying which:
	//
	//   - on the user built from the login request body, the plaintext the caller supplied;
	//   - on a user freshly constructed for a save, the plaintext about to be hashed;
	//   - on a user read back from storage, the bcrypt hash — or "", see below.
	//
	// Authentication depends on exactly that asymmetry:
	// bcrypt.CompareHashAndPassword(existingUser.GetPassword(), usr.GetPassword())
	// (laatoomodules/security/dev/plugins/common/common.go:84).
	//
	// Two hazards attach to the value:
	//
	// It is re-hashed on EVERY save, not only on create. DefaultUser.PreSave calls
	// bcrypt.GenerateFromPassword over whatever GetPassword returns
	// (defaultuser.go:96-103, :235-242), so saving a user whose field still holds a hash
	// double-hashes it and permanently invalidates the password. That is why PostLoad blanks the
	// field (defaultuser.go:105-109) and why ClearPassword exists.
	//
	// It is a live credential in memory. laatoomodules/security/dev/plugins/common/common.go:83
	// and laatoomodules/scripts/dev/plugins/common/common.go:80 log both the supplied plaintext
	// and the stored hash at Trace level. Never log, serialise or put this value on a claim.
	GetPassword() string

	// ClearPassword blanks the in-memory password field. DefaultUser sets it to ""
	// (defaultuser.go:114-116).
	//
	// Call it as soon as a comparison is done and before the object can be returned to a caller or
	// saved: the authenticated user is handed straight back in the login response
	// (common/common.go:86, then SetAuthToken at :95-124), and PreSave would otherwise re-hash the
	// hash. Because it is only a memory operation it does not touch what is stored — clearing is
	// not deleting the password.
	//
	// Several services reach it structurally rather than through this interface, as
	// entity.(interface{ ClearPassword() }) (firebaseusermgmt/.../FirebaseCreateUserService.go:103,
	// FirebaseUpdateUserService.go:153, FirebaseRegisterAdminService.go:127). Those assertions are
	// comma-ok and skip silently when the user type lacks the method, so a user object that does
	// not implement it leaves password material in a response with nothing logged.
	ClearPassword()
}
