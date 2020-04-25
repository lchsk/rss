var m = require("mithril");

var Logout = {
  view: node => {
    m.request({
      method: "POST",
      url: "http://localhost:8000/api/logout",
      withCredentials: true
    })
      .then(result => {})
      .catch(e => {});

    m.route.set("/index");
  }
};

module.exports = Logout;