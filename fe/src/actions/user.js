const m = require("mithril");

const Config = require("../config");
const { checkAuthAndExtract } = require("./request");

const UserAuthState = {
  UNKNOWN: "unknown",
  SIGNED_IN: "signed_in",
  SIGNED_OUT: "signed_out"
};

const User = {
  AuthState: UserAuthState,
  data: {},
  channels: [],
  channelsByCategory: {},
  authState: UserAuthState.UNKNOWN,

  defaultSuccess: () => {
    User.authState = User.AuthState.SIGNED_IN;
  },
  loadChannels: () => {
    return m
      .request({
        method: "GET",
        url: Config.api_url + "/users/current/channels",
        withCredentials: true,
        responseType: "json",
        extract: checkAuthAndExtract
      })
      .then(result => {
        User.defaultSuccess();

        const { response } = result;

        User.channels = response["user_channels"];

        for (let i = 0; i < User.channels.length; i++) {
          const channel = User.channels[i];
          User.channelsByCategory[channel]["category_id"] = {
            categoryTitle: channel["category_title"]
          };
        }
      })
      .catch(e => {
          console.log("load channels", e);
      });
  },
  load: () => {
    // if (User.authState === UserAuthState.SIGNED_OUT) {
    // return;
    // }
    // else if (User.authState === UserAuthState.SIGNED_IN) {
    //   return;
    // }

    return m
      .request({
        method: "GET",
        url: Config.api_url + "/users/current",
        withCredentials: true,
        responseType: "json",
        extract: function(xhr, options) {
          if (xhr.status === 401) {
            return false;
          }

          return true;
        }
      })
      .then(result => {
        // User.defaultSuccess();

        const { response } = result;
        User.data = { u: response };
      })
      .catch(e => {
        User.data = { error: e };
      });
  }
};

module.exports = User;
