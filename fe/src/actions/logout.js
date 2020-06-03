var m = require("mithril");

const Config = require("../config");
const User = require("./user");

var Logout = {
  view: node => {
    m.request({
      method: "POST",
      url: Config.api_url + "/logout",
      withCredentials: true,
      responseType: "json",
      extract: checkAuthAndExtract
    })
      .then(result => {
        User.defaultSuccess();
      })
      .catch(e => {
        User.authState = User.AuthState.SIGNED_OUT;
      });

    m.route.set("/index");
  }
};

module.exports = Logout;
