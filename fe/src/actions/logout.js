var m = require("mithril");

const Config = require("../config");

var Logout = {
  view: node => {
    m.request({
      method: "POST",
      url: Config.api_url + "/logout",
      withCredentials: true
    })
      .then(result => {})
      .catch(e => {});

    m.route.set("/index");
  }
};

module.exports = Logout;
