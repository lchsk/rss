import m from "mithril";

const SignUpComponent = require('./signup');
const LoginComponent = require('./login');
const LogoutComponent = require('./logout');
const User = require('./user');

const UserComponent = {
  oninit: node => {
    User.load();
  },
  view: node => {}
};

m.route.prefix("#!");

m.route(document.body, "/", {
  "/": {
    render: function() {}
  },
  "/login": {
    render: function() {
      return m(LoginComponent);
    }
  },
  "/logout": {
    render: function() {
      return m(LogoutComponent);
    }
  },
  "/signup": {
    render: () => {
      return m(SignUpComponent);
    }
  },
  "/user": {
    render: function() {
      return m(UserComponent);
    }
  }
});
