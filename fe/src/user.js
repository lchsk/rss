var m = require("mithril");

const UserAuthState = {
  UNKNOWN: "unknown",
  SIGNED_IN: "signed_in",
  SIGNED_OUT: "signed_out"
};

var User = {
  AuthState: UserAuthState,
  data: {},
  channels: [],
  authState: UserAuthState.UNKNOWN,

  loadChannels: () => {
    console.log("Loading channels");
    return m
      .request({
        method: "GET",
        url: "http://localhost:8000/api/users/current/channels",
        withCredentials: true
      })
      .then(result => {
        User.channels = result["user_channels"];
        console.log("channels loaded", result);
        // m.redraw();
        // User.data = { u: result };
        // User.authState = User.AuthState.SIGNED_IN;
      })
      .catch(e => {
        console.log(e);
        // User.data = { error: e };
        // TODO: Check for status
        // User.authState = User.AuthState.SIGNED_OUT;
        // m.route.set("/login");
      });
  },
  load: () => {
    if (User.authState === UserAuthState.SIGNED_OUT) {
      // if (m.route.get() !== "/login" && m.route.get() !== "/signup") {
      // m.route.set("/login");
      // }
      return;
    } else if (User.authState === UserAuthState.SIGNED_IN) {
      return;
    }

    m.request({
      method: "GET",
      url: "http://localhost:8000/api/users/current",
      withCredentials: true
    })
      .then(result => {
        User.data = { u: result };
        User.authState = User.AuthState.SIGNED_IN;
      })
      .catch(e => {
        User.data = { error: e };
        // TODO: Check for status
        User.authState = User.AuthState.SIGNED_OUT;
        // m.route.set("/login");
      });
  }
};

module.exports = User;
