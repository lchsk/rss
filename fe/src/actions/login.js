var m = require("mithril");
var { getErrorMessage, getSingleError } = require("../common/error");

const Config = require("../config");
const { checkAuthAndExtract } = require("./request");

const User = require("./user");

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
        User.authState = User.AuthState.SIGNED_IN;
        this.setError("");
        m.route.set("/");
      })
      .catch(e => {
        User.authState = User.AuthState.SIGNED_OUT;
        this.setError(getSingleError(e.message));
      });
  }
};

module.exports = Login;
