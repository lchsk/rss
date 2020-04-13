import m from "mithril";

import { SignUpComponent } from "./signup";
import { LoginComponent } from "./login";

const UserComponent = {
  oninit: node => {
    m.request({
      method: "GET",
      url: "http://localhost:8000/api/users/1",
      withCredentials: true
    })
      .then(result => {
        console.log(result);
      })
      .catch(e => {
        console.log(e);
      });
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
