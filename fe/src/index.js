import m from "mithril";

import { SignUpComponent } from "./signup";
import { LoginComponent } from "./login";
import { LogoutComponent } from "./logout";
import { User } from "./user"

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
