package auth

import "laatoo.io/sdk/datatypes"

// UserAccount is the profile entity hanging off a User — the display-facing half of an identity,
// kept separate from the credential-bearing User record.
//
// NOTHING IN THE PLATFORM IMPLEMENTS THIS INTERFACE. Both account entities that ship —
// useraccount.UserAccount
// (laatoomodules/security/dev/plugins/useraccount/src/server/go/autogen_UserAccount.go) and
// designeruser.UserAccount — are pure codegen entities whose only method beyond ReadAll/WriteAll
// and Config is the GetId inherited from data.StorageInfo. Neither declares Gender or a username
// field at all, so neither could satisfy this interface without a schema change.
//
// The consequence is in User.GetUserAccount: DefaultUser asserts
// usr.Account.Entity.(auth.UserAccount) unguarded (defaultuser.go:143). Today that line is
// unreachable because Entity is never populated, so the method returns nil instead. Expanding the
// account reference is what would turn a nil return into a panic.
//
// Document any new account type against the accessors below before wiring it in; the behaviour
// described for each is the interface's intent, not observed behaviour of a live implementation.
type UserAccount interface {
	datatypes.Serializable

	// GetId returns the account's own storage id.
	//
	// Note the convention clash worth resolving before relying on either: the login path loads
	// the profile with accountService.GetById(ctx, usr.GetId(), "")
	// (laatoomodules/security/dev/plugins/common/common.go:104), which requires the account row's
	// id to BE the user id — while the shipped useraccount.UserAccount entity carries a separate
	// UserId field for the same linkage.
	GetId() string

	// GetEmail returns the email held on the profile.
	//
	// A different field on a different entity from User.GetEmail; nothing keeps the two in sync.
	GetEmail() string

	// GetFullName returns the profile's display name.
	//
	// Its only caller is laatoomodules/security/dev/plugins/mfa/.../mfaemailtask.go:130, which
	// invokes it on the result of User.GetUserAccount() with no nil check — and that result is
	// nil on any token-authenticated request.
	GetFullName() string

	// GetPicture returns a reference to the profile image.
	GetPicture() string

	// GetGender returns the gender recorded on the profile.
	GetGender() string

	// GetUsernameField returns the NAME of the field holding the login name, not the login name —
	// the same contract as User.GetUsernameField, where the mechanism and its silent failure mode
	// are documented.
	GetUsernameField() string

	// GetUserName returns the login name as held on the profile.
	GetUserName() string
}
