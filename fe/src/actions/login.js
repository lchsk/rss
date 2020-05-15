var m = require("mithril");
var { getErrorMessage, getSingleError } = require("../common/error");

const Config = require("../config");

const User = require("./user");
// const getLoadingView = require("./loading");

var Login = {
  current: {},
  setError: function(error) {
    this.current.error = error;
  },
  submit: function() {
    return m
      .request({
        method: "POST",
        url: Config.api_url + "/authentication",
        data: Login.current,
        withCredentials: true
      })
      .then(result => {
        this.setError("");
        m.route.set("/index");
      })
      .catch(e => {
        this.setError(getSingleError(e.message));
      });
  }
};

module.exports = Login;
