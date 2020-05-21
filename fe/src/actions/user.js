var m = require("mithril");

const Config = require("../config");

const UserAuthState = {
  UNKNOWN: "unknown",
  SIGNED_IN: "signed_in",
  SIGNED_OUT: "signed_out"
};

var User = {
  AuthState: UserAuthState,
  data: {},
  channels: [],
  channelsByCategory: {},
  authState: UserAuthState.UNKNOWN,

  loadChannels: () => {
    return m
      .request({
        method: "GET",
        url: Config.api_url + "/users/current/channels",
        withCredentials: true
      })
      .then(result => {
        User.channels = result["user_channels"];

        for (let i = 0; i < User.channels.length; i++) {
          const channel = User.channels[i];
          User.channelsByCategory[channel.category_id] = {
            categoryTitle: channel.category_title,
          };
        }
      })
      .catch(e => {
        console.log(e);
      });
  },
  load: () => {
    if (User.authState === UserAuthState.SIGNED_OUT) {
      return;
    } else if (User.authState === UserAuthState.SIGNED_IN) {
      return;
    }

    m.request({
      method: "GET",
      url: Config.api_url + "/users/current",
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
