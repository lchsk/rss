import m from "mithril";

const UserAuthState = {
  UNKNOWN: "unknown",
  SIGNED_IN: "signed_in",
  SIGNED_OUT: "signed_out"
};

var User = {
  AuthState: UserAuthState,
  data: {},
  authState: UserAuthState.UNKNOWN,

  load: () => {
    if (User.authState !== UserAuthState.UNKNOWN) {
      return;
    }

    m.request({
      method: "GET",
      url: "http://localhost:8000/api/users/1",
      withCredentials: true
    })
      .then(result => {
        User.data = { u: result };
        User.authState = User.AuthState.SIGNED_IN;
      })
      .catch(e => {
        User.data = { error: e };
        User.authState = User.AuthState.SIGNED_OUT;
      });
  }
};

export { User };
