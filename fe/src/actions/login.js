var m = require("mithril");
var { getErrorMessage, getSingleError } = require("./error");

const User = require("./user");
const getLoadingView = require("./loading");

var Login = {
  current: {},
  setError: function(error) {
    this.current.error = error;
  },
  submit: function() {
    return m
      .request({
        method: "POST",
        url: "http://localhost:8000/api/authentication",
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

module.exports = LoginComponent;
