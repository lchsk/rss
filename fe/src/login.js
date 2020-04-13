import m from "mithril";

import { getErrorMessage, getSingleError } from "./error";

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

var LoginComponent = {
  view: node => {
    return m(
      "form",
      {
        onsubmit: e => {
          e.preventDefault();
          Login.submit();
        }
      },
      [
        m("div#login-error", Login.current.error),
        m("input[type=text][placeholder=Email]", {
          oninput: m.withAttr("value", value => {
            Login.current.email = value;
          })
        }),
        m("input[type=password][placeholder=Password]", {
          oninput: m.withAttr("value", value => {
            Login.current.password = value;
          })
        }),
        m("button[type=submit]", "Sign in")
      ]
    );
  }
};

export var LoginComponent;
